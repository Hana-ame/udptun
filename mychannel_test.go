package main

import (
	"fmt"
	"testing"
	"time"
)

type ic chan int

func (c ic) RecvChan() chan int {
	return chan int(c)
}

func TestXxx(t *testing.T) {
	mc := make(ic)
	rc := mc.RecvChan()
	// 0xc000088310
	// 0xc000088310
	println(mc)
	println(rc)
	// 0xc00005a070
	// 0xc000058740
	println(&mc)
	println(&rc)
	go func() {
		mc <- 1
	}()

	r := <-rc
	fmt.Println(r)
}

func TestMerge(t *testing.T) {
	sc0 := make(RawChannel)
	sc1 := make(RawChannel)
	mc := NewMergeChannel(sc0, sc1)

	go testHelper.ReadChan(mc)

	sc0 <- NewFrame(0, 0, 0, 0, 0, 0, []byte("1"))
	sc1 <- NewFrame(0, 0, 0, 0, 0, 0, []byte("2"))
	sc0 <- NewFrame(0, 0, 0, 0, 0, 0, []byte("3"))
	sc1 <- NewFrame(0, 0, 0, 0, 0, 0, []byte("4"))

	time.Sleep(5 * time.Second)
}
