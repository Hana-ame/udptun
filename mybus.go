package main

import (
	"fmt"
	"io"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ERR_CLOSED Error = "closed"
)

type MyBusWriter interface {
	SendFrame(MyFrame) error

	io.Closer
}

type MyBusReader interface {
	RecvFrame() (MyFrame, error)

	io.Closer
}

type MyBus interface {
	MyBusReader
	MyBusWriter

	io.Closer
}

type MyPongBus struct {
	localAddr Addr
	fc        chan MyFrame
}

func (b *MyPongBus) SendFrame(f MyFrame) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	b.fc <- f
	return
}

func (b *MyPongBus) RecvFrame() (f MyFrame, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	var ok bool
	f, ok = <-b.fc
	if ok {
		f.SetCommand(Pong)
		src := f.Source()
		f.SetSource(b.localAddr)
		f.SetDestination(src)
	}
	return
}

func (b *MyPongBus) Close() (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	close(b.fc)
	return
}

func NewPongBus() *MyPongBus {
	return &MyPongBus{fc: make(chan MyFrame)}
}

type MyNullBus struct{ closed bool }

func (b *MyNullBus) SendFrame(f MyFrame) (err error) {
	if b.closed {
		return ERR_CLOSED
	}
	return
}

func (b *MyNullBus) RecvFrame() (f MyFrame, err error) {
	if b.closed {
		return nil, ERR_CLOSED
	}
	return
}

func (b *MyNullBus) Close() (err error) {
	if b.closed {
		return ERR_CLOSED
	}
	b.closed = true
	return
}

func NewNullBus() *MyNullBus {
	return &MyNullBus{}
}
