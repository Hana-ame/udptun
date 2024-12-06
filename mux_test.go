package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/Hana-ame/udptun/Tools/debug"
)

func clientCopy(src, dst FrameHandler) {
	defer src.Close()
	defer dst.Close()
	for {
		f, err := src.Poll()
		if err != nil {
			debug.E("copy", err)
			break
		}
		// debug.I("copy", SprintFrame(f))
		if f.Command() == ClientRequest {
			f.SetCommand(ServerAccept)
		}
		if f.Command() == Close {
			f.SetPort(f.Port() ^ 1)
		}
		err = dst.Push(f)
		if err != nil {
			debug.E("copy", err)
			break
		}
	}
}

func send(c FrameHandler, tag string) {
	for {
		e := c.Push(NewFrame(0, 0, 0, Data, 0, 0, []byte(tag)))
		if e != nil {
			debug.E("recv", e)
			return
		}
		time.Sleep(time.Second)
	}
}

func recv(c FrameHandler, tag string) {
	for {
		f, e := c.Poll()
		if e != nil {
			debug.E("recv", e)
			return
		}
		debug.I(tag, SprintFrame(f))
	}
}

func conn(client *Client, tag string) *ClientConn {
	c, _ := client.Dial()
	go send(c, tag)
	go recv(c, tag)
	return c
}

func TestClient(t *testing.T) {
	mux := NewMux()
	client := NewClient(0, 0, mux)
	// debug.T("client", client)
	// rif := client.RouterInterface()
	// f, e := rif.Poll()
	// if e != nil {
	// 	debug.E("rif", e)
	// }
	// PrintFrame(f)
	go clientCopy(client.RouterInterface(), client.RouterInterface())
	// debug.E("Copy", err)

	conn(client, "1")
	c2 := conn(client, "2")
	conn(client, "3")
	c4 := conn(client, "4")
	conn(client, "5")

	c2.Close()
	time.Sleep(5 * time.Second)
	c4.Close()

	// fmt.Println(client.Mux)

	time.Sleep(time.Second * 10)
}

func testServer(server *Server) {
	for {
		// c1, _ := server.Accept()
		// send(c1, "server send")
		// recv(c1, "server recv")
		// c2, _ := server.Accept()
		// recv(c2, "server recv")
		c, _ := server.Accept()
		go Copy(c, c)
	}
}

func TestClientServer(t *testing.T) {
	smux := NewMux()
	server := NewServer(33, 44, smux)
	go testServer(server)

	cmux := NewMux()
	client := NewClient(222, 111, cmux)

	go Copy(client.RouterInterface(), server.RouterInterface())
	go Copy(server.RouterInterface(), client.RouterInterface())

	conn(client, "1")
	c2 := conn(client, "2")
	conn(client, "3")
	c4 := conn(client, "4")
	conn(client, "5")

	c2.Close()
	time.Sleep(5 * time.Second)
	c4.Close()

	fmt.Println(client.Mux)

	time.Sleep(time.Second * 10)

}
