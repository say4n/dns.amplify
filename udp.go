// Source: https://github.com/dimalinux/spoofsourceip/blob/e1554cd99d5fd7b5d3ba199fba4a3acc5308d5db/udpspoof/udp.go

package main

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// UdpFrameOptions is a struct that holds all the information needed to create a UDP frame.
// @property sourceIP - The source IP address of the UDP packet.
// @property {uint16} sourcePort - The source port of the UDP packet.
// @property sourceMac - The MAC address of the sender.
// @property {bool} isIPv6 - Whether the frame is IPv6 or IPv4.
// @property {[]byte} payloadBytes - The bytes of the payload to be sent.
type UdpFrameOptions struct {
	sourceIP, destIP     net.IP
	sourcePort, destPort uint16
	sourceMac, destMac   net.HardwareAddr
	isIPv6               bool
	payloadBytes         []byte
}

// "A SerializableNetworkLayer is a gopacket.NetworkLayer that can be serialized to
// a buffer."
//
// The first thing to notice is that the type is an interface.  This means that any
// type that implements the interface can be used as a SerializableNetworkLayer.
// The second thing to notice is that the interface has two methods:
// gopacket.NetworkLayer and SerializeTo. The first method is a method that is
// already defined in the gopacket library. The second method is a method that we
// are defining. This is the method that will be called by the
// serialization code to serialize the layer.
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
