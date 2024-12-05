package wsreverse

import (
	"testing"
	"time"

	"github.com/Hana-ame/udptun/Tools/debug"
	mymux "github.com/Hana-ame/udptun/Tools/my_mux"
)

func TestBuffer(t *testing.T) {
	b := mymux.NewGBNBuffer(8)
	go func() {
		var i byte = 0
		for {
			time.Sleep(time.Second)
			debug.I("offer", i)
			b.Offer([]byte{i})
			debug.I("offer", i)
			i++
		}
	}()
	go func() {
		for i := 0; i < 2; i++ {
			id, data, ok := b.Read()
			debug.I("read", id, data, ok)
		}
		time.Sleep(time.Second * 8)
		for i := 0; i < 2; i++ {
			id, data, ok := b.Read()
			debug.I("read", id, data, ok)
		}
		time.Sleep(time.Second * 8)
		for i := 0; i < 2; i++ {
			id, data, ok := b.Read()
			debug.I("read", id, data, ok)
		}
	}()
	go func() {
		time.Sleep(time.Second * 8)
		time.Sleep(time.Second * 8)
		b.SetTail(4)
	}()

	time.Sleep(time.Hour)
}

func TestBuffer2(t *testing.T) {
	b := mymux.NewGBNBuffer(8)
	go func() {
		var i byte = 0
		for {
			debug.I("offer", i)
			b.Offer([]byte{i})
			debug.I("offer", i)
			i++
		}
	}()
	go func() {
		for {
			id, data, ok := b.Read()
			debug.I("read", id, data, ok)
		}
	}()
	go func() {
		for i := 1; i < 100; i++ {
			time.Sleep(time.Second / 10)
			b.SetTail(uint8(i * 8))
		}

	}()

	time.Sleep(time.Hour)
}
