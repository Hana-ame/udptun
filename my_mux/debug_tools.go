package mymux

import "github.com/Hana-ame/udptun/Tools/debug"

type Helper string

func (h Helper) ReadBus(r BusReader) {
	for {
		f, e := r.RecvFrame()
		if e != nil {
			debug.W(h, e)
		}
		debug.I(h, (f).String())
	}
}

func (h Helper) Copy(r BusReader, w BusWriter) {
	for {
		f, e := r.RecvFrame()
		if e != nil {
			debug.W(h, e)
		}
		e = w.SendFrame(f)
		if e != nil {
			debug.W(h, e)
		}
	}
}
func (h Helper) CopyWithHandler(r BusReader, w BusWriter, handler func(Frame) Frame) {
	for {
		f, e := r.RecvFrame()
		if e != nil {
			debug.W(h, e)
		}
		f = handler(f)
		e = w.SendFrame(f)
		if e != nil {
			debug.W(h, e)
		}
	}
}
