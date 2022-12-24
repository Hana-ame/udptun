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
	go helper.Server("127.0.0.1:8888")
	var m map[string]string
	var a []string
	m = helper.Post("http://127.0.0.1:8888/", "name", "value")
	m = helper.Post("http://127.0.0.1:8888/", "name1", "value")
	m = helper.Post("http://127.0.0.1:8888/", "name1", "value2")
	m = helper.Post("http://127.0.0.1:8888/", "name1", "value3")
	helper.Delete("http://127.0.0.1:8888/", "name")
	m = helper.Post("http://127.0.0.1:8888/", "name1", "value3")
	a = helper.Get("http://127.0.0.1:8888/", "name1")
	fmt.Println(m, a)
	helper.Append("http://127.0.0.1:8888/", "name1", "value4")
	helper.Append("http://127.0.0.1:8888/", "name", "value4")
	a = helper.Get("http://127.0.0.1:8888/", "name1")
	fmt.Println(m, a)
	a = helper.Get("http://127.0.0.1:8888/", "name")
	fmt.Println(m, a)

	helper.Delete("http://127.0.0.1:8888/", "name1")
	m = helper.Post("http://127.0.0.1:8888/", "name1", "value")
	a = helper.Get("http://127.0.0.1:8888/", "name1")
	fmt.Println(m, a)
}

func main1() {
	flag.StringVar(&dst, "d", "", "destination host")
	flag.StringVar(&address, "httpHost", "", "http listen address")

	p = NewPortal(dst)

	flag.Parse()

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
