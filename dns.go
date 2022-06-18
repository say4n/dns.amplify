package main

import (
	"encoding/binary"
	"log"
	"net"
)

// See RFC 1035: https://datatracker.ietf.org/doc/html/rfc1035
type DNSHeader struct {
	Xid     uint16 // Randomly chosen identifier.
	Flags   uint16 // Bit-mask to indicate request/response.
	Qdcount uint16 // Number of questions.
	Ancount uint16 // Number of answers.
	Nscount uint16 // Number of authority records.
	Arcount uint16 // Number of additional records.
}

// See RFC 1035: https://datatracker.ietf.org/doc/html/rfc1035
type DNSQuestion struct {
	Name     []byte // Query domain name.
	Dnstype  uint16 // The QTYPE (1 = A, 255 = all)
	Dnsclass uint16 // The QCLASS (1 = IN) (IN = Internet)
}

type DNSMessage struct {
	Header   DNSHeader
	Question DNSQuestion
}

func (m *DNSMessage) ToByteSlice() []byte {
	var messageBytes []byte
	buffer := make([]byte, 2)

	// // Process header section.
	binary.BigEndian.PutUint16(buffer, m.Header.Xid)
	messageBytes = append(messageBytes, buffer...)

	binary.BigEndian.PutUint16(buffer, m.Header.Flags)
	messageBytes = append(messageBytes, buffer...)

	binary.BigEndian.PutUint16(buffer, m.Header.Qdcount)
	messageBytes = append(messageBytes, buffer...)

	binary.BigEndian.PutUint16(buffer, m.Header.Ancount)
	messageBytes = append(messageBytes, buffer...)

	binary.BigEndian.PutUint16(buffer, m.Header.Nscount)
	messageBytes = append(messageBytes, buffer...)

	binary.BigEndian.PutUint16(buffer, m.Header.Arcount)
	messageBytes = append(messageBytes, buffer...)

	// Process question section.
	messageBytes = append(messageBytes, m.Question.Name...)

	binary.BigEndian.PutUint16(buffer, m.Question.Dnstype)
	messageBytes = append(messageBytes, buffer...)

	binary.BigEndian.PutUint16(buffer, m.Question.Dnsclass)
	messageBytes = append(messageBytes, buffer...)

	return messageBytes
}

func GenerateDNSMessage(domain string) DNSMessage {
	var sizes []int
	var index int

	for index = 0; index < len(domain); index++ {
		if domain[index] == '.' {
			if len(sizes) > 0 {
				sizes = append(sizes, index-sizes[len(sizes)-1]-1)
			} else {
				sizes = append(sizes, index)
			}
		}
	}
	sizes = append(sizes, index-sizes[len(sizes)-1]-1)

	log.Println("sizes:", sizes)

	var parts []byte
	parts = append(parts, byte(sizes[0]))
	offset := 1

	for index := 0; index < len(domain); index++ {
		value := domain[index]

		if value == '.' {
			parts = append(parts, byte(sizes[offset]))
			offset++
		} else {
			parts = append(parts, value)
		}
	}
	parts = append(parts, 0)

	log.Printf("%s is %x.\n", domain, parts)

	header := DNSHeader{
		Xid:     0xbeef, // Randomly chosen ID. ;)
		Flags:   0x0100, // Q=0, RD=1.
		Qdcount: 0x1,    // Sending one query.
		Ancount: 0x0,
		Nscount: 0x0,
		Arcount: 0x0,
	}

	question := DNSQuestion{
		Name:     parts,
		Dnstype:  0xff, // RR type.
		Dnsclass: 0x1,  // Querying on the internet (IN).
	}

	return DNSMessage{
		Header:   header,
		Question: question,
	}
}

func PerformDNSRequest(dnsServerAddress, targetAddress string, query DNSMessage) error {
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
