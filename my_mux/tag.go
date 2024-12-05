package mymux

import (
	"encoding/binary"
	"fmt"
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
