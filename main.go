package main

import (
	"fmt"
	"time"
)

func main() {
	p := NewPortal("")
	go func() {
		for {
			time.Sleep(5 * time.Second)
			fmt.Println(p)
		}
	}()
	laddr := "0.0.0.0:3456"
	c := NewUDPMux(laddr, "127.0.0.1:9999", p)
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
	time.Sleep(time.Minute)

	c.closed = true

	time.Sleep(time.Minute)

}

// func main() {
// 	utils.UDPEcho(":9999")
// }
