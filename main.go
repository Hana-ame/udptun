package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/hana-ame/udptun/helper"
	"github.com/hana-ame/udptun/utils"
)

var p *Portal

var dst string
var address string
var name string
var helperAddr string
var mode string
var lm *utils.LockedMap

func main() {
	flag.StringVar(&dst, "d", "", "destination host")
	flag.StringVar(&address, "a", "", "http listen address")
	flag.StringVar(&name, "n", "", "name")
	flag.StringVar(&helperAddr, "h", "", "helpserver address")
	flag.StringVar(&mode, "m", "udp6", "udp4/udp6")
	flag.Parse()
	lm = utils.NewLockedMap()
	p = NewPortal(dst)
	if dst != "" {
		// as server
		go func() {
			// arr := make([]string, 0)
			cnt := 0
			var arr []string
			for {
				if laddr := p.GetLocalAddr(mode); laddr != "" {
					helper.Post(helperAddr, name, laddr)
					// time.Sleep(5 * time.Second)
				}
				if a := helper.Get(helperAddr, name); a != nil {
					arr = append(arr, a...)
					cnt = 0
				} else {
					if cnt += 1; cnt > 90 {
						cnt = 0
						arr = nil
					}
				}
				p.Ping(arr)
				time.Sleep(time.Second)
			}
		}()
	} // dst != ""

	if address != "" && dst == "" {
		// add pointer and as client
		r := mux.NewRouter()
		r.HandleFunc("/", handleRoot)
		r.HandleFunc("/{peer}", handlePeer)
		http.ListenAndServe(address, r)

	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	m := helper.Post(helperAddr, "", "")
	json.NewEncoder(w).Encode(m)
}

func handlePeer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	peer, ok := vars["peer"]
	if !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	if r.Method == "POST" {
		m := helper.Post(helperAddr, "", "")
		if _, ok := m[peer]; ok {
			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "can not read body", http.StatusBadRequest)
			}
			// remote addr, local addr
			args := strings.Split(string(body), ",")
			c := NewUDPMux(args[1], args[0], p)
			lm.Put(peer, c)
		} else {
			http.Error(w, "not found peer", http.StatusNotFound)
		}
	} else if r.Method == "DELETE" {
		if c, ok := lm.Get(peer); ok {
			c.(*UDPMux).Close()
			lm.Remove(peer)
		}
	} else if r.Method == "GET" {
		if c, ok := lm.Get(peer); ok {
			json.NewEncoder(w).Encode(c.(*UDPMux))
		}
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
