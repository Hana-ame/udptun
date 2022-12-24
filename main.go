package main

import "fmt"

func proxy(addr string, p Portal) *UDPMux {
	// listen udp local

	// create UDPMux

	// return
	return nil
}

func main() {
	b := make(PortalBuf, 10)
	b[0] = 1
	b[1] = 2
	b[2] = 3
	b[3] = 4
	fmt.Println(b)
	fmt.Println(b.Tag())
	copy(b.Tag(), b.Data(0))
	fmt.Println(b.Raw(4).Data(0))
	fmt.Println(b.Tag())
	b[3] = 5
	b[0] = 0
	fmt.Println(b.Raw(4).Data(0))
	fmt.Println(b.Tag())

}
