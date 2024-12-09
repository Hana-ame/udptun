package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"testing"
	"time"

	"github.com/Hana-ame/udptun/Tools/debug"
)

func DebugCopy(dst, src FrameHandler) {
	defer src.Close()
	defer dst.Close()
	for {
		f, err := src.Poll()
		if err != nil {
			debug.E("copy", err)
			time.Sleep(time.Second * 5)
			continue
		}
		err = dst.Push(f)
		if err != nil {
			debug.E("copy", err)
			time.Sleep(time.Second * 5)
			continue
		}
	}
}

// startServer 启动一个简单的 TCP 服务器
func startServer() {
	ln, err := net.Listen("tcp", "0.0.0.0:7890")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer ln.Close()
	fmt.Println("Server listening on localhost:7890")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

// handleConnection 处理客户端连接
func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("client accepted")
	// 读取数据
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading from connection:", err)
			return
		}

		// fmt.Printf("Server received: %d\n", n)

		// 回送相同的数据
		_, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Println("Error writing to connection:", err)
		}
	}
}

func TestTCPEndpoin(t *testing.T) {
	go startServer()

	time.Sleep(time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:7890")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("client connected", err)

	message := []byte("Hello, TCP Server!")
	fmt.Printf("Client sending: %s\n", message)

	// 写入数据到服务器
	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("Error writing to server:", err)
		return
	}

	// 读取服务器的响应
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from server:", err)
		return
	}

	fmt.Printf("Client received: %s\n", string(buffer[:n]))

}

func TestTCPEndpointWithConn(t *testing.T) {
	go startServer()

	time.Sleep(time.Second)

	conn, err := net.Dial("tcp", "localhost:7890")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	ep := &IOReadWriteCloserEndpoint{ReadWriteCloser: conn}

	c := NewConn()

	go DebugCopy(c.RouterInterface(), ep)
	go DebugCopy(ep, c.RouterInterface())

	writer := FrameWriter{FrameHandler: c.ApplicatonInterface(), MTU: 100000}
	reader := FrameReader{FrameHandler: c.ApplicatonInterface()}

	cnt := 50
	go func() {
		for i := 0; i < cnt; i++ {
			writer.Write([]byte("1234567890qwertyuiopasdfghjklzxcvbnm"))
		}
	}()
	go func() {
		buf := make([]byte, 1000)
		for i := 0; i < cnt; i++ {
			n, _ := reader.Read(buf)
			fmt.Printf("%s\n", buf[:n])
		}
	}()
	time.Sleep(3 * time.Second)

}
