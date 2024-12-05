package mymux

import (
	"strconv"
	"testing"
	"time"

	"github.com/Hana-ame/udptun/Tools/debug"
)

// // not tested.
// func TestBus(t *testing.T) {
// 	muxReader, busWriter := NewPipe()
// 	busReader, muxWriter := NewPipe()

// 	mux := NewMuxServer(muxWriter, 1)
// 	go mux.ReadDaemon(muxReader)

// 	bus := NewBusFromPipe(busReader, busWriter)

// 	router := NewRouter()
// 	_ = bus
// 	_ = mux
// 	_ = router
// }

// 正常情况下能够传输的。
func TestReliableBus(t *testing.T) {
// 	a, b := NewDebugPipeBusPair("123")
// 	ra := NewReliableBus(a, 4)
// 	rb := NewReliableBus(b, 4)
// 	go func() {
// 		// 创建一个字节数组
// 		for i := 22; i < 1234; i++ {
// 			// 将整数转换为字符串
// 			str := strconv.Itoa(i)

// 			ra.SendFrame(NewFrame(0, 0, 0, Disorder, 0, 0, []byte(str)))
// 			time.Sleep(time.Second)
// 		}
// 	}()
// 	go func() {
// 		for {
// 			f, e := rb.RecvFrame()
// 			debug.I("rb", SprintFrame(f), e)
// 			time.Sleep(time.Second)
// 		}
// 	}()
// 	select {}
// }

