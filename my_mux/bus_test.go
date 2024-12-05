package mymux

import (
	"testing"

	"github.com/Hana-ame/udptun/Tools/debug"
)

func TestLoopBus(t *testing.T) {
	rb, wb := NewPipe()
	b := &RawBus{rb, wb}

	go Helper("111").ReadBus(b)

	b.SendFrame(NewDataFrame(0, 1, 2, 3, 4, []byte("1")))
	b.SendFrame(NewDataFrame(0, 1, 2, 3, 4, []byte("2")))
	b.SendFrame(NewDataFrame(0, 1, 2, 3, 4, []byte("3")))
	b.SendFrame(NewDataFrame(0, 1, 2, 3, 4, []byte("4")))
	b.SendFrame(NewDataFrame(0, 1, 2, 3, 4, []byte("1")))
}

func TestHandlerBus(t *testing.T) {
	b := NewBusWithHandler(func(f Frame) Frame {
		src, dst := f.Source(), f.Destination()
		f.SetDestination(src)
		f.SetSource(dst)
		return f
	})

	go Helper("1").ReadBus(b)

	var e error
	e = b.SendFrame(NewDataFrame(0, 1, 2, 3, 4, []byte("1")))
	debug.I("e", e)
	e = b.SendFrame(NewDataFrame(0, 1, 2, 3, 4, []byte("2")))
	debug.I("e", e)
	e = b.SendFrame(NewDataFrame(0, 1, 2, 3, 4, []byte("3")))
	debug.I("e", e)
	e = b.SendFrame(NewDataFrame(0, 1, 2, 3, 4, []byte("4")))
	debug.I("e", e)
	e = b.SendFrame(NewDataFrame(0, 1, 2, 3, 4, []byte("1")))
	debug.I("e", e)
}
