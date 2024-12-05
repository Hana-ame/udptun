package remoteclient

// import (
// 	"io"
// 	"log"
// 	"net"
// 	"net/http"

// 	"github.com/Hana-ame/udptun/Tools/debug"
// 	mymux "github.com/Hana-ame/udptun/Tools/my_mux"
// 	wsreverse "github.com/Hana-ame/udptun/Tools/my_mux/example/ws-reverse"
// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/websocket"
// )

// // // 定义一个客户端结构体
// // type Client struct {
// // 	conn *websocket.Conn
// // }

// // // 定义一个客户端集合
// // var clients = make(map[*Client]bool)
// // var broadcast = make(chan string)

// var Conn = wsreverse.NewConn(
// 	nil,
// )

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// // 处理WebSocket连接
// func HandleWebSocket(c *gin.Context) {
// 	const Tag = "HandleWebSocket"
// 	// 升级HTTP连接为WebSocket连接
// 	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		debug.E(Tag, err.Error())
// 	}
// 	// defer ws.Close()

// 	Conn.SetConn(ws)
// 	debug.I(Tag, "accept a new conn and set")
// }

// func Server(conn *wsreverse.Conn) mymux.Bus {
// 	const Tag = "ws server"
// 	cbus, sbus := mymux.NewPipeBusPair()
// 	// conn -> bus
// 	go func() {
// 		for {
// 			_, data, err := conn.ReadMessage()
// 			if err != nil {
// 				debug.E(Tag, err.Error())
// 				conn.WaitOnError()
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

// func init() {
// 	const Tag = "remote-client"
// 	bus := Server(Conn)

// 	var addr mymux.Addr = 5
// 	client := mymux.NewClient(bus, addr)

// 	listener, err := net.Listen("tcp", "127.24.10.4:8080")
// 	debug.E(Tag, err.Error())

// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			debug.E(Tag, err.Error())
// 			continue
// 		}

// 		muxc, err := client.Dial(5)
// 		if err != nil {
// 			conn.Close()
// 			debug.E(Tag, err.Error())
// 			continue
// 		}

// 		go handle(muxc, conn)

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
