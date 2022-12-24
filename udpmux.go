package main

import (
	"log"
	"net"

	"github.com/hana-ame/udptun/utils"
)

// a udpConn mux
type UDPMux struct {
	*net.UDPConn
	dstAddr *net.UDPAddr
	connMap *utils.LockedMap

	portal *Portal

	closed bool
}

func NewUDPMux(listen string, dst string, portal *Portal) *UDPMux {
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
		closed:  false,
	}

	return c
}

func (c *UDPMux) ReadFromPortal(data []byte) {
	if v, ok := c.connMap.Get(string(data[0:2])); ok {
		if fc, ok := v.(*fakeUDPConn); ok {
			// TODO if want to ake portal package edit here
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
	c.portal.router.Put(c.dstAddr.String(), c.ReadFromPortal)
	buf := make([]byte, 1500)
	for !c.closed {
		n, addr, err := c.ReadFromUDP(buf[2:]) // TODO if want to make protal packege edit here
		if err != nil {
			log.Println(err)
			continue
		}
		// tag := itoa.Itoa(addr.Port)
		tag := string([]byte{byte(addr.Port / 256), byte(addr.Port % 256)}) // Big
		if v, ok := c.connMap.Get(tag); ok {
			if fc, ok := v.(*fakeUDPConn); ok {
				fc.WriteToDst(buf[:n+2]) // TODO
			} else {
				log.Println("fakeConn is", fc, " not fakeUDPConn")
				continue
			}
		} else {
			// create new fakeConn
			if fc := NewFakeUDPConn(addr, c.UDPConn, c.dstAddr, c.portal.UDPConn); fc != nil {
				c.connMap.Put(tag, fc)
				fc.WriteToDst(buf[:n+2]) // TODO
			} else {
				log.Println("fakeConn is nil")
				continue
			}
		}
	}
	c.Close()
}

func (c *UDPMux) Close() {
	c.closed = true
	c.portal.router.Remove(c.dstAddr.String())
	c.UDPConn.Close()
}
