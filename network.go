// Package tsuniqid - Network utilities for machine identification
package tsuniqid

import (
	"errors"
	"net"
)

// getLocalIP retrieves the first available non-loopback IPv4 address from network interfaces.
// This function iterates through all network interfaces and returns the first valid local IP address.
//
// Returns:
//   - net.IP: The first available local IPv4 address
//   - error: An error if no valid IP address is found
func getLocalIP() (net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		// Skip interfaces that are down
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		// Skip loopback interfaces
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addresses, err := iface.Addrs()
		if err != nil {
			continue // Skip this interface if we can't get addresses
		}

		for _, addr := range addresses {
			ip := extractIPFromAddr(addr)
			if ip != nil {
				return ip, nil
			}
		}
	}

	return nil, errors.New("no valid local IP address found")
}

// extractIPFromAddr extracts an IPv4 address from a network address.
// This function handles both *net.IPNet and *net.IPAddr types and filters out
// loopback addresses and IPv6 addresses.
//
// Parameters:
//   - addr: The network address to extract IP from
//
// Returns:
//   - net.IP: The extracted IPv4 address, or nil if not valid
func extractIPFromAddr(addr net.Addr) net.IP {
	var ip net.IP

	// Handle different address types
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	default:
		return nil
	}

	// Filter out invalid addresses
	if ip == nil || ip.IsLoopback() {
		return nil
	}

	// Convert to IPv4 and filter out IPv6 addresses
	ipv4 := ip.To4()
	if ipv4 == nil {
		return nil
	}

	return ipv4
}
