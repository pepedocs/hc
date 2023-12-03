package internal

import (
	"net"
)

// Gets free/unused network ports
func GetFreePorts(numPorts int) ([]int, error) {
	var ports []int

	for idx := 0; idx < numPorts; idx++ {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			return ports, err
		}
		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return ports, err
		}

		defer listener.Close()
		ports = append(ports, listener.Addr().(*net.TCPAddr).Port)
	}
	return ports, nil

}
