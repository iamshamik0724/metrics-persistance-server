package message

type Message struct {
	Version        string
	MessageType    string
	SequenceNumber int
	Timestamp      uint64
	Payload        *Payload
}

type Payload struct {
	Route        string
	Method       string
	StatusCode   int
	ResponseTime float64
}
