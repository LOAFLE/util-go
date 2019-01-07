package cidr

import (
	"fmt"
	"net"
)

func NewCIDRRanger(cidr string) (CIDRRanger, error) {
	_, nIPNet, err := net.ParseCIDR(cidr)
	if nil != err {
		return nil, err
	}

	switch len(nIPNet.IP) {
	case net.IPv4len:
		cr := &cidrRangeIPv4{
			cidrNet: nIPNet,
		}
		return cr, nil
	case net.IPv6len:
		cr := &cidrRangeIPv6{
			cidrNet: nIPNet,
		}
		return cr, nil
	default:
		return nil, fmt.Errorf("Net: not supported IP length")
	}
}

type CIDRRanger interface {
	Contains(ip net.IP) bool
	First() net.IP
	Last() net.IP
	Range() []net.IP
	// (!Contains(startIP) || !Contains(endIP)) return error
	// (startIP > endIP) return error
	// (nil != startIP && nil != endIP) return (startIP ~ endIP) + include - exclude
	// (nil == startIP || nil == endIP) return include - exclude
	Ranges(startIP net.IP, endIP net.IP, include []net.IP, exclude []net.IP) ([]net.IP, error)
	Broadcast() net.IP
	Network() net.IP
	Next(ip net.IP) net.IP
	Previous(ip net.IP) net.IP
}
