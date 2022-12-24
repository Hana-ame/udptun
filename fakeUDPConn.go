package main

import (
	"log"
	"net"
	"time"
)

type fakeUDPConn struct {
	srcAddr *net.UDPAddr
	srcConn *net.UDPConn

	dstAddr *net.UDPAddr
	dstConn *net.UDPConn

	lastactivity int64
	timeout      int64
	closed       bool
	close        func()
}

func NewFakeUDPConn(
	srcAddr *net.UDPAddr, srcConn *net.UDPConn,
	dstAddr *net.UDPAddr, dstConn *net.UDPConn,
	tag string, timeout int64, close func(),
) *fakeUDPConn {
	fc := &fakeUDPConn{
		srcAddr:      srcAddr,
		srcConn:      srcConn,
		dstAddr:      dstAddr,
		dstConn:      dstConn,
		lastactivity: time.Now().Unix(),
		timeout:      timeout,
		closed:       false,
		close:        close,
	}

	go fc.Run()

	return fc
}

// raw
// to portal, with tag
func (c *fakeUDPConn) WriteToDst(b []byte) (int, error) {
	c.lastactivity = time.Now().Unix()
	// log.Println("WriteToDst ", len(b), c.dstConn.LocalAddr().String(), "->", c.dstAddr.String())
	return c.dstConn.WriteToUDP(b, c.dstAddr)
}

// raw
// to udp, only data
func (c *fakeUDPConn) WriteToSrc(b []byte) (int, error) {
	c.lastactivity = time.Now().Unix()
	// log.Println("WriteToSrc ", len(b), c.srcConn.LocalAddr().String(), "->", c.srcAddr.String())
	return c.srcConn.WriteToUDP(b, c.srcAddr)
}

func (c *fakeUDPConn) Run() {
	for {
		time.Sleep(10 * time.Second)
		if time.Now().Unix()-c.lastactivity > c.timeout {
			c.Close()
			return
		}
	}
}

func (c *fakeUDPConn) Close() {
	log.Println("fakeUDPConn close")
	c.closed = true
	c.close()
}
