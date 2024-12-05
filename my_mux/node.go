package mymux

import (
	"fmt"
	"sync"

	"github.com/Hana-ame/udptun/Tools/debug"
)

type BusAccpeter struct {
	Bus
	port byte
	*sync.Cond

	closed bool
}

func (a *BusAccpeter) AcceptableBus(b Bus, port byte) {
	a.L.Lock()
	for a.Bus != nil && !a.closed {
		a.Wait()
	}
	a.Bus = b
	a.port = port
	a.L.Unlock()
	a.Signal()
}

func (a *BusAccpeter) Accept() (b Bus, port byte) {
	a.L.Lock()
	for a.Bus == nil && !a.closed {
		a.Wait()
	}
	b = a.Bus
	port = a.port
	a.Bus = nil
	a.L.Unlock()
	a.Signal()
	return
}

func (a *BusAccpeter) Close() error {
	a.closed = true
	a.Broadcast()
	return nil
}

func NewBusAccepter() *BusAccpeter {
	a := &BusAccpeter{
		Cond: sync.NewCond(&sync.Mutex{}),
	}
	return a
}

type Node struct {
	*Mux

	*BusAccpeter

	f Frame
	*sync.Cond

	closed bool
}

func (n *Node) SendFrame(f Frame) error {
	bus, ok := n.Get(f.Port())
	if f.Command() == Close {
		if ok {
			bus.Close()
		}
		return nil
	}
	// 如果请求就新建链接
	if f.Command() == Request {
		if ok {
			bus.Close()
		}
		bus, ret := NewBusPair()
		n.AddBus(f.Port(), bus)
		n.AcceptableBus(ret, f.Port())
		return nil
	}
	// 这里没有的话返回close
	if !ok {
		n.L.Lock()
		for n.f != nil && !n.closed {
			n.Wait()
		}
		n.f = NewFrame(f.Destination(), f.Source(), f.Port(), Close, 0, 0, nil)
		n.L.Unlock()
		n.Broadcast()
		return nil
	}

	return bus.SendFrame(f)
}

func (n *Node) RecvFrame() (f Frame, err error) {
	n.L.Lock()
	for n.f == nil && !n.closed {
		n.Wait()
	}
	f = n.f
	n.f = nil
	n.L.Unlock()
	n.Broadcast()
	return
}

func (n *Node) Close() error {
	if n.closed {
		return ERR_CLOSED
	}
	n.BusAccpeter.Close()
	n.Mux.Close()
	n.closed = true
	n.Broadcast()
	return nil
}

func (n *Node) Dial(dst, src Addr, port byte) (Bus, error) {
	bus, ok := n.Get(port)
	if ok {
		return nil, fmt.Errorf("exist")
	}
	n.L.Lock()
	for n.f != nil && !n.closed {
		n.Wait()
	}
	n.f = NewFrame(src, dst, port, Request, 0, 0, nil)
	n.L.Unlock()
	n.Broadcast()
	bus, ret := NewBusPair()
	n.AddBus(port, bus)

	return NewBusOnClose(ret, func() error {
		return n.RemoveBus(port)
	}), nil
}

func NewNode() *Node {
	n := &Node{
		Mux:         NewMux(),
		BusAccpeter: NewBusAccepter(),
		Cond:        sync.NewCond(&sync.Mutex{}),
	}
	go func() {
		for {
			f, e := n.Mux.RecvFrame()
			if e != nil {
				debug.W("node", e)
				n.Close()
			}
			n.L.Lock()
			for n.f != nil && !n.closed {
				n.Wait()
			}
			n.f = f
			n.L.Unlock()
			n.Broadcast()
		}
	}()
	return n
}
