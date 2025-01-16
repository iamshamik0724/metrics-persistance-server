package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"metrics-persistance-server/internal/message"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mutex      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte, 256), // Buffered channel to avoid blocking
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (h *Hub) AddClient(conn *websocket.Conn) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.clients[conn] = true
}

func (h *Hub) RemoveClient(conn *websocket.Conn) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	err := conn.Close()
	if err != nil {
		log.Println("Error closing WebSocket connection:", err)
	}
	delete(h.clients, conn)
}

func (h *Hub) BroadcastMessage(message []byte) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for client := range h.clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Broadcast error:", err)
			client.Close()
			delete(h.clients, client)
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.AddClient(conn)
		case conn := <-h.unregister:
			h.RemoveClient(conn)
		case message := <-h.broadcast:
			h.BroadcastMessage(message)
		}
	}
}

func (h *Hub) SendMessage(message []byte) {
	h.broadcast <- message
}

func (h *Hub) SendMetricsMessage(msg *message.Message) {
	metricMessage := &MetricsBroadcastMessage{
		Route:        fmt.Sprintf("%s - %s", msg.Payload.Method, msg.Payload.Route),
		Timestamp:    msg.Timestamp,
		ResponseTime: msg.Payload.ResponseTime,
		Status:       msg.Payload.StatusCode,
	}
	data, err := json.Marshal(metricMessage)
	if err != nil {
		log.Println("Error marshaling message:", err)
		return
	}
	h.broadcast <- data
}

func (h *Hub) Shutdown() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	for client := range h.clients {
		client.Close()
		delete(h.clients, client)
	}
	close(h.broadcast)
	close(h.register)
	close(h.unregister)
}
