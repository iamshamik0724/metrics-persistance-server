package message

func CreateHeartbeatMessage() []byte {
	return CreateMessage(HeartBeat, []byte{})
}
