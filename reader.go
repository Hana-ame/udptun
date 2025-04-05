// 失去任何引用时是否会回收？

package main

import (
	"fmt"
)

type FrameReader struct {
	// sync.Mutex
	// errCh  chan error
	rcvCh  chan frame
	closed bool
}

func NewFrameReader() *FrameReader {
	return &FrameReader{
		// Mutex: sync.Mutex{},
		// errCh: make(chan error),
		rcvCh: make(chan frame),
	}
}

func (r *FrameReader) ReadFrame() (frame, error) {
	f, ok := <-r.rcvCh
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

func (r *FrameReader) WriteFrame(f frame) (err error) {
	defer func() {
		recover()
		err = fmt.Errorf("closed")
	}()
	if r.closed {
		return fmt.Errorf("closed")
	}
	r.rcvCh <- f
	return nil
}

// func (r *FrameReader) WriteError(err error) error {
// 	if r.closed {
// 		return fmt.Errorf("closed")
// 	}
// 	r.errCh <- err
// 	return nil
// }

func (r *FrameReader) Close() error {
	if r.closed {
		return nil
	}
	r.closed = true
	close(r.rcvCh)
	// r.errCh <- fmt.Errorf("closed")
	return nil
}

type Reader struct {
}
