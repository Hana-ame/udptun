package mymux

import (
	"encoding/binary"
	"io"
	"net"

	"github.com/Hana-ame/udptun/Tools/debug"
)

// 这个是放Node之间连接的，其实加上router的话可能没法用
type TCPNode struct {
	reading bool
	writing bool
	Conn    net.Conn

	f    Frame
	Node // 假设 Node 是一个定义好的接口或结构体
}

func (n *TCPNode) SetConn(c net.Conn) {
	n.Conn = c
}

func (n *TCPNode) SetReading(f bool) {
	n.reading = f
}

func (n *TCPNode) SetWriting(f bool) {
	n.writing = f
}

func (n *TCPNode) ReadCopy() error {
	defer n.SetReading(false)
	n.reading = true
	for n.reading {
		size := make([]byte, 2)
		_, err := n.Conn.Read(size)
		if err != nil {
			return err
		}

		buffer := make([]byte, binary.BigEndian.Uint16(size))
		_, err = n.Conn.Read(buffer)
		if err != nil {
			return err
		}

		err = n.SendFrame(buffer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *TCPNode) WriteCopy() (err error) {
	defer n.SetWriting(false)
	n.writing = true
	for n.writing {
		if n.f == nil {
			n.f, err = n.RecvFrame()
			if err != nil {
				return err
			}
		}

		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, uint16(len(n.f)))

		_, err = n.Conn.Write(size)
		if err != nil {
			return err
		}
		_, err = n.Conn.Write(n.f)
		if err != nil {
			return err
		}
		n.f = nil
	}
	return nil
}

func NewTCPListener(addr string, node *Node, dst, src Addr) (net.Listener, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	go func() {
		defer listener.Close()
		var port byte = 0
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			var bus Bus
			for bus == nil {
				port++
				bus, err = node.Dial(dst, src, port)
				debug.W("listener", err)
			}
			fc := NewFrameConn(bus, src, dst, port)
			go Copy(fc, conn)
			go Copy(conn, fc)
		}
	}()

	return listener, nil
}

// 直接在node里面用、
func NewTCPDialer(addr string, node *Node, dst, src Addr) error {
	for {
		bus, port := node.Accept()
		fc := NewFrameConn(bus, src, dst, port)
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return err
		}
		go Copy(fc, conn)
		go Copy(conn, fc)
	}
	return nil
}

func Copy(dst io.WriteCloser, src io.ReadCloser) {
	defer dst.Close()
	defer src.Close()

	io.Copy(dst, src)
}
