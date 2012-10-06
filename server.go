package rfc868

import (
	"fmt"
	"net"
	"sync"
	"time"
	"bytes"
	"encoding/binary"
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

	// In all this, we hope to avoid memory allocation. Allocation
	// shouldn't happen since we are Write()ing 4 bytes -- the
	// size of int32 -- to a bytes.Buffer that is initialized from
	// a byte slice that is already exactly that size, and we call
	// bytes.Buffer.Reset() before writing. I have checked that
	// this holds, but a sanity check might be in order (with a
	// panic() if not passed).
	buf := bytes.NewBuffer(th.data)
	buf.Reset()
	// We must send back time encoded as a Big endian 32 bit (signed) int
	err := binary.Write(buf, binary.BigEndian,
		int32(time.Since(th.rfc868_epoch)/1000000000))
	if err != nil {
		panic(err)
	}

	// fmt.Println(buf.Len(), cap(buf.Bytes()))

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
