package main

const tagLength = 2 // 定义标签的长度为 2 字节

// PortalBuf 是一个字节切片，用于存储带标签的数据
type PortalBuf []byte

// Tag 返回 PortalBuf 的前 tagLength 字节作为标签
func (p PortalBuf) Tag() PortalBuf {
	if len(p) >= tagLength {
		return (PortalBuf)(p)[:tagLength] // 返回标签部分
	}
	return nil // 如果长度不足，返回 nil
}

// AddTag 向 PortalBuf 添加标签，支持 string、[]byte 和 int 类型
func (p PortalBuf) AddTag(tag any) PortalBuf {
	if len(p) >= tagLength { // 确保 PortalBuf 足够长以存储标签
		switch t := tag.(type) {
		case string:
			// 将字符串转换为字节并添加标签
			if len(t) >= tagLength {
				copy(p[:tagLength], []byte(t)[:tagLength]) // 复制前 2 字节
			}
		case []byte:
			// 确保字节切片长度足够
			if len(t) >= tagLength {
				copy(p[:tagLength], t[:tagLength]) // 复制前 2 字节
			}
		case int:
			// 将整数转换为两个字节
			p[0] = byte(t % 256) // 低字节
			p[1] = byte(t / 256) // 高字节
		default:
			return nil // 不支持的类型，返回 nil
		}
		return (PortalBuf)(p) // 返回修改后的 PortalBuf
	}
	return nil // 如果长度不足，返回 nil
}

// Data 返回数据部分，根据 n 的值返回不同长度的数据
// n <= 0 返回所有数据部分
// n > 0 返回长度为 n 的数据部分
func (p PortalBuf) Data(n int) PortalBuf {
	if len(p) > tagLength { // 确保有足够的长度
		if n > 0 {
			if len(p) >= tagLength+n {
				return (PortalBuf)(p)[tagLength : tagLength+n] // 返回指定长度的数据
			}
			return (PortalBuf)(p)[tagLength:] // 如果数据长度不足，返回所有数据
		}
		return (PortalBuf)(p)[tagLength:] // 返回所有数据
	}
	return nil // 如果长度不足，返回 nil
}

// Raw 返回 PortalBuf 的前 n 个字节
func (p PortalBuf) Raw(n int) PortalBuf {
	if n > len(p) {
		n = len(p) // 确保不超过 PortalBuf 的长度
	}
	return (p)[:n] // 返回前 n 字节
}

// DataAndTag 返回包含标签和数据部分的切片
// 返回长度为 tagLength+n 的切片
func (p PortalBuf) DataAndTag(n int) PortalBuf {
	if n > len(p)-tagLength {
		n = len(p) - tagLength // 确保不会超出范围
	}
	return (p)[:n+tagLength] // 返回标签和数据部分
}
