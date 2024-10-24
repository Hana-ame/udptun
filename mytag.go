package main

import (
	"encoding/binary"
	"fmt"
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
