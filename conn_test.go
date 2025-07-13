package main

import (
	"fmt"
	"testing"
	"time"
)

func TestConn(t *testing.T) {
	c := NewConn()
	go Copy(c.RouterInterface(), c.RouterInterface())
	// go Copy(c.RouterInterface(), c.RouterInterface())
	writer := FrameWriter{FrameHandler: c.ApplicatonInterface(), MTU: 100000}
	reader := FrameReader{FrameHandler: c.ApplicatonInterface()}
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
