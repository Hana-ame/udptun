package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
)

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

// const (
// 	ANY_ADDR Addr = 65535
// )

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
	Data  Command = 0      // 数据命令
	Close Command = 1      // 关闭命令
	Aloha Command = 1 << 2 // Aloha 命令

	Pong Command = 1 << 3
	// UDP 用于无序数据
	Disorder            Command = 1 << 4
	DisorderAcknowledge Command = 1<<4 | 1
	// MUX 命令
	Request Command = 1<<6 | 1 // 请求命令
	Accept  Command = 1<<6 | 2 // 确认命令
)

// MyFrame 定义了一个帧类型，基于字节切片。
// src, src, dst, dst, port, cmd, seqn, ackn, data...
type MyFrame []byte

// NewCtrlFrame 创建一个新的控制帧。
func NewCtrlFrame(source, destination Addr, port uint8, command Command, sequenceNumber, acknowledgeNumber uint8) MyFrame {
	f := make(MyFrame, FrameHeadLength) // 创建一个指定长度的帧
	f.SetSource(source)
	f.SetDestination(destination)
	f.SetPort(port)
	f.SetCommand(command)
	f.SetSequenceNumber(sequenceNumber)
	f.SetAcknowledgeNumber(acknowledgeNumber)
	return f
}

// NewDataFrame 创建一个新的数据帧。
func NewDataFrame(source, destination Addr, port uint8, sequenceNumber, acknowledgeNumber uint8, data []byte) MyFrame {
	f := make(MyFrame, FrameHeadLength+len(data)) // 创建一个包含数据长度的帧
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
func NewFrame(source, destination Addr, port uint8, command Command, sequenceNumber, acknowledgeNumber uint8, data []byte) MyFrame {
	f := make(MyFrame, FrameHeadLength+len(data)) // 创建一个包含数据长度的帧
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
func PrintFrame(f MyFrame) {
	if len(f) < FrameHeadLength {
		return // 如果帧长度小于头部长度，则返回
	}
	log.Printf("%d->%d:%d,%s, %s\n",
		f.Source(), f.Destination(),
		f.Port(), f.Command().String(),
		f.Data())
}

// SprintFrame 返回帧的详细信息字符串。
func SprintFrame(f MyFrame) string {
	if len(f) < FrameHeadLength {
		return "" // 如果帧长度小于头部长度，则返回空字符串
	}
	return fmt.Sprintf("%d->%d:%d,%s,[%d,%d], {%s}",
		f.Source(), f.Destination(),
		f.Port(), f.Command().String(),
		f.SequenceNumber(), f.AcknowledgeNumber(),
		f.Data())
}

// Tag 方法获取帧的标签。
// src, src, dst, dst, port
func (f MyFrame) Tag() MyTag {
	var tag MyTag
	copy(tag[:], f[:TagLength]) // 复制前 TagLength 字节作为标签
	return tag
}

// 获取源地址
func (f MyFrame) Source() Addr {
	return Addr(binary.BigEndian.Uint16(f[0:2])) // 使用大端字节序读取源地址
}

// 获取目的地址
func (f MyFrame) Destination() Addr {
	return Addr(binary.BigEndian.Uint16(f[2:4])) // 使用大端字节序读取目的地址
}

// 获取端口
func (f MyFrame) Port() uint8 {
	return f[4] // 端口在第 5 个字节
}

// 获取命令类型
func (f MyFrame) Command() Command {
	return Command(f[5]) // 命令在第 6 个字节
}

// 获取序列号
func (f MyFrame) SequenceNumber() uint8 {
	return f[6] // 序列号在第 7 个字节
}

// 获取确认号
func (f MyFrame) AcknowledgeNumber() uint8 {
	return f[7] // 确认号在第 8 个字节
}

// 获取数据内容
func (f MyFrame) Data() []byte {
	return f[FrameHeadLength:] // 数据内容从第 9 个字节开始
}

// 设置源地址
func (f MyFrame) SetSource(source Addr) {
	binary.BigEndian.PutUint16(f[0:2], uint16(source)) // 使用大端字节序写入源地址
}

// 设置目的地址
func (f MyFrame) SetDestination(destination Addr) {
	binary.BigEndian.PutUint16(f[2:4], uint16(destination)) // 使用大端字节序写入目的地址
}

// 设置端口
func (f MyFrame) SetPort(port uint8) {
	f[4] = port // 设置端口
}

// 设置命令
func (f MyFrame) SetCommand(command Command) {
	f[5] = byte(command) // 设置命令
}

// 设置序列号
func (f MyFrame) SetSequenceNumber(sequence uint8) {
	f[6] = sequence // 设置序列号
}

// 设置确认号
func (f MyFrame) SetAcknowledgeNumber(acknowledge uint8) {
	f[7] = acknowledge // 设置确认号
}

// 设置数据内容
func (f MyFrame) SetData(data []byte) int {
	if len(data) > 0 {
		return copy(f[FrameHeadLength:], data) // 将数据复制到帧的适当位置
	}
	return 0
}
