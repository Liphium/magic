package integration

import (
	"fmt"
	"net"
)

func ScanForOpenPort(start, end uint) (uint, error) {
	for port := start; port <= end; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			ln.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no open port found in range %d-%d", start, end)
}
