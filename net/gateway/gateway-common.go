package gateway

import (
	"errors"
	"net"
	"strings"
)

var errNoGateway = errors.New("no gateway found")

func parseWindowsRoutePrint(output []byte) (net.IP, string, error) {
	// Windows route output format is always like this:
	// ===========================================================================
	// Active Routes:
	// Network Destination        Netmask          Gateway       Interface  Metric
	//           0.0.0.0          0.0.0.0      192.168.1.1    192.168.1.100     20
	// ===========================================================================
	// I'm trying to pick the active route,
	// then jump 2 lines and pick the third IP
	// Not using regex because output is quite standard from Windows XP to 8 (NEEDS TESTING)
	lines := strings.Split(string(output), "\n")
	for idx, line := range lines {
		if strings.HasPrefix(line, "Active Routes:") {
			if len(lines) <= idx+2 {
				return nil, "", errNoGateway
			}

			fields := strings.Fields(lines[idx+2])
			if len(fields) < 3 {
				return nil, "", errNoGateway
			}

			ip := net.ParseIP(fields[2])
			if ip != nil {
				return ip, fields[3], nil
			}
		}
	}
	return nil, "", errNoGateway
}

func parseLinuxIPRoute(output []byte) (net.IP, string, error) {
	// Linux '/usr/bin/ip route show' format looks like this:
	// default via 192.168.178.1 dev wlp3s0  metric 303
	// 192.168.178.0/24 dev wlp3s0  proto kernel  scope link  src 192.168.178.76  metric 303
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 && fields[0] == "default" {
			ip := net.ParseIP(fields[2])
			if ip != nil {
				return ip, fields[4], nil
			}
		}
	}

	return nil, "", errNoGateway
}

func parseLinuxRoute(output []byte) (net.IP, string, error) {
	// Linux route out format is always like this:
	// Kernel IP routing table
	// Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
	// 0.0.0.0         192.168.1.1     0.0.0.0         UG    0      0        0 eth0
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "0.0.0.0" {
			ip := net.ParseIP(fields[1])
			if ip != nil {
				return ip, fields[7], nil
			}
		}
	}

	return nil, "", errNoGateway
}

func parseDarwinRouteGet(output []byte) (net.IP, string, error) {
	// Darwin route out format is always like this:
	//    route to: default
	// destination: default
	//        mask: default
	//     gateway: 192.168.1.1
	//   interface: tun0
	//       flags: <UP,GATEWAY,DONE,STATIC,PRCLONING>
	// lines := strings.Split(string(output), "\n")
	// for _, line := range lines {
	// 	fields := strings.Fields(line)
	// 	if len(fields) >= 2 && fields[0] == "gateway:" {
	// 		ip := net.ParseIP(fields[1])
	// 		if ip != nil {
	// 			return ip, "", nil
	// 		}
	// 	}
	// }

	// Darwin route out format is always like this:
	// Internet:
	// Destination        Gateway            Flags        Refs      Use   Netif Expire
	// default            192.168.10.254     UGSc          194        0     en3
	// 127                127.0.0.1          UCS             0      429     lo0
	// 127.0.0.1          127.0.0.1          UH              1   587632     lo0
	// 169.254            link#7             UCS             0        0     en3
	// 192.168.10         link#7             UCS             4        0     en3
	// 192.168.10.1       0:11:32:7f:20:61   UHLWIi          1      202     en3   1065
	// 224.0.0/4          link#7             UmCS            3        0     en3
	// 224.0.0.251        1:0:5e:0:0:fb      UHmLWI          0     2325     en3
	// 239.192.152.143    1:0:5e:40:98:8f    UHmLWI          0    22892     en3
	// 239.255.255.250    1:0:5e:7f:ff:fa    UHmLWI          0    15988     en3
	// 255.255.255.255/32 link#7             UCS             0        0     en3

	// Internet6:
	// Destination                             Gateway                         Flags         Netif Expire
	// default                                 fe80::%utun0                    UGcI          utun0
	// default                                 fe80::%utun1                    UGcI          utun1
	// default                                 fe80::%utun2                    UGcI          utun2
	// default                                 fe80::%utun3                    UGcI          utun3
	outputLines := strings.Split(string(output), "\n")
	for _, line := range outputLines {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "default" {
			ip := net.ParseIP(fields[1])
			if ip != nil {
				return ip, fields[5], nil
			}
		}
	}
	return nil, "", errNoGateway
}

func parseBSDSolarisNetstat(output []byte) (net.IP, string, error) {
	// netstat -rn produces the following on FreeBSD:
	// Routing tables
	//
	// Internet:
	// Destination        Gateway            Flags      Netif Expire
	// default            10.88.88.2         UGS         em0
	// 10.88.88.0/24      link#1             U           em0
	// 10.88.88.148       link#1             UHS         lo0
	// 127.0.0.1          link#2             UH          lo0
	//
	// Internet6:
	// Destination                       Gateway                       Flags      Netif Expire
	// ::/96                             ::1                           UGRS        lo0
	// ::1                               link#2                        UH          lo0
	// ::ffff:0.0.0.0/96                 ::1                           UGRS        lo0
	// fe80::/10                         ::1                           UGRS        lo0
	// ...
	outputLines := strings.Split(string(output), "\n")
	for _, line := range outputLines {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "default" {
			ip := net.ParseIP(fields[1])
			if ip != nil {
				return ip, fields[3], nil
			}
		}
	}

	return nil, "", errNoGateway
}
