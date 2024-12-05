package mymux

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/Hana-ame/udptun/Tools/debug"
	"github.com/gorilla/websocket"
)

const (
	ERR_BUS_CLOSED Error = "my bus already closed" // 总线关闭错误信息
)

const (
	ERR_CONN_CLOSED Error = "my frame conn closed"
)

func ErrorIsClosed(e error) bool {
	err := Error(e.Error())
	return err == ERR_BUS_CLOSED || err == ERR_PIPE_CLOSED || err == ERR_CONN_CLOSED
}

type MyFrameConn struct {
	MyBus

	localAddr  Addr
	remoteAddr Addr
	port       uint8

	closed bool

	MTU int // for body
}

func NewFrameConn(bus MyBus, localAddr, remoteAddr Addr, port uint8) *MyFrameConn {
	c := &MyFrameConn{
		MyBus: bus,

		localAddr:  localAddr,
		remoteAddr: remoteAddr,
		port:       port,

		MTU: 1024,
	}
	return c
}

func (c *MyFrameConn) WriteFrame(p []byte) (n int, err error) {
	const Tag = "MyFrameConn.WriteFrame"
	debug.T(Tag, c.localAddr, "->", c.remoteAddr, ":", c.port, string(p))
	if c.closed {
		debug.D(Tag, c.localAddr, "->", c.remoteAddr, ":", c.port, "conn closed")
		err = fmt.Errorf("closed")
		return
	}
	if len(p) > c.MTU {
		p = p[:c.MTU]
	}
	f := NewFrame(c.localAddr, c.remoteAddr, c.port, Disorder, 0, 0, p)

	n = len(p)
	err = c.MyBus.SendFrame(f)
	return
}

// 需要大于MTU
// 从ReadBuf里面取到纯净的Data
func (c *MyFrameConn) ReadFrame() ([]byte, error) {
	const Tag = "MyFrameConn.ReadFrame"
	if c.closed {
		return nil, (ERR_CONN_CLOSED)
	}

	f, err := c.MyBus.RecvFrame()
	if err != nil {
		return nil, err
	}
	debug.T(Tag, c.localAddr, "<-", c.remoteAddr, ":", c.port, f.Command().String())
	if f.Command() == Close {
		defer c.Close()
		return nil, (ERR_CONN_CLOSED)
	}
	debug.T(Tag, c.localAddr, "<-", c.remoteAddr, ":", c.port, string(f.Data()))
	return f.Data(), nil
}

// close
func (c *MyFrameConn) Close() error {
	const Tag = "MyFrameConn.Close"
	debug.D(Tag, c.localAddr, "<-", c.remoteAddr, ":", c.port, "closing")
	defer debug.D(Tag, c.localAddr, "<-", c.remoteAddr, ":", c.port, "closed")

	if c.closed {
		return (ERR_CONN_CLOSED)
	}
	c.SendFrame(NewCtrlFrame(c.localAddr, c.remoteAddr, c.port, Close, 0, 0))
	// time.Sleep(time.Second) // it seems that close cannot send, so sleep and
	c.MyBus.Close()
	// c.MyMux.PrintMap() // debug 加了这句client Close不能
	c.closed = true
	return nil
}

// for net.Conn interface
func (c *MyFrameConn) LocalAddr() Addr {
	return c.localAddr
}
func (c *MyFrameConn) RemoteAddr() Addr {
	return c.remoteAddr
}

func (c *MyFrameConn) SetDeadline(t time.Time) error {
	return fmt.Errorf("todo")
}
func (c *MyFrameConn) SetReadDeadline(t time.Time) error {
	return fmt.Errorf("todo")
}
func (c *MyFrameConn) SetWriteDeadline(t time.Time) error {
	return fmt.Errorf("todo")
}

// 插口，专门把FreamConn转换为io.Streamer
type MyFrameConnStreamer struct {
	*MyFrameConn

	rb []byte
}

func (c *MyFrameConnStreamer) Write(p []byte) (n int, err error) {
	return c.WriteFrame(p)
}

func (c *MyFrameConnStreamer) Read(p []byte) (n int, err error) {
	if len(c.rb) == 0 {
		c.rb, err = c.ReadFrame()
		if err != nil {
			return
		}
	}
	n = copy(p, c.rb)
	c.rb = c.rb[n:]
	return
}

// 这里开始没什么关系，可能用到TCP的东西再说。
type MyConn struct {
	MyBus

	// MyTag

	localAddr  Addr
	remoteAddr Addr
	Port       uint8

	sequenceNumber uint8 // 即将发送的frame的Seq number
	requestingSeq  uint8 // 对方要求的最近的Seq num

	ReadBuf     chan MyFrame // 先做简单的
	nextReadSeq uint8        // 己方维护的自卷积要求的最近的Seq num

	MTU        int
	WindowSize int // 用于更新acknowledgeNumber

	closed bool
}

func NewConn(mux MyBus, frameTag MyTag, localAddr, remoteAddr Addr, port uint8) *MyConn {
	conn := &MyConn{
		MyBus: mux,
		// MyTag:          frameTag,
		localAddr:      localAddr,
		remoteAddr:     remoteAddr,
		Port:           port,
		sequenceNumber: 0,
		requestingSeq:  0,
		ReadBuf:        make(chan MyFrame),
		nextReadSeq:    0,
		MTU:            1024,
		WindowSize:     32,
		closed:         false,
	}
	return conn
}

// c.localAddr, c.remoteAdr, c.port
func (c *MyConn) Tag() MyTag {
	var tag MyTag
	binary.BigEndian.PutUint16(tag[0:2], uint16(c.remoteAddr))
	binary.BigEndian.PutUint16(tag[2:4], uint16(c.localAddr))
	tag[4] = c.Port
	return tag
}

// 会限制不能大于MTU
// 封装成DataFrame从Mux发送
func (c *MyConn) Write(p []byte) (n int, err error) {
	if c.closed {
		err = fmt.Errorf("closed")
		return
	}
	if len(p) > c.MTU {
		p = p[:c.MTU]
	}
	f := NewDataFrame(c.localAddr, c.remoteAddr, c.Port, c.sequenceNumber, c.nextReadSeq, p)

	n = len(p)
	err = c.MyBus.SendFrame(f)
	return
}

// 需要大于MTU
// 从ReadBuf里面取到纯净的Data
func (c *MyConn) Read(p []byte) (n int, err error) {
	if c.closed {
		err = fmt.Errorf("closed")
		return
	}
	f := <-c.ReadBuf

	if f.Command() == Close {
		return 0, io.EOF
	}

	// 不是Close也不是其他frame，DataFrame根据状态来的
	// 更新最后收到的帧
	if f.AcknowledgeNumber()-c.requestingSeq < uint8(c.WindowSize) {
		c.requestingSeq = f.AcknowledgeNumber()
	}
	c.nextReadSeq = f.SequenceNumber() // 这个需要稍后改一下。

	n = copy(p, f.Data())
	return
}

func (c *MyConn) Close() error {
	// debug
	const Tag = "MyConn.Close"
	debug.T(Tag, "closing")
	defer debug.T(Tag, "closed")

	if c.closed {
		return fmt.Errorf("closed")
	}
	// 给ReadBuf发送一个Close的CtrlFrame，读到就直接EOF
	c.ReadBuf <- MyFrame(NewCtrlFrame(0, 0, 0, Close, 0, 0))
	c.SendFrame(NewCtrlFrame(c.localAddr, c.remoteAddr, c.Port, Close, c.sequenceNumber, c.nextReadSeq))
	// c.MyBus.RemoveConn(c)
	c.MyBus.Close()
	// c.MyMux.PrintMap() // debug 加了这句client Close不能
	return nil
}

// for mux
// 从这里接受Frame到缓冲区
func (c *MyConn) PutFrame(f MyFrame) {
	// 及时Close
	if f.Command() == Close {
		c.Close()
		return
	}

	c.ReadBuf <- f
}

// for net.Conn interface
func (c *MyConn) LocalAddr() Addr {
	return c.localAddr
}
func (c *MyConn) RemoteAddr() Addr {
	return c.remoteAddr
}

func (c *MyConn) SetDeadline(t time.Time) error {
	return fmt.Errorf("todo")
}
func (c *MyConn) SetReadDeadline(t time.Time) error {
	return fmt.Errorf("todo")
}
func (c *MyConn) SetWriteDeadline(t time.Time) error {
	return fmt.Errorf("todo")
}

// MyConnBus 用于 TCP 连接的总线结构。
type MyConnBus struct {
	net.Conn

	sync.Mutex // 仅允许一个读取守护进程读取。
}

// RecvFrame 从连接中接收一帧数据。
func (b *MyConnBus) RecvFrame() (MyFrame, error) {
	// 获取帧长度
	l := make([]byte, 2)
	_, err := b.Read(l)
	if err != nil {
		return nil, err
	}
	pl := binary.BigEndian.Uint16(l)
	// 获取帧内容
	f := make([]byte, pl)
	_, err = b.Read(f)
	return MyFrame(f), err
}

// SendFrame 发送一帧数据到连接。
func (b *MyConnBus) SendFrame(f MyFrame) error {
	l := make([]byte, 2)
	binary.BigEndian.PutUint16(l, uint16(len(f)))
	if _, err := b.Write(l); err != nil {
		return err
	}
	if _, err := b.Write(f); err != nil {
		return err
	}
	return nil
}

// MyWsBus 用于 WebSocket 连接的总线结构。
type MyWsBusReader struct {
	*websocket.Conn

	sync.Mutex
}

// MyWsBus 用于 WebSocket 连接的总线结构。
type MyWsBusWriter struct {
	*websocket.Conn

	sync.Mutex
}

type MyWsBus struct {
	*websocket.Conn

	*MyWsBusReader
	*MyWsBusWriter

	sync.Mutex // 仅允许一个读取守护进程读取。
}

func NewWsBus(c *websocket.Conn) MyBus {
	// func NewWsBus(c *websocket.Conn) *MyWsBus {
	return &MyWsBus{
		Conn:          c,
		MyWsBusReader: &MyWsBusReader{Conn: c},
		MyWsBusWriter: &MyWsBusWriter{Conn: c},
	}
}

// RecvFrame 从 WebSocket 连接接收一帧数据。
func (b *MyWsBusReader) RecvFrame() (MyFrame, error) {
	b.Lock()
	defer b.Unlock()
	_, f, err := b.ReadMessage()
	return MyFrame(f), err
}

// SendFrame 通过 WebSocket 发送一帧数据。
func (b *MyWsBusWriter) SendFrame(f MyFrame) error {
	b.Lock()
	defer b.Unlock()
	err := b.WriteMessage(websocket.BinaryMessage, f)
	return err
}

func (b *MyWsBus) Close() error {
	return b.Conn.Close()
}

// MyPipeBus 本地管道总线结构。
type MyPipeBus struct {
	MyBusReader
	MyBusWriter

	closed bool // 标记总线是否已关闭

	sync.Mutex
}

// Close 关闭总线，释放相关资源。
func (b *MyPipeBus) Close() error {
	// const Tag = "MyPipeBus.Close"
	if b.closed {
		// 		debug.E(Tag, "already closed")
		return ERR_BUS_CLOSED
	}
	b.closed = true
	b.MyBusReader.Close() // 关闭读取器
	b.MyBusWriter.Close() // 关闭写入器
	return nil
}

// NewBusFromPipe 创建一个新的管道总线实例。
func NewBusFromPipe(reader MyBusReader, writer MyBusWriter) *MyPipeBus {
	return &MyPipeBus{
		MyBusReader: reader,
		MyBusWriter: writer,
	}
}

// NewPipeBusPair 创建一对本地管道总线。
func NewPipeBusPair() (*MyPipeBus, *MyPipeBus) {
	a2bReader, b2aWriter := NewPipe()              // 创建 a 到 b 的读写管道
	b2aReader, a2bWriter := NewPipe()              // 创建 b 到 a 的读写管道
	a2bBus := NewBusFromPipe(a2bReader, a2bWriter) // 创建 a 到 b 的总线
	b2aBus := NewBusFromPipe(b2aReader, b2aWriter) // 创建 b 到 a 的总线
	return a2bBus, b2aBus
}

// NewDebugPipeBusPair 创建一对带调试信息的本地管道总线。
func NewDebugPipeBusPair(tag string) (*MyPipeBus, *MyPipeBus) {
	a2bReader, b2aWriter := NewDebugPipe(tag)      // 创建带调试信息的管道
	b2aReader, a2bWriter := NewDebugPipe(tag)      // 创建带调试信息的管道
	a2bBus := NewBusFromPipe(a2bReader, a2bWriter) // 创建 a 到 b 的总线
	b2aBus := NewBusFromPipe(b2aReader, b2aWriter) // 创建 b 到 a 的总线
	return a2bBus, b2aBus
}

// 在正常情况下能够传输 见test
// 如果有什么问题遇到的时候再来debug
type ReliableBus struct {
	MyBus

	f      MyFrame
	e      error
	nextId uint8

	*Buffer
	request uint8

	*sync.Cond
}

func NewReliableBus(b MyBus, size uint8) *ReliableBus {
	rb := &ReliableBus{
		MyBus: b,

		Buffer: NewGBNBuffer(size),

		Cond: sync.NewCond(&sync.Mutex{}),
	}

	go rb.ReadDaemon()
	go rb.WriteDeamon()
	go rb.AcknowledgeDeamon()
	return rb
}

func (b *ReliableBus) SendFrame(f MyFrame) error {
	if f.Command() == Disorder || f.Command() == DisorderAcknowledge {
		b.Offer(f)
		return nil
	}
	return b.MyBus.SendFrame(f)
}
func (b *ReliableBus) RecvFrame() (MyFrame, error) {
	b.L.Lock()
	for !(b.f != nil || b.closed) {
		b.Wait()
	}
	if b.closed {
		b.L.Unlock()
		return b.f, ERR_BUS_CLOSED
	}
	f, e := b.f, b.e
	b.f, b.e = nil, nil

	b.L.Unlock()
	b.Broadcast()
	return f, e
}

func (b *ReliableBus) ReadDaemon() {
	// const Tag = "ReliableBus.ReadDaemon"
	for {
		f, e := b.MyBus.RecvFrame()
		// 		debug.T(Tag, "recv Frame", SprintFrame(f))
		b.L.Lock()
		for !(b.f == nil || b.closed) {
			b.Wait()
		}
		if b.closed {
			b.L.Unlock()
			return
		}
		if f.Command() == Disorder {
			// 如果是disorder，那么在bus处处理。
			if f.SequenceNumber() == b.nextId {
				b.f, b.e = f, e
				b.nextId++
				// 				debug.T(Tag, "b.nextid = ", b.nextId)
			}
		}
		if f.Command() == DisorderAcknowledge || f.Command() == Disorder {
			// 			debug.T(Tag, b.request, " should set to ", f.AcknowledgeNumber())
			if b.request-f.AcknowledgeNumber() > b.size {
				// 				debug.T(Tag, b.request, " set to ", f.AcknowledgeNumber())
				b.request = f.AcknowledgeNumber()
				b.Buffer.SetTail(b.request)
			}
		} else {
			b.f, b.e = f, e
		}

		b.L.Unlock()
		b.Broadcast()
	}
}
func (b *ReliableBus) WriteDeamon() {
	// const Tag = "ReliableBus.WriteDeamon"

	for {
		id, data, ok := b.Buffer.Read() // 在buffer里的一定是disorder
		if !ok {
			if b.closed {
				return
			}
			// 			debug.E(Tag, id, data, ok)
			continue
		}
		f := MyFrame(data)
		f.SetSequenceNumber(id)
		f.SetAcknowledgeNumber(b.nextId)

		e := b.MyBus.SendFrame(f)
		if e != nil {
			if b.closed {
				return
			}
			// 			debug.E(Tag, e.Error())
			continue
		}
	}
}

func (b *ReliableBus) AcknowledgeDeamon() {
	for {
		time.Sleep(time.Second)
		f := NewFrame(0, 0, 0, DisorderAcknowledge, 0, b.nextId, nil)
		e := b.MyBus.SendFrame(f)
		if e != nil {
			if b.closed {
				return
			}
			// 				debug.E(Tag, "send frame error", e.Error())
			continue
		}
		// 			debug.T(Tag, "requesting", b.request)
		b.Buffer.SetTail(b.request) // 顺手在这里设置了
	}
}

func (b *ReliableBus) Close() error {
	e := b.MyBus.Close()
	b.Buffer.Close()
	b.closed = true
	b.Broadcast()
	return e
}

func (b *ReliableBus) Lock() {
	b.MyBus.Lock()
}

func (b *ReliableBus) UnLock() {
	b.MyBus.Unlock()
}
