package main

type Conn struct {
	// 相对于application而言的
	ConnReader *Pipe
	ConnWriter *Pipe
}

func NewConn() *Conn {
	return &Conn{
		ConnReader: NewPipe(),
		ConnWriter: NewPipe(),
	}
}

func (c *Conn) ApplicatonInterface() FrameHandlerInterface {
	return FrameHandlerInterface{
		push:  c.ConnWriter.Push,
		poll:  c.ConnReader.Poll,
		close: c.ConnWriter.Close,
	}
}
func (c *Conn) RouterInterface() FrameHandlerInterface {
	return FrameHandlerInterface{
		push:  c.ConnReader.Push,
		poll:  c.ConnWriter.Poll,
		close: c.ConnReader.Close,
	}
}
