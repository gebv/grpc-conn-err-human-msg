package grpcx

import (
	"net"
	"time"
)

// TCPConnectOK returns an error if it was not possible to connect to the server.
func TCPConnectOK(timeout time.Duration, addr string) error {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
