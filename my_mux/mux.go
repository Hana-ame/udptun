package mymux

import (
	"fmt"
	"sync"

	tools "github.com/Hana-ame/udptun/Tools"
	"github.com/Hana-ame/udptun/Tools/debug"
)

// port, bus
type Mux struct {
	*tools.ConcurrentHashMap[byte, Bus]

	f Frame
	*sync.Cond

	closed bool
}

func NewMux() *Mux {
	m := &Mux{
		ConcurrentHashMap: tools.NewConcurrentHashMap[byte, Bus](),
		Cond:              sync.NewCond(&sync.Mutex{}),
	}
	return m
}

func (m *Mux) SendFrame(f Frame) error {
	bus, ok := m.Get(f.Port())
	if !ok {
		return fmt.Errorf("not exist")
	}
	return bus.SendFrame(f)
}
func (m *Mux) RecvFrame() (Frame, error) {
	m.L.Lock()
	for m.f == nil && !m.closed {
		m.Wait()
	}
	f := m.f
	m.f = nil
	m.L.Unlock()
	m.Signal()
	return f, nil
}
func (m *Mux) Close() error {
	if m.closed {
		return ERR_CLOSED
	}
	m.ConcurrentHashMap.ForEach(func(_ byte, bus Bus) {
		bus.Close()
	})
	m.closed = true
	m.Broadcast()
	return nil
}
func (m *Mux) AddBus(port byte, bus Bus) error {
	m.Put(port, bus)
	go func() {
		defer bus.Close()
		for {
			f, err := bus.RecvFrame()
			if err != nil {
				debug.I("mux", err)
				return
			}
			m.L.Lock()
			for m.f != nil && !m.closed {
				m.Wait()
			}
			m.f = f
			m.L.Unlock()
			m.Signal()
		}
	}()
	return nil
}
func (m *Mux) RemoveBus(port byte) error {
	bus, ok := m.Get(port)
	if !ok {
		return ERR_CLOSED
	}
	bus.Close()
	m.Remove(port)

	return nil
}

// 	for {
// 		f, err := b.RecvFrame()
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		// 怎么处理草泥马
// 	}
// }
