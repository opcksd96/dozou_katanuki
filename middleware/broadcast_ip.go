// middleware/broadcast_ip.go (100行以下)
package middleware

import (
	"fmt"
	"net"
)

func GetLocalIPv4s() []string {
	var ips []string
	ifaces, err := net.Interfaces()
	if err != nil { return []string{"127.0.0.1"} }
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 { continue }
		addrs, err := iface.Addrs()
		if err != nil { continue }
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet: ip = v.IP
			case *net.IPAddr: ip = v.IP
			}
			if ip == nil || ip.IsLoopback() { continue }
			if ip4 := ip.To4(); ip4 != nil { ips = append(ips, ip4.String()) }
		}
	}
	if len(ips) == 0 { ips = append(ips, "127.0.0.1") }
	return ips
}

func GetLocalSubnets() []string {
	seen := make(map[string]bool)
	var subnets []string
	ifaces, err := net.Interfaces()
	if err != nil { return []string{"127.0.0.1/32"} }
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 { continue }
		addrs, err := iface.Addrs()
		if err != nil { continue }
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if ipNet.IP == nil || ipNet.IP.IsLoopback() { continue }
				if ip4 := ipNet.IP.To4(); ip4 != nil {
					maskedIP := ip4.Mask(ipNet.Mask)
					ones, _ := ipNet.Mask.Size()
					cidr := fmt.Sprintf("%s/%d", maskedIP.String(), ones)
					if !seen[cidr] { seen[cidr] = true; subnets = append(subnets, cidr) }
				}
			}
		}
	}
	if !seen["127.0.0.1/32"] { subnets = append(subnets, "127.0.0.1/32") }
	return subnets
}
