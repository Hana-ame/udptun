package mymux

import (
	"github.com/gorilla/websocket"
)

// 这个是放Node之间连接的，其实加上router的话可能没法用
type WebSocketNode struct {
	reading bool
	writing bool
	Conn    *websocket.Conn

	f    Frame
	Node // 假设 Node 是一个定义好的接口或结构体
}

func (n *WebSocketNode) SetConn(c *websocket.Conn) {
	n.Conn = c
}

func (n *WebSocketNode) SetReading(f bool) {
	n.reading = f
}

func (n *WebSocketNode) SetWriting(f bool) {
	n.writing = f
}

func (n *WebSocketNode) ReadCopy() error {
	defer n.SetReading(false)
	n.reading = true
	for n.reading {
		_, p, err := n.Conn.ReadMessage()
		if err != nil {
			return err
		}
		err = n.SendFrame(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *WebSocketNode) WriteCopy() (err error) {
	defer n.SetWriting(false)
	n.writing = true
	for n.writing {
		if n.f == nil {
			n.f, err = n.RecvFrame()
			if err != nil {
				return err
			}
		}

		err = n.Conn.WriteMessage(websocket.BinaryMessage, n.f)
		if err != nil {
			return err
		}
		n.f = nil
	}
	return nil
}
