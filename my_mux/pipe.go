package mymux

import (
	"io"
	"sync"

	log "github.com/Hana-ame/udptun/Tools/debug"
)

type IPipe[T any] interface {
	Send(data T) (err error)
	Recv() (data T, err error)

	io.Closer
}

// Pipe 定义了一个管道结构，用于在读取和写入之间传递数据。
type Pipe struct {
	*sync.Cond // 条件变量，用于同步

	f      Frame // 存储要传递的帧
	closed bool  // 管道是否关闭
}

// SendFrame 发送帧到管道。
func (p *Pipe) SendFrame(f Frame) (err error) {
	p.L.Lock() // 锁定互斥锁
	// 当帧不为空且管道未关闭时，等待
	for p.f != nil && !p.closed {
		p.Wait()
	}
	if p.closed {
		err = ERR_CLOSED
	}
	p.f = f
	p.L.Unlock()

	p.Broadcast()
	return
}

func (p *Pipe) RecvFrame() (f Frame, err error) {
	p.L.Lock()
	for p.f == nil && !p.closed {
		p.Wait()
	}

	if p.closed {
		err = ERR_CLOSED
	}
	f = p.f   // 获取帧
	p.f = nil // 清空帧
	p.L.Unlock()

	p.Broadcast() // 唤醒等待的协程
	return f, err
}

// Close 关闭管道并广播唤醒所有等待的协程。
func (p *Pipe) Close() error {
	p.L.Lock()
	for p.f != nil && !p.closed {
		p.Wait()
	}
	p.closed = true
	p.L.Unlock()
	p.Broadcast()
	return nil
}

// NewPipe 创建一个新的管道，返回读取和写入接口。
func NewPipe() (BusReader, BusWriter) {
	pipe := &Pipe{Cond: sync.NewCond(&sync.Mutex{})}
	return pipe, pipe // 返回同一个管道的读写接口
}

// edit before use.
func NewDebugPipe(tag string) (BusReader, BusWriter) {
	r := &Pipe{Cond: sync.NewCond(&sync.Mutex{})}
	w := &Pipe{Cond: sync.NewCond(&sync.Mutex{})}
	go func() {
		for {
			f, e := w.RecvFrame() // 从写管道接收帧
			if e != nil {
				log.E(tag, e) // 记录错误
				return
			} else if len(f) < FrameHeadLength {
				log.W(tag, "length = ", len(f)) // 记录警告
			} else {
				log.D(tag, SprintFrame(f)) // 打印帧内容
				e = r.SendFrame(f)         // 发送帧到读管道
				if e != nil {
					log.E(tag, e) // 记录错误
					return
				}
			}
		}
	}()
	return r, w // 返回读写管道
}
