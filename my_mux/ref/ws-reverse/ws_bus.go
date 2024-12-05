// 没用到

package wsreverse

import (
	"time"

	"github.com/Hana-ame/udptun/Tools/debug"
	mymux "github.com/Hana-ame/udptun/Tools/my_mux"
	"github.com/gorilla/websocket"
)

// [0, id]  接收前
// [data]
// [1，id]  这个是response

const (
	CLIENT_RECV_SHOULD_BE_ID = iota
	CLIENT_RECV_SHOULD_BE_DATA
)

func NewWsClient(dst string, timeout time.Duration) mymux.Bus {
	const Tag = "NewWsClient"
	// cbus will return to client mux to use
	cbus, sbus := mymux.NewPipeBusPair()

	dialer := func() *websocket.Conn {
		// 	const Tag = "dialer"
		for {
			conn, _, err := websocket.DefaultDialer.Dial(dst, nil)
			if err == nil {
				return conn
			}
			debug.E(Tag, err.Error())
		}
	}
	// create conn from websocket.conn
	conn := NewConn(dialer())

	var size uint8 = 64
	buffer := mymux.NewGBNBuffer(uint8(size))

	// onerr := false

	// read channel
	go func() {
		const Tag = "client read channel"
		var next uint8 = 0
		var nextState = CLIENT_RECV_SHOULD_BE_ID
	loop:
		for {
			msgType, data, err := conn.ReadMessage()
			if err != nil {
				debug.E(Tag, err.Error())
				conn.SetConn(dialer())
				continue
			}
			switch msgType {
			case websocket.TextMessage:
				switch data[0] {
				case 1: // it is a acknowledge
					buffer.SetRead(data[1]) // should be moded by size
					buffer.SetTail(data[1]) // should be moded by size
				case 0: //
					if nextState != CLIENT_RECV_SHOULD_BE_ID {
						conn.WriteMessage(websocket.TextMessage, []byte{1, next})
						nextState = CLIENT_RECV_SHOULD_BE_ID
						continue loop
					}
					if next != data[1] {
						conn.WriteMessage(websocket.TextMessage, []byte{1, next})
						nextState = CLIENT_RECV_SHOULD_BE_ID
						continue loop
					}
					nextState = CLIENT_RECV_SHOULD_BE_DATA
				}
			case websocket.BinaryMessage:
				if nextState != CLIENT_RECV_SHOULD_BE_DATA {
					conn.WriteMessage(websocket.TextMessage, []byte{1, next})
					nextState = CLIENT_RECV_SHOULD_BE_ID
					continue loop
				}
				e := sbus.SendFrame(data)
				if e != nil {
					debug.E(Tag, e.Error())
					continue
				}

				nextState = CLIENT_RECV_SHOULD_BE_ID
				next++

				conn.WriteMessage(websocket.TextMessage, []byte{1, next})
			}
		}
	}()

	// write channel
	// bus side
	go func() {
		const Tag = "client write channel bus side"
		for {
			f, e := sbus.RecvFrame()
			if e != nil {
				debug.E(Tag, e.Error())
				continue
			}

			buffer.Offer(f)
		}
	}()

	// conn side
	go func() {
		const Tag = "client write channel bus side"
		for {
			id, data, e := buffer.Read()
			if !e {
				debug.E(Tag, id, data, "not existed")
			}
			conn.WriteMessage(websocket.TextMessage, []byte{0, id})
			conn.WriteMessage(websocket.BinaryMessage, data)
		}
	}()

	return cbus
}

func NewWsServer(conn *Conn) mymux.Bus {
	const Tag = "NewWsServer"
	// cbus will return to client mux to use
	cbus, sbus := mymux.NewPipeBusPair()

	var size uint8 = 64
	buffer := mymux.NewGBNBuffer(uint8(size))

	// read channel
	go func() {
		const Tag = "server read channel"
		var next uint8 = 0
		var nextState = CLIENT_RECV_SHOULD_BE_ID
	loop:
		for {
			msgType, data, err := conn.ReadMessage()
			if err != nil {
				debug.E(Tag, err.Error())
				conn.WaitOnError()
				continue
			}
			switch msgType {
			case websocket.TextMessage:
				switch data[0] {
				case 1: // it is a acknowledge
					buffer.SetRead(data[1]) // should be moded by size
					buffer.SetTail(data[1]) // should be moded by size
				case 0: //
					if nextState != CLIENT_RECV_SHOULD_BE_ID {
						conn.WriteMessage(websocket.TextMessage, []byte{1, next})
						nextState = CLIENT_RECV_SHOULD_BE_ID
						continue loop
					}
					if next != data[1] {
						conn.WriteMessage(websocket.TextMessage, []byte{1, next})
						nextState = CLIENT_RECV_SHOULD_BE_ID
						continue loop
					}
					nextState = CLIENT_RECV_SHOULD_BE_DATA
				}
			case websocket.BinaryMessage:
				if nextState != CLIENT_RECV_SHOULD_BE_DATA {
					conn.WriteMessage(websocket.TextMessage, []byte{1, next})
					nextState = CLIENT_RECV_SHOULD_BE_ID
					continue loop
				}
				e := sbus.SendFrame(data)
				if e != nil {
					debug.E(Tag, e.Error())
					continue
				}

				nextState = CLIENT_RECV_SHOULD_BE_ID
				next++

				conn.WriteMessage(websocket.TextMessage, []byte{1, next})
			}
		}
	}()

	// write channel
	// bus side
	go func() {
		const Tag = "client write channel bus side"
		for {
			f, e := sbus.RecvFrame()
			if e != nil {
				debug.E(Tag, e.Error())
				continue
			}

			buffer.Offer(f)
		}
	}()

	// conn side
	go func() {
		const Tag = "client write channel bus side"
		for {
			id, data, e := buffer.Read()
			if !e {
				debug.E(Tag, id, data, "not existed")
			}
			conn.WriteMessage(websocket.TextMessage, []byte{0, id})
			conn.WriteMessage(websocket.BinaryMessage, data)
		}
	}()

	return cbus
}
