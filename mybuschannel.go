package main

import (
	"encoding/binary"

	"github.com/Hana-ame/udptun/Tools/debug"
)

// 双向的。

type BusChannel interface {
	RecvChannel
	SendChannel
}

type RawBusChannel struct {
	RecvChannel
	SendChannel
}

func NewLoopBusChannel() BusChannel {
	c := make(RawChannel)
	return RawBusChannel{
		c,
		c,
	}
}

func NewPongBusChannel() BusChannel {
	sc := make(RawChannel)
	rc := make(RawChannel)
	go func() {
		for {
			f, ok := <-sc
			if !ok {
				close(rc)
				return
			}
			f.SetCommand(Pong)
			rc <- f
		}
	}()
	return &RawBusChannel{
		SendChannel: sc,
		RecvChannel: rc,
	}
}

type AlohaBusChannel struct {
	RecvChannel
	SendChannel

	localAddr Addr
	bus       BusChannel
}

func NewAlohaBusChannel(laddr Addr, bus BusChannel) BusChannel {
	sc := make(RawChannel)
	c := make(RawChannel)
	rc := NewMergeChannel(c, bus.RecvChan())
	go func() {
		defer func() {
			e := recover()
			if e != nil {
				debug.W("AlohaBusChannel", e)
			}
		}()
		defer CloseCh(sc)
		defer CloseCh(c)
		defer CloseCh(rc)
		for {
			f, ok := <-sc
			if !ok {
				return
			}
			if f.Command() == Aloha {
				data := make([]byte, 2)
				binary.BigEndian.PutUint16(data, uint16(laddr))
				c <- NewFrame(laddr, f.Source(), f.Port(), Aloha, 0, 0, data)
			} else {
				bus.SendChan() <- f
			}

		}
	}()
	return &AlohaBusChannel{
		SendChannel: sc,
		RecvChannel: rc,
		localAddr:   laddr,
		bus:         bus,
	}
}
