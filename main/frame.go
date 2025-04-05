package main

// 通过最后一位判断是发送方还是接收方
type addr [4]byte

type frame []byte

func (f frame) Src() addr {
	return addr(f[0:4])
}
func (f frame) Dst() addr {
	return addr(f[4:8])
}
