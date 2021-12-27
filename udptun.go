package main

import (
	"encoding/binary"
	"flag"
	"log"
	"net"
	"time"

	"github.com/pkg/errors"
)

var pool map[int]*net.UDPConn
var poolc map[int]*net.Addr

// ./udptun.exe -l ":6000" -r "127.0.0.1:4000"
func main() {
	var mode string
	var laddr string
	var raddr string
	flag.StringVar(&mode, "mode", "raw", "client: send with work as client (mux) or\nwork as server (demux)\nnone: do nothing")
	flag.StringVar(&laddr, "l", ":6000", "local addr")
	flag.StringVar(&raddr, "r", "127.0.0.1:40000", "remote addr")
	flag.Parse()

	if mode == "raw" {
		raw(laddr, raddr)
	} else if mode == "server" {
		server2(laddr, raddr)
		// server(laddr, raddr)
	} else if mode == "client" {
		client2(laddr, raddr)
		// client(laddr, raddr)
	}
}

// 添加打洞
// 改写raddr
func server2(laddr, raddr string) { // 接受，去头，传送
	pool = make(map[int]*net.UDPConn)

	// udpaddr, err := net.ResolveUDPAddr("udp", laddr)
	// if err != nil {
	// 	log.Println(errors.WithStack(err))
	// }
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		log.Println(errors.WithStack(err))
	}

	// stun
	s, err := GetAddr(conn)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(s)
	NewNode(laddr, s)
	go func(path string, conn *net.UDPConn) {
		for {
			n := GetNode(path)
			n.PingPeer(conn)
			time.Sleep(time.Second)
		}
	}(laddr, conn)

	serveraddr, err := net.ResolveUDPAddr("udp", raddr)
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	listenServer(conn, serveraddr)
}

func server(laddr, raddr string) { // 接受，去头，传送
	pool = make(map[int]*net.UDPConn)

	udpaddr, err := net.ResolveUDPAddr("udp", laddr)
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	conn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		log.Println(errors.WithStack(err))
	}

	serveraddr, err := net.ResolveUDPAddr("udp", raddr)
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	listenServer(conn, serveraddr)
}

func listenServer(lc *net.UDPConn, rep *net.UDPAddr) {
	for {
		buf := make([]byte, 1500)
		n, addr, err := lc.ReadFrom(buf)
		if err != nil {
			log.Println(errors.WithStack(err))
		}
		if n < 4 {
			continue
		}

		// 如果没有则新建，按照tag来
		tag := int(binary.BigEndian.Uint32(buf[:4]))
		if pool[tag] == nil {
			udpaddr := &net.UDPAddr{
				IP:   net.IPv4(127, 27, buf[0], buf[1]),
				Port: tag & 0xffff,
			}
			c, err := net.ListenUDP("udp", udpaddr)
			if err != nil {
				log.Println(errors.WithStack(err))
			}
			pool[tag] = c
			go newConnServer(lc, c, addr, tag) // 监听的conn，新的conn，新建时收到的地址，tag
		}

		// 通过port索引并发送
		_, err = pool[tag].WriteToUDP(buf[4:n], rep)
		if err != nil {
			log.Println(errors.WithStack(err))
		}
	}
}

// 监听的conn，新的conn，新建时收到的地址，tag
// 目的：向Server发起的连接接收到的数据打上tag之后回传新建时收到的地址
func newConnServer(lc, conn *net.UDPConn, addr net.Addr, tag int) {
	defer func() {
		pool[tag] = nil
		delete(pool, tag)
		err := conn.Close()
		if err != nil {
			log.Println(errors.WithStack(err))
		}
	}()
	for {
		buf := make([]byte, 1500)
		n, _, err := conn.ReadFrom(buf[4:])
		if err != nil {
			log.Println(errors.WithStack(err))
		}
		// 贴入端口信息
		binary.BigEndian.PutUint32(buf[:4], uint32(tag))
		// 回传
		_, err = lc.WriteTo(buf[:n+4], addr)
		if err != nil {
			log.Println(errors.WithStack(err))
		}

	}
}

// 添加打洞
func client2(laddr, raddr string) { // 接受，加头，传送
	poolc = make(map[int]*net.Addr)

	udpaddr, err := net.ResolveUDPAddr("udp", laddr)
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	conn, err := net.ListenUDP("udp", udpaddr) // 本地映射
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	pc, err := net.ListenUDP("udp", nil) // 集束
	if err != nil {
		log.Println(errors.WithStack(err))
	}

	// stun
	s, err := GetAddr(pc)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(s)
	JoinNode(raddr, s)
	go func(path string, conn *net.UDPConn) {
		for {
			n := GetNode(path)
			n.PingHost(conn)
			time.Sleep(time.Second)
		}
	}(raddr, pc)
	n := GetNode(raddr)

	remoteaddr, err := net.ResolveUDPAddr("udp", n.Endpoint)
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	go listenClient(conn, pc, remoteaddr)
	listenClientR(conn, pc)
}

func client(laddr, raddr string) { // 接受，加头，传送
	poolc = make(map[int]*net.Addr)

	udpaddr, err := net.ResolveUDPAddr("udp", laddr)
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	conn, err := net.ListenUDP("udp", udpaddr) // 本地映射
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	pc, err := net.ListenUDP("udp", nil) // 集束
	if err != nil {
		log.Println(errors.WithStack(err))
	}

	remoteaddr, err := net.ResolveUDPAddr("udp", raddr)
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	go listenClient(conn, pc, remoteaddr)
	listenClientR(conn, pc)
}

// 接受（raw），加tag，传送走到rep
func listenClient(lc, pc *net.UDPConn, rep *net.UDPAddr) {
	for {
		buf := make([]byte, 1500)
		n, addr, err := lc.ReadFrom(buf[4:])
		if err != nil {
			log.Println(errors.WithStack(err))
		}
		// log.Println("addr:", addr)

		// 如果没有则新建，按照tag来
		tag := int(binary.BigEndian.Uint32(
			[]byte{
				addr.(*net.UDPAddr).IP[14],
				addr.(*net.UDPAddr).IP[15],
				byte(addr.(*net.UDPAddr).Port >> 8),
				byte(addr.(*net.UDPAddr).Port & 0xff),
			},
		))
		// 贴入端口信息
		binary.BigEndian.PutUint32(buf[:4], uint32(tag))
		if poolc[tag] == nil {
			log.Println("添加tag", buf[:4])
			poolc[tag] = &addr
		}

		// 传送走
		_, err = pc.WriteToUDP(buf[:n+4], rep)
		if err != nil {
			log.Println(errors.WithStack(err))
		}
	}
}

// pc 接受，通过lc发送，由tag决定方向
func listenClientR(lc, pc *net.UDPConn) {
	for {
		buf := make([]byte, 1500)
		n, _, err := pc.ReadFrom(buf)
		if err != nil {
			log.Println(errors.WithStack(err))
		}
		if n < 4 {
			continue
		}
		tag := int(binary.BigEndian.Uint32(buf[:4]))
		if poolc[tag] == nil {
			log.Println("不存在tag")
			log.Println(buf[:4])
			continue
		}

		// 传送走
		_, err = lc.WriteTo(buf[4:n], *poolc[tag])
		if err != nil {
			log.Println(errors.WithStack(err))
		}
	}
}

func raw(laddr, raddr string) {
	pool = make(map[int]*net.UDPConn)

	udpaddr, err := net.ResolveUDPAddr("udp", laddr)
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	conn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		log.Println(errors.WithStack(err))
	}

	rep, err := net.ResolveUDPAddr("udp", raddr)
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	// network := "udp4"
	// if udpaddr.IP.To4() == nil {
	// 	network = "udp"
	// }

	listen(conn, rep)
}

func listen(lc *net.UDPConn, rep *net.UDPAddr) {
	// buf := make([]byte, 1500)
	for {
		buf := make([]byte, 1500)
		n, addr, err := lc.ReadFrom(buf)
		if err != nil {
			log.Println(errors.WithStack(err))
		}

		// debug
		// fmt.Println("From ", addr, " len=", n)
		// fmt.Println(buf[:n])
		// fmt.Println(addr.(*net.UDPAddr).Port)

		// 如果没有则新建，按照port来
		if pool[addr.(*net.UDPAddr).Port] == nil {
			udpaddr := &net.UDPAddr{
				IP:   net.IPv4(127, 27, 7, 1),
				Port: addr.(*net.UDPAddr).Port,
			}
			nc, err := net.ListenUDP("udp", udpaddr)
			if err != nil {
				log.Println(errors.WithStack(err))
			}
			pool[addr.(*net.UDPAddr).Port] = nc
			go newConnDeamon(lc, nc, addr) // 监听的conn，新的conn，新建时收到的地址
		}

		// 通过port索引并发送
		_, err = pool[addr.(*net.UDPAddr).Port].WriteToUDP(buf[:n], rep)
		if err != nil {
			log.Println(errors.WithStack(err))
		}
		// debug
		// fmt.Println("From ", pool[addr.(*net.UDPAddr).Port].LocalAddr().String(), "To ", rep, " len=", n)

	}
}

// 监听的conn，新的conn，新建时收到的地址
func newConnDeamon(lc, conn *net.UDPConn, la net.Addr) {
	// buf := make([]byte, 1500)
	for {
		buf := make([]byte, 1500)
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			log.Println(errors.WithStack(err))
		}
		// debug
		// fmt.Println("回传：From ", addr, " len=", n)
		// fmt.Println(buf[:n])
		// fmt.Println(addr.(*net.UDPAddr).Port)

		// 回传
		_, err = lc.WriteTo(buf[:n], la)
		if err != nil {
			log.Println(errors.WithStack(err))
		}
		// fmt.Println("回传：From ", lc.LocalAddr(), "To ", la, " len=", n)

	}
}
