package main

import "fmt"

type ReadChan[T any] (chan T)
type WriteChan[T any] (chan T)

func (c ReadChan[T]) Read() (T, bool) {
	v, ok := <-c
	return v, ok
}
func (c WriteChan[T]) Write(data T) {
	c <- data
}

type IChan[T any] interface {
	ReadChan[T]
	WriteChan[T]
}

type Chan struct {
	ReadChan[frame]
	WriteChan[frame]
	// status
	closed bool
	// config
	src uint32
	dst uint32
}

func (c *Chan) Close() error {
	c.closed = true
	return nil
}

func (c *Chan) Write(data []byte) error {
	if c.closed {
		return fmt.Errorf("chan: write to a closed chan")
	}
	c.closed = true
	return nil
}

func (c *Chan) Read() ([]byte, error) {
	data := <-c.ReadChan

	return data, nil
}
