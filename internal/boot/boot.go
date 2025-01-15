package boot

import (
	"context"
	"fmt"
	"log"
	"metrics-persistance-server/config"
	"metrics-persistance-server/internal/message"
	"metrics-persistance-server/internal/metrics"
	"metrics-persistance-server/internal/metrics/repo"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitilizeDb() {}

func InitializeUdpConnection(config *config.Config, ctx context.Context) error {

	//Initialize DB
	db, dbErr := InitializeDb(config)

	if dbErr != nil {
		log.Println("Error Initializing database")
		return dbErr
	}

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

	//Initialize repo
	apiMetricRepo := repo.NewApiMetricRepository(db)

	//Initialize Service
	apiMetricService := metrics.NewApiMetricService(apiMetricRepo)

	//Initialize consumer channel
	go startConsumerChannel(ctx, conn, apiMetricService)

	//Initialize Gin Router
	router := gin.Default()

	//Initialize handlers
	metricsHandler := metrics.NewHandler(apiMetricService)

	//Initialize Routes
	router.GET("/metrics", metricsHandler.GetMetrics)

	router.Run(":8085")

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

func startConsumerChannel(ctx context.Context, conn *net.UDPConn, metricService metrics.IService) {
	log.Println("Starting consumer channel...")
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping consumer channel...")
			return
		default:
			buffer := make([]byte, 1024)
			_, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				log.Printf("Error reading UDP message: %v\n", err)
				continue
			}
			//log.Printf("Received message from %v: %s %d %s\n", addr, string(buffer[:n]), len(buffer), string(buffer))
			message, errMessage := message.ParseMessage(buffer)
			if errMessage != nil {
				log.Printf("Error while parsing message")
			}
			//log.Println(message)
			//log.Println(message.Payload)
			recordMetricError := metricService.RecordMetric(message)
			if recordMetricError != nil {
				log.Println("Error while recording metric message. Message: ", message, "Payload:", message.Payload)
			}
		}
	}
}

func InitializeDb(config *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.DatabaseConfig.Host, config.DatabaseConfig.User, config.DatabaseConfig.Password, config.DatabaseConfig.Database, config.DatabaseConfig.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error opening database connection: %v\n", err)
		return nil, err
	}
	log.Println("Database initialized successfully")
	return db, nil
}
