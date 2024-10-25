package main

import "testing"

func TestLoopBusChannel(t *testing.T) {
	lb := NewLoopBusChannel()

	go testHelper.ReadChan(lb)

	sc := lb.SendChan()

	sc <- NewFrame(0, 0, 0, 0, 0, 0, []byte("1"))
	sc <- NewFrame(0, 0, 0, 0, 0, 0, []byte("2"))
	sc <- NewFrame(0, 0, 0, 0, 0, 0, []byte("3"))
	sc <- NewFrame(0, 0, 0, 0, 0, 0, []byte("4"))
}

func TestPongBusChannel(t *testing.T) {
	lb := NewPongBusChannel()

	go testHelper.ReadChan(lb)

	sc := lb.SendChan()

	sc <- NewFrame(0, 0, 0, 0, 0, 0, []byte("1"))
	sc <- NewFrame(0, 0, 0, 0, 0, 0, []byte("2"))
	sc <- NewFrame(0, 0, 0, 0, 0, 0, []byte("3"))
	sc <- NewFrame(0, 0, 0, 0, 0, 0, []byte("4"))
}

func TestAlohaBusChannel(t *testing.T) {
	bb := NewPongBusChannel()
	lb := NewAlohaBusChannel(123, bb)

	go testHelper.ReadChan(lb)

	sc := lb.SendChan()

	sc <- NewFrame(0, 0, 0, 0, 0, 0, []byte("1"))
	sc <- NewFrame(0, 0, 0, 0, 0, 0, []byte("1"))
	sc <- NewFrame(0, 0, 0, Aloha, 0, 0, []byte("1"))
	sc <- NewFrame(0, 0, 0, 0, 0, 0, []byte("2"))
	sc <- NewFrame(0, 0, 0, 0, 0, 0, []byte("3"))
	sc <- NewFrame(0, 0, 0, Aloha, 0, 0, []byte("1"))
	sc <- NewFrame(0, 0, 0, 0, 0, 0, []byte("4"))
}
