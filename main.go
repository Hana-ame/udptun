package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/hana-ame/udptun/helper"
)

var p *Portal

var dst string
var address string

func main() {
	flag.StringVar(&dst, "d", "", "destination host")
	flag.StringVar(&address, "httpHost", "", "http listen address")
	flag.Parse()
	p = NewPortal(dst)

	if address != "" {
		// http
	}
}

func helper8888() {
	helper.Server("127.0.0.1:8888")
}

func server9999() {
	p := NewPortal("127.0.0.1:9999")
	go func() {
		for {
			time.Sleep(5 * time.Second)
			fmt.Println(p.connMap.Size(), p.router.Size())
			p.connMap.Iter(func(k string, v any) {
				fmt.Println((k))
			})
		}
	}()
	time.Sleep(time.Minute)
	time.Sleep(time.Minute)
}

func client3333() {
	p := NewPortal("")
	// go func() {
	// 	for {
	// 		time.Sleep(5 * time.Second)
	// 		fmt.Println(p.connMap.Size(), p.router.Size())
	// 		p.router.Iter(func(k string, v any) {
	// 			fmt.Println((k))
	// 		})
	// 	}
	// }()
	laddr := "0.0.0.0:3333"
	c := NewUDPMux(laddr, "127.0.0.1:4444", p)
	// go func() {
	// 	for {
	// 		time.Sleep(5 * time.Second)
	// 		fmt.Println(c)
	// 		fmt.Println(c.connMap.Size())
	// 		c.connMap.Iter(
	// 			func(key string, value any) {
	// 				fmt.Println([]byte(key))
	// 			},
	// 		)
	// 	}
	// }()
	time.Sleep(time.Minute * 10)

	c.Close()

	time.Sleep(time.Minute)

}

// func main() {
// 	utils.UDPEcho(":9999")
// }
