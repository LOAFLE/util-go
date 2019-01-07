package gateway

import (
	"net"
	"os/exec"
)

func DiscoverGateway() (ip net.IP, iface string, err error) {
	routeCmd := exec.Command("netstat", "-rn")
	output, err := routeCmd.CombinedOutput()
	if err != nil {
		return nil, "", err
	}

	return parseBSDSolarisNetstat(output)
}
