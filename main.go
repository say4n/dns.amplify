package main

import (
	"log"
	"net"
)

func performDNSRequest(dnsServerAddress, targetAddress string, query DNSMessage) error {
	raddr, err := net.ResolveUDPAddr("udp", dnsServerAddress)
	if err != nil {
		log.Fatalf("Encountered error `%s` while trying to resolve address `%s`.\n", err, dnsServerAddress)
	}

	var laddr *net.UDPAddr

	if targetAddress != "" {
		laddr, err = net.ResolveUDPAddr("udp", targetAddress)
		if err != nil {
			log.Fatalf("Encountered error `%s` while trying to resolve address `%s`.\n", err, targetAddress)
		}
	} else {
		laddr = nil
	}

	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		log.Fatalf("Encountered error `%s` trying to start UDP communication.", err)
	}

	defer conn.Close()

	log.Println("targetIP:", conn.LocalAddr().String())

	messageBytes := query.ToByteSlice()
	log.Printf("messageBytes: %x\n", messageBytes)

	conn.Write(messageBytes)

	return nil
}

func main() {
	query := "fb.me"
	queryMessage := GenerateDNSMessage(query)

	err := performDNSRequest("1.0.0.1:53", "", queryMessage)
	if err != nil {
		panic(err)
	}
}
