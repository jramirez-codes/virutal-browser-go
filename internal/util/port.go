package util

import (
	"fmt"
	"net"
)

func GetPort() (int, error) {
	for port := 3000; port <= 3999; port++ {
		addr := fmt.Sprintf(":%d", port)

		ln, err := net.Listen("tcp", addr)
		if err == nil {
			ln.Close() // free it immediately
			return port, nil
		}
	}

	return 0, fmt.Errorf("no available port found in range 3000–3999")
}
