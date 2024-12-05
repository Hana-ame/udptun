package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
)

// src src dst dst port
const TagLength = 5

// src src dst dst port
type Tag [TagLength]byte

// src src dst dst port
func NewTag(src, dst Addr, port uint8) Tag {
	var tag Tag
	binary.BigEndian.PutUint16(tag[0:2], uint16(src))
	binary.BigEndian.PutUint16(tag[2:4], uint16(dst))
	tag[4] = port // 设置端口
	return tag
}

// src src dst dst port
func (f Tag) Tag() Tag {
	return f
}

// src src dst dst port
func (f Tag) String() string {
	src := binary.BigEndian.Uint16(f[0:2])
	dst := binary.BigEndian.Uint16(f[2:4])
	port := f[4]

	return fmt.Sprintf("%d->%d:%d", src, dst, port)
}

// Command 定义了一个命令类型，基于 uint8。
type Command uint8

// String 方法将 Command 转换为字符串表示。
func (cmd Command) String() string {
	switch cmd {
	case Data:
		return "Data"
	case Close:
		return "Close"
	case Aloha:
		return "Aloha"
	case Request:
		return "Request"
	case Accept:
		return "Acknowledge"
	case Disorder:
		return "Disorder"
	case DisorderAcknowledge:
		return "DisorderAcknowledge"
	default:
		return "Unknown"
	}
}

// Addr 定义了一个地址类型，基于 uint16。
type Addr uint16

// NetWork 方法返回地址的网络类型。
func (a Addr) NetWork() string {
	return "mymux"
}

// String 方法将 Addr 转换为字符串表示。
func (a Addr) String() string {
	return strconv.Itoa(int(a))
}

// 常量定义
const (
	FrameHeadLength int = 8 // 帧头长度

	// 命令定义
	Data  Command = 0 // 数据命令
	Close Command = 1 // 关闭命令
	Aloha Command = 4 // Aloha 命令
	// UDP 用于无序数据
	Disorder            Command = 1 << 4
	DisorderAcknowledge Command = 1<<4 | 1
	// MUX 命令
	Request Command = 1<<6 | 1 // 请求命令
	Accept  Command = 1<<6 | 2 // 确认命令
)

// Frame 定义了一个帧类型，基于字节切片。
// src, src, dst, dst, port, cmd, seqn, ackn, data...
type Frame []byte

// NewCtrlFrame 创建一个新的控制帧。
func NewCtrlFrame(source, destination Addr, port uint8, command Command, sequenceNumber, acknowledgeNumber uint8) Frame {
	f := make(Frame, FrameHeadLength) // 创建一个指定长度的帧
	f.SetSource(source)
	f.SetDestination(destination)
	f.SetPort(port)
	f.SetCommand(command)
	f.SetSequenceNumber(sequenceNumber)
	f.SetAcknowledgeNumber(acknowledgeNumber)
	return f
}

// NewDataFrame 创建一个新的数据帧。
func NewDataFrame(source, destination Addr, port uint8, sequenceNumber, acknowledgeNumber uint8, data []byte) Frame {
	f := make(Frame, FrameHeadLength+len(data)) // 创建一个包含数据长度的帧
	f.SetSource(source)
	f.SetDestination(destination)
	f.SetPort(port)
	f.SetCommand(Data)
	f.SetSequenceNumber(sequenceNumber)
	f.SetAcknowledgeNumber(acknowledgeNumber)
	f.SetData(data)
	return f
}

// NewFrame 创建一个新的帧，可以指定数据和命令。
func NewFrame(source, destination Addr, port uint8, command Command, sequenceNumber, acknowledgeNumber uint8, data []byte) Frame {
	f := make(Frame, FrameHeadLength+len(data)) // 创建一个包含数据长度的帧
	f.SetSource(source)
	f.SetDestination(destination)
	f.SetPort(port)
	f.SetCommand(command)
	f.SetSequenceNumber(sequenceNumber)
	f.SetAcknowledgeNumber(acknowledgeNumber)
	f.SetData(data)
	return f
}

// PrintFrame 打印帧的详细信息。
func PrintFrame(f Frame) {
	if len(f) < FrameHeadLength {
		return // 如果帧长度小于头部长度，则返回
	}
	log.Printf("%d->%d:%d,%s, %s\n",
		f.Source(), f.Destination(),
		f.Port(), f.Command().String(),
		f.Data())
}

// SprintFrame 返回帧的详细信息字符串。
func SprintFrame(f Frame) string {
	if len(f) < FrameHeadLength {
		return "" // 如果帧长度小于头部长度，则返回空字符串
	}
	return fmt.Sprintf("%d->%d:%d,%s,[%d,%d], {%s}",
		f.Source(), f.Destination(),
		f.Port(), f.Command().String(),
		f.SequenceNumber(), f.AcknowledgeNumber(),
		f.Data())
}

func (f Frame) String() string {
	return SprintFrame(f)
}

// Tag 方法获取帧的标签。
// src, src, dst, dst, port
func (f Frame) Tag() Tag {
	var tag Tag
	copy(tag[:], f[:TagLength]) // 复制前 TagLength 字节作为标签
	return tag
}

// 获取源地址
func (f Frame) Source() Addr {
	return Addr(binary.BigEndian.Uint16(f[0:2])) // 使用大端字节序读取源地址
}

// 获取目的地址
func (f Frame) Destination() Addr {
	return Addr(binary.BigEndian.Uint16(f[2:4])) // 使用大端字节序读取目的地址
}

// 获取端口
func (f Frame) Port() uint8 {
	return f[4] // 端口在第 5 个字节
}

// 获取命令类型
func (f Frame) Command() Command {
	return Command(f[5]) // 命令在第 6 个字节
}

// 获取序列号
func (f Frame) SequenceNumber() uint8 {
	return f[6] // 序列号在第 7 个字节
}

// 获取确认号
func (f Frame) AcknowledgeNumber() uint8 {
	return f[7] // 确认号在第 8 个字节
}

// 获取数据内容
func (f Frame) Data() []byte {
	return f[FrameHeadLength:] // 数据内容从第 9 个字节开始
}

// 设置源地址
func (f Frame) SetSource(source Addr) {
	binary.BigEndian.PutUint16(f[0:2], uint16(source)) // 使用大端字节序写入源地址
}

// 设置目的地址
func (f Frame) SetDestination(destination Addr) {
	binary.BigEndian.PutUint16(f[2:4], uint16(destination)) // 使用大端字节序写入目的地址
}

// 设置端口
func (f Frame) SetPort(port uint8) {
	f[4] = port // 设置端口
}

// 设置命令
func (f Frame) SetCommand(command Command) {
	f[5] = byte(command) // 设置命令
}

// 设置序列号
func (f Frame) SetSequenceNumber(sequence uint8) {
	f[6] = sequence // 设置序列号
}

// 设置确认号
func (f Frame) SetAcknowledgeNumber(acknowledge uint8) {
	f[7] = acknowledge // 设置确认号
}

// 设置数据内容
func (f Frame) SetData(data []byte) int {
	return copy(f[FrameHeadLength:], data) // 将数据复制到帧的适当位置
}
