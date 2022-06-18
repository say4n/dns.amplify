package main

import (
	"log"
	"net"
	"regexp"
	"strings"

	"github.com/google/gopacket/pcap"
	"github.com/jackpal/gateway"
	"github.com/mostlygeek/arp"
)

func performDNSRequest(query DNSMessage) {
	messageBytes := query.ToByteSlice()
	log.Printf("messageBytes: %x\n", messageBytes)

	defaultGatewayIP, _ := gateway.DiscoverGateway()
	log.Printf("defaultGatewayIP: %#v\n", defaultGatewayIP.String())

	defaultGatewayMac, _ := net.ParseMAC(standardizeMACFormat(arp.Search(defaultGatewayIP.String())))
	log.Printf("defaultGatewayMac: %#v\n", defaultGatewayMac)

	udpFrameOptions := UdpFrameOptions{
		sourceIP:     net.IPv4(127, 0, 0, 1),
		destIP:       net.IPv4(1, 1, 1, 1),
		sourcePort:   4000,
		destPort:     53,
		sourceMac:    getMacAddrForInterface("en0"),
		destMac:      defaultGatewayMac,
		isIPv6:       false,
		payloadBytes: query.ToByteSlice(),
	}

	frameBytes, err := CreateSerializedUDPFrame(udpFrameOptions)
	if err != nil {
		log.Fatal(err)
	}

	handle, err := pcap.OpenLive("en0", 1024, false, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	if err := handle.WritePacketData(frameBytes); err != nil {
		log.Fatal(err)
	}
	log.Println("DNS request sent.")
}

// standardizeMACFormat fixes dash-separated MAC addresses from Windows ipconfig
// and macOS arp results which don't include leading zeros (:9: instead of :09:)
func standardizeMACFormat(macAddr string) string {
	macAddr = strings.Replace(macAddr, "-", ":", -1)
	return regexp.MustCompile(`(\b)(\d)(\b)`).ReplaceAllString(macAddr, "${1}0${2}${3}")
}

func getMacAddrForInterface(iface string) net.HardwareAddr {
	ifas, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, ifa := range ifas {
		if ifa.Name == iface {
			return ifa.HardwareAddr
		}
	}
	return nil
}

func main() {
	query := "fb.me"
	queryMessage := GenerateDNSMessage(query)
	performDNSRequest(queryMessage)
}
