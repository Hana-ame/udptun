单纯测试一下是否是严格的nat类型。
返回一个值时说明不是严格类型，有nat希望。
python版本的不知道放到哪里去了。


```golang

func main() {
	conn, err := net.ListenUDP(`udp`, nil)
	if err != nil {
		log.Fatal(`NewWrapper : `, err.Error())
	}
	go func() {
		time.Sleep(time.Second)
		conn.Close()
	}()
	buf := make([]byte, 100)
	_, _, err = conn.ReadFrom(buf)
	if err != nil {
		log.Println(err)
	}

}

```

会返回已经关闭的conn错误




```bash

nc -uvlp 10000


```


```golang

func main() {
	pool := NewPortalPool(5, 5)
	for i := 0; i < 5; i++ {
		p := NewPortal("udp")
		pool.Add(p)
	}
	fmt.Println(pool)
	// time.Sleep(time.Second * 5)
	fmt.Println(pool)
	// time.Sleep(time.Second * 120)

	p := pool.Pick()
	s, err := GetAddr(p.Conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s)
	fmt.Println(p.Conn.LocalAddr())
	fmt.Println(pool)

	for i := 0; i < 4; i++ {
		p := pool.Pick()
		var paddr *string
		laddr := "localhost:9999"
		if i == 0 {
			addr := "localhost:10000"
			paddr = &addr
		}
		if i == 2 {
			addr := "localhost:10001"
			paddr = &addr
		}
		p.Set(paddr, &laddr, nil)
	}

	time.Sleep(time.Second * 60)
}

```

server这边大大概没什么问题


```powershell

.\ncat.exe -vlup 10000
.\ncat.exe -vlup 10001


```

debian的会断开
https://superuser.com/questions/1008348/netcat-keep-listening-for-connection-in-debian



```golang
func main() {
	ps := NewPortalServer("localhost:9999")
	ps.NewPortal()
	ps.NewPortal()
	ps.NewPortal()
	ps.NewPortal()

	var paddr string

	paddr = "localhost:10000"
	ps.ActivePortal(&paddr)
	paddr = "localhost:10001"
	ps.ActivePortal(&paddr)

	time.Sleep(time.Second * 90)
}

func main() {
	ps := NewPortalServer("localhost:9999")
	ps.NewPortal()
	ps.NewPortal()
	ps.NewPortal()
	ps.NewPortal()

	fmt.Println(ps.Pool)

	var paddr string

	paddr = "localhost:10000"
	ps.ActivePortal(&paddr,ps.LocalAddr,nil)
	paddr = "localhost:10001"
	ps.ActivePortal(&paddr,ps.LocalAddr,nil)

	time.Sleep(time.Second * 90)
}
```
PortalServer没啥问题




```golang 
// stun-tester.go

func _main() {
	Conn, err := net.ListenPacket("udp", fmt.Sprintf("0.0.0.0:%d", 12321))
	if err != nil {
		log.Fatal("sb")
		return
	}
	// fmt.Println(GetAddr(Conn))

	// addr, err := net.ResolveUDPAddr("udp", "34.145.70.165:12421")
	// if err != nil {
	// 	log.Printf("error : %v", err)
	// 	return
	// }
	// fmt.Println("1")
	// time.Sleep(time.Second * 3)
	// Conn.WriteTo([]byte{0}, addr)
	// Conn.WriteTo([]byte{0}, addr)
	// Conn.WriteTo([]byte{0}, addr)
	// Conn.WriteTo([]byte{0}, addr)
	// Conn.WriteTo([]byte{0}, addr)
	// fmt.Println("2")

	// buffer := make([]byte, 2048)
	// for {
	// 	fmt.Println(3)
	// 	n, addr, err := Conn.ReadFrom(buffer)
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 	}
	// 	fmt.Printf("packet-received: bytes=%d from=%s\n", n, addr.String())
	// 	// portalproxy.PrintHex(buffer[:n])
	// }
	s, _ := GetAddr(Conn)
	fmt.Println(s)
}
```

```golang 
func __1_main() {
	p1 := &Portal{}
	p2 := &Portal{}
	p3 := &Portal{}
	fmt.Println([]*Portal{p1, p2, p3})

	pool := NewPortalPool(1, 1)
	fmt.Println(pool)

	pool.Add(p1)
	fmt.Println(pool)
	pool.Add(p2)
	fmt.Println(pool)
	pool.Add(p3)
	fmt.Println(pool)

	fmt.Println(pool.Pick())
	fmt.Println(pool)
	fmt.Println(pool.Pick())
	fmt.Println(pool)
	fmt.Println(pool.Pick())
	fmt.Println(pool)

	pool.Add(p1)
	fmt.Println(pool)
	pool.Add(p2)
	fmt.Println(pool)
	pool.Add(p3)
	fmt.Println(pool)

	fmt.Println(pool.Pick())
	fmt.Println(pool)
	fmt.Println(pool.Pick())
	fmt.Println(pool)
	fmt.Println(pool.Pick())
	fmt.Println(pool)

	fmt.Println(pool.Pick())
	fmt.Println(pool)
	fmt.Println(pool.Pick())
	fmt.Println(pool)
	fmt.Println(pool.Pick())
	fmt.Println(pool)

	pool.Add(p1)
	fmt.Println(pool)
	fmt.Println(pool.Pick())
	fmt.Println(pool)

}

```



```golang
func main() {
	pc := NewPortalClient("localhost:10000")
	pc.NewPortal()
	pc.NewPortal()
	pc.NewPortal()
	pc.NewPortal()

	fmt.Println(pc.Pool)

	var paddr string = ""

	paddr = "localhost:9999"
	pc.ActivePortal(&paddr, nil, pc.Mux)
	pc.ActivePortal(&paddr, nil, pc.Mux)
	pc.ActivePortal(&paddr, nil, pc.Mux)
	pc.ActivePortal(&paddr, nil, pc.Mux)

	// p := pc.Mux.Pool.m[0]
	// go func() {
	// 	for {
	// 		fmt.Println(p)
	// 		time.Sleep(time.Second * 5)
	// 	}
	// }()
	for {
		fmt.Println(pc.Pool)
		fmt.Println(pc.Mux)
		fmt.Println(pc.Mux.Pool)
		time.Sleep(time.Second * 2)
		fmt.Println("==============")
	}

	time.Sleep(time.Second * 30)
	fmt.Println("==============")

	fmt.Println(pc.Pool)
	fmt.Println(pc.Mux)
	fmt.Println(pc.Mux.Pool)

	fmt.Println("==============")
	time.Sleep(time.Second * 60)
	fmt.Println("==============")

	fmt.Println(pc.Pool)
	fmt.Println(pc.Mux)
	fmt.Println(pc.Mux.Pool)

	fmt.Println("==============")
	time.Sleep(time.Second * 90)
}


```

PortalClient大概没什么问题。


```golang

	// getPool = pc.Pool
	// putPool = pc.Mux.Pool


```


突然不对了，回个档看看


```bash 

# server
./udptun.exe -s --a=localhost:9999 --p=test 

# client
./udptun.exe --p=test --a=:10000
./udptun.exe --p=test 

```


```cmd

.\hysteria-windows-amd64.exe 
.\hysteria-windows-amd64.exe -c server.json server

```


google cloudshell不通。



gcp
```bash
python3 -m http.server 80

#~/Toys/stun-tester
$ git pull origin
#~/Toys/stun-tester
$ go build -o ./release/udptun .

#~/Toys/stun-tester
$ ./release/udptun  -s --a=localhost:9999 --p=test


cd ~/Toys/stun-tester
./release/udptun  -s --a=localhost:9999 --p=test

#:~
$ cd ./Toys/stun-tester/release/
#~/Toys/stun-tester/release
$ ./hysteria-linux-amd64 server

```


Pool应该加入回调让Pool自己填充更好点。
概率性开不起来不想debug了。又要重写，草。