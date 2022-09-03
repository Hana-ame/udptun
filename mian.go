package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var c *net.UDPConn
var peerAddr *net.UDPAddr

// var getPool *PortalPool
// var putPool *PortalPool
var pc *PortalClient

var m map[string]*Portal
var mu sync.Mutex

var addr string
var host string
var path string

func main() {

	addr = "localhost:9999"
	host = "http://127.0.233.0:8080"
	path = "test"

	flag.StringVar(&addr, "a", ":10000", "address, udp")
	flag.StringVar(&host, "h", "https://moonchan.xyz/233", "help server host, no /")
	flag.StringVar(&path, "p", "test", "path, should be specified")

	isServer := flag.Bool("s", false, "work as server when specified")

	flag.Parse()
	fmt.Println(addr, host, path, *isServer)

	if *isServer {
		Server(addr)
	} else {
		Client(addr)
	}
	// go Server(addr)

	// addr = ":10000"
	// go Client(addr)

	// time.Sleep(time.Hour)
}

func ActiveClientPortal() {
	if pc.Pool.Cnt() < pc.Pool.mlen {
		go pc.NewPortal()
	}
	if peerAddr == nil {
		return
	}
	p := pc.Pool.Pick()
	if p == nil {
		return
	}

	// fmt.Println("ActiveClientPortal", p, pc.Pool, pc.Mux.Pool, peerAddr, c)
	// fmt.Println("****** ActiveClientPortal *******", "\n", peerAddr, "\n", p.LocalAddr, "\n", "******")
	// fmt.Println(m)
	c.WriteToUDP([]byte(p.LocalAddr), peerAddr)
	mu.Lock()
	m[p.LocalAddr] = p
	mu.Unlock()

	// fmt.Println("****** ActiveClientPortal *******", "\n", p, "\n", pc.Pool, "\n", pc.Mux.Pool, "\n", "******")

	// for i := putPool.mlen - putPool.Cnt(); i > 0; i-- {
	// 	ActiveClientPortal()
	// }
	// time.Sleep(time.Second)
	if pc.Mux.Pool.mlen-1 > pc.Mux.Pool.Cnt() {
		go ActiveClientPortal()
	}
}

func Client(listenAddr string) {
	stopFlag := false
	pc = NewPortalClient(listenAddr)
	for i := 0; i < 5; i++ {
		go pc.NewPortal()
	}
	pc.NewPortal()

	go debug(pc)
	// getPool = pc.Pool
	// putPool = pc.Mux.Pool
	pc.Mux.RecvConnCallBack = ActiveClientPortal

	var err error
	c, err = net.ListenUDP("udp", nil)
	if err != nil {
		log.Println(err)
		return
	}
	localAddr, err := GetAddr(c)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("!!!!!!!", localAddr)
	go func(host string, path string, data string) { // not tested, should be work. connect between
		for !stopFlag {
			res := PostAddr(host, path, data, 3600*24*365)
			fmt.Println("!!!!!!", res)
			if res != nil && len(res) == 2 {
				// fmt.Println("!!!!!!", res)
				for _, s := range res {
					udpaddr, err := net.ResolveUDPAddr("udp", s)
					if err != nil {
						log.Println(err)
					}
					_, err = c.WriteToUDP([]byte{}, udpaddr)
					if err != nil {
						log.Println(err)
					}
					if s != data {
						peerAddr = udpaddr
						log.Println("peerAddr : ", peerAddr)
					}
				}
				time.Sleep(time.Second * 3)
			}
		}
	}(host, path, localAddr)

	m = make(map[string]*Portal)
	for peerAddr == nil {
		time.Sleep(time.Second)
	}

	buf := make([]byte, 2048)
	log.Println("before recv first pack")
	l, raddr, err := c.ReadFromUDP(buf)
	log.Println("recv first pack", raddr, "len = ", l)

	ActiveClientPortal()
	ActiveClientPortal()
	ActiveClientPortal()

	fmt.Println("===============")
	fmt.Println(pc.Mux.Pool)
	fmt.Println(pc.Pool)
	fmt.Println(pc.Mux.m)
	fmt.Println("===============")

	ActiveClientPortal()
	ActiveClientPortal()

	for {
		// fmt.Println("recv")
		l, raddr, err := c.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		if l == 0 {
			continue
		}
		msg := string(buf[:l])
		fmt.Println("==== msg: ====", "\n", msg, "\n", "====")
		msgs := strings.Split(msg, "\n")
		_, err = net.ResolveUDPAddr("udp", msgs[1])
		if err != nil {
			log.Println(err)
			continue
		}
		_, err = net.ResolveUDPAddr("udp", msgs[0])
		if err != nil {
			log.Println(err)
			continue
		}

		mu.Lock()
		p := m[msgs[0]]
		if p != nil {
			p.Set(&msgs[1], nil, pc.Mux) // set peer only, local addr will set by mux
			p.Mux.Pool.Add(p)
			delete(m, msgs[0])
		}
		mu.Unlock()

		_ = raddr
	}

}

func Server(forwardAddr string) {
	stopFlag := false
	ps := NewPortalServer(forwardAddr)
	for i := 0; i < 5; i++ {
		go ps.NewPortal()
	}

	var err error
	c, err = net.ListenUDP("udp", nil)
	if err != nil {
		log.Println(err)
		return
	}
	localAddr, err := GetAddr(c)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("!!!!!!!", localAddr)
	go func(host string, path string, data string) { // not tested, should be work. connect between
		for !stopFlag {
			res := PostAddr(host, path, data, 3600*24*365)
			if res != nil {
				for _, s := range res {
					udpaddr, err := net.ResolveUDPAddr("udp", s)
					if err != nil {
						log.Println(err)
					}
					c.WriteToUDP([]byte{}, udpaddr)
				}
				time.Sleep(time.Second * 3)
			}
		}
	}(host, path, localAddr)

	// loop
	// handle the recv message
	buf := make([]byte, 2048)
	log.Println("before recv first pack")
	l, raddr, err := c.ReadFromUDP(buf)
	log.Println("recv first pack", raddr, "len = ", l)
	if err != nil {
		log.Println(err)
	}
	for {
		l, raddr, err := c.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		if l == 0 {
			continue
		}
		msg := string(buf[:l])

		fmt.Println("=== msg ===")
		fmt.Println(msg)
		fmt.Println("===")

		_, err = net.ResolveUDPAddr("udp", msg)
		if err != nil {
			log.Println(err)
			continue
		}
		p := ps.ActivePortal(&msg, ps.LocalAddr, nil)
		if p == nil {
			log.Println(msg, "no portal avaliable")
			continue
		}
		// fmt.Println(ps.Pool)
		// fmt.Println(p)

		res := msg + "\n" + p.LocalAddr
		fmt.Println("=== res ===")
		fmt.Println(res)
		fmt.Println(raddr)
		fmt.Println("===")
		c.WriteToUDP([]byte(res), raddr)
		c.WriteToUDP([]byte(res), raddr)
		c.WriteToUDP([]byte(res), raddr)
	}
}

func debug(pc *PortalClient) {
	for {
		time.Sleep(time.Second * 10)
		fmt.Println("===============")
		fmt.Println(pc.Mux.Pool)
		fmt.Println(pc.Pool)
		fmt.Println(pc.Mux.m)
		fmt.Println("===============")
	}
}
