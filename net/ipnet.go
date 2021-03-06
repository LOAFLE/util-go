package net

import (
	"net"
)

// Sub class net.IPNet so that we can add JSON marshalling and unmarshalling.
type IPNet struct {
	net.IPNet
}

// Version returns the IP version for an IPNet, or 0 if not a valid IP net.
func (i *IPNet) Version() int {
	if i.IP.To4() != nil {
		return 4
	} else if len(i.IP) == net.IPv6len {
		return 6
	}
	return 0
}

// IsNetOverlap is a utility function that returns true if the two subnet have an overlap.
func (i *IPNet) IsNetOverlap(n net.IPNet) bool {
	return n.Contains(i.IP) || i.Contains(n.IP)
}

// Network returns the masked IP network.
func (i *IPNet) Network() *IPNet {
	_, n, _ := ParseCIDR(i.String())
	return n
}

func ParseCIDR(c string) (*IP, *IPNet, error) {
	netIP, netIPNet, e := net.ParseCIDR(c)
	if netIPNet == nil || e != nil {
		return nil, nil, e
	}
	ip := &IP{netIP}
	ipnet := &IPNet{*netIPNet}

	// The base golang net library always uses a 4-byte IPv4 address in an
	// IPv4 IPNet, so for uniformity in the returned types, make sure the
	// IP address is also 4-bytes - this allows the user to safely assume
	// all IP addresses returned by this function use the same encoding
	// mechanism (not strictly required but better for testing and debugging).
	if ip4 := ip.IP.To4(); ip4 != nil {
		ip.IP = ip4
	}

	return ip, ipnet, nil
}

// Parse a CIDR or an IP address and return the IP, CIDR or error.  If an IP address
// string is supplied, then the CIDR returned is the fully masked IP address (i.e /32 or /128)
func ParseCIDROrIP(c string) (*IP, *IPNet, error) {
	// First try parsing as a CIDR.
	ip, cidr, err := ParseCIDR(c)
	if err == nil {
		return ip, cidr, nil
	}

	// That failed, so try parsing as an IP.
	ip = &IP{}
	if err2 := ip.UnmarshalText([]byte(c)); err2 == nil {
		if ip4 := ip.IP.To4(); ip4 != nil {
			ip.IP = ip4
		}
		n := ip.Network()
		return ip, n, nil
	}

	// That failed too, return the original error.
	return nil, nil, err
}

// String returns a friendly name for the network.  The standard net package
// implements String() on the pointer, which means it will not be invoked on a
// struct type, so we re-implement on the struct type.
func (i *IPNet) String() string {
	ip := &i.IPNet
	return ip.String()
}

// MustParseNetwork parses the string into a IPNet.  The IP address in the
// IPNet is masked.
func MustParseNetwork(c string) IPNet {
	_, cidr, err := ParseCIDR(c)
	if err != nil {
		panic(err)
	}
	return *cidr
}

// MustParseCIDR parses the string into a IPNet.  The IP address in the
// IPNet is not masked.
func MustParseCIDR(c string) IPNet {
	ip, cidr, err := ParseCIDR(c)
	if err != nil {
		panic(err)
	}
	n := IPNet{}
	n.IP = ip.IP
	n.Mask = cidr.Mask
	return n
}
