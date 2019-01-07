package gateway

import (
	"net"
	"os/exec"
)

func DiscoverGateway() (ip net.IP, iface string, err error) {
	output, err := exec.Command("cmd.exe", "/C", "chcp 437 && route print 0.0.0.0").CombinedOutput()
	if err != nil {
		panic(err)
	}
	if err != nil {
		return nil, "", err
	}

	return parseWindowsRoutePrint(output)
}
