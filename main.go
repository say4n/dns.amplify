package main

import (
	"log"
	"net"
)

func connect(dnsServerAddress, targetAddress string) {
	raddr, err := net.ResolveUDPAddr("udp", dnsServerAddress)
	if err != nil {
		log.Fatalf("Encountered error `%s` while trying to resolve address `%s`.\n", err, dnsServerAddress)
	}

	laddr, err := net.ResolveUDPAddr("udp", targetAddress)
	if err != nil {
		log.Fatalf("Encountered error `%s` while trying to resolve address `%s`.\n", err, targetAddress)
	}

	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		log.Fatalf("Encountered error `%s` trying to start UDP communication.", err)
	}

	defer conn.Close()
}

func main() {
	connect("dns.google:53", "127.0.0.1:2000")
}
