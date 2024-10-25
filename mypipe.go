package main

import (
	"sync"

	"github.com/Hana-ame/udptun/Tools/debug"
)

// MySyncPipe 定义了一个管道结构，用于在读取和写入之间传递数据。
type MySyncPipe struct {
	*sync.Cond // 条件变量，用于同步
	sync.Mutex // 互斥锁，确保只有一个读取协程在运行

	f      MyFrame // 存储要传递的帧
	closed bool    // 管道是否关闭
}

// SendFrame 发送帧到管道。
func (p *MySyncPipe) SendFrame(f MyFrame) (err error) {
	p.L.Lock() // 锁定互斥锁
	// 当帧不为空且管道未关闭时，等待
	for (p.f != nil) && !p.closed {
		p.Wait()
	}
	if p.closed {
		err = ERR_CLOSED
	}
	p.f = f // 设置帧
	p.L.Unlock()

	p.Broadcast() // 唤醒等待的协程
	return
}

// RecvFrame 从管道接收帧。
func (p *MySyncPipe) RecvFrame() (f MyFrame, err error) {
	p.L.Lock() // 锁定互斥锁
	// 当帧为空且管道未关闭时，等待
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
	return
}

// Close 关闭管道并广播唤醒所有等待的协程。
func (p *MySyncPipe) Close() (err error) {
	p.L.Lock() // 锁定互斥锁
	// 要sending最后一个package，这是为了正常关闭。
	for p.f != nil && !p.closed {
		p.Wait()
	}
	if p.closed {
		err = ERR_CLOSED
	}
	p.closed = true
	p.L.Unlock()

	p.Broadcast() // 广播信号
	return
}

// NewSyncPipe 创建一个新的管道，返回读取和写入接口。
func NewSyncPipe() (MyBusReader, MyBusWriter) {
	pipe := &MySyncPipe{Cond: sync.NewCond(&sync.Mutex{})}
	return pipe, pipe // 返回同一个管道的读写接口
}

func NewDebugPipe(tag string) (MyBusReader, MyBusWriter) {
	reader := &MySyncPipe{Cond: sync.NewCond(&sync.Mutex{})}
	writer := &MySyncPipe{Cond: sync.NewCond(&sync.Mutex{})}
	go func() {
		for {
			f, e := writer.RecvFrame()
			if e != nil {
				debug.E(tag, e) // 记录错误
				continue
			}
			debug.D(tag, SprintFrame(f)) // 打印帧内容
			e = reader.SendFrame(f)      // 发送帧到读管道
			if e != nil {
				debug.E(tag, e) // 记录错误
			}

		}
	}()
	return reader, writer // 返回同一个管道的读写接口
}
