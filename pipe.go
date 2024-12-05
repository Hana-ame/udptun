package main

import (
	"fmt"
	"sync"
)

type Pipe struct {
	*sync.Cond

	f      Frame
	closed bool
}

func NewPipe() *Pipe {
	return &Pipe{
		Cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (p *Pipe) Push(f Frame) (err error) {
	p.L.Lock() // 锁定互斥锁
	// 帧未读且未关闭时等待, 这个条件要等待
	for p.f != nil && !p.closed {
		p.Wait()
	}
	if p.closed {
		err = fmt.Errorf("pipe is closed")
	} else {
		p.f = f
	}
	p.L.Unlock()

	p.Broadcast()
	return
}

func (p *Pipe) Poll() (f Frame, err error) {
	p.L.Lock()
	// 帧存在且未关闭时等待, 这个条件要等待
	for p.f == nil && !p.closed {
		p.Wait()
	}
	if p.closed && p.f == nil {
		err = fmt.Errorf("pipe is closed")
	}
	f = p.f   // 获取帧
	p.f = nil // 清空帧
	p.L.Unlock()

	p.Broadcast() // 唤醒等待的协程
	return f, err
}

func (p *Pipe) Close() error {
	p.closed = true

	p.Broadcast()
	return nil
}

func Copy(dst, src FrameHandler) error {
	defer src.Close()
	defer dst.Close()
	for {
		f, err := src.Poll()
		if err != nil {
			return err
		}
		err = dst.Push(f)
		if err != nil {
			return err
		}
	}
}
