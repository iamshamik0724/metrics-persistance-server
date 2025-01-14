package boot

import (
	"context"
	"log"
	"metrics-persistance-server/config"
	"metrics-persistance-server/internal/message"
	"net"
	"time"
)

func InitilizeDb() {}

func InitializeUdpConnection(config *config.Config, ctx context.Context) error {

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

	go func() {
		<-ctx.Done()
		log.Println("Context canceled. Closing UDP connection...")
		conn.Close()
	}()

	//Send connection Init message
	errInitMessage := sendInitMessage(conn)
	if errInitMessage != nil {
		log.Println("Failed to send Init Message. Server Start failed")
		return err
	}

	// Initialize heartbeat
	go startHeartbeat(ctx, conn, config.UdpServer.HeartBeatInterval)

	//Initialize consumer channel
	go startConsumerChannel(ctx, conn)

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

func startHeartbeat(ctx context.Context, conn *net.UDPConn, heartBeatInterval int) {
	ticker := time.NewTicker(time.Duration(heartBeatInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping heartbeat...")
			return
		case <-ticker.C:
			log.Println("Sending heartbeat...")
			msg := message.CreateHeartbeatMessage()
			if _, err := conn.Write(msg); err != nil {
				log.Printf("Error sending heartbeat: %v\n", err)
			}
		}
	}
}

func startConsumerChannel(ctx context.Context, conn *net.UDPConn) {
	log.Println("Starting consumer channel...")
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping consumer channel...")
			return
		default:
			// Example: Read from UDP and process messages
			buffer := make([]byte, 1024)
			n, addr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				log.Printf("Error reading UDP message: %v\n", err)
				continue
			}
			log.Printf("Received message from %v: %s %d %s\n", addr, string(buffer[:n]), len(buffer), string(buffer))
			message, errMessage := message.ParseMessage(buffer)
			if errMessage != nil {
				log.Printf("Error while parsing message")
			}
			log.Println(message)
			log.Println(message.Payload)
		}
	}
}
