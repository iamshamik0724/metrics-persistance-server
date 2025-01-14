package message

const (
	ConnectionInit string = "ConnectionInit"
	HeartBeat      string = "HeartBeat"
	MetricsData    string = "MetricsData"
)

func GetMessageType(messageType byte) string {
	switch messageType {
	case 0x01:
		return ConnectionInit
	case 0x02:
		return HeartBeat
	case 0x03:
		return MetricsData
	}
	return ""
}

func GetByteForMessageType(messageType string) byte {
	switch messageType {
	case ConnectionInit:
		return 0x01
	case HeartBeat:
		return 0x02
	case MetricsData:
		return 0x03
	}
	return 0x00
}
