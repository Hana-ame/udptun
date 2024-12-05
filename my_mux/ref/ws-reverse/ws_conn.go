package wsreverse

import (
	"sync"

	"github.com/gorilla/websocket"
)

// 会无限重试websocket.Conn
// 通过再loader中定义新Conn的生成方式
// 取代websocket.Conn的位置，表现为websocket.Conn

type ConnWriter struct {
	*websocket.Conn
	sync.Mutex
}

func (w *ConnWriter) WriteMessage(messageType int, data []byte) error {
	// const Tag = "ConnWriter.WriteMessage"
	w.Lock()
	defer w.Unlock()

	return w.Conn.WriteMessage(messageType, data)
}

type ConnReader struct {
	*websocket.Conn
	sync.Mutex
}

func (r *ConnReader) ReadMessage() (messageType int, data []byte, err error) {
	// const Tag = "ConnWriter.WriteMessage"
	r.Lock()
	defer r.Unlock()

	return r.Conn.ReadMessage()
}

// 带锁的websockt.Conn，不能同时写或者同时读
// 可以使用setConn重置
type Conn struct {
	*ConnReader
	*ConnWriter

	// accept or connect ws
	// loader func() *websocket.Conn

	*sync.Cond

	onError bool
	closed  bool
}

func NewConn(c *websocket.Conn) *Conn {
	return &Conn{
		ConnReader: &ConnReader{Conn: c},
		ConnWriter: &ConnWriter{Conn: c},
		Cond:       sync.NewCond(&sync.Mutex{}),
	}
}

// 大概有问题。
func (c *Conn) WaitOnError() {
	c.L.Lock()
	c.onError = true
	for !(!c.onError || c.closed) {
		c.Wait()
	}
	c.L.Unlock()
}

func (c *Conn) SetConn(conn *websocket.Conn) {
	c.L.Lock()
	for !(c.onError || c.closed) {
		c.Wait()
	}
	c.ConnWriter.Conn.Close()
	c.ConnReader.Conn.Close()

	c.ConnReader = &ConnReader{Conn: conn}
	c.ConnWriter = &ConnWriter{Conn: conn}

	c.onError = false

	c.L.Unlock()
	c.Broadcast()
}

func (c *Conn) Close() {
	c.closed = true

	c.ConnWriter.Conn.Close()
	c.ConnReader.Conn.Close()

	c.Broadcast()
}
