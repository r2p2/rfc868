package rfc868

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type TimeHandle struct {
	data         []byte
	rfc868_epoch time.Time

	mtx sync.RWMutex
}

func NewTimeHandle() (*TimeHandle, error) {
	epoch, err := time.Parse("2006-01-02 15:04:05", "1900-01-01 00:00:00")
	if err != nil {
		return nil, err
	}

	return &TimeHandle{
		make([]byte, 4),
		epoch,
		sync.RWMutex{},
	}, nil
}

func (th *TimeHandle) update() {
	th.mtx.Lock()

	to_byte(uint(time.Since(th.rfc868_epoch)/1000000000), &th.data)

	th.mtx.Unlock()
}

func update_service(th *TimeHandle) {
	var s1 time.Duration = 1000000000
	for {
		th.update()
		time.Sleep(s1)
	}
}

func (th *TimeHandle) send(udpconn *net.UDPConn, caddr *net.UDPAddr) error {
	th.mtx.RLock()

	_, err := udpconn.WriteToUDP(th.data, caddr)
	if err != nil {
		return err
	}

	th.mtx.RUnlock()

	return nil
}

func ServeTime(addr string) error {
	tx := make([]byte, 4)
	timehandle, err := NewTimeHandle()
	go update_service(timehandle)

	if err != nil {
		return err
	}

	udpaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	udpconn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		return err
	}

	//fps := fpsCounter()
	for {
		_, caddr, err := udpconn.ReadFromUDP(tx)
		if err != nil {
			fmt.Println("error: " + err.Error())
			continue
		}

		err = timehandle.send(udpconn, caddr)
		if err != nil {
			fmt.Println("error: " + err.Error())
			continue
		}

		//fmt.Println("rps", fps())
	}

	err = udpconn.Close()
	if err != nil {
		return err
	}

	return nil
}
