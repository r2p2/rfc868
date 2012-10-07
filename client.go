package rfc868

import (
	"net"
	"time"
	"encoding/binary"
	"fmt"
)

type client struct {
	addr *net.UDPAddr
	conn *net.UDPConn
	data []byte
}

var epoch time.Time

func init() {
	var err error
	epoch, err = time.Parse("2006-01-02 15:04:05", "1900-01-01 00:00:00")
	if err != nil {
		panic(err)
	}
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
func (c *client) RequestTime() (time.Time, error) {
	_, err := c.conn.Write([]byte{})
	if err != nil {
		return epoch, err
	}

	var n uint32
	err = binary.Read(c.conn, binary.BigEndian, &n)
	if err != nil {
		return epoch, err
	}

	if err != nil {
		return epoch, err
	}
	t := epoch.Add(time.Duration(n) * time.Second)

	return t, nil
}

// If you don't want to handle the client object
// use this function.
func RequestTime(addr string) (time.Time, error) {
	c, err := NewClient(addr)
	if err != nil {
		return epoch, err
	}

	t, err := c.RequestTime()
	if err != nil {
		return epoch, err
	}

	return t, nil
}
