package utils

import (
	"log"
	"net"
)

func UDPEcho(laddr string) {
	la, _ := net.ResolveUDPAddr(`udp`, laddr)
	lc, _ := net.ListenUDP(`udp`, la)
	buf := make([]byte, 2048)
	for {
		l, addr, err := lc.ReadFrom(buf)
		log.Println(`UDPEcho:`, l, addr)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		lc.WriteTo(buf[:l], addr)
	}
}
