package main

import (
	"log"
	"net"
	"strings"
)

func getIP() string {
	c, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	unFormattedIP := c.LocalAddr()
	ip := strings.Split(unFormattedIP.String(), ":")

	return ip[0]
}
