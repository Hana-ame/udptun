// 设计坏了，server只能固定address，浮动port，相当于只做了一个conn的分流
// 先写一个监听所有addr，能分流的server来

package main

import (
	"fmt"

	tools "github.com/Hana-ame/udptun/Tools"
)

type Sucker struct {
	FrameHandler

	FramePushHandler
}

func NewSucker(i FrameHandler, dst FramePushHandler) *Sucker {
	go func() {
		defer i.Close()
		for {
			f, err := i.Poll()
			if err != nil {
				return
			}
			err = dst.Push(f)
			if err != nil {
				return
			}
		}
	}()
	return &Sucker{FrameHandler: i, FramePushHandler: dst}
}

// 只能自己建立监听。
// 必须对地址严查。

type Router struct {
	FrameHandler // 默认路由
	*tools.ConcurrentHashMap[Addr, FrameHandler]
}

func NewRouter(i FrameHandler) *Router {
	return &Router{
		FrameHandler:      i,
		ConcurrentHashMap: tools.NewConcurrentHashMap[Addr, FrameHandler](),
	}
}

func (r *Router) Push(f Frame) error {
	return r.GetOrDefault(f.Destination(), r.FrameHandler).Push(f)
}

func (r *Router) Poll() (f Frame, err error) {
	return nil, fmt.Errorf("not supported")
}

func (r *Router) Close() error {
	return fmt.Errorf("todo")
}
