package main

import "encoding/binary"

type flag byte

// 通过最后一位判断是发送方还是接收方
type addr [4]byte

type frame []byte

func (f frame) src() addr {
	return addr(f[0:4])
}
func (f frame) dst() addr {
	return addr(f[4:8])
}

func (f frame) seq() uint16 {
	return binary.BigEndian.Uint16(f[8:10])
}

func (f frame) ack() uint16 {
	return binary.BigEndian.Uint16(f[8:10])
}

func (f frame) flags() flag {
	return flag(f[10])
}
