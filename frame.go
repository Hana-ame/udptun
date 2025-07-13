package main

import "encoding/binary"

const (
	HEADER_OFFSET = 12

	SYN byte = 0b0000_0001
	FIN byte = 0b0000_0010
)

type flag []byte

func (flag flag) isSYN() bool {
	return flag[0] == SYN
}

func (flag flag) setSYN(value bool) {
	if value {
		flag[0] = flag[0] | SYN
	} else {
		flag[0] = ^((^flag[0]) & (^SYN))
	}
}

func (flag flag) isFIN() bool {
	return flag[0] == FIN
}

func (flag flag) setFIN(value bool) {
	if value {
		flag[0] = flag[0] | FIN
	} else {
		flag[0] = ^((^flag[0]) & (^FIN))
	}
}

type frame []byte

func (f frame) src() uint32 {
	return binary.BigEndian.Uint32(f[0:4])
}

func (f frame) dst() uint32 {
	return binary.BigEndian.Uint32(f[4:8])
}

func (f frame) reserved() byte {
	return f[8]
}

func (f frame) flag() flag {
	return flag(f[9:10])
}

func (f frame) dataLength() uint16 {
	return binary.BigEndian.Uint16(f[10:12])
}

func (f frame) data() []byte {
	return f[HEADER_OFFSET:]
}
