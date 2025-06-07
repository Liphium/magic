package integration

import (
	"fmt"
	"net"
)

// Scan a range of ports for an open port
func ScanForOpenPort(start, end uint) (uint, error) {
	for port := start; port <= end; port++ {
		if !ScanPort(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no open IPv4 port found in range %d-%d", start, end)
}

// Scan an individual port. Returns true when the connection succeeds.
func ScanPort(port uint) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err == nil {
		conn.Close()
		return true
	}
	return false
}
