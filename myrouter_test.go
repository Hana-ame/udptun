package main

import (
	"testing"

	"github.com/Hana-ame/udptun/Tools/debug"
)

func TestRoute(t *testing.T) {
	pong0 := NewPongBus(99)
	pong1 := NewPongBus(1)
	pong2 := NewPongBus(2)

	r := NewRouter(99, pong0)
	r.Add(pong1)
	r.Add(pong2)

	debug.I(r.ConcurrentHashMap)

	DebugHelperReader("png1", pong1)
}
