package mymux

import (
	"testing"
	"time"

	"github.com/Hana-ame/udptun/Tools/debug"
)

func TestMux(t *testing.T) {
	mux := NewMux()

	b1 := NewBusWithHandler(func(f Frame) Frame {
		src, dst := f.Source(), f.Destination()
		f.SetDestination(src)
		f.SetSource(dst)
		return f
	})
	b2 := NewBusWithHandler(func(f Frame) Frame {
		src, dst := f.Source(), f.Destination()
		f.SetDestination(src)
		f.SetSource(dst)
		return f
	})
	b3 := NewBusWithHandler(func(f Frame) Frame {
		src, dst := f.Source(), f.Destination()
		f.SetDestination(src)
		f.SetSource(dst)
		return f
	})
	b, b0 := NewBusPair()
	mux.AddBus(1, b1)
	mux.AddBus(2, b2)
	mux.AddBus(3, b3)
	mux.AddBus(0, b0)

	go Helper("111").ReadBus(b)
	go Helper("222").ReadBus(mux)

	var e error
	e = mux.SendFrame(NewDataFrame(0, 1, 1, 3, 4, []byte("1")))
	debug.I("e", e)
	e = mux.SendFrame(NewDataFrame(0, 2, 2, 3, 4, []byte("2")))
	debug.I("e", e)
	e = mux.SendFrame(NewDataFrame(0, 3, 3, 3, 4, []byte("3")))
	debug.I("e", e)
	e = mux.SendFrame(NewDataFrame(0, 4, 4, 3, 4, []byte("4")))
	debug.I("e", e)
	e = mux.SendFrame(NewDataFrame(0, 0, 0, 3, 4, []byte("1")))
	debug.I("e", e)
	e = b.SendFrame(NewDataFrame(0, 0, 0, 3, 4, []byte("1")))
	debug.I("e", e)

	time.Sleep(time.Second * 5)
}
