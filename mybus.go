package main

import (
	"encoding/binary"
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

func Copy(writer MyBusWriter, reader MyBusReader) error {
	for {
		f, err := reader.RecvFrame()
		if err != nil {
			return err
		}
		if err := writer.SendFrame(f); err != nil {
			return err
		}
	}
}

type MyBus interface {
	MyBusReader
	MyBusWriter

	io.Closer
}
type MyBusWrapper struct {
	MyBusReader
	MyBusWriter

	closed bool
}

func (b *MyBusWrapper) Close() (err error) {
	if b.closed {
		return ERR_CLOSED
	}
	er := b.MyBusReader.Close()
	ew := b.MyBusWriter.Close()
	if er != nil {
		err = er
	}
	if ew != err {
		err = ew
	}
	b.closed = true
	return
}

func NewBus(reader MyBusReader, writer MyBusWriter) *MyBusWrapper {
	return &MyBusWrapper{
		MyBusReader: reader,
		MyBusWriter: writer,
	}
}

// 会响应aloha，但其实应该填充到C里面
type MyPongBus struct {
	localAddr Addr
	C         chan MyFrame
}

func (b *MyPongBus) SendFrame(f MyFrame) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	if f.Command() == Aloha {
		data := make([]byte, 2)
		binary.BigEndian.PutUint16(data, uint16(b.localAddr))
		f = NewFrame(f.Source(), f.Destination(), f.Port(), f.Command(), f.SequenceNumber(), f.AcknowledgeNumber(), data)
	} else {
		f.SetCommand(Pong)
	}
	f.SetDestination(f.Source())
	f.SetSource(b.localAddr)

	b.C <- f
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
	f, ok = <-b.C
	if !ok {
		return nil, fmt.Errorf("%v", ok)
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
	close(b.C)
	return
}

func NewPongBus(localAddr Addr) *MyPongBus {
	return &MyPongBus{localAddr: localAddr, C: make(chan MyFrame, 1)}
}

// drop to void
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
