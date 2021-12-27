package main

import (
	"bytes"
	"encoding/binary"
	"errors"
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

/*
func _main() {
	Conn, err := net.ListenPacket("udp", fmt.Sprintf("0.0.0.0:%d", 12321))
	if err != nil {
		log.Fatal("sb")
		return
	}
	fmt.Println(GetAddr(Conn))

	addr, err := net.ResolveUDPAddr("udp", "34.145.70.165:12421")
	if err != nil {
		log.Printf("error : %v", err)
		return
	}
	fmt.Println("1")
	time.Sleep(time.Second * 3)
	Conn.WriteTo([]byte{0}, addr)
	Conn.WriteTo([]byte{0}, addr)
	Conn.WriteTo([]byte{0}, addr)
	Conn.WriteTo([]byte{0}, addr)
	Conn.WriteTo([]byte{0}, addr)
	fmt.Println("2")

	buffer := make([]byte, 2048)
	for {
		fmt.Println(3)
		n, addr, err := Conn.ReadFrom(buffer)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("packet-received: bytes=%d from=%s\n", n, addr.String())
		// portalproxy.PrintHex(buffer[:n])
	}
}
*/
// copy 111223
func GetAddr(conn net.PacketConn) (string, error) {
	addr1, err := GetAddress(conn, "stun1.l.google.com:19302")
	if err != nil {
		log.Printf("error : %v", err)
		return "", err
	}
	addr2, err := GetAddress(conn, "stun2.l.google.com:19302")
	if err != nil {
		log.Printf("error : %v", err)
		return "", err
	}
	if addr1 == addr2 {
		return addr1, nil
	}
	return "", errors.New("严格的nat类型")
}

// copy 111223
func GetAddress(conn net.PacketConn, server string) (string, error) {
	addr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		log.Printf("error : %v", err)
		return "", err
	}
	sp := stunPack{
		1,
		0,
		0x2112A442,
		[12]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
	}
	conn.WriteTo(sp.Bytes(), addr)
	buffer := make([]byte, 2048)
	n, _, err := conn.ReadFrom(buffer)
	if err != nil {
		log.Fatal(err.Error())
		return "", err
	}

	return xorAddr(buffer[n-6 : n]), nil
}

// copy 111223
func xorAddr(b []byte) string {
	if len(b) != 6 {
		return ""
	}
	port := binary.BigEndian.Uint16(b[0:2])
	port ^= 0x2112
	ip := net.IPv4(b[2]^0x21, b[3]^0x12, b[4]^0xa4, b[5]^0x42)
	return fmt.Sprintf("%s:%d", ip, port)
}

// BigEncoder 111223
func BE(v ...interface{}) []byte {
	buf := bytes.Buffer{}
	for _, i := range v {
		// fmt.Println(i)
		buf.Write(BEbytes(i))
	}
	return buf.Bytes()
}

// tested 11223
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
