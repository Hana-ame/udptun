package mymux

import (
	"io"

	"github.com/Hana-ame/udptun/Tools/debug"
)

// BusWriter 接口定义了发送帧的功能，并包含关闭功能。
type BusWriter interface {
	SendFrame(Frame) error

	io.Closer
}

// BusReader 接口定义了接收帧的功能，提供锁功能以确保线程安全，并包含关闭功能。
type BusReader interface {
	RecvFrame() (Frame, error)

	io.Closer
}

type Bus interface {
	BusReader
	BusWriter

	io.Closer
}

type RawBus struct {
	BusReader
	BusWriter
}

//	func (b *RawBus) RecvFrame() (MyFrame, error) {
//		return b.MyBusReader.RecvFrame()
//	}
//
//	func (b *RawBus) SendFrame(f MyFrame) error {
//		return b.MyBusWriter.SendFrame(f)
//	}
func (b *RawBus) Close() error {
	re := b.BusReader.Close()
	rw := b.BusWriter.Close()
	if re != nil {
		return re
	}
	if rw != nil {
		return rw
	}
	return nil
}

func NewBusWithHandler(handler func(f Frame) Frame) Bus {
	// reader, iw := NewDebugPipe("1")
	// ir, writer := NewDebugPipe("2")
	reader, iw := NewPipe()
	ir, writer := NewPipe()
	// go ReadBus(ir)
	// _ = iw
	go func() {
		defer ir.Close()
		defer iw.Close()

		for {
			f, e := ir.RecvFrame()
			if e != nil {
				debug.T("handler", e)
				return
			}
			f = handler(f)
			if e := iw.SendFrame(f); e != nil {
				debug.T("handler", e)
				return
			}
		}
	}()
	return &RawBus{reader, writer}
}

func NewBusPair() (Bus, Bus) {
	r0, w1 := NewPipe()
	r1, w0 := NewPipe()
	return &RawBus{r0, w0}, &RawBus{r1, w1}
}

type BusOnClose struct {
	Bus
	onClose func() error
}

func NewBusOnClose(bus Bus, onClose func() error) *BusOnClose {
	return &BusOnClose{
		Bus:     bus,
		onClose: onClose,
	}
}

func (b *BusOnClose) Close() error {
	be := b.Close()
	oe := b.onClose()
	if oe != nil {
		return oe
	}
	if be != nil {
		return be
	}
	return nil
}
