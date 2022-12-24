package main

import (
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
	return c.dstConn.WriteToUDP(b, c.dstAddr)
}

// raw
// to udp, only data
func (c *fakeUDPConn) WriteToSrc(b []byte) (int, error) {
	c.lastactivity = time.Now().Unix()
	return c.srcConn.WriteToUDP(b, c.srcAddr)
}

// unused
func (c *fakeUDPConn) Run() {
	for {
		time.Sleep(time.Second)
		if time.Now().Unix()-c.lastactivity > c.timeout {
			c.closed = true
			c.close()
			return
		}
	}
}

// unused
// func (c *fakeUDPConn) Close() {
// 	c.closed = true
// }
