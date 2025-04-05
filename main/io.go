package main

import "fmt"

type FrameWriter struct {
	FrameHandler

	MTU int
}

func (w *FrameWriter) Write(p []byte) (n int, err error) {
	if len(p) > int(w.MTU) {
		return 0, fmt.Errorf("framewriter, too large")
	}
	f := NewFrame(0, 0, 0, 0, 0, 0, p)
	err = w.Push(f)
	return len(p), err
}

type FrameReader struct {
	FrameHandler
}

func (r *FrameReader) Read(p []byte) (n int, err error) {
	f, err := r.Poll()
	if err != nil {
		return
	}
	n = copy(p, f.Data())
	return n, err
}

type FrameReadWirteCloser struct {
	FrameReader
	FrameWriter
}

func (c *FrameReadWirteCloser) Close() error {
	return c.FrameWriter.Close()
}
