package message

import "time"

type Message struct {
	Version        string
	MessageType    string
	SequenceNumber int
	Timestamp      time.Time
	Payload        *Payload
}

type Payload struct {
	Route        string
	Method       string
	StatusCode   int
	ResponseTime float64
}
