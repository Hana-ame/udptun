package mymux

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

// TestMutex 测试互斥锁的工作原理。
func TestMutex(t *testing.T) {
	reader, writer := NewPipe()           // 创建读写管道
	bus := NewBusFromPipe(reader, writer) // 从管道创建总线
	// reader.Lock()                         // 锁定读取器
	bus.Lock()           // 锁定总线
	fmt.Println("never") // 这行代码不会被执行，因为锁定后没有解锁
}

// TestCond 测试条件变量的工作原理。
func TestCond(t *testing.T) {
	var b bool
	var c sync.Cond = *sync.NewCond(&sync.Mutex{}) // 创建条件变量

	// 启动第一个 goroutine
	go func() {
		c.L.Lock() // 锁定条件变量的互斥锁
		for !b {   // 当 b 为 false 时等待
			c.Wait()                 // 等待条件变量
			log.Println("waiting 1") // 打印日志
		}
		c.L.Unlock() // 解锁

		time.Sleep(time.Second) // 等待 1 秒
		c.Signal()              // 发出信号，唤醒一个等待的 goroutine
		time.Sleep(time.Second) // 再等待 1 秒
		b = false               // 设置 b 为 false
	}()

	// 启动第二个 goroutine
	go func() {
		c.L.Lock() // 锁定条件变量的互斥锁
		for !b {   // 当 b 为 false 时等待
			c.Wait()                 // 等待条件变量
			log.Println("waiting 0") // 打印日志
		}
		c.L.Unlock() // 解锁

		time.Sleep(time.Second) // 等待 1 秒
		b = false               // 设置 b 为 false
		time.Sleep(time.Second) // 再等待 1 秒
		c.Signal()              // 发出信号，唤醒一个等待的 goroutine
	}()

	time.Sleep(time.Second * 2) // 主 goroutine 等待 2 秒
	b = true                    // 设置 b 为 true
	c.Broadcast()               // 广播信号，唤醒所有等待的 goroutine
	c.L.Lock()                  // 锁定条件变量的互斥锁
	for b {                     // 当 b 为 true 时等待
		c.Wait()                 // 等待条件变量
		log.Println("waiting 2") // 打印日志
	}
	c.L.Unlock()      // 解锁
	log.Println("OK") // 打印完成消息
}
