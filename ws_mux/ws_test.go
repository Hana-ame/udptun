package wsmux

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestClient(t *testing.T) {
	// WebSocket 服务器的 URL
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8080", Path: "/ws"}
	log.Println(u)

	// 建立 WebSocket 连接
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	log.Println("dial")

	// 发送消息
	msg := "Hello, server!"
	err = c.WriteMessage(websocket.BinaryMessage, []byte(msg))
	if err != nil {
		log.Println("write:", err)
	}
	log.Println("write")

	// 接收消息
	// for {
	_, message, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		// break
	}
	log.Printf("recv: %s\n", message)
	time.Sleep(1 * time.Second)
	// }

}

func TestMux(t *testing.T) {
	// WebSocket 服务器的 URL
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8080", Path: "/ws"}
	log.Println(u)

	// 建立 WebSocket 连接
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	log.Println("dial")

	mux := NewWsMux(c, MuxSeqClient)
	go mux.ReadDaemon(c)

	go func() {
		c := mux.Accept()
		pkg := c.ReadPackage()
		log.Println(pkg)
	}()

	// 发送消息
	subConn, _ := mux.Dial()
	subConn.Write([]byte("你是不是傻逼"))
	// 接收消息
	pkg := subConn.ReadPackage()
	log.Printf("recv: %s\n", pkg.Message)
	time.Sleep(1 * time.Second)
	// }

}

func TestMuxMul(t *testing.T) {
	// WebSocket 服务器的 URL
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8080", Path: "/ws"}
	log.Println(u)

	// 建立 WebSocket 连接
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	log.Println("dial")

	mux := NewWsMux(c, MuxSeqClient)
	go mux.ReadDaemon(c)

	go func() {
		c := mux.Accept()
		pkg := c.ReadPackage()
		log.Println(pkg)
	}()

	go dialSendRecv(mux, []byte("1"), 5)
	go dialSendRecv(mux, []byte("2"), 5)
	go dialSendRecv(mux, []byte("3"), 5)

	time.Sleep(time.Second * 200)
}

func dialSendRecv(mux *WsMux, msg []byte, times int) {
	c, _ := mux.Dial()
	go func() {
		for {
			pkg := c.ReadPackage()
			log.Println("recv:", pkg, msg)
		}
	}()
	for i := 0; i < times; i++ {
		c.Write(msg)
		log.Println("send:", i, msg)
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}
}

func TestCopy(t *testing.T) {
	// 打开源文件
	srcFile, err := os.Open("source.txt")
	if err != nil {
		fmt.Println("打开源文件失败:", err)
		return
	}
	defer srcFile.Close()

	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8080", Path: "/ws"}
	log.Println(u)

	// 建立 WebSocket 连接
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	log.Println("dial")

	mux := NewWsMux(c, MuxSeqClient)
	go mux.ReadDaemon(c)

	go func() {
		for {
			c := mux.Accept()
			go func(c *WsMuxConn) {
				// 创建目标文件
				dstFile, err := os.Create("destination" + strconv.Itoa(int(c.ID)) + ".txt")
				if err != nil {
					fmt.Println("创建目标文件失败:", err)
					return
				}
				defer dstFile.Close()
				buffer := make([]byte, 1024)
				io.CopyBuffer(dstFile, c, buffer)
			}(c)
		}
	}()

	cs := make([]*WsMuxConn, 3)
	for i := 0; i < 3; i++ {
		go func(i int) {
			cs[i], _ = mux.Dial()
			buf := make([]byte, 1024) // 1KB缓冲区
			file, err := os.Open("source.txt")
			if err != nil {
				fmt.Println("打开文件失败:", err)
				return
			}
			defer file.Close()

			io.CopyBuffer(cs[i], file, buf)
		}(i)
	}
	// 创建缓冲区
	// 循环读取并写入

	fmt.Println("文件复制完成")

	time.Sleep(time.Minute * 2)

}
