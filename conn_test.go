package main

import (
	"fmt"
	"testing"
	"time"
)

func TestWrite(t *testing.T) {
	ch := NewFrameChan()
	go func() {
		for {
			f, e := ch.ReadFrame()
			if e != nil {
				fmt.Println(2, e)
				return
			}
			fmt.Println(f, e)
		}
	}()
	ch.WriteFrame(Frame(1, 2, 3, 4, 5, []byte("hello1")))
	time.Sleep(time.Second / 2)
	ch.WriteFrame(Frame(1, 2, 3, 4, 5, []byte("hello2")))
	time.Sleep(time.Second / 2)
	ch.WriteFrame(Frame(1, 2, 3, 4, 5, []byte("hello3")))
	time.Sleep(time.Second / 2)
	ch.WriteFrame(Frame(1, 2, 3, 4, 5, []byte("hello4")))
	time.Sleep(time.Second / 2)

	e := ch.Close()
	fmt.Println(1, e, ch.closed)

	time.Sleep(time.Second / 2)

	e = ch.WriteFrame(Frame(1, 2, 3, 4, 5, []byte("hello1")))
	fmt.Println(3, e)

}

func TestConn(t *testing.T) {
	c := NewFrameConn(111, 222)
	go func() {
		for {
			f, e := c.WriteChan.ReadFrame()
			if e != nil {
				fmt.Println("loop", e)
				return
			}
			e = c.ReadChan.WriteFrame(f)
			if e != nil {
				fmt.Println("loop", e)
				return
			}
		}
	}()

	go func() {
		for {
			f, e := c.ReadFrame()
			if e != nil {
				fmt.Println("read", e)
				return
			}
			fmt.Println(f, e)
		}
	}()

	c.WriteFrame(Frame(1, 2, 3, 4, 5, []byte("hello1")))
	time.Sleep(time.Second / 2)
	c.WriteFrame(Frame(1, 2, 3, 4, 5, []byte("hello2")))
	time.Sleep(time.Second / 2)
	c.WriteFrame(Frame(1, 2, 3, 4, 5, []byte("hello3")))
	time.Sleep(time.Second / 2)
	c.WriteFrame(Frame(1, 2, 3, 4, 5, []byte("hello4")))
	time.Sleep(time.Second / 2)

	e := c.Close()
	fmt.Println(1, e, c.WriteChan, c.ReadChan)

	time.Sleep(time.Second / 2)

	e = c.WriteFrame(Frame(1, 2, 3, 4, 5, []byte("hello1")))
	fmt.Println(3, e)

}

func TestConnReader(t *testing.T) {
	c := NewFrameConn(111, 222)
	go func() {
		for {
			f, e := c.WriteChan.ReadFrame()
			if e != nil {
				fmt.Println("loop", e)
				return
			}
			e = c.ReadChan.WriteFrame(f)
			if e != nil {
				fmt.Println("loop", e)
				return
			}
		}
	}()

	r := c.Reader()
	go func() {
		buf := make([]byte, 1024)
		for {
			n, e := r.Read(buf)
			if e != nil {
				fmt.Println("read", e)
				return
			}
			fmt.Println(string(buf[:n]))
		}
	}()

	w := c.Writer()

	w.Write([]byte("我是傻逼1"))
	w.Write([]byte("我是傻逼2"))
	w.Write([]byte("我是傻逼3"))
	w.Write([]byte("我是傻逼4"))
	w.Write([]byte("我是傻逼5"))

	time.Sleep(time.Second / 2)

	w.Close()

	time.Sleep(time.Second / 2)

	n, e := w.Write([]byte("我是傻逼6"))
	fmt.Println(n, e)

	time.Sleep(time.Second / 2)

}

func TestConnReader2(t *testing.T) {
	c := NewFrameConn(111, 222)
	go func() {
		for {
			f, e := c.WriteChan.ReadFrame()
			if e != nil {
				fmt.Println("loop", e)
				return
			}
			e = c.ReadChan.WriteFrame(f)
			if e != nil {
				fmt.Println("loop", e)
				return
			}
		}
	}()

	go func() {
		for {
			f, e := c.ReadFrame()
			if e != nil {
				fmt.Println("read", e)
				return
			}
			fmt.Println(f, e)
		}
	}()

	w := c.Writer()

	w.Write([]byte("我是傻逼1"))
	w.Write([]byte("我是傻逼2"))
	w.Write([]byte("我是傻逼3"))
	w.Write([]byte("我是傻逼4"))
	w.Write([]byte("我是傻逼5"))

	time.Sleep(time.Second / 2)

	w.Close()

	time.Sleep(time.Second / 2)

	n, e := w.Write([]byte("我是傻逼6"))
	fmt.Println(n, e)

	time.Sleep(time.Second / 2)

}
