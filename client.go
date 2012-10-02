package rfc868

import (
	"net"
)

func RequestTime(addr string) (uint, error) {
	rx := make([]byte, 4)

	udpaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return 0, err
	}

	udpconn, err := net.DialUDP("udp", nil, udpaddr)
	if err != nil {
		return 0, err
	}

	_, err = udpconn.Write([]byte{})
	if err != nil {
		return 0, err
	}

	_, err = udpconn.Read(rx)
	if err != nil {
		return 0, err
	}

	err = udpconn.Close()
	if err != nil {
		return 0, err
	}

	return to_uint(rx), nil
}
