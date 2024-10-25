package main

import (
	"testing"

	"github.com/Hana-ame/udptun/Tools/debug"
)

func DebugHelperReader(tag string, reader MyBusReader) error {
	for {
		f, e := reader.RecvFrame()
		if e != nil {
			return e
		}

		debug.I(tag, SprintFrame(f))
	}
}

func TestErr(t *testing.T) {
	reader1, writer2 := NewSyncPipe()
	reader2, writer1 := NewSyncPipe()
	bus1 := NewBus(reader1, writer1)
	bus2 := NewBus(reader2, writer2)
	go Copy(bus2.MyBusWriter, bus2.MyBusReader)
	go DebugHelperReader("bus1", bus1)
	bus1.SendFrame(NewFrame(
		0, 1, 2, Data, 0, 0, []byte("1111"),
	))
	bus1.SendFrame(NewFrame(
		0, 1, 2, Data, 0, 0, []byte("1111"),
	))
	bus1.SendFrame(NewFrame(
		0, 1, 2, Data, 0, 0, []byte("1111"),
	))
}
