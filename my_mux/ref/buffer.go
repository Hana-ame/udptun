// not tested.

package mymux

import (
	"sync"
)

type record struct {
	id   uint8
	data []byte
}

// buffer

// for go back n
type Buffer struct {
	size uint8

	buffer []*record
	valid  []bool

	head uint8
	tail uint8
	rptr uint8

	*sync.Cond

	closed bool
}

func NewGBNBuffer(size uint8) *Buffer {
	buf := &Buffer{
		size:   size,
		buffer: make([]*record, size),
		valid:  make([]bool, size),
		Cond:   sync.NewCond(&sync.Mutex{}),
	}
	return buf
}

func (b *Buffer) Offer(data []byte) {
	b.L.Lock()
	// 有效 且 未closed 时等待
	for (b.valid[b.head%b.size]) && !b.closed {
		b.Wait()
	}
	b.buffer[b.head%b.size] = &record{(b.head), data}
	b.valid[b.head%b.size] = true
	b.head++
	b.L.Unlock()
	b.Broadcast()
}

func (b *Buffer) Read() (uint8, []byte, bool) {
	b.L.Lock()
	// 无效 且 可以读取(即rptr小于head) 且 未closed 时等待
	for ((!b.valid[b.rptr%b.size]) || (b.rptr == b.head)) && !b.closed {
		b.Wait()
	}
	record, ok := b.buffer[b.rptr%b.size], b.valid[b.rptr%b.size]
	b.rptr++
	b.L.Unlock()
	b.Broadcast()
	return record.id, record.data, ok

}

func (b *Buffer) SetTail(tail uint8) {
	b.L.Lock()
	for tail-b.tail > 0 {
		b.valid[b.tail%b.size] = false
		b.tail++
	}
	b.L.Unlock()
	b.Broadcast()
}

func (b *Buffer) SetRead(rptr uint8) {
	b.L.Lock()
	b.rptr = rptr
	b.L.Unlock()
	b.Broadcast()
}

func (b *Buffer) Close() {
	b.L.Lock()
	b.closed = true
	b.L.Unlock()
	b.Broadcast()
}

// window no larger than 127
type MyBuffer struct {
	size   uint8
	buffer [][]byte
	valid  []bool

	head uint8
	tail uint8

	sync.Cond
	closed bool
}

func NewBuffer(size uint8) *MyBuffer {
	buf := &MyBuffer{
		size:   size,
		buffer: make([][]byte, size),
		valid:  make([]bool, size),
	}
	return buf
}

func (b *MyBuffer) IsEmpty() bool {
	return b.head == b.tail
}
func (b *MyBuffer) IsFull() bool {
	return b.head-b.tail == b.size
}

func (b *MyBuffer) Offer(data []byte) {
	for (b.valid[b.head]) && !b.closed {
		b.Wait()
	}
	b.buffer[b.head] = data
	b.valid[b.head] = true
	b.head++
	b.Broadcast()
}
func (b *MyBuffer) Poll() []byte {
	for !(!b.valid[b.tail] || b.tail != b.head || b.closed) {
		b.Wait()
	}
	data := b.buffer[b.tail]
	b.valid[b.tail] = false
	b.tail++
	b.Broadcast()
	return data
}

func (b *MyBuffer) SetTail(tail uint8) {
	for b.tail < tail {
		b.valid[b.tail] = false
		b.tail++
	}

	// b.head = tail + b.size
	b.Broadcast()
}

// func (b *MyBuffer) SetHead(head uint8) {
// 	b.head = head
// 	b.tail = head - b.size
// 	b.Broadcast()
// }

func (b *MyBuffer) Put(p uint8, data []byte) {
	for !(p-b.tail < b.size && !b.valid[p%b.size]) && !b.closed {
		b.Wait()
	}
	b.buffer[p%b.size] = data
	b.valid[p%b.size] = true
	b.Broadcast()
}
func (b *MyBuffer) Get(p uint8) []byte {
	for !(p-b.tail < b.size && b.valid[p%b.size]) && !b.closed {
		b.Wait()
	}
	data := b.buffer[p%b.size]
	b.Broadcast()
	return data
}

func (b *MyBuffer) Close() {
	b.closed = true
	b.Broadcast()
}
