package converter

import (
	"fmt"
	"net"
)

func IPToInt(ip net.IP) (int32, error) {
	switch len(ip) {
	case net.IPv4len:
		return IPv4ToInt(ip), nil
	case net.IPv6len:
		return IPv6ToInt(ip), nil
	default:
		return 0, fmt.Errorf("Net: not supported IP length")
	}
}
