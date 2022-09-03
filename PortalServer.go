package main

type PortalServer struct {
	LocalAddr *string
	Pool      *PortalPool
}

func NewPortalServer(addr string) *PortalServer {
	pool := NewPortalPool(5, 10)
	return &PortalServer{
		LocalAddr: &addr,
		Pool:      pool,
	}
}

func (s *PortalServer) NewPortal() {
	p := NewPortal("udp")
	s.Pool.Add(p)
}

// paddr: address from peer,
// laddr: default(nil) is s.LocalAddr
// mux  : never used
func (s *PortalServer) ActivePortal(paddr *string, laddr *string, mux *Multiplexer) (p *Portal) {
	p = s.Pool.Pick()
	if p == nil {
		return
	}
	if laddr == nil {
		laddr = s.LocalAddr
	}
	p.Set(paddr, laddr, nil)
	// if mux != nil {
	// 	s.Pool.Add(p)
	// }
	if s.Pool.cnt < s.Pool.mlen {
		go s.NewPortal()
	}
	return
}
