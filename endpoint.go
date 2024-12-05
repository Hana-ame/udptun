package main

import (
	"encoding/binary"
	"io"
	"sync"

	"github.com/gorilla/websocket"
)

type IOReadWriteCloserEndpoint struct {
	io.ReadWriteCloser
}

func (e *IOReadWriteCloserEndpoint) Push(f Frame) error {
	// debug.I("push", SprintFrame(f))
	l := make([]byte, 2)
	binary.BigEndian.PutUint16(l, uint16(len(f)))
	if _, err := e.Write(l); err != nil {
		return err
	}
	if _, err := e.Write(f); err != nil {
		return err
	}
	return nil
}

func (e *IOReadWriteCloserEndpoint) Poll() (Frame, error) {
	l := make([]byte, 2)
	if _, err := e.Read(l); err != nil {
		return nil, err
	}
	v := binary.BigEndian.Uint16(l)
	// debug.I("poll", "v", v)
	f := make(Frame, v)
	if _, err := e.Read(f); err != nil {
		return nil, err
	}
	// debug.I("poll", SprintFrame(f))
	return f, nil
}

type WebSocketReader struct {
	*websocket.Conn

	sync.Mutex
}

func (e *WebSocketReader) Poll() (Frame, error) {
	e.Lock()
	defer e.Unlock()

	_, f, err := e.ReadMessage()

	return f, err
}

type WebSocketWriter struct {
	*websocket.Conn

	sync.Mutex
}

func (e *WebSocketWriter) Push(f Frame) error {
	e.Lock()
	defer e.Unlock()

	err := e.WriteMessage(websocket.BinaryMessage, f)

	return err
}

type WebSocketEndpoint struct {
	WebSocketReader
	WebSocketWriter
}

func (e *WebSocketEndpoint) Close() error {
	re := e.WebSocketReader.Close()
	we := e.WebSocketWriter.Close()
	if re != nil {
		return re
	}
	if we != nil {
		return we
	}
	return nil
}

func NewWebSocketEndpoint(c *websocket.Conn) *WebSocketEndpoint {
	return &WebSocketEndpoint{
		WebSocketReader: WebSocketReader{Conn: c},
		WebSocketWriter: WebSocketWriter{Conn: c},
	}
}
