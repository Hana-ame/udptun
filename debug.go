package main

import "fmt"

func Debug(s ...interface{}) {
	// return
	fmt.Println(s...)
	// for i := range s {
	// fmt.Println(i)
	// }
}
