package integration

import (
	"fmt"
	"net"
)

func ScanForOpenPort(start, end uint) (uint, error) {
	for port := start; port <= end; port++ {
		ln, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return port, nil
		}
		ln.Close()
	}
	return 0, fmt.Errorf("no open IPv4 port found in range %d-%d", start, end)
}
