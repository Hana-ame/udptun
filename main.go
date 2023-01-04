package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
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

var isHelperServer bool

func main() {
	flag.StringVar(&dst, "d", "", "destination host")
	flag.StringVar(&address, "a", "", "http listen address")
	flag.StringVar(&name, "n", "", "name")
	flag.StringVar(&helperAddr, "h", "", "helpserver address")
	flag.StringVar(&mode, "m", "udp6", "udp4/udp6")
	flag.BoolVar(&isHelperServer, "isHelpServer", false, "work as help server")

	flag.Parse()

	if isHelperServer {
		err := helper.Server(helperAddr)
		log.Fatal(err)
		return
	}

	lm = utils.NewLockedMap()
	p = NewPortal(dst)
	if dst != "" {
		// as server
		// arr := make([]string, 0)
		cnt := 0
		var arr []string
		for {
			// every 10 seconds renew the server node.
			if cnt%10 == 0 {
				if laddr := p.GetLocalAddr(mode); laddr != "" {
					helper.Post(helperAddr, name, laddr)
					// time.Sleep(5 * time.Second)
				}
			}
			// every second get peers and ping.
			if a := helper.Get(helperAddr, name); a != nil {
				log.Println(a)
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
	} // dst != ""

	if address != "" && dst == "" {
		// add pointer and as client
		r := mux.NewRouter()
		r.HandleFunc("/", handleRoot)
		r.HandleFunc("/{peer}", handlePeer)
		err := http.ListenAndServe(address, r)
		log.Fatal(err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		m := helper.Post(helperAddr, "", "")
		json.NewEncoder(w).Encode(m)
	} else if r.Method == "POST" {
		defer r.Body.Close()
		var m map[string]string
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		local, ok := m["local"]
		if !ok {
			http.Error(w, "bad request, no local field", http.StatusBadRequest)
			return
		}
		destination, ok := m["destination"]
		if !ok {
			http.Error(w, "bad request, no destination field", http.StatusBadRequest)
			return
		}

		c := NewUDPMux(local, destination, p, "peer")
		lm.Put(destination, c)
	}
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
				return
			}
			// args := {remote addr, local addr}
			arg := string(body)
			helper.Append(helperAddr, peer, p.GetLocalAddr(mode))
			c := NewUDPMux(arg, m[peer], p, peer)

			// fmt.Println(helperAddr, peer, p.GetLocalAddr(mode))
			// fmt.Println(arg, m[peer], p)
			// fmt.Println(c)

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
			// fmt.Println(c)
			json.NewEncoder(w).Encode(c)
		} else {
			http.Error(w, "not found peer", http.StatusNotFound)
		}
	}
}
