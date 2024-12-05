package main

type ConnConfig struct {
	localAddr  Addr
	remoteAddr Addr
	port       uint8
}

type ConnReader struct {
	*ConnConfig

	*Pipe // 如果卡住了, 需要改成Buffer
}

type ConnWriter struct {
	*ConnConfig

	*Pipe // 如果卡住了, 需要改成Buffer
}

func (c *ConnWriter) Push(f Frame) error {
	f.SetSource(c.localAddr)
	f.SetDestination(c.remoteAddr)
	f.SetPort(c.port)
	return c.Pipe.Push(f)
}

type Conn struct {
	ConnReader
	ConnWriter
}

func NewConn(local, destination Addr, port uint8) *Conn {
	config := &ConnConfig{
		localAddr:  local,
		remoteAddr: destination,
		port:       port,
	}
	return &Conn{
		ConnReader: ConnReader{ConnConfig: config, Pipe: NewPipe()},
		ConnWriter: ConnWriter{ConnConfig: config, Pipe: NewPipe()},
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
