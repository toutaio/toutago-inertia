// Package realtime provides WebSocket-based real-time updates for Inertia.js applications.
package realtime

import (
	"context"
	"encoding/json"
	"io"
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
	//nolint:unused // reserved for future use
	maxMessageSize = 512 * 1024 // 512 KB
)

// defaultUpgrader is the default WebSocket upgrader configuration.
// This is a package-level variable but is treated as immutable.
//
//nolint:gochecknoglobals // WebSocket upgrader is effectively a constant configuration.
var defaultUpgrader = websocket.Upgrader{
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
	defer c.cleanupConnection(ticker)

	for {
		select {
		case message, ok := <-c.send:
			if !c.handleOutgoingMessage(message, ok) {
				return
			}
		case <-ticker.C:
			if !c.sendPing() {
				return
			}
		}
	}
}

// cleanupConnection closes the ticker and connection when writePump exits.
func (c *Client) cleanupConnection(ticker *time.Ticker) {
	ticker.Stop()
	if c.conn != nil {
		c.conn.Close()
	}
}

// handleOutgoingMessage processes an outgoing message from the send channel.
func (c *Client) handleOutgoingMessage(message []byte, ok bool) bool {
	if c.conn != nil {
		_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	}

	if !ok {
		return c.sendCloseMessage()
	}

	if c.conn == nil {
		return false
	}

	return c.writeMessageWithQueued(message)
}

// sendCloseMessage sends a close message to the WebSocket.
func (c *Client) sendCloseMessage() bool {
	if c.conn != nil {
		_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	}
	return false
}

// writeMessageWithQueued writes a message and any queued messages.
func (c *Client) writeMessageWithQueued(message []byte) bool {
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return false
	}

	_, _ = w.Write(message)

	// Add queued messages to the current websocket message
	c.writeQueuedMessages(w)

	return w.Close() == nil
}

// writeQueuedMessages writes all queued messages from the send channel.
func (c *Client) writeQueuedMessages(w io.WriteCloser) {
	n := len(c.send)
	for range n {
		_, _ = w.Write([]byte{'\n'})
		_, _ = w.Write(<-c.send)
	}
}

// sendPing sends a ping message to keep the connection alive.
func (c *Client) sendPing() bool {
	if c.conn == nil {
		return true
	}

	_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(websocket.PingMessage, nil) == nil
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
			h.shutdown()
			return
		case client := <-h.register:
			h.handleRegister(client)
		case client := <-h.unregister:
			h.handleUnregister(client)
		case message := <-h.broadcast:
			h.handleBroadcast(message)
		}
	}
}

// shutdown closes all client connections.
func (h *Hub) shutdown() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		close(client.send)
	}
}

// handleRegister registers a new client and adds it to its subscribed channels.
func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true
	h.addClientToChannels(client)
}

// addClientToChannels adds a client to all its subscribed channels.
func (h *Hub) addClientToChannels(client *Client) {
	client.mu.RLock()
	defer client.mu.RUnlock()

	for channel := range client.channels {
		if _, ok := h.channels[channel]; !ok {
			h.channels[channel] = make(map[*Client]bool)
		}
		h.channels[channel][client] = true
	}
}

// handleUnregister removes a client and cleans up its channel subscriptions.
func (h *Hub) handleUnregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; !ok {
		return
	}

	delete(h.clients, client)
	close(client.send)
	h.removeClientFromAllChannels(client)
}

// removeClientFromAllChannels removes a client from all channels.
func (h *Hub) removeClientFromAllChannels(client *Client) {
	for channel := range h.channels {
		if clients, ok := h.channels[channel]; ok {
			delete(clients, client)
			if len(clients) == 0 {
				delete(h.channels, channel)
			}
		}
	}
}

// handleBroadcast processes a broadcast message.
func (h *Hub) handleBroadcast(message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	if message.Channel == "*" {
		h.broadcastToAll(data)
	} else {
		h.broadcastToChannel(message.Channel, data)
	}
}

// broadcastToAll sends a message to all connected clients.
func (h *Hub) broadcastToAll(data []byte) {
	for client := range h.clients {
		h.sendToClient(client, data)
	}
}

// broadcastToChannel sends a message to all clients in a specific channel.
func (h *Hub) broadcastToChannel(channel string, data []byte) {
	clients, ok := h.channels[channel]
	if !ok {
		return
	}

	for client := range clients {
		h.sendToClient(client, data)
	}
}

// sendToClient sends data to a client, unregistering if the buffer is full.
func (h *Hub) sendToClient(client *Client, data []byte) {
	select {
	case client.send <- data:
	default:
		// Client buffer full, close it
		go func(c *Client) {
			h.unregister <- c
		}(client)
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
	conn, err := defaultUpgrader.Upgrade(w, r, nil)
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
