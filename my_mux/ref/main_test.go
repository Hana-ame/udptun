package mymux

import (
	"testing"
	"time"

	"github.com/Hana-ame/udptun/Tools/debug"
)

// func TestXxx(t *testing.T) {
// 	a2bReader, b2bWriter := io.Pipe()
// 	b2bReader, b2aWriter := io.Pipe()
// 	b2aReader, a2aWriter := io.Pipe()
// 	a2aReader, a2bWriter := io.Pipe()

// 	go func() {
// 		buf := make([]byte, 1500)
// 		for {
// 			n, _ := b2bReader.Read(buf)
// 			// log.Printf("==========pipe b:")
// 			// PrintFrame(buf[:n])
// 			b2bWriter.Write(buf[:n])
// 		}
// 	}()
// 	go func() {
// 		buf := make([]byte, 1500)
// 		for {
// 			n, _ := a2aReader.Read(buf)
// 			// log.Printf("==========pipe a:")
// 			// PrintFrame(buf[:n])
// 			a2aWriter.Write(buf[:n])
// 		}
// 	}()

// 	aBus := NewReaderWriterBus(a2bReader, a2bWriter)
// 	bBus := NewReaderWriterBus(b2aReader, b2aWriter)

// 	aMux := NewMuxServer(aBus, 5)
// 	go aMux.ReadDaemon(aBus)

// 	bMux := NewMuxClient(bBus, 0)
// 	go bMux.ReadDaemon(bBus)

// 	go handleServer(aMux)
// 	time.Sleep(3 * time.Second)

// 	go handleClient(bMux)
// 	time.Sleep(9 * time.Second)
// 	go handleClient(bMux)
// 	var a uint
// 	_ = a
// 	time.Sleep(60 * time.Second)
// }

// var handleServer = func(server *MyMuxServer) {
// 	go server.ReadDaemon()
// 	log.Println("handleServer")
// 	for {
// 		c := server.Accept()
// 		go handleServerConn(c)
// 	}
// }

// var handleClient = func(client *MyMuxClient) {
// 	log.Println("handleClient")
// 	c, e := client.Dial(5)
// 	if e != nil {
// 		debug.E("handleClient", e.Error())
// 	}
// 	go func() {
// 		buf := make([]byte, 1500)
// 		for {
// 			n, e := c.Read(buf)
// 			if e != nil {
// 				debug.E("handleClient", e.Error())
// 			}
// 			debug.I("handleClient", c.Tag(), n, "client recv", string(buf[:n]))
// 		}
// 	}()
// 	c.Close()
// 	for i := 0; i < 5; i++ {
// 		// i := -1
// 		c.Write([]byte(fmt.Sprintf("来自client %d", i)))
// 		time.Sleep(time.Second)
// 	}
// 	time.Sleep(time.Minute)
// }

// var handleAcceptedConn = func(c *MyFrameConn) {
// 	for {
// 		f, e := c.ReadFrame()
// 		if e != nil {
// 			debug.E("handleAcceptedConn", e.Error())
// 			if ErrorIsClosed(e){
// 				return
// 			}
// 			continue
// 		}
// 		debug.I("serve recv", string(f))
// 	}
// }

// handleServerConn := func(c *MyConn) {
// 	log.Println("handleConn")
// 	go func() {
// 		buf := make([]byte, 1500)
// 		for {
// 			n, e := c.Read(buf)
// 			if e != nil {
// 				debug.E("handleConn", e.Error())
// 			}

// 			debug.I("handleServerConn", c.Tag(), n, "server recv", string(buf[:n]))

// 			c.Write([]byte(fmt.Sprintf("反弹 %s", buf[:n])))
// 		}
// 	}()

//		for i := 0; i < 5; i++ {
//			// i := -1
//			c.Write([]byte(fmt.Sprintf("来自server %d", i)))
//			time.Sleep(time.Second)
//		}
//		time.Sleep(time.Minute)
//	}
func TestClient(t *testing.T) {
	handleClientConn := func(c *MyFrameConn) {
		debug.T("handleClientConn", "")
		for {
			f, e := c.ReadFrame()
			if e != nil {
				debug.E("handleClientConn", e.Error())
				if ErrorIsClosed(e) {
					return
				}
				continue
			}
			debug.I("client recv", string(f))
		}
	}

	handleServer := func(server *MyServer) {
		const Tag = "handleServer"
		debug.T(Tag, "")

		handleAcceptedConn := func(c *MyFrameConn) {
			const Tag = "handleAcceptedConn"
			debug.T(Tag, "")
			for {
				f, e := c.ReadFrame()
				if e != nil {
					debug.E(Tag, e.Error())
					if ErrorIsClosed(e) {
						return
					}
					continue
				}
				debug.I("serve recv", string(f))
			}
		}
		go server.ReadDeamon()
		for {
			c := server.Accpet()
			go handleAcceptedConn(c)

			time.Sleep(time.Second)
			c.WriteFrame([]byte("from server 1"))
			time.Sleep(time.Second)
			c.WriteFrame([]byte("from server 2"))
		}
	}

	handleClient := func(client *MyClient) {
		const Tag = "handleClient"
		debug.T(Tag, "")
		go client.ReadDaemon()
	}

	// 这是一对bus，至少应该是正常传输的。
	cb, sb := NewDebugPipeBusPair("bus")
	// 这是server，listen在bus上然后地址是0
	server := NewServer(sb, 0)
	go handleServer(server)

	client := NewClient(cb, 2)
	go handleClient(client)

	dialAndEcho := func(client *MyClient) {
		const Tag = "dialAndEcho"
		c, e := client.Dial(0)
		if e != nil {
			debug.E(Tag, e.Error())
			t.Error(e)
		}
		go handleClientConn(c)
		// time.Sleep(time.Second)
		_, e = c.WriteFrame([]byte("from client 11"))
		if e != nil {
			t.Error(e)

		}
		_, e = c.WriteFrame([]byte("from client 12"))
		if e != nil {
			t.Error(e)

		}
	}
	go dialAndEcho(client)
	// go dialAndEcho(client)
	// go dialAndEcho(client)

	time.Sleep(time.Minute)
}

// func TestClient(t *testing.T) {
// 	cb, sb := NewBusPipe()
// 	server := NewServer(sb, 0)
// 	go server.ReadDeamon()
// 	go func() {
// 		for {
// 			c := server.Accpet()
// 			go handleAcceptedConn(c)

// 			time.Sleep(time.Second)
// 			c.WriteFrame([]byte("from server 1"))
// 			time.Sleep(time.Second)
// 			c.WriteFrame([]byte("from server 2"))
// 		}
// 	}()
// 	{
// 		client := NewClient(cb, 1, 2, 3)
// 		go handleClientConn(client)
// 		client.WriteFrame([]byte("from client 11"))
// 		client.WriteFrame([]byte("from client 12"))
// 	}
// 	{
// 		client := NewClient(cb, 1, 2, 4)
// 		go handleClientConn(client)
// 		client.WriteFrame([]byte("from client 21"))
// 		client.WriteFrame([]byte("from client 22"))
// 		client.Close()
// 		client.WriteFrame([]byte("from client 25"))

// 	}
// 	time.Sleep(time.Minute)

// }

func TestClient2(t *testing.T) {
	handleClientConn := func(c *MyFrameConn) {
		debug.T("handleClientConn", "initial")
		for {
			f, e := c.ReadFrame()
			if e != nil {
				debug.E("handleClientConn", e.Error())
				if ErrorIsClosed(e) {
					return
				}
				continue
			}
			debug.I("client recv", string(f))
		}
	}

	handleServer := func(server *MyServer) {
		const Tag = "handleServer"
		debug.T(Tag, "initial")

		handleAcceptedConn := func(c *MyFrameConn) {
			const Tag = "handleAcceptedConn"
			debug.T(Tag, "initial")
			for {
				f, e := c.ReadFrame()
				if e != nil {
					debug.E(Tag, e.Error())
					if ErrorIsClosed(e) {
						return
					}
					continue
				}
				debug.I("serve recv", string(f))

				n, e := c.WriteFrame(f)
				if e != nil {
					debug.E(Tag, n, e.Error())
					if ErrorIsClosed(e) {
						return
					}
					continue
				}
			}
		}
		go server.ReadDeamon()
		for {
			c := server.Accpet()
			go handleAcceptedConn(c)

			time.Sleep(time.Second)
			c.WriteFrame([]byte("from server 1"))
			time.Sleep(time.Second)
			c.WriteFrame([]byte("from server 2"))
		}
	}

	handleClient := func(client *MyClient) {
		const Tag = "handleClient"
		debug.T(Tag, "initial")
		go client.ReadDaemon()
	}

	// 这是一对bus，至少应该是正常传输的。
	cb, sb := NewDebugPipeBusPair("bus")
	// 这是server，listen在bus上然后地址是0
	server := NewServer(sb, 0)
	go handleServer(server)

	client := NewClient(cb, 1)
	go handleClient(client)

	dialAndEcho := func(client *MyClient, append string) {
		const Tag = "dialAndEcho"
		c, e := client.Dial(0)
		if e != nil {
			debug.E(Tag, e.Error())
			t.Error(e)
		}
		go handleClientConn(c)
		// time.Sleep(time.Second)
		_, e = c.WriteFrame([]byte("from client 1" + append))
		if e != nil {
			t.Error(e)

		}
		_, e = c.WriteFrame([]byte("from client 2" + append))
		if e != nil {
			t.Error(e)

		}
	}
	go dialAndEcho(client, "aaa")
	go dialAndEcho(client, "bbb")
	// go dialAndEcho(client)
	// go dialAndEcho(client)

	time.Sleep(time.Minute)
}

func TestClient3(t *testing.T) {
	handleClientConn := func(c *MyFrameConn) {
		debug.T("handleClientConn", "initial")
		for {
			f, e := c.ReadFrame()
			if e != nil {
				debug.E("handleClientConn", e.Error())
				if ErrorIsClosed(e) {
					return
				}
				continue
			}
			debug.I("client recv", string(f))
		}
	}

	handleServer := func(server *MyServer) {
		const Tag = "handleServer"
		debug.T(Tag, "initial")

		handleAcceptedConn := func(c *MyFrameConn) {
			const Tag = "handleAcceptedConn"
			debug.T(Tag, "initial")
			for {
				f, e := c.ReadFrame()
				if e != nil {
					debug.E(Tag, e.Error())
					if ErrorIsClosed(e) {
						return
					}
					continue
				}
				debug.I("serve recv", string(f))

				n, e := c.WriteFrame(f)
				if e != nil {
					debug.E(Tag, n, e.Error())
					if ErrorIsClosed(e) {
						return
					}
					continue
				}
			}
		}
		go server.ReadDeamon()
		for {
			c := server.Accpet()
			go handleAcceptedConn(c)
			go func() {
				const Tag = "ServerSending"
				time.Sleep(time.Second)
				debug.T(Tag, "closed?", c.closed)
				_, e := c.WriteFrame([]byte("from server 1"))
				if e != nil {
					debug.E(Tag, e.Error())
					return
				}
				time.Sleep(time.Second)
				debug.T(Tag, "closed?", c.closed)
				_, e = c.WriteFrame([]byte("from server 2"))
				if e != nil {
					debug.E(Tag, e.Error())
					return
				}
				time.Sleep(time.Second)
				debug.T(Tag, "closed?", c.closed)
				_, e = c.WriteFrame([]byte("from server 3"))
				if e != nil {
					debug.E(Tag, e.Error())
					return
				}
				time.Sleep(time.Second)
				debug.T(Tag, "closed?", c.closed)
				_, e = c.WriteFrame([]byte("from server 4"))
				if e != nil {
					debug.E(Tag, e.Error())
					return
				}
			}()
		}
	}

	handleClient := func(client *MyClient) {
		const Tag = "handleClient"
		debug.T(Tag, "initial")
		go client.ReadDaemon()
	}

	// 这是一对bus，至少应该是正常传输的。
	cb, sb := NewDebugPipeBusPair("bus")
	// 这是server，listen在bus上然后地址是0
	server := NewServer(sb, 0)
	go handleServer(server)

	client := NewClient(cb, 1)
	go handleClient(client)

	dialAndEcho := func(client *MyClient, append string) {
		const Tag = "dialAndEcho"
		c, e := client.Dial(0)
		if e != nil {
			debug.E(Tag, e.Error())
			t.Error(e)
		}
		go handleClientConn(c)
		time.Sleep(time.Second)
		// _, e = c.WriteFrame([]byte("from client 1" + append))
		// if e != nil {
		// 	debug.E(Tag, e.Error())
		// 	return
		// }
		// time.Sleep(time.Second)
		// _, e = c.WriteFrame([]byte("from client 2" + append))
		// if e != nil {
		// 	debug.E(Tag, e.Error())
		// 	return
		// }
		// time.Sleep(time.Second)
		c.Close()
		// _, e = c.WriteFrame([]byte("from client 3" + append))
		// if e != nil {
		// 	debug.E(Tag, e.Error())
		// 	return
		// }
		// _, e = c.WriteFrame([]byte("from client 4" + append))
		// if e != nil {
		// 	debug.E(Tag, e.Error())
		// 	return
		// }
	}
	// go dialAndEcho(client, "aaa")
	go dialAndEcho(client, "bbb")
	// go dialAndEcho(client)
	// go dialAndEcho(client)

	time.Sleep(time.Minute)
}
