package cidr

import (
	"fmt"
	"net"

	"git.loafle.net/overflow/util-go/net/converter"
)

type cidrRangeIPv4 struct {
	cidrNet *net.IPNet
}

func (cr *cidrRangeIPv4) Contains(ip net.IP) bool {
	return cr.cidrNet.Contains(ip)
}

func (cr *cidrRangeIPv4) First() net.IP {
	nIP := cr.Network()
	return cr.Next(nIP)
}

func (cr *cidrRangeIPv4) Last() net.IP {
	bIP := cr.Broadcast()
	return cr.Previous(bIP)
}

func (cr *cidrRangeIPv4) Range() []net.IP {
	fIP := cr.First()
	if nil == fIP {
		return nil
	}
	lIP := cr.Last()
	if nil == lIP {
		return nil
	}
	fNum := converter.IPv4ToInt(fIP.To4())
	lNum := converter.IPv4ToInt(lIP.To4())

	r := make([]net.IP, 0)
	for i := fNum; i <= lNum; i++ {
		r = append(r, converter.IntToIPv4(i))
	}
	return r
}

// (!Contains(startIP) || !Contains(endIP)) return nil
// (startIP > endIP) return nil
// (nil != startIP && nil != endIP) return (startIP ~ endIP) + include - exclude
// (nil == startIP || nil == endIP) return include - exclude
func (cr *cidrRangeIPv4) Ranges(startIP net.IP, endIP net.IP, include []net.IP, exclude []net.IP) ([]net.IP, error) {

	res := make(map[int32]bool)

	if nil != startIP && nil != endIP {
		if !cr.Contains(startIP) {
			return nil, fmt.Errorf("CIDR Range: CIDR not contains start ip[%v]", startIP)
		}
		if !cr.Contains(endIP) {
			return nil, fmt.Errorf("CIDR Range: CIDR not contains end ip[%v]", endIP)
		}
		sNum := converter.IPv4ToInt(startIP.To4())
		eNum := converter.IPv4ToInt(endIP.To4())
		if sNum > eNum {
			return nil, fmt.Errorf("CIDR Range: Start IP[%v] must smaller then End IP[%v]", startIP, endIP)
		}
		for i := sNum; i <= eNum; i++ {
			res[i] = true
		}
	}

	if nil != include {
		for _, in := range include {
			iNum := converter.IPv4ToInt(in.To4())
			if _, ok := res[iNum]; !ok {
				res[iNum] = true
			}
		}
	}

	if nil != exclude {
		for _, ex := range exclude {
			iNum := converter.IPv4ToInt(ex.To4())
			if _, ok := res[iNum]; ok {
				delete(res, iNum)
			}
		}
	}

	r := make([]net.IP, 0)
	for k, _ := range res {
		r = append(r, converter.IntToIPv4(k))
	}

	return r, nil
}

func (cr *cidrRangeIPv4) Broadcast() net.IP {
	ip := cr.cidrNet.IP.To4()
	bIP := net.IPv4(0, 0, 0, 0).To4()
	for i := 0; i < len(bIP); i++ {
		bIP[i] = ip[i] | ^cr.cidrNet.Mask[i]
	}
	return bIP
}

func (cr *cidrRangeIPv4) Network() net.IP {
	ip := cr.cidrNet.IP.To4()
	return ip.Mask(cr.cidrNet.Mask)
}

func (cr *cidrRangeIPv4) Next(ip net.IP) net.IP {
	nNum := converter.IPv4ToInt(ip.To4()) + 1
	nIP := converter.IntToIPv4(nNum)
	if cr.Contains(nIP) {
		return nIP
	}
	return nil
}

func (cr *cidrRangeIPv4) Previous(ip net.IP) net.IP {
	nNum := converter.IPv4ToInt(ip.To4()) - 1
	nIP := converter.IntToIPv4(nNum)
	if cr.Contains(nIP) {
		return nIP
	}
	return nil
}
