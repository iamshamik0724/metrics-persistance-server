package message

import (
	"encoding/binary"
	"errors"
	"log"
	"math"
	"math/rand"
	"time"
)

func CreateMessage(messageType string, payload []byte) []byte {

	const baseSize = 14

	message := make([]byte, baseSize+len(payload))

	//Add version
	message[0] = 0x01

	//addd message type
	message[1] = GetByteForMessageType(messageType)

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

func ParseMessage(buffer []byte) (*Message, error) {
	if len(buffer) < 14 {
		return nil, errors.New("buffer is too short")
	}

	version := string(buffer[0])
	messageType := GetMessageType(buffer[1])
	seqNum := int(binary.BigEndian.Uint32(buffer[2:6]))
	timestamp := binary.BigEndian.Uint64(buffer[6:14])
	payload := buffer[14:]
	messagePayload, parseErr := ParsePayload(payload, messageType)
	if parseErr != nil {
		log.Println("Error while parsing message payload")
		return nil, parseErr
	}

	return &Message{
		Version:        version,
		MessageType:    messageType,
		SequenceNumber: seqNum,
		Payload:        messagePayload,
		Timestamp:      timestamp,
	}, nil
}

func ParsePayload(payload []byte, messageType string) (*Payload, error) {
	switch messageType {
	case MetricsData:
		return parseMetricsData(payload)
	case ConnectionInit, HeartBeat:
		// For ConnectionInit and HeartBeat, we expect empty payloads
		if len(payload) != 0 {
			return nil, errors.New("unexpected non-empty payload")
		}
		return nil, nil
	default:
		return nil, errors.New("unknown message type")
	}
}

// parseMetricsData parses the payload for MetricsData type
func parseMetricsData(payload []byte) (*Payload, error) {
	if len(payload) < 10 {
		return nil, errors.New("payload is too short for MetricsData")
	}

	statusCode := int(binary.BigEndian.Uint16(payload[0:2]))
	responseTime := binary.BigEndian.Uint64(payload[2:10])
	method := GetMethod(payload[10])

	// The rest of the payload is the route
	route := payload[11:]
	routeString := string(trimNullBytes(route))

	// Convert responseTime from uint64 to float64 (assuming it's a raw double precision float in binary)
	// You can use any conversion method based on how the float is encoded
	respTime := math.Float64frombits(responseTime)

	return &Payload{
		Route:        routeString,
		Method:       method,
		StatusCode:   statusCode,
		ResponseTime: respTime,
	}, nil
}

func trimNullBytes(route []byte) []byte {
	i := len(route) - 1
	for i >= 0 && route[i] == 0x00 {
		i--
	}
	return route[:i+1]
}
