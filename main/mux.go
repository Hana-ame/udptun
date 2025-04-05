package main

import (
	"fmt"

	tools "github.com/Hana-ame/udptun/Tools"
	"github.com/Hana-ame/udptun/Tools/debug"
)

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

// 单向
type PortMux struct {
	*tools.ConcurrentHashMap[uint8, *PortConn] // port, router interface

	*Pipe
}

func (m *PortMux) Push(f Frame) error {
	if portMux, ok := m.Get(f.Port()); ok {
		return portMux.Push(f)
	}
	return m.Pipe.Push(f)
}

func (m *PortMux) Close() error {
	m.ConcurrentHashMap.ForEach(func(key uint8, value *PortConn) {
		defer value.Close()
	})
	return m.Pipe.Close()
}

func NewPortMux(pipe *Pipe) *PortMux {
	return &PortMux{
		ConcurrentHashMap: tools.NewConcurrentHashMap[uint8, *PortConn](),
		Pipe:              pipe,
	}
}

type PortClient struct {
	local  addr
	remote addr

	port uint8

	*PortMux
	*Pipe
}

func NewPortClient(local, remote addr, mux *PortMux) *PortClient {
	client := &PortClient{
		local:  local,
		remote: remote,

		PortMux: mux,
		Pipe:    NewPipe(),
	}

	handler := func() error {
		defer client.Close()
		for {
			// 所有的f都是map中没有储存过的port响应
			f, err := mux.Poll()
			if err != nil {
				debug.E("client's mux handler", err)
				return client.Close()
			}

			// 不响应已经关闭了的连接的Close请求
			if f.Command() == Close {
				continue
			}

			// 其他情况告知这里已经没有了
			src := f.Source()
			dst := f.Destination()
			if dst != local {
				continue
			}
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
	for !c.PutIfAbsent(conn.port, conn) {
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
	local  addr
	remote addr

	AcceptChan chan *PortConn

	*AddrMux
	*Pipe
}

// 单向
type AddrMux struct {
	*tools.ConcurrentHashMap[addr, *PortMux] // port, router interface

	*Pipe
}

func (m *AddrMux) Push(f Frame) error {
	if portMux, ok := m.Get(f.Source()); ok {
		return portMux.Push(f)
	}
	return m.Pipe.Push(f)
}

func (m *AddrMux) Close() error {
	m.ConcurrentHashMap.ForEach(func(key addr, value *PortMux) {
		defer value.Close()
	})
	return m.Pipe.Close()
}

func NewAddrMux() *AddrMux {
	return &AddrMux{
		ConcurrentHashMap: tools.NewConcurrentHashMap[addr, *PortMux](),
		Pipe:              NewPipe(),
	}
}

func NewPortServer(local, remote addr, mux *AddrMux) *PortServer {
	server := &PortServer{
		local:  local,
		remote: remote,

		AcceptChan: make(chan *PortConn, 5),

		AddrMux: mux,
		Pipe:    NewPipe(),
	}

	handler := func() error {
		defer server.Close()
		for {
			// 处理所有不在map中的请求
			f, err := mux.Poll()
			if err != nil {
				debug.E("client's mux handler", err)
				return server.Close()
			}

			// 不响应已经关闭了的连接的Close请求
			if f.Command() == Close {
				continue
			}

			// 是请求的情况下，进行响应
			if f.Command() == ClientRequest {
				if !mux.Contains(f.Source()) {
					portMux := NewPortMux(server.AddrMux.Pipe)
					server.AddrMux.Put(f.Source(), portMux)
				}
				portMux, ok := mux.Get(f.Source())
				if ok && !portMux.Contains(f.Port()) {
					// 仅在是Request并且map中空缺port的情况下创建
					conn := server.PortConn((f.Source()), f.Port())
					server.AcceptChan <- conn
					server.Pipe.Push(NewFrame(local, remote, f.Port(), ServerAccept, 0, 0, []byte{}))
					portMux.Put(f.Port(), conn)
				}
			}

			// 其他情况告知这里已经没有了
			src := f.Source()
			dst := f.Destination()
			if dst != local {
				continue
			}
			err = server.Pipe.Push(NewFrame(dst, src, f.Port(), Close, 0, 0, []byte{}))
			if err != nil {
				debug.E("server's mux handler", err)
				return server.Close()
			}

		}
	}

	go handler()

	return server
}
func (s *PortServer) PortConn(src addr, port uint8) *PortConn {
	portMux, ok := s.Get(src)
	if !ok {
		return nil
	}
	return &PortConn{
		FramePushHandler: s.ApplicatonInterface(),
		PortMux:          portMux,
		Pipe:             NewPipe(),
		port:             port,
	}
}

func (s *PortServer) Accept() (*PortConn, error) {
	conn := <-s.AcceptChan
	return conn, nil
}

func (s *PortServer) Close() error {
	s.AddrMux.Close()
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
		push:  s.AddrMux.Push,
		poll:  s.Pipe.Poll,
		close: s.Close,
	}
}
