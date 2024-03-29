package main

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/hana-ame/udptun/utils"
)

type Portal struct {
	*net.UDPConn
	// stun
	// stunServer string
	localAddr string

	// portal -> src
	//
	// map [addr.String()] func(PortalBuf)
	// where addr is remote portal's address
	// func(PortalBuf) is provided by udpmux.
	router *utils.LockedMap

	// portal -> dst
	// if dst is not nil, it will be used to send udp packets to dst
	dst *net.UDPAddr
	// map [addr.String() + tag] *fakeUDPConn
	// send data for Conn here
	connMap *utils.LockedMap
}

// "" means not accept remote
func NewPortal(dst string) *Portal {
	// addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:4444") // !!!!debug
	// c, err := net.ListenUDP("udp", addr)                   // !!!!debug
	c, err := net.ListenUDP("udp", nil)
	for err != nil {
		log.Println("what?", err)
		time.Sleep(2 * time.Second)
		c, err = net.ListenUDP("udp", nil)
	}

	var dstAddr *net.UDPAddr = nil
	if dst != "" {
		dstAddr, err = net.ResolveUDPAddr("udp", dst)
		if err != nil {
			log.Fatal(err)
			return nil
		}
	}
	p := &Portal{
		UDPConn:   c,
		localAddr: "",
		router:    utils.NewLockedMap(),
		dst:       dstAddr,
		connMap:   utils.NewLockedMap(),
	}

	go p.Run()

	return p
}

func (p *Portal) Ping(arr []string) {
	for _, v := range arr {
		addr, err := net.ResolveUDPAddr("udp", v)
		if err != nil {
			log.Println(v, " can not resolve ", err)
			continue
		}
		p.UDPConn.WriteToUDP([]byte{}, addr)
	}
}

// stun.
// empty stunServer for IPv6
// for IPv4, first stunServer is "stun.l.google.com:19302".
// and set stunServer to "udp4" for latest result
// for IPv6, set the stunServer to "udp" or "udp6"
func (p *Portal) GetLocalAddr(stunServer string) string {
	if p.localAddr == "" && !strings.HasPrefix(stunServer, "udp") {
		go func() {
			for {
				utils.StunRequest(stunServer, p.UDPConn)
				time.Sleep(30 * time.Second)
			}
		}()
		time.Sleep(5 * time.Second)
	} else if stunServer != "udp4" {
		t := utils.GetOutboundIPv6(p.UDPConn)
		if t != "" {
			p.localAddr = t
		}
	}
	return p.localAddr
}

// go
func (p *Portal) Run() {
	// portal -X-> stun
	stunServer := "142.251.2.127:19302" // google
	p.GetLocalAddr(stunServer)

	// src --> portal -X-> portal --> dst
	// dst --> portal -X-> portal --> src
	for {
		buf := make(PortalBuf, 1500)
		// fmt.Println("reading")
		n, addr, err := p.ReadFromUDP(buf)
		// if n == 0 {
		// 	continue
		// }
		// when reading error
		if err != nil {
			log.Println(err)
			continue
		}
		// addrString is the address of remote Portal
		addrString := addr.String()
		if addrString == stunServer {
			// stun -X-> portal
			p.localAddr, err = utils.StunResolve(buf.Raw(n))
			if err != nil {
				log.Println("error when recv from stunServer", err)
			}
		} else if v, ok := p.router.Get(addrString); ok {
			// dst --> portal -X-> portal --> src
			if handler, ok := v.(func(PortalBuf)); ok {
				// handler is UDPMux.ReadFromPortal
				handler(buf.Raw(n))
			} else {
				log.Println("invalid router") // never
			}
		} else if p.dst != nil {
			// src --> portal -X-> portal -X-> dst
			if n < tagLength {
				log.Println("never n < tagLength the pack from antoher portal")
				return
			}
			tag := string(buf.Tag())
			// fmt.Println(addrString + tag)
			if value, ok := p.connMap.Get(addrString + tag); ok {
				if fc, ok := value.(*fakeUDPConn); ok {
					fc.WriteToSrc(buf.Raw(n).Data(0))
				} else {
					log.Println(fc, ok, value)  // never
					log.Println("invalid conn") // never
				}
			} else {
				// didn't get conn
				// create a new conn
				c, err := net.ListenUDP("udp", nil)
				if err != nil {
					log.Println("DailUDP failed") // never
					continue
				}
				fc := NewFakeUDPConn(
					p.dst, c,
					addr, p.UDPConn,
					addrString+tag, 90, func() {
						c.Close()
						p.connMap.Remove(addrString + tag)
					})
				// dst -X-> portal --> portal --> src
				go handleUDPConn(fc, p, tag)
				p.connMap.Put(addrString+tag, fc)
				fc.WriteToSrc(buf.Raw(n).Data(0))
			}
		}
	} //for
}

// dst -X-> portal --> portal --> src
func handleUDPConn(fc *fakeUDPConn, p *Portal, tag any) {
	defer fc.srcConn.Close()
	buf := make(PortalBuf, 1500)
	buf.AddTag(tag)
	for !fc.closed {
		n, err := fc.srcConn.Read(buf.Data(0))
		if err != nil {
			log.Println(err)
			continue
		}
		fc.WriteToDst(buf.DataAndTag(n))
	}
}
