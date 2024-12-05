package wsmux

import (
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type WsMuxConn struct {
	sync.Mutex

	*WsMux

	ID       uint16 // WsMux.SeqN
	SeqN     uint16
	ReadChan chan *WsPackage

	MTU int

	closed bool
}

func NewWsConn(id uint16, w *WsMux) *WsMuxConn {
	conn := &WsMuxConn{
		WsMux:    w,
		ID:       id,
		ReadChan: make(chan *WsPackage, 32),

		MTU: 1024,
	}

	return conn
}

func (c *WsMuxConn) PutPackage(pkg *WsPackage) bool {
	select {
	case c.ReadChan <- pkg:
		return true
	default:
		return false
	}
}

// 喷了这边怎么实现EOF啊
func (c *WsMuxConn) ReadPackage() *WsPackage {
	return <-c.ReadChan
}
func (c *WsMuxConn) WritePackage(pkg *WsPackage) error {
	if c.closed {
		return fmt.Errorf("WsMuxConn is closed")
	}

	if pkg == nil {
		return fmt.Errorf("pkg is nil")
		// pkg = &WsPackage{ID: c.ID, SeqN: c.SeqN, Message: []byte{}}
	}
	err := c.WriteMessage(websocket.BinaryMessage, pkg.ToBytes())

	return err
}

func (c *WsMuxConn) Read(p []byte) (n int, err error) {
	pkg := c.ReadPackage()
	// log.Println("read", len(pkg.Message)) // debug
	if len(pkg.Message) == 0 {
		err = io.EOF
		c.Close()
	}
	return copy(p, pkg.Message), err
}
func (c *WsMuxConn) Write(p []byte) (n int, err error) {
	c.Lock()
	// log.Println("write", len(p), p[:min(len(p), 10)])
	defer c.Unlock()
	pkg := &WsPackage{
		ID:      c.ID,
		SeqN:    c.SeqN,
		Message: p,
	}
	err = c.WritePackage(pkg)
	if err != nil {
		c.SeqN++
	}
	return len(pkg.Message), err
}

func (c *WsMuxConn) Close() error {
	c.Lock()
	log.Println("close")
	defer c.Unlock()

	if c.closed {
		return nil
	}
	c.WsMux.DeleteConn(c.ID)
	c.WritePackage(&WsPackage{ID: c.ID, SeqN: 0, Message: []byte{}})
	c.closed = true
	return nil
}
