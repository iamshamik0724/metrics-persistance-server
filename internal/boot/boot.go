package boot

import (
	"context"
	"fmt"
	"log"
	"metrics-persistance-server/config"
	"metrics-persistance-server/internal/message"
	"metrics-persistance-server/internal/metrics"
	"metrics-persistance-server/internal/metrics/repo"
	"metrics-persistance-server/internal/websocket"
	"net"
	"net/http"
	"time"

	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitilizeDb() {}

func Initialize(config *config.Config, ctx context.Context, stopChannel chan struct{}) error {

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

	//Send connection Init message
	errInitMessage := sendInitMessage(conn)
	if errInitMessage != nil {
		log.Println("Failed to send Init Message. Server Start failed")
		return err
	}

	var wg sync.WaitGroup

	//Waitgroup for heartbeat and consumer channel closures
	wg.Add(2)

	// Initialize heartbeat
	stopHeartbeatChannel := make(chan struct{})
	go startHeartbeat(conn, config.UdpServer.HeartBeatInterval, stopHeartbeatChannel, &wg)

	//Initialize repo
	apiMetricRepo := repo.NewApiMetricRepository(db)

	//Initialize Service
	apiMetricService := metrics.NewApiMetricService(apiMetricRepo)

	//Initialize Gin Router
	router := gin.Default()

	//Initialize and run web socket hub
	hub := websocket.NewHub()
	go hub.Run()

	//Initialize handlers
	metricsHandler := metrics.NewHandler(apiMetricService)
	webSocketHandler := websocket.NewHandler(hub)

	//Initialize consumer channel
	stopConsumerChannel := make(chan struct{})
	go startConsumerChannel(conn, apiMetricService, hub, stopConsumerChannel, &wg)

	//Initialize Routes
	router.GET("/ws", webSocketHandler.HandleConnections)
	router.GET("/metrics", metricsHandler.GetMetrics)

	srv := &http.Server{
		Addr:    ":8085",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	go func() {
		<-stopChannel
		close(stopHeartbeatChannel)
		close(stopConsumerChannel)
		wg.Wait()
		log.Println("Context canceled. Closing UDP connection...")
		conn.Close()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("HTTP server Shutdown: %v", err)
		} else {
			log.Println("HTTP server stopped")
		}
	}()

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

func startHeartbeat(conn *net.UDPConn, heartBeatInterval int, stopChannel chan struct{}, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Duration(heartBeatInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stopChannel:
			log.Println("Stopping heartbeat channel due to stop signal...")
			wg.Done()
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

func startConsumerChannel(conn *net.UDPConn, metricService metrics.IService, hub *websocket.Hub, stopChannel chan struct{}, wg *sync.WaitGroup) {
	log.Println("Starting consumer channel...")
	for {
		select {
		case <-stopChannel:
			log.Println("Stopping Consumer channel due to stop signal...")
			wg.Done()
			return
		default:
			buffer := make([]byte, 1024)
			_, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				log.Printf("Error reading UDP message: %v\n", err)
				continue
			}
			message, errMessage := message.ParseMessage(buffer)
			if errMessage != nil {
				log.Printf("Error while parsing message")
			}
			recordMetricError := metricService.RecordMetric(message)
			if recordMetricError != nil {
				log.Println("Error while recording metric message. Message: ", message, "Payload:", message.Payload)
			}
			hub.SendMetricsMessage(message)
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
