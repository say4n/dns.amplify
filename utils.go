package main

import (
	"log"
	"net"
	"regexp"
	"strings"
)

// StandardizeMACFormat fixes dash-separated MAC addresses from Windows ipconfig
// and macOS arp results which don't include leading zeros (:9: instead of :09:)
func StandardizeMACFormat(macAddr string) string {
	macAddr = strings.Replace(macAddr, "-", ":", -1)
	return regexp.MustCompile(`(\b)(\d)(\b)`).ReplaceAllString(macAddr, "${1}0${2}${3}")
}

// GetMacAddrForInterface returns the MAC address of the interface with the given name.
func GetMacAddrForInterface(iface string) net.HardwareAddr {
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
