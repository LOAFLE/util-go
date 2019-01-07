package cidr

import "net"

type cidrRangeIPv6 struct {
	cidrNet *net.IPNet
}

func (cr *cidrRangeIPv6) Contains(ip net.IP) bool {
	return cr.cidrNet.Contains(ip)
}

func (cr *cidrRangeIPv6) First() net.IP {
	return nil
}

func (cr *cidrRangeIPv6) Last() net.IP {
	return nil
}

func (cr *cidrRangeIPv6) Range() []net.IP {
	return nil
}

func (cr *cidrRangeIPv6) Ranges(startIP net.IP, endIP net.IP, include []net.IP, exclude []net.IP) ([]net.IP, error) {
	return nil, nil
}

func (cr *cidrRangeIPv6) Broadcast() net.IP {
	return nil
}

func (cr *cidrRangeIPv6) Network() net.IP {
	return nil
}

func (cr *cidrRangeIPv6) Next(ip net.IP) net.IP {
	return nil
}

func (cr *cidrRangeIPv6) Previous(ip net.IP) net.IP {
	return nil
}
