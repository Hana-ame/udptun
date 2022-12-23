package main

import (
	"log"
	"net"
	"time"

	"github.com/hana-ame/udptun/utils"
)

// a udpConn mux
type UDPMux struct {
	*net.UDPConn
	dstAddr *net.UDPAddr
	connMap *utils.LockedMap

	portal *portal
}

func newUDPMux(listen string, dst string, portal *portal) *UDPMux {
	addr, err := net.ResolveUDPAddr("udp", listen)
	if err != nil {
		panic(err)
	}
	dstAddr, err := net.ResolveUDPAddr("udp", dst)
	if err != nil {
		panic(err)
	}
	lc, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

	c := &UDPMux{
		UDPConn: lc,
		dstAddr: dstAddr,
		connMap: utils.NewLockedMap(),
		portal:  portal,
	}

	portal.router.Put(dst, c.ReadFromPortal)

	return c
}

func (c *UDPMux) ReadFromPortal(data []byte) {
	if v, ok := c.connMap.Get(string(data[0:2])); ok {
		if fc, ok := v.(*fakeUDPConn); ok {
			fc.WriteToSrc(data)
		} else {
			log.Println("value not *fakeUDPConn")
		}
	} else {
		log.Println("connMap do not have key:", data[0:2])
	}
}

// will only recv from local
func (c *UDPMux) Run() {
	for {
		buf := make([]byte, 1500)
		n, addr, err := c.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		// tag := itoa.Itoa(addr.Port)
		tag := string([]byte{byte(addr.Port / 256), byte(addr.Port % 256)}) // Big
		if v, ok := c.connMap.Get(tag); ok {
			if fc, ok := v.(*fakeUDPConn); ok {
				fc.WriteToDst(buf[:n])
			} else {
				log.Println("fakeConn is", fc, " not fakeUDPConn")
				continue
			}
		} else {
			// create new fakeConn
			if fc := NewFakeUDPConn(addr, c.UDPConn, c.dstAddr, c.portal.UDPConn); fc != nil {
				c.connMap.Put(tag, fc)
				fc.WriteToDst(buf[:n])
			} else {
				log.Println("fakeConn is nil")
				continue
			}
		}
	}
}

func (c *UDPMux) Close() {
	c.portal.router.Remove(c.dstAddr.String())
	c.UDPConn.Close()
}

type fakeUDPConn struct {
	srcAddr *net.UDPAddr
	srcConn *net.UDPConn

	dstAddr *net.UDPAddr
	dstConn *net.UDPConn

	lastactivity int64
}

func NewFakeUDPConn(srcAddr *net.UDPAddr, srcConn *net.UDPConn, dstAddr *net.UDPAddr, dstConn *net.UDPConn) *fakeUDPConn {
	return &fakeUDPConn{
		srcAddr:      srcAddr,
		srcConn:      srcConn,
		dstAddr:      dstAddr,
		dstConn:      dstConn,
		lastactivity: time.Now().Unix(),
	}
}

// do here
func (c *fakeUDPConn) WriteToDst(b []byte) (int, error) {
	c.lastactivity = time.Now().Unix()
	return c.dstConn.WriteToUDP(b, c.dstAddr)
}

func (c *fakeUDPConn) WriteToSrc(b []byte) (int, error) {
	c.lastactivity = time.Now().Unix()
	return c.srcConn.WriteToUDP(b, c.srcAddr)
}
