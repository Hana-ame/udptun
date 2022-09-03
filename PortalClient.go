package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
// MTU = 2048
)

type Multiplexer struct {
	Conn *net.UDPConn // 占用的port

	Pool *PortalPool
	m    map[string]*Portal
	mu   sync.RWMutex

	RecvConnCallBack func()
}

func NewMultiplexer(ptype string, addr *string) *Multiplexer {
	udpaddr, err := net.ResolveUDPAddr(ptype, *addr)
	if err != nil {
		log.Println(err)
		return nil
	}
	c, err := net.ListenUDP(ptype, udpaddr)
	if err != nil {
		log.Println(err)
		return nil
	}

	m := make(map[string]*Portal)

	pool := NewPortalPool(5, 10)

	mux := &Multiplexer{
		Conn: c,
		m:    m,
		Pool: pool,
	}
	go mux.Start()

	return mux
}

func (m *Multiplexer) Start() {
	fmt.Println("well, did you worked?")
	for {
		buffer := make([]byte, MTU)
		l, addr, err := m.Conn.ReadFrom(buffer)
		if err != nil {
			log.Println(err.Error())
			log.Println(err.Error())
			log.Println(err.Error())
			continue
		}
		// m.addrfrom = &addr
		// log.Println(`[Mux]receive from `, addr, `len=`, l)
		// _ = l
		m.handlePacket(l, buffer, &addr)
	}
}

func (m *Multiplexer) handlePacket(l int, b []byte, addr *net.Addr) {
	addrString := (*addr).String()

	m.mu.RLock()
	if m.m[addrString] == nil { // 处理新链接
		m.mu.RUnlock()

		p := m.Pool.Pick() // pick a unused portal to handle packets from this address

		go m.RecvConnCallBack() // would it work?

		for p == nil || p.status == DYING {
			go m.RecvConnCallBack() // would it work?
			fmt.Println("peer ")
			fmt.Println(m.Pool)
			time.Sleep(time.Second)
			p = m.Pool.Pick() // pick a unused portal to handle packets from this address
		}
		fmt.Println(addrString)
		p.Set(nil, &addrString, nil)

		m.mu.Lock()
		m.m[addrString] = p
		m.mu.Unlock()

		m.mu.RLock()
	} // m.mp[addrString] != nil
	p := m.m[addrString]
	m.mu.RUnlock()

	// fmt.Println(p) // 不在

	p.RecvPacketFromOthers(l, b)
}

func (m *Multiplexer) Remove(addrString string) {
	m.mu.Lock()
	delete(m.m, addrString)
	m.mu.Unlock()
}

type PortalClient struct {
	Pool *PortalPool

	LocalAddr *string
	Mux       *Multiplexer
}

func NewPortalClient(addr string) *PortalClient {
	pool := NewPortalPool(5, 10)
	mux := NewMultiplexer("udp", &addr)
	return &PortalClient{
		LocalAddr: &addr,
		Pool:      pool,
		Mux:       mux,
	}
}

func (c *PortalClient) NewPortal() {
	p := NewPortal("udp")
	c.Pool.Add(p)
	if c.Pool.Cnt() < c.Pool.mlen {
		c.NewPortal()
	}
}

// not used!
// paddr: address from peer,
// laddr: always nil
// mux  : always c.Mux
func (c *PortalClient) ActivePortal(paddr *string, laddr *string, mux *Multiplexer) (p *Portal) {
	p = c.Pool.Pick()
	if c.Pool.Cnt() < c.Pool.mlen {
		c.NewPortal()
	}
	if p == nil {
		return
	}
	p.Set(paddr, nil, mux)
	if mux != nil {
		c.Mux.Pool.Add(p)
	} else {
		log.Println("This should never occure")
	}
	if c.Pool.cnt < c.Pool.mlen {
		go c.NewPortal()
	}
	return
}
