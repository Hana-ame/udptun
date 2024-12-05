package mymux

import (
	"io"
	"testing"
	"time"

	"github.com/Hana-ame/udptun/Tools/debug"
)

func read(reader io.Reader) {
	b := make([]byte, 2048)
	for {
		n, err := reader.Read(b)
		if err != nil {
			debug.E("read", err)
		}
		debug.T("read", string(b[:n]))
	}

}

func TestFC(t *testing.T) {
	b0, b1 := NewBusPair()
	go Helper("").CopyWithHandler(b1, b1, func(f Frame) Frame {
		debug.T("f", f.String())
		return f
	})

	fc := NewFrameConn(b0, 0, 0, 0)
	go read(fc)
	fc.Write([]byte("a"))
	fc.Write([]byte("sdf2124df"))
	fc.Write([]byte("sdfd55y65f"))
	fc.Write([]byte("sdfd345435f"))

	time.Sleep(time.Second)
}

// 注意port
// 没有往对面发送Close指令
func TestFC2(t *testing.T) {
	n0, n1 := NewNode(), NewNode()
	go Helper("").Copy(n1, n0)
	go Helper("").Copy(n0, n1)

	go func() {
		for {
			b, port := n1.Accept()
			fc := NewFrameConn(b, 0, 0, port)
			go read(fc)
		}
	}()
	{
		b, err := n0.Dial(0, 0, 0)
		if err != nil {
			t.Fatal(err)
		}
		fc := NewFrameConn(b, 0, 0, 0)
		fc.Write([]byte("123"))
		fc.Write([]byte("456"))
		fc.Write([]byte("789"))
		fc.Write([]byte("000"))
	}
	{
		b, err := n0.Dial(0, 0, 2)
		if err != nil {
			t.Fatal(err)
		}
		fc := NewFrameConn(b, 0, 0, 2)
		fc.Write([]byte("123"))
		fc.Write([]byte("456"))
		fc.Write([]byte("789"))
		fc.Write([]byte("000"))

		fc.Close()
	}
	{
		b, err := n0.Dial(0, 0, 1)
		if err != nil {
			t.Fatal(err)
		}
		fc := NewFrameConn(b, 0, 0, 1)
		fc.Write([]byte("123"))
		fc.Write([]byte("456"))
		fc.Write([]byte("789"))
		fc.Write([]byte("000"))
	}
}
