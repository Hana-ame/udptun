package main

// import (
// 	"flag"
// 	"io"
// 	"log"
// 	"net"

// 	"github.com/Hana-ame/udptun/Tools/debug"
// 	mymux "github.com/Hana-ame/udptun/Tools/my_mux"
// 	wsreverse "github.com/Hana-ame/udptun/Tools/my_mux/example/ws-reverse"
// 	"github.com/gorilla/websocket"
// )

// var Conn *wsreverse.Conn

// func Client(conn *wsreverse.Conn, dst string) mymux.Bus {
// 	const Tag = "ws client"
// 	cbus, sbus := mymux.NewPipeBusPair()
// 	// conn -> bus
// 	go func() {
// 		for {
// 			_, data, err := conn.ReadMessage()
// 			if err != nil {
// 				debug.E(Tag, err.Error())
// 				var ws *websocket.Conn
// 				for err != nil {
// 					ws, _, err = websocket.DefaultDialer.Dial(dst, nil)
// 				}
// 				conn.SetConn(ws)
// 				continue
// 			}
// 			sbus.SendFrame(data)
// 		}
// 	}()
// 	// bus -> conn
// 	go func() {
// 		for {
// 			f, e := sbus.RecvFrame()
// 			if e != nil {
// 				debug.E(Tag, e.Error())
// 			}

// 			e = conn.WriteMessage(websocket.BinaryMessage, f)
// 			for e != nil {
// 				debug.E(Tag, e.Error())
// 				conn.WaitOnError()
// 				e = conn.WriteMessage(websocket.BinaryMessage, f)
// 			}
// 		}
// 	}()
// 	return cbus
// }

// func main() {
// 	// 定义命令行参数
// 	wsUrl := flag.String("ws", "ws://file.moonchan.xyz/ws/server", "")
// 	localUrl := flag.String("l", "localhost:8080", "")

// 	// 解析命令行参数
// 	flag.Parse()
// 	ws, _, err := websocket.DefaultDialer.Dial(*wsUrl, nil)
// 	if err != nil {
// 		debug.F("main", err.Error())
// 	}
// 	Conn = wsreverse.NewConn(ws)

// 	bus := Client(Conn, *wsUrl)

// 	rb := mymux.NewReliableBus(bus, 64)

// 	server := mymux.NewServer(rb, 5)
// 	go server.ReadDeamon()

// 	for {
// 		nc := server.Accpet()

// 		lc, e := net.Dial("tcp", *localUrl)
// 		if e != nil {
// 			nc.Close()
// 			debug.E("dial tcp", *localUrl, " ", e.Error())
// 		}

// 		handle(nc, lc)
// 	}
// }

// func handle(nc *mymux.FrameConn, lc net.Conn) {
// 	sc := &mymux.MyFrameConnStreamer{MyFrameConn: nc}

// 	// 从 lc 读取数据并写入到 sc
// 	go func() {
// 		defer lc.Close() // 确保连接在结束时关闭
// 		if _, err := io.Copy(sc, lc); err != nil {
// 			log.Println("Error copying from lc to sc:", err)
// 			return
// 		}
// 	}()

// 	// 从 sc 读取数据并写入到 lc
// 	go func() {
// 		defer nc.Close() // 确保 nc 在结束时关闭
// 		if _, err := io.Copy(lc, sc); err != nil {
// 			log.Println("Error copying from sc to lc:", err)
// 			return
// 		}
// 	}()

// }
