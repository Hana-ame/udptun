package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Node map[string]int64

// type Site map[string]Node

// func GetAddr(host string, path string, timeout int64) []string {

// }
func PostAddr(host string, path string, data string, timeout int64) []string {
	client := http.Client{Timeout: time.Second * 5}
	// fmt.Println(host + "/" + path)
	resp, err := client.Post(host+"/"+path, "application/json",
		bytes.NewBuffer([]byte(`"`+data+`"`)))
	if err != nil {
		log.Println(err)
		return nil
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	// bodyString := string(bodyBytes)
	// fmt.Println(bodyString)
	n := &Node{}
	err = json.Unmarshal(bodyBytes, n)
	if err != nil {
		log.Println(err)
		return nil
	}

	res := make([]string, 0)
	for k, v := range *n {
		fmt.Println(v)
		fmt.Println(time.Now().UnixNano() - timeout*1_000_000)
		if v > time.Now().UnixNano()-timeout*1_000_000 { // timeout: s -> nano
			res = append(res, k)
		}
	}
	fmt.Println("=======")

	return res
}
