package main

import (
	"fmt"

	tools "github.com/Hana-ame/udptun/Tools"
	"github.com/Hana-ame/udptun/Tools/debug"
)

type Mux struct {
	*tools.ConcurrentHashMap[uint8, FramePushCloserHandler] // port, router interface

	*Pipe
}

func (m *Mux) Push(f Frame) error {
	return m.GetOrDefault(f.Port(), m.Pipe).Push(f)
}

func (m *Mux) Close() error {
	m.ConcurrentHashMap.ForEach(func(key uint8, value FramePushCloserHandler) {
		defer value.Close()
	})
	return m.Pipe.Close()
}

func NewMux() *Mux {
	return &Mux{
		ConcurrentHashMap: tools.NewConcurrentHashMap[uint8, FramePushCloserHandler](),
		Pipe:              NewPipe(),
	}
}

type ClientConn struct {
	FramePushHandler

	*Mux

	*Pipe

	port uint8

	requested    bool
	acknowledged bool
	closed       bool
}

// 不要调用多次
// 调用多次可能的情况是最开始的几个frame会丢失掉
func (c *ClientConn) Request() error {
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

func (c *ClientConn) Close() error {
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

func (c *ClientConn) Poll() (Frame, error) {
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

func (c *ClientConn) Push(f Frame) error {
	f.SetPort(c.port)
	return c.FramePushHandler.Push(f)
}

func (c *ClientConn) ApplicatonInterface() FrameHandlerInterface {
	return FrameHandlerInterface{
		push:  c.Push,
		poll:  c.Poll,
		close: c.Close,
	}
}
func (c *ClientConn) RouterInterface() FramePushCloserHandler {
	return FrameHandlerInterface{
		push:  c.Pipe.Push,
		poll:  nil, // 直接推到Client/Server
		close: c.Close,
	}
}

type Client struct {
	local  Addr
	remote Addr

	port uint8

	*Mux
	*Pipe
}

func NewClient(local, remote Addr, mux *Mux) *Client {
	client := &Client{
		local:  local,
		remote: remote,

		Mux:  mux,
		Pipe: NewPipe(),
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

func (c *Client) ClientConn() *ClientConn {
	if c.port == 0 {
		c.port++
	}
	return &ClientConn{
		FramePushHandler: c.ApplicatonInterface(),
		Mux:              c.Mux,
		Pipe:             NewPipe(),
		port:             c.port,
	}
}

func (c *Client) Dial() (*ClientConn, error) {
	conn := c.ClientConn()
	for !c.PutIfAbsent(conn.port, conn.RouterInterface()) {
		c.port++
		conn = c.ClientConn()
	}

	if err := conn.Request(); err != nil {
		defer conn.Close()
		return nil, err
	}

	return conn, nil
}
func (c *Client) Close() error {
	c.Mux.Close()
	return c.Pipe.Close()
}

func (c *Client) Push(f Frame) error {
	f.SetSource(c.local)
	f.SetDestination(c.remote)
	return c.Pipe.Push(f)
}

func (c *Client) ApplicatonInterface() FramePushHandler {
	return FrameHandlerInterface{
		push:  c.Push,
		poll:  nil, // 从Mux得到map直连
		close: nil, // 不需要
	}
}
func (c *Client) RouterInterface() FrameHandler {
	return FrameHandlerInterface{
		push:  c.Mux.Push,
		poll:  c.Pipe.Poll,
		close: c.Close,
	}
}

type Server struct {
	local  Addr
	remote Addr

	AcceptChan chan *ClientConn

	*Mux
	*Pipe
}

func NewServer(local, remote Addr, mux *Mux) *Server {
	server := &Server{
		local:  local,
		remote: remote,

		AcceptChan: make(chan *ClientConn, 5),

		Mux:  mux,
		Pipe: NewPipe(),
	}

	handler := func() error {
		for {
			f, err := mux.Poll()
			if err != nil {
				debug.E("client's mux handler", err)
				return server.Close()
			}
			// 不响应已经关闭了的连接的Close请求
			if f.Command() == Close {
				continue
			}

			// 可能需要修改一下
			if f.Command() == ClientRequest {
				port := f.Port()
				if mux.Contains(port) {
					err = server.Pipe.Push(NewFrame(local, remote, f.Port(), Close, 0, 0, []byte{}))

					if err != nil {
						debug.E("client's mux handler", err)
						return server.Close()
					}
				} else { // 可以accept
					conn := server.ClientConn(port)
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
func (s *Server) ClientConn(port uint8) *ClientConn {
	return &ClientConn{
		FramePushHandler: s.ApplicatonInterface(),
		Mux:              s.Mux,
		Pipe:             NewPipe(),
		port:             port,
	}
}

func (s *Server) Accept() (*ClientConn, error) {
	conn := <-s.AcceptChan
	return conn, nil
}

func (s *Server) Close() error {
	s.Mux.Close()
	return s.Pipe.Close()
}

func (s *Server) Push(f Frame) error {
	f.SetSource(s.local)
	f.SetDestination(s.remote)
	return s.Pipe.Push(f)
}

func (s *Server) ApplicatonInterface() FramePushHandler {
	return FrameHandlerInterface{
		push:  s.Push,
		poll:  nil, // Mux直连
		close: nil, // 待定
	}
}
func (s *Server) RouterInterface() FrameHandler {
	return FrameHandlerInterface{
		push:  s.Mux.Push,
		poll:  s.Pipe.Poll,
		close: s.Close,
	}
}
