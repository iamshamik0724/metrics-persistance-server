package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	Hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		Hub: hub,
	}
}

func (h *Handler) HandleConnections(c *gin.Context) {
	conn, err := h.Upgrade(c.Writer, c.Request)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}
	h.Hub.AddClient(conn)

	successMessage := []byte("Successfully connected to WebSocket")
	err = conn.WriteMessage(websocket.TextMessage, successMessage)
	if err != nil {
		log.Printf("Error sending success message to client: %v", err)
	}
}

func (h *Handler) Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
