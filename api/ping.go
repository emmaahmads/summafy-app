package api

import (
	"log"
	"net/http"
	"time"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Validate origins
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{"http://localhost:8081"}

		// In production, use environment variables instead of hardcoded values
		return slices.Contains(allowedOrigins, origin)
	},
}

// HandlerWebSocket handles WebSocket connections for the ping mechanism
func (server *Server) HandlerWebSocket(c *gin.Context) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	done := make(chan struct{})
	defer close(done)

	// Set read deadline to detect disconnections
	if err := conn.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
		log.Printf("Failed to set read deadline: %v", err)
		return
	}
	conn.SetPongHandler(func(string) error {
		if err := conn.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
			log.Printf("Failed to extend read deadline: %v", err)
		}
		return nil
	})

	// Start ping-pong mechanism
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					log.Printf("Error sending ping: %v", err)
					return
				}
			case <-done:
				return
			}
		}
	}()

	// Handle incoming messages (mainly pong responses)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// If we receive any message, just log it for debugging
		log.Printf("Received message: %s", message)
	}
}

// HandlerKeepAlive is an HTTP endpoint for keepalive checks
// This is used as a fallback when WebSocket connection fails
func (server *Server) HandlerKeepAlive(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
