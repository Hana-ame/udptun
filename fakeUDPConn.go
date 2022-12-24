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
	closed       bool
}

func NewFakeUDPConn(srcAddr *net.UDPAddr, srcConn *net.UDPConn, dstAddr *net.UDPAddr, dstConn *net.UDPConn) *fakeUDPConn {
	return &fakeUDPConn{
		srcAddr:      srcAddr,
		srcConn:      srcConn,
		dstAddr:      dstAddr,
		dstConn:      dstConn,
		lastactivity: time.Now().Unix(),
		closed:       false,
	}
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

func (c *fakeUDPConn) Close() {
	c.closed = true
}
