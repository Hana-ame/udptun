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

	router utils.LockedMap

	connMap utils.LockedMap
}

func (p *portal) getLocalAddr(IPv4 bool) string {
	if IPv4 {
		fmt.Println(p.localAddr, p.stunServer)
		if p.localAddr == "" && p.stunServer != "" {
			go func() {
				for {
					// fmt.Println("send")
					utils.StunRequest(p.stunServer, p.UDPConn)
					time.Sleep(1 * time.Second)
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
		log.Println(addr.String())
		s, err := utils.StunResolve(buf[:n])
		log.Println(s)
	}
}

func main() {
	lc, _ := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 2345})

	p := &portal{
		UDPConn: lc,
	}
	go p.run()
	r := p.getLocalAddr(true)
	fmt.Println(r)
	fmt.Println(p.UDPConn.LocalAddr().String())

	time.Sleep(time.Second * 9)
	r = p.getLocalAddr(true)
	fmt.Println(r)
}

func runServer(src string, dst string) {
	lc, err := net.Listen("udp", src)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	_ = lc

}
