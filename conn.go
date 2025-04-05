package main

import (
	"io"
	"time"
)

type FrameConn struct {
	src       Addr
	dst       Addr
	seq       uint16
	ack       uint16
	WriteChan Chan
	ReadChan  Chan
}

func NewFrameConn(src, dst Addr) *FrameConn {
	return &FrameConn{
		src, dst, 0, 0,
		NewFrameChan(),
		NewFrameChan(),
	}
}

func (c *FrameConn) WriteFrame(f frame) (err error) {
	return c.WriteChan.WriteFrame(f)
}

func (c *FrameConn) ReadFrame() (frame, error) {
	frame, err := c.ReadChan.ReadFrame()
	if err == nil && frame.Flags().Close() {
		err = io.EOF
	}
	return frame, err
}

func (c *FrameConn) Close() error {
	c.WriteFrame(Frame(c.src, c.dst, c.seq, c.ack, FLAG_CLOSE, nil))
	go func() {
		time.Sleep(time.Minute)
		c.ReadChan.Close()
	}()
	return c.WriteChan.Close()
	// return DoubleError(c.WriteChan.Close(), c.ReadChan.Close(), "write", "read")
}

type FrameConnReader struct {
	buf []byte
	*FrameConn
}

func (r *FrameConnReader) Read(buf []byte) (n int, err error) {
	if len(r.buf) == 0 {
		// 从通道获取新帧
		frame, err := r.ReadFrame()
		if err != nil {
			return 0, err
		}
		r.ack = frame.seq()  // 验证序列号TODO
		r.buf = frame.data() // 假设frame有Data()方法返回[]byte
		if frame.Flags().Close() {
			r.FrameConn.ReadChan.Close()
			return 0, io.EOF
		}
	}

	// 拷贝数据到用户缓冲区
	n = copy(buf, r.buf)
	r.buf = r.buf[n:] // 更新剩余数据

	return n, nil
}

func (r *FrameConn) Reader() *FrameConnReader {
	return &FrameConnReader{nil, r}
}

type FrameConnWriter struct {
	// buf []byte
	*FrameConn
}

func (w *FrameConnWriter) Write(buf []byte) (n int, err error) {
	if len(buf) > MTU {
		buf = buf[:MTU]
	}
	err = w.WriteFrame(Frame(w.src, w.dst, w.seq, w.ack, 0, buf))
	if err != nil {
		return
	}
	w.seq++
	return len(buf), err
}

func (w *FrameConn) Writer() *FrameConnWriter {
	return &FrameConnWriter{w}
}

type Conn struct {
	*FrameConnWriter
	*FrameConnReader
}

func (c *Conn) Close() error {
	return DoubleError(c.FrameConnWriter.Close(), c.FrameConnReader.Close(), "write", "read")
}
