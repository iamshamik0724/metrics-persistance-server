package message

import (
	"encoding/binary"
	"math/rand"
	"time"
)

func CreateMessage(messageType string, payload []byte) []byte {

	const baseSize = 14

	message := make([]byte, baseSize+len(payload))

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
	sequence := rand.Uint32()
	binary.BigEndian.PutUint32(message[2:6], sequence)

	// add time
	timestamp := time.Now().UnixNano()
	binary.BigEndian.PutUint64(message[6:14], uint64(timestamp))

	//add payload
	if len(payload) > 0 {
		copy(message[baseSize:], payload)
	}

	return message
}
