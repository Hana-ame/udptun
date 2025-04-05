// 目前不一定能互相close的都

package main

import (
	"fmt"
	"testing"
	"time"
)

func TestSegment(t *testing.T) {
	seg := make(AddrSegments, 0)
	fmt.Println(seg)
	seg = append(seg, AddrSegment{5, 1 << 31})
	fmt.Println(seg)
	// seg = append(seg, AddrSegment{3, 1 << 31})
	fmt.Println(seg)
	seg = append(seg, AddrSegment{1, 1 << 31})
	fmt.Println(seg)
	seg.Sort()
	fmt.Println(seg)
	{
		index := seg.BinarySearch(0)
		fmt.Println(index)
	}
	{
		index := seg.BinarySearch(1)
		fmt.Println(index)
	}
	{
		index := seg.BinarySearch(2)
		fmt.Println(index)
	}
	{
		index := seg.BinarySearch(3)
		fmt.Println(index)
	}
	{
		index := seg.BinarySearch(4)
		fmt.Println(index)
	}
	{
		index := seg.BinarySearch(5)
		fmt.Println(index)
	}
	{
		index := seg.BinarySearch(6)
		fmt.Println(index)
	}
}

func TestSimpeRoute(t *testing.T) {
	defCh := NewFrameChan()
	c12 := NewFrameConn(1, 2)
	c21 := NewFrameConn(2, 1)
	r := NewRouter(0, defCh)
	go r.Serve(1, c12.WriteChan, c12.ReadChan)
	go r.Serve(2, c21.WriteChan, c21.ReadChan)
	go func() {
		for {
			f, e := defCh.ReadFrame()
			if e != nil {
				fmt.Println("read", e)
				return
			}
			fmt.Println("def", f, e)
		}
	}()

	go func() {
		for {
			f, e := c21.ReadFrame()
			if e != nil {
				fmt.Println("read", e)
				return
			}
			fmt.Println("c21", f, e)
		}
	}()
	go func() {
		for {
			f, e := c12.ReadFrame()
			if e != nil {
				fmt.Println("read", e)
				return
			}
			fmt.Println("c12", f, e)
		}
	}()

	c21.WriteFrame(Frame(2, 1, 0, 0, 0, []byte("123")))
	c12.WriteFrame(Frame(1, 2, 0, 0, 0, []byte("321")))
	c21.WriteFrame(Frame(2, 0, 0, 0, 0, []byte("1243")))
	c12.WriteFrame(Frame(1, 0, 0, 0, 0, []byte("3241")))
	time.Sleep(time.Second)
	c12.Close()
	c12.ReadChan.Close()
	time.Sleep(time.Second)
	c21.WriteFrame(Frame(2, 1, 0, 0, 0, []byte("12345")))
	time.Sleep(time.Second * 20)
}
