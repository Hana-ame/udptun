package main

import (
	"encoding/binary"

	tools "github.com/Hana-ame/udptun/Tools"
)

type MyRouter struct {
	localAddr Addr
	*tools.ConcurrentHashMap[Addr, MyBus]

	bus MyBus
}

func (r *MyRouter) Route(f MyFrame) error {
	return r.GetOrDefault(f.Destination(), r.bus).SendFrame(f)
}

func (r *MyRouter) Read(remoteAddr Addr, bus MyBus) error {

	for {
		f, err := bus.RecvFrame()
		if err != nil {
			return err
		}
		if err := r.Route(f); err != nil {
			return err
		}
	}
}

func (r *MyRouter) Aloha(bus MyBus) error {
	// 用于存储所有键的切片
	// 使用 ForEach 方法获取所有键
	keys := make([]byte, r.Size()*2+2)
	binary.BigEndian.PutUint16(keys[:], uint16(r.localAddr))
	i := 1
	r.ForEach(func(key Addr, _ MyBus) {
		binary.BigEndian.PutUint16(keys[i*2:], uint16(key))
		i++
	})

	f := NewFrame(r.localAddr, 0, 0, Aloha, 0, 0, keys)
	if err := bus.SendFrame(f); err != nil {
		return err
	}

	return nil
}

func (r *MyRouter) Add(bus MyBus) error {
	if err := r.Aloha(bus); err != nil {
		return err
	}

	var remoteAddr Addr
	for {
		f, err := bus.RecvFrame()
		if err != nil {
			return err
		}
		if f.Command() == Aloha {
			data := f.Data()
			if len(data) < 2 || len(data)%2 != 0 {
				continue
			}
			for i := 0; i < len(data); i += 2 {
				addr := binary.BigEndian.Uint16(data[i:])
				r.Put(Addr(addr), bus)
			}

			remoteAddr = f.Source()
			break
		}
	}

	go r.Read(remoteAddr, bus)

	return nil
}
