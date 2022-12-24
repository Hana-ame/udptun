package main

import (
	"log"
	"net"

	"github.com/hana-ame/udptun/utils"
)

// a udpConn mux
type UDPMux struct {
	// listen on this conn
	*net.UDPConn

	// the address of the other portal
	dstAddr *net.UDPAddr

	// map[tag string]fc fakeUDPConn.
	//
	// tag is both: (1) the port of local conn, (2) the first 2 bytes of portalBuf
	connMap *utils.LockedMap

	// the portal it use
	portal *Portal

	// is closed
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

	go c.Run()

	return c
}

func (c *UDPMux) ReadFromPortal(buf PortalBuf) {
	// dst --> portal --> portal -X-> src
	tag := string(buf.Tag())
	if v, ok := c.connMap.Get(tag); ok {
		if fc, ok := v.(*fakeUDPConn); ok {
			// buf is what recv from Portal
			fc.WriteToSrc(buf.Data(0))
		} else {
			log.Println("value not *fakeUDPConn")
		}
	} else {
		log.Println("connMap do not have key:", buf.Tag())
	}
}

// will only recv from local
func (c *UDPMux) Run() {
	c.portal.router.Put(c.dstAddr.String(), c.ReadFromPortal)
	buf := make(PortalBuf, 1500)
	for !c.closed {
		// read data from local conn
		// src -X-> portal -X-> portal --> dst
		n, addr, err := c.ReadFromUDP(buf.Data(0))
		if err != nil {
			log.Println(err)
			continue
		}
		// tag := itoa.Itoa(addr.Port)
		// tag := ([]byte{byte(addr.Port / 256), byte(addr.Port % 256)}) // Big
		buf.AddTag(addr.Port)
		tag := string(buf.Tag())

		if v, ok := c.connMap.Get(tag); ok {
			if fc, ok := v.(*fakeUDPConn); ok {
				fc.WriteToDst(buf.DataAndTag(n))
			} else {
				log.Println("fakeConn is", fc, " not fakeUDPConn")
				continue
			}
		} else {
			// create new fakeConn
			if fc := NewFakeUDPConn(
				addr, c.UDPConn,
				c.dstAddr, c.portal.UDPConn,
				tag, 5, func() {
					c.connMap.Remove(tag)
				}); fc != nil {
				c.connMap.Put(tag, fc)
				fc.WriteToDst(buf.DataAndTag(n))
			} else {
				log.Println("fakeConn is nil")
				continue
			}
		}
	}
	c.Close()
}

// TODO
func (c *UDPMux) Close() {
	if c.closed {
		return
	}
	c.closed = true
	c.UDPConn.Close()
	c.portal.router.Remove(c.dstAddr.String())
}
