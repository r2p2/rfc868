package rfc868

import (
	"net"
)

type client struct {
	addr *net.UDPAddr
	conn *net.UDPConn
	data []byte
}

// Create a new client object and establishes a connection
// for later usage.
func NewClient(addr string) (*client, error) {
	udpaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	udpconn, err := net.DialUDP("udp", nil, udpaddr)
	if err != nil {
		return nil, err
	}

	return &client{
		udpaddr,
		udpconn,
		make([]byte, 4),
	}, nil
}

// Uses the established connection to request the current
// time from the server.
func (c *client) RequestTime() (uint, error) {
	_, err := c.conn.Write([]byte{})
	if err != nil {
		return 0, err
	}

	_, err = c.conn.Read(c.data)
	if err != nil {
		return 0, err
	}

	return to_uint(c.data), nil
}

// If you don't want to handle the client object
// use this function.
func RequestTime(addr string) (uint, error) {
	c, err := NewClient(addr)
	if err != nil {
		return 0, err
	}

	time, err := c.RequestTime()
	if err != nil {
		return 0, err
	}

	return time, nil
}
