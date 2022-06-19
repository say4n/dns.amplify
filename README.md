# dns.amplify
a tiny proof of concept to understand dns amplification attacks

# what even is this codebase?

## main.go

It is the entry point to the codebase.
The request URL, the forged source IP, the DNS server to send the request to are
all defined her.

## dns.go

This houses the various data structures to craft a DNS query message.
The crafted message is then sent over UDP to the DNS server.

## udp.go

This file is a copy of [this file](https://github.com/dimalinux/spoofsourceip/blob/e1554cd99d5fd7b5d3ba199fba4a3acc5308d5db/udpspoof/udp.go) with more verbose documentation.
It creates a UDP frame with a given payload (our DNS query).

## utils.go

This file defines various utility functions used by the main.go file.

# notes
- [udp socket programming: dns](https://w3.cs.jmu.edu/kirkpams/OpenCSF/Books/csf/html/UDPSockets.html)
- [ip spoofing does not work on macOS?](https://dev.to/conner/ip-spoofing-theory-and-implementation-ep6)
- [spoofing UDP source IP in golang](https://github.com/dimalinux/spoofsourceip)
