package main

import (
	// "fmt"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/hana-ame/udptun/utils"
)

type portal struct {
	*net.UDPConn
	stunServer string
	localAddr  string

	router utils.LockedMap // map[addr.String()]func([]byte)

	serverDst *net.UDPAddr
	connMap   utils.LockedMap
}

func (p *portal) getLocalAddr(isIPv4 bool) string {
	if isIPv4 {
		if p.localAddr == "" && p.stunServer != "" {
			go func() {
				for {
					utils.StunRequest(p.stunServer, p.UDPConn)
					time.Sleep(5 * time.Second)
				}
			}()
		}
		return p.localAddr
	} else {
		return utils.GetOutboundIPv6(p.UDPConn)
	}
}

func (p *portal) run() {
	p.stunServer = "142.251.2.127:19302"

	for {
		buf := make([]byte, 1500)
		fmt.Println("reading")
		n, addr, err := p.ReadFrom(buf)
		if err != nil {
			log.Println(err)
			continue
		}

		if addr.String() == p.stunServer {
			p.localAddr, err = utils.StunResolve(buf[:n])
			if err != nil {
				log.Println("error when recv from stunServer", err)
				continue
			}
		} else if value, ok := p.router.Get(addr.String()); ok {
			if handler, ok := value.(func([]byte)); ok {
				handler(buf[:n])
			} else {
				log.Println("invalid router")
				continue
			}
		} else {
			log.Println("not supposed to be seen@protal run for loop")
		}
	}
}

var p *portal

func renewAddr(p *portal, isIPv4 bool) {
	for {
		r := p.getLocalAddr(isIPv4)
		if r != "" {
			//TODO: Renew
		}
		time.Sleep(5 * time.Second)
	}
}

func main() {

	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:1900")
	lc, _ := net.ListenUDP("udp", addr)
	NewU

	// initial the portal
	// lc, _ := net.ListenUDP("udp", nil)
	// p = &portal{
	// 	UDPConn: lc,
	// }
	// go p.run()
	// go renewAddr(p, true)

}
