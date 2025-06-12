package tsuniqid

import (
	"net"
	"testing"
)

// TestGetLocalIP tests the getLocalIP function to ensure it returns a valid IP address.
func TestGetLocalIP(t *testing.T) {
	ip, err := getLocalIP()

	// The function should either return a valid IP or an error
	if err != nil {
		t.Logf("getLocalIP returned error (this may be expected in some environments): %v", err)
		return
	}

	// If no error, IP should be valid
	if ip == nil {
		t.Error("getLocalIP returned nil IP without error")
		return
	}

	// IP should be IPv4
	if ip.To4() == nil {
		t.Errorf("getLocalIP returned non-IPv4 address: %v", ip)
	}

	// IP should not be loopback
	if ip.IsLoopback() {
		t.Errorf("getLocalIP returned loopback address: %v", ip)
	}

	t.Logf("getLocalIP returned: %v", ip)
}

// TestExtractIPFromAddr tests the extractIPFromAddr function with various address types.
func TestExtractIPFromAddr(t *testing.T) {
	testCases := []struct {
		name     string
		addr     net.Addr
		expected bool // whether we expect a valid IP
	}{
		{
			name:     "Valid IPv4 IPNet",
			addr:     &net.IPNet{IP: net.ParseIP("192.168.1.100"), Mask: net.CIDRMask(24, 32)},
			expected: true,
		},
		{
			name:     "Valid IPv4 IPAddr",
			addr:     &net.IPAddr{IP: net.ParseIP("10.0.0.1")},
			expected: true,
		},
		{
			name:     "IPv6 address (should be filtered out)",
			addr:     &net.IPNet{IP: net.ParseIP("2001:db8::1"), Mask: net.CIDRMask(64, 128)},
			expected: false,
		},
		{
			name:     "Loopback address (should be filtered out)",
			addr:     &net.IPNet{IP: net.ParseIP("127.0.0.1"), Mask: net.CIDRMask(8, 32)},
			expected: false,
		},
		{
			name:     "Nil IP in IPNet",
			addr:     &net.IPNet{IP: nil, Mask: net.CIDRMask(24, 32)},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ip := extractIPFromAddr(tc.addr)

			if tc.expected {
				if ip == nil {
					t.Errorf("Expected valid IP, got nil")
				} else if ip.To4() == nil {
					t.Errorf("Expected IPv4 address, got %v", ip)
				} else if ip.IsLoopback() {
					t.Errorf("Expected non-loopback address, got %v", ip)
				}
			} else {
				if ip != nil {
					t.Errorf("Expected nil IP, got %v", ip)
				}
			}
		})
	}
}

// TestExtractIPFromAddr_UnsupportedType tests extractIPFromAddr with unsupported address types.
func TestExtractIPFromAddr_UnsupportedType(t *testing.T) {
	// Test with an unsupported address type
	unsupportedAddr := &net.TCPAddr{IP: net.ParseIP("192.168.1.1"), Port: 8080}

	ip := extractIPFromAddr(unsupportedAddr)
	if ip != nil {
		t.Errorf("Expected nil for unsupported address type, got %v", ip)
	}
}
