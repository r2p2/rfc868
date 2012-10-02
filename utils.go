package rfc868

import "time"

func to_byte(time uint, data *[]byte) {
	(*data)[0] = byte(time)
	(*data)[1] = byte(time >> 8)
	(*data)[2] = byte(time >> 16)
	(*data)[3] = byte(time >> 24)
}

func to_uint(data []byte) (time uint) {
	time |= uint(data[0])
	time |= uint(data[1]) << 8
	time |= uint(data[2]) << 16
	time |= uint(data[3]) << 24

	return
}

func fpsCounter() func() uint32 {
	var nextFPS, currentFPS uint32
	var timestamp time.Time
	fpsDelay, _ := time.ParseDuration("1s")

	nextFPS = 1
	currentTime := time.Now()
	timestamp = currentTime.Add(fpsDelay)

	return func() uint32 {
		if now := time.Now(); timestamp.Sub(now) <= 0 {
			currentFPS = nextFPS
			nextFPS = 1
			timestamp = now.Add(fpsDelay)
		} else {
			nextFPS++
		}
		return currentFPS
	}
}
