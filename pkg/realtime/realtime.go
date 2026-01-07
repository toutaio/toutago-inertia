// Package realtime provides WebSocket-based real-time updates for Inertia.js applications.
package realtime

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512 * 1024 // 512 KB
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		return true // Allow all origins in development
	},
}

// Message represents a WebSocket message.
type Message struct {
	Channel string      `json:"channel"`
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
}

// Client represents a WebSocket client connection.
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	channels map[string]bool
	mu       sync.RWMutex
}

// Subscribe adds the client to a channel.
func (c *Client) Subscribe(channel string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.channels[channel] = true
}

// Unsubscribe removes the client from a channel.
func (c *Client) Unsubscribe(channel string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.channels, channel)
}

// IsSubscribed checks if client is subscribed to a channel.
func (c *Client) IsSubscribed(channel string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.channels[channel]
}

// readPump pumps messages from the WebSocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		if c.conn != nil {
			c.conn.Close()
		}
	}()

	if c.conn != nil {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		c.conn.SetPongHandler(func(string) error {
			_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
	}

	for {
		if c.conn == nil {
			return
		}

		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle subscription/unsubscription messages
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "subscribe":
			c.Subscribe(msg.Channel)
		case "unsubscribe":
			c.Unsubscribe(msg.Channel)
		}
	}
}

// writePump pumps messages from the hub to the WebSocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if c.conn != nil {
			c.conn.Close()
		}
	}()

	for {
		select {
		case message, ok := <-c.send:
			if c.conn != nil {
				_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			}
			if !ok {
				// Hub closed the channel
				if c.conn != nil {
					_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				}
				return
			}

			if c.conn == nil {
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write([]byte{'\n'})
				_, _ = w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			if c.conn != nil {
				_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}
}

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	clients    map[*Client]bool
	channels   map[string]map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		channels:   make(map[string]map[*Client]bool),
	}
}

// Run starts the hub's message processing loop.
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// Close all client connections
			h.mu.Lock()
			for client := range h.clients {
				close(client.send)
			}
			h.mu.Unlock()
			return

		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			// Add client to their subscribed channels
			client.mu.RLock()
			for channel := range client.channels {
				if _, ok := h.channels[channel]; !ok {
					h.channels[channel] = make(map[*Client]bool)
				}
				h.channels[channel][client] = true
			}
			client.mu.RUnlock()
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// Remove client from all channels
				for channel := range h.channels {
					if clients, ok := h.channels[channel]; ok {
						delete(clients, client)
						if len(clients) == 0 {
							delete(h.channels, channel)
						}
					}
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			data, err := json.Marshal(message)
			if err != nil {
				h.mu.RUnlock()
				continue
			}

			// Broadcast to all if channel is "*"
			if message.Channel == "*" {
				for client := range h.clients {
					select {
					case client.send <- data:
					default:
						// Client buffer full, close it
						go func(c *Client) {
							h.unregister <- c
						}(client)
					}
				}
			} else {
				// Broadcast to specific channel
				if clients, ok := h.channels[message.Channel]; ok {
					for client := range clients {
						select {
						case client.send <- data:
						default:
							// Client buffer full, close it
							go func(c *Client) {
								h.unregister <- c
							}(client)
						}
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast sends a message to all clients subscribed to a channel.
func (h *Hub) Broadcast(msg *Message) {
	h.broadcast <- msg
}

// Publish is a helper method to broadcast a message.
func (h *Hub) Publish(channel, msgType string, data interface{}) {
	h.Broadcast(&Message{
		Channel: channel,
		Type:    msgType,
		Data:    data,
	})
}

// HandleWebSocket handles WebSocket connection upgrades.
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	client := &Client{
		hub:      h,
		conn:     conn,
		send:     make(chan []byte, 256),
		channels: make(map[string]bool),
	}

	h.register <- client

	// Allow collection of memory referenced by the caller
	go client.writePump()
	go client.readPump()

	return nil
}

// UpdateChannelMembership updates a client's channel subscriptions.
func (h *Hub) UpdateChannelMembership(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Remove from old channels
	for channel, clients := range h.channels {
		if !client.IsSubscribed(channel) {
			delete(clients, client)
			if len(clients) == 0 {
				delete(h.channels, channel)
			}
		}
	}

	// Add to new channels
	client.mu.RLock()
	for channel := range client.channels {
		if _, ok := h.channels[channel]; !ok {
			h.channels[channel] = make(map[*Client]bool)
		}
		h.channels[channel][client] = true
	}
	client.mu.RUnlock()
}
