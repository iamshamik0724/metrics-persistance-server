package boot

import (
	"log"
	"metrics-persistance-server/config"
	"metrics-persistance-server/internal/message"
	"net"
)

func InitilizeDb() {}

func InitlizeUdpConnection(config *config.Config) error {

	// Initialize UDP
	udpAddr, err := net.ResolveUDPAddr("udp", config.UdpServer.Address)
	if err != nil {
		log.Println("Error resolving address:", err)
		return err
	}

	// Create a UDP connection
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Println("Error connecting to server:", err)
		return err
	}

	//TODO: handle connection close on server shutdown

	//Send connection Init message
	errInitMessage := sendInitMessage(conn)
	if errInitMessage != nil {
		panic("Failed to send Init Message. Server Start failed")
	}

	// Initialize heartbeat

	//Initialize consumer channel

	return nil

}

func sendInitMessage(conn *net.UDPConn) error {

	message := message.CreateConnectionMessage()

	_, err := conn.Write(message)
	if err != nil {
		log.Println("Error while sending connection initialize message")
		return err
	}

	return nil

}
