package mymux

import (
	"strconv"
	"testing"
	"time"

	"github.com/Hana-ame/udptun/Tools/debug"
)

// Helper(strconv.Itoa(1)).ReadBus(b)

// pass
func TestNode(t *testing.T) {
	n0 := NewNode()
	n1 := NewNode()
	go Helper("").Copy(n0, n1)
	go Helper("").Copy(n1, n0)

	go func() {
		i := 1
		for {
			b, _ := n0.Accept()
			go Helper(strconv.Itoa(i)).ReadBus(b)
			i++
		}
	}()

	go func() {
		b, e := n1.Dial(0, 0, 0)
		if e != nil {
			debug.E("aaa", e)
		}
		b.SendFrame(NewDataFrame(0, 0, 0, 3, 4, []byte("1")))
		b.SendFrame(NewDataFrame(0, 0, 0, 3, 4, []byte("2")))
		b.SendFrame(NewDataFrame(0, 0, 0, 3, 4, []byte("3")))
		b.SendFrame(NewDataFrame(0, 0, 0, 3, 4, []byte("4")))
	}()

	go func() {
		b, e := n1.Dial(0, 0, 1)
		if e != nil {
			debug.E("aaa", e)
		}
		b.SendFrame(NewDataFrame(0, 0, 1, 3, 4, []byte("1")))
		b.SendFrame(NewDataFrame(0, 0, 1, 3, 4, []byte("2")))
		b.SendFrame(NewDataFrame(0, 0, 1, 3, 4, []byte("3")))
		b.SendFrame(NewDataFrame(0, 0, 1, 3, 4, []byte("4")))
	}()
	time.Sleep(time.Second * 3)

}

// pass
func TestNode2(t *testing.T) {
	n1 := NewNode()

	go Helper("").CopyWithHandler(n1, n1, func(f Frame) Frame {
		if f.Command() == Request {
			f.SetCommand(Accept)
		}
		// debug.I("", f.String())
		return f
	})

	go func() {
		b, e := n1.Dial(0, 0, 0)
		if e != nil {
			debug.E("aaa", e)
		}
		go Helper("").ReadBus(b)
		b.SendFrame(NewDataFrame(0, 0, 0, 3, 4, []byte("1")))
		b.SendFrame(NewDataFrame(0, 0, 0, 3, 4, []byte("2")))
		b.SendFrame(NewDataFrame(0, 0, 0, 3, 4, []byte("3")))
		b.SendFrame(NewDataFrame(0, 0, 0, 3, 4, []byte("4")))
	}()

	go func() {
		b, e := n1.Dial(0, 0, 1)
		if e != nil {
			debug.E("aaa", e)
		}
		go Helper("").ReadBus(b)
		b.SendFrame(NewDataFrame(0, 0, 1, 3, 4, []byte("1")))
		b.SendFrame(NewDataFrame(0, 0, 1, 3, 4, []byte("2")))
		b.SendFrame(NewDataFrame(0, 0, 1, 3, 4, []byte("3")))
		b.SendFrame(NewDataFrame(0, 0, 1, 3, 4, []byte("4")))
	}()
	time.Sleep(time.Second * 4)

}

func TestNodeReader(t *testing.T) {
	n0 := NewNode()
	go func() {
		Helper("111").ReadBus(n0)
	}()

	b, e := n0.Dial(0, 0, 0)
	if e != nil {
		debug.E("e", e)
	}

	b.SendFrame(NewDataFrame(0, 0, 0, 3, 4, []byte("1")))
	b.SendFrame(NewDataFrame(0, 0, 0, 3, 4, []byte("2")))
	b.SendFrame(NewDataFrame(0, 0, 0, 3, 4, []byte("3")))
	b.SendFrame(NewDataFrame(0, 0, 0, 3, 4, []byte("4")))

	time.Sleep(time.Second)
}

// pass
func TestNodeEcho2(t *testing.T) {
	n0 := NewNode()
	n1 := NewNode()
	go Helper("").Copy(n0, n1)
	go Helper("").Copy(n1, n0)

	go func() {
		i := 1
		for {
			b, _ := n0.Accept()
			go Helper(strconv.Itoa(i)).CopyWithHandler(b, b, func(f Frame) Frame {
				// debug.I("", f.String())
				return f
			})
			// i++
		}
	}()

	go func() {
		b, e := n1.Dial(0, 0, 0)
		if e != nil {
			debug.E("aaa", e)
		}
		go Helper(strconv.Itoa(1)).ReadBus(b)
		b.SendFrame(NewDataFrame(0, 0, 0, 3, 4, []byte("1")))
		b.SendFrame(NewDataFrame(0, 0, 0, 3, 4, []byte("2")))
		b.SendFrame(NewDataFrame(0, 0, 0, 3, 4, []byte("3")))
		b.SendFrame(NewDataFrame(0, 0, 0, 3, 4, []byte("4")))
	}()

	go func() {
		b, e := n1.Dial(0, 0, 1)
		if e != nil {
			debug.E("aaa", e)
		}
		go Helper(strconv.Itoa(2)).ReadBus(b)
		b.SendFrame(NewDataFrame(0, 0, 1, 3, 4, []byte("1")))
		b.SendFrame(NewDataFrame(0, 0, 1, 3, 4, []byte("2")))
		b.SendFrame(NewDataFrame(0, 0, 1, 3, 4, []byte("3")))
		b.SendFrame(NewDataFrame(0, 0, 1, 3, 4, []byte("4")))
	}()
	time.Sleep(time.Second * 13)

}

// pzss
func TestNodeEcho(t *testing.T) {
	n0 := NewNode()
	n1 := NewNode()
	go Helper("").Copy(n0, n1)
	go Helper("").Copy(n1, n0)

	go func() {
		i := 1
		for {
			b, _ := n0.Accept()
			go Helper(strconv.Itoa(i)).Copy(b, b)
			// i++
		}
	}()

	go func() {
		b, e := n1.Dial(0, 0, 1)
		if e != nil {
			debug.E("aaa", e)
		}
		go Helper(strconv.Itoa(2)).ReadBus(b)
		b.SendFrame(NewDataFrame(0, 0, 1, 3, 4, []byte("1")))
		b.SendFrame(NewDataFrame(0, 0, 1, 3, 4, []byte("2")))
		b.SendFrame(NewDataFrame(0, 0, 1, 3, 4, []byte("3")))
		b.SendFrame(NewDataFrame(0, 0, 1, 3, 4, []byte("4")))
	}()
	time.Sleep(time.Second * 13)

}
