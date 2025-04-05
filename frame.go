package main

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

type flag byte

const (
	FLAG_CLOSE flag = 1 << iota

	FLAG_DATA flag = 0
)

func (flag flag) Close() bool { return (flag & FLAG_CLOSE) != 0 }

func (flag flag) String() string { return strconv.Itoa(int(flag)) }

// 通过最后一位判断是发送方还是接收方
type Addr uint32

func (addr Addr) IsClient() bool { return (addr & 1) != 0 }
func (addr Addr) IsServer() bool { return (addr & 1) == 0 }
func (addr Addr) tag() uint32    { return (uint32(addr) & (^uint32(1))) }
func (addr Addr) String() string { return strconv.Itoa(int(addr)) }

const (
	MTU = 1024
)

type frame []byte

func (f frame) src() Addr {
	return Addr(binary.BigEndian.Uint32(f[0:4]))
}
func (f frame) Dst() Addr {
	return Addr(binary.BigEndian.Uint32(f[4:8]))
}

func (f frame) seq() uint16 {
	return binary.BigEndian.Uint16(f[8:10])
}

func (f frame) ack() uint16 {
	return binary.BigEndian.Uint16(f[10:12])
}

func (f frame) Flags() flag {
	return flag(f[12])
}

func (f frame) data() []byte {
	return f[13:]
}

func (f frame) dataLength() int {
	return len(f) - 13
}

func Frame(src, dst Addr, seq, ack uint16, flag flag, data []byte) frame {
	f := make([]byte, 13+len(data))
	binary.BigEndian.PutUint32(f[0:4], uint32(src))
	binary.BigEndian.PutUint32(f[4:8], uint32(dst))
	binary.BigEndian.PutUint16(f[8:10], seq)
	binary.BigEndian.PutUint16(f[10:12], ack)
	f[12] = byte(flag)
	copy(f[13:], data)
	return f
}

func (f frame) String() string {
	return fmt.Sprintf("src:%d,dst:%d,seq:%d,ack:%d,flag:%d,dataLength:%d",
		f.src(), f.Dst(), f.seq(), f.ack(), f.Flags(), f.dataLength())
}
