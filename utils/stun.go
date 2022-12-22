package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

type stunPack struct {
	MessageType          uint16
	MessageLength        uint16
	MessageCookie        uint32
	MessageTransactionID [12]uint8
}

func (p *stunPack) Bytes() []byte {
	return BE(
		p.MessageType,
		p.MessageLength,
		p.MessageCookie,
		p.MessageTransactionID,
	)
}

// BigEncoder
func BE(v ...interface{}) []byte {
	buf := bytes.Buffer{}
	for _, i := range v {
		// fmt.Println(i)
		buf.Write(BEbytes(i))
	}
	return buf.Bytes()
}

func BEbytes(v interface{}) []byte {
	if va, ok := v.(uint16); ok {
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, va)
		return b
	}
	if va, ok := v.(uint32); ok {
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, va)
		return b
	}
	if va, ok := v.([12]byte); ok {
		return va[:]
	}
	if va, ok := v.([]byte); ok {
		return va[:]
	}
	return nil
}

func xorAddr(b []byte) (string, error) {
	if len(b) != 6 {
		return "", fmt.Errorf("xorAddr: len(b)!= 6")
	}
	port := binary.BigEndian.Uint16(b[0:2])
	port ^= 0x2112
	ip := net.IPv4(b[2]^0x21, b[3]^0x12, b[4]^0xa4, b[5]^0x42)
	return fmt.Sprintf("%s:%d", ip, port), nil
}

func StunRequest(server string, conn *net.UDPConn) {
	addr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		log.Printf("error : %v", err)
		return
	}
	sp := stunPack{
		1,
		0,
		0x2112A442,
		[12]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
	}
	conn.WriteTo(sp.Bytes(), addr)
}

func StunResolve(data []byte) (string, error) {
	n := len(data)
	if n-6 < 0 {
		return "", fmt.Errorf("invalid stun message")
	}
	return xorAddr(data[n-6 : n])
}

func GetOutboundIPv6(lc *net.UDPConn) string {
	conn, err := net.Dial("udp", "[2606:4700:4700::1111]:53")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return fmt.Sprintf("[%s]:%d", localAddr.IP.String(), lc.LocalAddr().(*net.UDPAddr).Port)

}
