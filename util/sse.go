package util

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

// SSEClient represents a single SSE client
type SSEClient struct {
	Channel chan string
	Level   string
}

// SSEServer holds the list of clients
type SSEServer struct {
	Clients map[*SSEClient]bool
	mu      sync.Mutex // untuk aman mengakses Clients
	Add     chan *SSEClient
	Remove  chan *SSEClient
	Message chan Message
}

// Message struct to include message content and level
type Message struct {
	Content string
	Level   string
}

// NewSSEServer initializes a new SSE server
func NewSSEServer() *SSEServer {
	return &SSEServer{
		Clients: make(map[*SSEClient]bool),
		Add:     make(chan *SSEClient),
		Remove:  make(chan *SSEClient),
		Message: make(chan Message),
	}
}

// Run starts the SSE server loop
func (s *SSEServer) Run() {
	for {
		select {
		case client := <-s.Add:
			s.mu.Lock()
			s.Clients[client] = true
			s.mu.Unlock()

		case client := <-s.Remove:
			s.mu.Lock()
			if _, ok := s.Clients[client]; ok {
				delete(s.Clients, client)
				close(client.Channel)
			}
			s.mu.Unlock()

		case message := <-s.Message:
			s.mu.Lock()
			for client := range s.Clients {
				if client.Level == message.Level {
					// kirim pesan (non-blocking jika channel penuh)
					select {
					case client.Channel <- message.Content:
					default:
						// jika channel penuh → bisa log atau abaikan
					}
				}
			}
			s.mu.Unlock()
		}
	}
}

// SSEHandler returns an echo.HandlerFunc for SSE streaming
func SSEHandler(server *SSEServer) echo.HandlerFunc {
	return func(c echo.Context) error {
		level := c.QueryParam("level")

		client := &SSEClient{
			Channel: make(chan string, 100), // beri buffer agar tidak mudah blocking
			Level:   level,
		}

		// Register client
		server.Add <- client

		// Pastikan client dihapus saat koneksi selesai
		defer func() {
			server.Remove <- client
		}()

		// Set header SSE yang benar
		c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
		c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

		// Optional: tambahkan header ini jika ingin lebih kompatibel
		// c.Response().Header().Set("X-Accel-Buffering", "no")

		flusher, ok := c.Response().Writer.(http.Flusher)
		if !ok {
			return c.String(http.StatusInternalServerError, "Streaming unsupported!")
		}

		for {
			select {
			case msg, ok := <-client.Channel:
				if !ok {
					return nil // channel ditutup → client keluar
				}

				// Format SSE yang benar
				fmt.Fprintf(c.Response().Writer, "data: %s\n\n", msg)
				flusher.Flush()

			case <-c.Request().Context().Done():
				// Client disconnect (lebih baik daripada CloseNotify di Echo)
				return nil
			}
		}
	}
}

// SendMessageHandler untuk mengirim pesan ke semua client dengan level tertentu
func SendMessageHandler(server *SSEServer) echo.HandlerFunc {
	return func(c echo.Context) error {
		type request struct {
			Message string `json:"message"`
			Level   string `json:"level"`
		}

		var req request
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid request body",
			})
		}

		if req.Message == "" || req.Level == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "message and level are required",
			})
		}

		// Kirim ke server
		server.Message <- Message{
			Content: req.Message,
			Level:   req.Level,
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status": "message sent",
		})
	}
}
