// 失去任何引用时是否会回收？

package main

import (
	"fmt"
)

type Chan interface {
	ReadFrame() (frame, error)
	WriteFrame(frame) error
	Close() error
}

type FrameChan struct {
	// sync.Mutex
	// errCh  chan error
	ch     chan frame
	closed bool
}

func NewFrameChan() *FrameChan {
	return &FrameChan{
		// Mutex: sync.Mutex{},
		// errCh: make(chan error),
		ch: make(chan frame),
	}
}

func (r *FrameChan) ReadFrame() (frame, error) {
	f, ok := <-r.ch
	if ok {
		return f, nil
	}
	return nil, fmt.Errorf("closed")
	// if r.closed {
	// 	return nil, fmt.Errorf("closed")
	// }
	// select {
	// case f, ok := <-r.rcvCh:
	// 	if ok {
	// 		return f, nil
	// 	}
	// 	r.errCh <- fmt.Errorf("rcv channel closed")
	// 	return nil, fmt.Errorf("rcv channel closed")
	// case err, ok := <-r.errCh:
	// 	if ok {
	// 		return nil, err
	// 	}
	// 	return nil, fmt.Errorf("err channel closed")
	// }
}

func (w *FrameChan) WriteFrame(f frame) (err error) {
	defer func() {
		if e := recover(); e == nil {
			return
		}
		err = fmt.Errorf("closed (recover)")
	}()
	if w.closed {
		return fmt.Errorf("closed")
	}
	w.ch <- f
	return nil
}

// func (r *FrameReader) WriteError(err error) error {
// 	if r.closed {
// 		return fmt.Errorf("closed")
// 	}
// 	r.errCh <- err
// 	return nil
// }

func (c *FrameChan) Close() error {
	if c.closed {
		return nil
	}
	c.closed = true
	close(c.ch)
	// r.errCh <- fmt.Errorf("closed")
	return nil
}
