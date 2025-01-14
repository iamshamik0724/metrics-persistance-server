package message

const (
	GET     byte = 0x01
	POST    byte = 0x02
	PUT     byte = 0x03
	DELETE  byte = 0x04
	PATCH   byte = 0x05
	OPTIONS byte = 0x06
	HEAD    byte = 0x07
	TRACE   byte = 0x08
	CONNECT byte = 0x09
)

// Convert Method byte value to string representation
func GetMethod(method byte) string {
	switch method {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	case PATCH:
		return "PATCH"
	case OPTIONS:
		return "OPTIONS"
	case HEAD:
		return "HEAD"
	case TRACE:
		return "TRACE"
	case CONNECT:
		return "CONNECT"
	default:
		return "UNKNOWN"
	}
}
