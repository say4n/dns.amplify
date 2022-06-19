package main

import (
	"encoding/binary"
	"log"
	"math/rand"
)

// See RFC 1035: https://datatracker.ietf.org/doc/html/rfc1035
// A DNSHeader is a struct with 6 fields, each of which is a 16-bit unsigned
// integer.
// @property {uint16} Xid - A random number that identifies the request.
// @property {uint16} Flags - The flags field is a bit-mask to indicate whether the
// message is a request or a response, whether recursion is desired, and whether
// the message has been truncated.
// @property {uint16} Qdcount - Number of questions.
// @property {uint16} Ancount - The number of answers in the response.
// @property {uint16} Nscount - Number of name server resource records in the
// authority records section.
// @property {uint16} Arcount - The number of additional records.
type DNSHeader struct {
	Xid     uint16 // Randomly chosen identifier.
	Flags   uint16 // Bit-mask to indicate request/response.
	Qdcount uint16 // Number of questions.
	Ancount uint16 // Number of answers.
	Nscount uint16 // Number of authority records.
	Arcount uint16 // Number of additional records.
}

// See RFC 1035: https://datatracker.ietf.org/doc/html/rfc1035
// A DNSQuestion is a DNS query for a domain name of a certain type and class.
// @property {[]byte} Name - The domain name we want to query.
// @property {uint16} Dnstype - The type of DNS record you're looking for.
// @property {uint16} Dnsclass - The class of the query. This is usually 1, which
// means the Internet.
type DNSQuestion struct {
	Name     []byte // Query domain name.
	Dnstype  uint16 // The QTYPE (1 = A, 255 = all)
	Dnsclass uint16 // The QCLASS (1 = IN) (IN = Internet)
}

// A DNSMessage is a DNSHeader and a DNSQuestion.
// @property {DNSHeader} Header - This is the header of the DNS message. It
// contains information about the message itself.
// @property {DNSQuestion} Question - The question section of the DNS message.
type DNSMessage struct {
	Header   DNSHeader
	Question DNSQuestion
}

// Converting the DNSMessage struct into a byte slice.
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

// GenerateDNSMessage takes a domain name, splits it into parts, and then creates a DNSMessage
// struct with the header and question fields filled out.
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

	header := DNSHeader{
		Xid:     uint16(rand.Uint32()), // Randomly chosen ID. ;)
		Flags:   0x0100,                // Q=0, RD=1.
		Qdcount: 0x1,                   // Sending one query.
		Ancount: 0x0,
		Nscount: 0x0,
		Arcount: 0x0,
	}

	log.Printf("DNSHeader: %#v", header)

	question := DNSQuestion{
		Name:     parts,
		Dnstype:  0xff, // RR type.
		Dnsclass: 0x1,  // Querying on the internet (IN).
	}

	log.Printf("DNSQuestion: %#v", question)

	return DNSMessage{
		Header:   header,
		Question: question,
	}
}
