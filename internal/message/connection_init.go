package message

func CreateConnectionMessage() []byte {
	return CreateMessage(ConnectionInit, []byte{})
}
