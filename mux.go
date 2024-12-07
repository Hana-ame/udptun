package main

import (
	"fmt"

	tools "github.com/Hana-ame/udptun/Tools"
	"github.com/Hana-ame/udptun/Tools/debug"
)

type PortMux struct {
	*tools.ConcurrentHashMap[uint8, FramePushCloserHandler] // port, router interface

	*Pipe
}

func (m *PortMux) Push(f Frame) error {
	return m.GetOrDefault(f.Port(), m.Pipe).Push(f)
}

func (m *PortMux) Close() error {
	m.ConcurrentHashMap.ForEach(func(key uint8, value FramePushCloserHandler) {
		defer value.Close()
	})
	return m.Pipe.Close()
}

func NewPortMux() *PortMux {
	return &PortMux{
		ConcurrentHashMap: tools.NewConcurrentHashMap[uint8, FramePushCloserHandler](),
		Pipe:              NewPipe(),
	}
}

type PortConn struct {
	FramePushHandler

	*PortMux

	*Pipe

	port uint8

	requested    bool
	acknowledged bool
	closed       bool
}

// 不要调用多次
// 调用多次可能的情况是最开始的几个frame会丢失掉
func (c *PortConn) Request() error {
	if c.requested {
		return nil
	}
	c.requested = true
	if c.acknowledged {
		return nil
	}
	if err := c.FramePushHandler.Push(NewFrame(0, 0, c.port, ClientRequest, 0, 0, []byte{})); err != nil {
		return err
	}
	for f, err := c.Poll(); !c.acknowledged; f, err = c.Poll() {
		if err != nil {
			return err
		}
		if f.Command() == ServerAccept {
			c.acknowledged = true
			break
		}
	}
	return nil
}

func (c *PortConn) Close() error {
	if c.closed {
		return nil
	}
	// 关闭Client
	c.Remove(c.port)
	// 关闭peer
	if err := c.FramePushHandler.Push(NewFrame(0, 0, c.port, Close, 0, 0, []byte{})); err != nil {
		return err
	}

	return c.Pipe.Close()
}

func (c *PortConn) Poll() (Frame, error) {
	f, err := c.Pipe.Poll()
	if err != nil {
		return f, err
	}
	if f.Command() == Close {
		defer c.Close()
		return f, fmt.Errorf("client conn: receive close")
	}
	return f, err
}

func (c *PortConn) Push(f Frame) error {
	f.SetPort(c.port)
	return c.FramePushHandler.Push(f)
}

// 是完整解耦合的，并且可以直接用。
func (c *PortConn) ApplicatonInterface() FrameHandlerInterface {
	return FrameHandlerInterface{
		push:  c.Push,
		poll:  c.Poll,
		close: c.Close,
	}
}
func (c *PortConn) RouterInterface() FramePushCloserHandler {
	return FrameHandlerInterface{
		push:  c.Pipe.Push,
		poll:  nil, // 耦合的，直推，不能从Mux取，会乱。
		close: c.Close,
	}
}

type PortClient struct {
	local  Addr
	remote Addr

	port uint8

	*PortMux
	*Pipe
}

func NewPortClient(local, remote Addr, mux *PortMux) *PortClient {
	client := &PortClient{
		local:  local,
		remote: remote,

		PortMux: mux,
		Pipe:    NewPipe(),
	}

	handler := func() error {
		for {
			// debug.T("clinet's mux handler", mux)
			f, err := mux.Poll()
			if err != nil {
				debug.E("client's mux handler", err)
				return client.Close()
			}
			// 不响应已经关闭了的连接的Close请求
			if f.Command() == Close {
				continue
			}
			src := f.Source()
			dst := f.Destination()
			// debug.T("clinet's mux handler", client.Pipe)
			err = client.Pipe.Push(NewFrame(dst, src, f.Port(), Close, 0, 0, []byte{}))
			if err != nil {
				debug.E("client's mux handler", err)
				return client.Close()
			}
		}
	}

	go handler()
	return client
}

func (c *PortClient) PortConn() *PortConn {
	if c.port == 0 {
		c.port++
	}
	return &PortConn{
		FramePushHandler: c.ApplicatonInterface(),
		PortMux:          c.PortMux,
		Pipe:             NewPipe(),
		port:             c.port,
	}
}

func (c *PortClient) Dial() (*PortConn, error) {
	conn := c.PortConn()
	for !c.PutIfAbsent(conn.port, conn.RouterInterface()) {
		c.port++
		conn = c.PortConn()
	}

	if err := conn.Request(); err != nil {
		defer conn.Close()
		return nil, err
	}

	return conn, nil
}
func (c *PortClient) Close() error {
	c.PortMux.Close()
	return c.Pipe.Close()
}

func (c *PortClient) Push(f Frame) error {
	f.SetSource(c.local)
	f.SetDestination(c.remote)
	return c.Pipe.Push(f)
}

func (c *PortClient) ApplicatonInterface() FramePushHandler {
	return FrameHandlerInterface{
		push:  c.Push,
		poll:  nil, // 从Mux得到map直连
		close: nil, // 不需要
	}
}

// 是完整解耦合的
func (c *PortClient) RouterInterface() FrameHandler {
	return FrameHandlerInterface{
		push:  c.PortMux.Push,
		poll:  c.Pipe.Poll,
		close: c.Close,
	}
}

type PortServer struct {
	local  Addr
	remote Addr

	AcceptChan chan *PortConn

	*PortMux
	*Pipe
}

func NewPortServer(local, remote Addr, mux *PortMux) *PortServer {
	server := &PortServer{
		local:  local,
		remote: remote,

		AcceptChan: make(chan *PortConn, 5),

		PortMux: mux,
		Pipe:    NewPipe(),
	}

	handler := func() error {
		for {
			// 处理所有不在
			f, err := mux.Poll()
			if err != nil {
				debug.E("client's mux handler", err)
				return server.Close()
			}
			// 不响应已经关闭了的连接的Close请求
			if f.Command() == Close {
				continue
			}

			// 是请求的情况下。
			if f.Command() == ClientRequest {
				port := f.Port()
				if mux.Contains(port) {
					// 如果port还存在就返回不让。
					err = server.Pipe.Push(NewFrame(local, remote, f.Port(), Close, 0, 0, []byte{}))
					if err != nil {
						debug.E("client's mux handler", err)
						return server.Close()
					}
				} else {
					// 如果可以创建那就创建了。
					conn := server.PortConn(port)
					server.Pipe.Push(NewFrame(local, remote, f.Port(), ServerAccept, 0, 0, []byte{}))
					server.AcceptChan <- conn
					server.Put(port, conn.RouterInterface())
				}
			} else if f.Command() != Close {
				err = server.Pipe.Push(NewFrame(local, remote, f.Port(), Close, 0, 0, []byte{}))

				if err != nil {
					debug.E("client's mux handler", err)
					return server.Close()
				}

			}
		}
	}

	go handler()

	return server
}
func (s *PortServer) PortConn(port uint8) *PortConn {
	return &PortConn{
		FramePushHandler: s.ApplicatonInterface(),
		PortMux:          s.PortMux,
		Pipe:             NewPipe(),
		port:             port,
	}
}

func (s *PortServer) Accept() (*PortConn, error) {
	conn := <-s.AcceptChan
	return conn, nil
}

func (s *PortServer) Close() error {
	s.PortMux.Close()
	return s.Pipe.Close()
}

func (s *PortServer) Push(f Frame) error {
	f.SetSource(s.local)
	f.SetDestination(s.remote)
	return s.Pipe.Push(f)
}

func (s *PortServer) ApplicatonInterface() FramePushHandler {
	return FrameHandlerInterface{
		push:  s.Push,
		poll:  nil, // Mux直连
		close: nil, // 待定
	}
}

// 是完整解耦合的
func (s *PortServer) RouterInterface() FrameHandler {
	return FrameHandlerInterface{
		push:  s.PortMux.Push,
		poll:  s.Pipe.Poll,
		close: s.Close,
	}
}
