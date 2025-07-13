package main

import (
	"fmt"
	"testing"
	"time"
)

func TestIO(t *testing.T) {
	pipe := NewPipe()
	writer := FrameWriter{FrameHandler: pipe, MTU: 100000}
	reader := FrameReader{FrameHandler: pipe}
	go func() {
		for {
			writer.Write([]byte("1234567890qwertyuiopasdfghjklzxcvbnm"))
		}
	}()
	go func() {
		buf := make([]byte, 1000)
		for {
			n, _ := reader.Read(buf)
			fmt.Printf("%s\n", buf[:n])
		}
	}()
	time.Sleep(30 * time.Second)
}
