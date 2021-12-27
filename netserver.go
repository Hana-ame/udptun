package main

import (
	"fmt"
	"net"
)

func UdpServer(port int) {
	pc, err := net.ListenPacket("udp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := make([]byte, 2048)
	for {
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("toolbox : udp from ", addr.String())
		fmt.Printf("toolbox : udp len=%d,%v\n", n, buf[:n])
	}
}
func TcpEcho(port int) {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err.Error())
			// handle error
		}
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer func() {
		c.Close()
		fmt.Println("toolbox : tcp close ", c.RemoteAddr().String())
	}()
	fmt.Println("toolbox : tcp from ", c.RemoteAddr().String())
	buf := make([]byte, 2048)
	for {
		n, err := c.Read(buf)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("toolbox : tcp len=%d,%v\n", n, buf[:n])
		_, err = c.Write(buf[:n])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
