// 只能适配MyConn了，
// 大概会弃用

package mymux

import (
	"encoding/binary"
	"fmt"
	"log"
	"sync"

	tools "github.com/Hana-ame/udptun/Tools"
)

// 定义标签的长度
const TagLength = 5

// MyTag 定义了一个标签类型，使用字节数组。
type MyTag [TagLength]byte

// NewTag 创建一个新的标签，包含远程地址、本地地址和端口信息。
func NewTag(src, dst Addr, port uint8) MyTag {
	var tag MyTag
	binary.BigEndian.PutUint16(tag[0:2], uint16(src))
	binary.BigEndian.PutUint16(tag[2:4], uint16(dst))
	tag[4] = port // 设置端口
	return tag
}

// Tag 方法返回标签本身。
func (f MyTag) Tag() MyTag {
	return f
}

func (f MyTag) String() string {
	src := binary.BigEndian.Uint16(f[0:2])
	dst := binary.BigEndian.Uint16(f[2:4])
	port := f[4]

	return fmt.Sprintf("%d->%d:%d", src, dst, port)
}

// MyMux 接口定义了一个多路复用器。
type MyMux interface {
	MyBus
	RemoveConn(*MyConn)
}

// MyMuxServer 定义了一个多路复用服务器。
type MyMuxServer struct {
	MyBusWriter
	localAddr Addr

	// 使用并发哈希表存储连接
	*tools.ConcurrentHashMap[MyTag, *MyConn]
	acceptedConnChannel chan *MyConn // 存储已接受的连接通道
}

// NewMuxServer 创建一个新的多路复用服务器实例。
func NewMuxServer(writer MyBusWriter, localAddr Addr) *MyMuxServer {
	mux := &MyMuxServer{
		MyBusWriter:         writer,
		localAddr:           localAddr,
		ConcurrentHashMap:   tools.NewConcurrentHashMap[MyTag, *MyConn](),
		acceptedConnChannel: make(chan *MyConn),
	}
	return mux
}

// RemoveConn 从哈希表中移除连接。
func (m *MyMuxServer) RemoveConn(c *MyConn) {
	m.Remove(c.Tag())
}

// Accept 从接受的连接通道中读取连接。
func (m *MyMuxServer) Accept() *MyConn {
	return <-m.acceptedConnChannel
}

// ReadDaemon 读取并处理连接的帧。
func (m *MyMuxServer) ReadDaemon(c MyBus) {
	c.Lock()
	defer c.Unlock()

	for {
		f, _ := c.RecvFrame()
		switch f.Command() {
		case Request:
			// 只响应目标地址为本地地址的请求
			if f.Destination() != m.localAddr {
				continue
			}
			// 创建新连接
			if _, exist := m.Get(f.Tag()); !exist {
				cBus, _ := NewPipeBusPair()
				c := NewConn(cBus, f.Tag(), f.Destination(), f.Source(), f.Port())
				m.Put(c.Tag(), c)
				m.acceptedConnChannel <- c // 将新连接发送到通道
			}
			// 发送确认帧
			m.SendFrame(NewCtrlFrame(f.Destination(), f.Source(), f.Port(), Accept, 0, 0))

		case Accept:
			continue

		default:
			// 其他情况直接转发帧
			if conn, exist := m.Get(f.Tag()); exist {
				conn.PutFrame(f)
			} else {
				if f.Command() == Close {
					continue // 如果是关闭命令则不返回关闭帧
				}
				m.SendFrame(NewCtrlFrame(f.Destination(), f.Source(), f.Port(), Close, 0, 0))
			}
		}
	}
}

// PrintMap 打印当前连接映射。
func (m *MyMuxServer) PrintMap() {
	log.Println("print mux map", m.localAddr)
	m.ConcurrentHashMap.ForEach(func(key MyTag, value *MyConn) {
		fmt.Println(key, value)
	})
}

// MyMuxClient 定义了一个多路复用客户端。
type MyMuxClient struct {
	MyBusWriter
	sync.Mutex

	localAddr                                Addr
	*tools.ConcurrentHashMap[MyTag, *MyConn] // 使用并发哈希表存储连接

	nextport uint8 // 下一个可用端口
}

// NewMuxClient 创建一个新的多路复用客户端实例。
func NewMuxClient(writer MyBusWriter, localAddr Addr) *MyMuxClient {
	mux := &MyMuxClient{
		MyBusWriter:       writer,
		localAddr:         localAddr,
		ConcurrentHashMap: tools.NewConcurrentHashMap[MyTag, *MyConn](),
	}
	return mux
}

// RemoveConn 从哈希表中移除连接。
func (m *MyMuxClient) RemoveConn(c *MyConn) {
	m.Remove(c.Tag())
}

// Dial 发起一个连接请求。
func (m *MyMuxClient) Dial(dst Addr) (*MyConn, error) {
	m.Lock()
	defer m.Unlock()

	if m.Size() > 254 {
		return nil, fmt.Errorf("no other ports") // 如果端口已用尽则返回错误
	}
	f := NewCtrlFrame(m.localAddr, dst, m.nextport, Request, 0, 0)
	// 查找可用端口
	for m.Contains(f.Tag()) || m.nextport == 0 {
		m.nextport++
		f.SetPort(m.nextport)
	}
	cBus, _ := NewPipeBusPair()
	c := NewConn(cBus, f.Tag(), m.localAddr, dst, m.nextport)

	// 发送连接请求
	m.SendFrame(f)
	m.PutIfAbsent(c.Tag(), c)

	m.nextport++

	return c, nil
}

// ReadDaemon 读取并处理连接的帧。
func (m *MyMuxClient) ReadDaemon(c MyBus) {
	c.Lock()
	defer c.Unlock()

	for {
		f, _ := c.RecvFrame()
		switch f.Command() {
		case Request:
			// 拒绝连接请求
			m.SendFrame(NewCtrlFrame(f.Destination(), f.Source(), f.Port(), Close, 0, 0))
			continue
		case Accept:
			continue
		default:
			// 其他情况直接转发帧
			if conn, exist := m.Get(f.Tag()); exist {
				conn.PutFrame(f)
			} else {
				if f.Command() == Close {
					continue // 如果是关闭命令则不返回关闭帧
				}
				m.SendFrame(NewCtrlFrame(f.Destination(), f.Source(), f.Port(), Close, 0, 0))
			}
		}
	}
}

// PrintMap 打印当前连接映射。
func (m *MyMuxClient) PrintMap() {
	log.Println("print mux map", m.localAddr)
	m.ConcurrentHashMap.ForEach(func(key MyTag, value *MyConn) {
		fmt.Println(key, value)
	})
}
