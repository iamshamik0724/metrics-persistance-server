package message

import (
	"encoding/binary"
	"time"
)

func CreateMessage(messageType string, payload []byte) []byte {

	message := make([]byte, 14)

	//Add version
	message[0] = 0x01

	//addd message type
	switch messageType {
	case ConnectionInit:
		message[1] = 0x01
	case HeartBeat:
		message[1] = 0x02
	case MetricsData:
		message[1] = 0x03
	}

	//add seq
	//ToDO: random sequence generation
	message[2] = 0x01
	message[3] = 0x03
	message[4] = 0x03
	message[5] = 0x03

	// add time
	currentTime := time.Now()
	seconds := currentTime.Unix()
	nanoseconds := currentTime.Nanosecond()
	timestamp := (int64(seconds) << 32) | int64(nanoseconds)
	var timestampBytes [8]byte
	binary.BigEndian.PutUint64(timestampBytes[:], uint64(timestamp))
	copy(message[6:], timestampBytes[:])

	//add payload
	copy(message[14:], payload[:])
	return message
}
