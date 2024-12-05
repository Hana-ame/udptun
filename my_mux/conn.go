package mymux

type FrameConn struct {
	Bus

	localAddr  Addr
	remoteAddr Addr
	port       uint8

	MTU int
}

func NewFrameConn(bus Bus, localAddr, remoteAddr Addr, port uint8) *FrameConn {
	c := &FrameConn{
		Bus: bus,

		localAddr:  localAddr,
		remoteAddr: remoteAddr,
		port:       port,

		MTU: 1024,
	}
	return c
}

// 没写大于MTU的处理
func (c *FrameConn) Write(p []byte) (n int, err error) {
	f := NewFrame(c.localAddr, c.remoteAddr, c.port, Disorder, 0, 0, p)
	// f := make(Frame, FrameHeadLength+len(p)) // 创建一个包含数据长度的帧
	// f.SetSource(c.localAddr)
	// f.SetDestination(c.remoteAddr)
	// f.SetPort(c.port)
	// f.SetCommand(Disorder)
	// f.SetSequenceNumber(sequenceNumber)
	// f.SetAcknowledgeNumber(acknowledgeNumber)
	n = len(p)
	err = c.SendFrame(f)
	return
}

// 没写小于MTU的处理
func (c *FrameConn) Read(p []byte) (n int, err error) {
	f, err := c.RecvFrame()
	if err != nil {
		return
	}
	n = copy(p, f.Data())
	return
}
