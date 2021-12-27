package main

import (
	"log"
	"strconv"
	"strings"
)

func GetTag(s string) byte {
	s_ip := GetIP(s)
	ip := strings.Split(s_ip, ".")
	tag, err := strconv.ParseUint(ip[3], 10, 8)
	if err != nil {
		log.Printf("error : %v, string = %s", err, s)
	}
	return byte(tag)
}

func GetIP(s string) string {
	s_ip_port := strings.Split(s, ":")
	return s_ip_port[0]
}

func GetPort(s string) uint16 {
	s_ip_port := strings.Split(s, ":")
	sport := s_ip_port[len(s_ip_port)-1]
	port, err := strconv.ParseUint(sport, 10, 16)
	if err != nil {
		log.Printf("error : %v, string = %s", err, s)
	}
	return uint16(port)
}
