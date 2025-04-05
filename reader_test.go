package main

import (
	"fmt"
	"testing"
	"time"
)

// 关闭之后再接受会返回 ok = false
func TestChannel(t *testing.T) {
	a := make(chan int)
	go func() {
		for {
			i, ok := <-a
			fmt.Println(i, ok)
		}
	}()
	a <- 1
	time.Sleep(time.Second)
	a <- 2
	time.Sleep(time.Second)
	close(a)
	time.Sleep(time.Second)
	a <- 3
}

func TestChannelClose(t *testing.T) {
	a := make(chan int, 2)
	a <- 1
	a <- 2
	close(a)
}

func TestGC(t *testing.T) {
	reader := NewFrameReader()
	for {
		reader.Close()
		reader = NewFrameReader()
	}
	reader.Close()
}
