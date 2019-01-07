package net

import (
	"net"
)

// Sub class net.IP so that we can add JSON marshalling and unmarshalling.
type IP struct {
	net.IP
}

// ParseIP returns an IP from a string
func ParseIP(ip string) *IP {
	addr := net.ParseIP(ip)
	if addr == nil {
		return nil
	}
	// Always return IPv4 values as 4-bytes to be consistent with IPv4 IPNet
	// representations.
	if addr4 := addr.To4(); addr4 != nil {
		addr = addr4
	}
	return &IP{addr}
}

// Version returns the IP version for an IP, or 0 if the IP is not valid.
func (i *IP) Version() int {
	if i.To4() != nil {
		return 4
	} else if len(i.IP) == net.IPv6len {
		return 6
	}
	return 0
}

// Network returns the IP address as a fully masked IPNet type.
func (i *IP) Network() *IPNet {
	// Unmarshaling an IPv4 address returns a 16-byte format of the
	// address, so convert to 4-byte format to match the mask.
	n := &IPNet{}
	if ip4 := i.IP.To4(); ip4 != nil {
		n.IP = ip4
		n.Mask = net.CIDRMask(net.IPv4len*8, net.IPv4len*8)
	} else {
		n.IP = i.IP
		n.Mask = net.CIDRMask(net.IPv6len*8, net.IPv6len*8)
	}
	return n
}

// MustParseIP parses the string into a IP.
func MustParseIP(i string) IP {
	var ip IP
	err := ip.UnmarshalText([]byte(i))
	if err != nil {
		panic(err)
	}
	// Always return IPv4 values as 4-bytes to be consistent with IPv4 IPNet
	// representations.
	if ip4 := ip.To4(); ip4 != nil {
		ip.IP = ip4
	}
	return ip
}
