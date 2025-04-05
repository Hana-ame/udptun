package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebsocketTunnel struct {
	sync.Mutex
	*websocket.Conn
}

func (t *WebsocketTunnel) ReadFrame() (frame, error) {
	t.Lock()
	defer t.Unlock()
	_, message, err := t.Conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (t *WebsocketTunnel) WriteFrame(f frame) error {
	t.Lock()
	defer t.Unlock()
	return t.Conn.WriteMessage(websocket.BinaryMessage, f)
}

func (t *WebsocketTunnel) Close() error {
	t.Lock()
	defer t.Unlock()
	if t.Conn == nil {
		return nil
	}
	err := t.Conn.Close()
	t.Conn = nil
	return err
}

func NewWebsocketTunnel(conn *websocket.Conn) (*WebsocketTunnel, error) {
	// conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	// if err != nil {
	// 	return nil, err
	// }

	return &WebsocketTunnel{
		Mutex: sync.Mutex{},
		Conn:  conn,
	}, nil
}
