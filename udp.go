// Source: https://github.com/dimalinux/spoofsourceip/blob/e1554cd99d5fd7b5d3ba199fba4a3acc5308d5db/udpspoof/udp.go

package main

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type UdpFrameOptions struct {
	sourceIP, destIP     net.IP
	sourcePort, destPort uint16
	sourceMac, destMac   net.HardwareAddr
	isIPv6               bool
	payloadBytes         []byte
}

type SerializableNetworkLayer interface {
	gopacket.NetworkLayer
	SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error
}

// CreateSerializedUDPFrame creates an Ethernet frame encapsulating our UDP
// packet for injection to the local network
func CreateSerializedUDPFrame(opts UdpFrameOptions) ([]byte, error) {

	buf := gopacket.NewSerializeBuffer()
	serializeOpts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	ethernetType := layers.EthernetTypeIPv4
	if opts.isIPv6 {
		ethernetType = layers.EthernetTypeIPv6
	}
	eth := &layers.Ethernet{
		SrcMAC:       opts.sourceMac,
		DstMAC:       opts.destMac,
		EthernetType: ethernetType,
	}
	var ip SerializableNetworkLayer
	if !opts.isIPv6 {
		ip = &layers.IPv4{
			SrcIP:    opts.sourceIP,
			DstIP:    opts.destIP,
			Protocol: layers.IPProtocolUDP,
			Version:  4,
			TTL:      32,
		}
	} else {
		ip = &layers.IPv6{
			SrcIP:      opts.sourceIP,
			DstIP:      opts.destIP,
			NextHeader: layers.IPProtocolUDP,
			Version:    6,
			HopLimit:   32,
		}
		ip.LayerType()
	}

	udp := &layers.UDP{
		SrcPort: layers.UDPPort(opts.sourcePort),
		DstPort: layers.UDPPort(opts.destPort),
		// we configured "Length" and "Checksum" to be set for us
	}
	udp.SetNetworkLayerForChecksum(ip)
	err := gopacket.SerializeLayers(buf, serializeOpts, eth, ip, udp, gopacket.Payload(opts.payloadBytes))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
