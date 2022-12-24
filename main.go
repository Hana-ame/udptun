package main

import (
	"time"
)

func main() {
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
	time.Sleep(time.Minute)

	c.Close()

	time.Sleep(time.Minute)

}

// func main() {
// 	utils.UDPEcho(":9999")
// }
