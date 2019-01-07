package converter

import (
	"net"
	"regexp"
	"strconv"
)

var (
	macStripRegexp = regexp.MustCompile(`[^a-fA-F0-9]`)
)

func MacToUint64(hwAddr net.HardwareAddr) (uint64, error) {
	mac := hwAddr.String()
	hex := macStripRegexp.ReplaceAllLiteralString(mac, "")

	return strconv.ParseUint(hex, 16, 64)
}
