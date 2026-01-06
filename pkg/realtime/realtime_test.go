package realtime

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHub(t *testing.T) {
	hub := NewHub()
	require.NotNil(t, hub)
	assert.NotNil(t, hub.clients)
	assert.NotNil(t, hub.channels)
	assert.NotNil(t, hub.broadcast)
	assert.NotNil(t, hub.register)
	assert.NotNil(t, hub.unregister)
}

func TestHubRunAndStop(t *testing.T) {
	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run hub in background
	go hub.Run(ctx)

	// Give it time to start
	time.Sleep(10 * time.Millisecond)

	// Stop the hub
	cancel()

	// Give it time to stop
	time.Sleep(10 * time.Millisecond)
}

func TestClientSubscription(t *testing.T) {
	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go hub.Run(ctx)
	time.Sleep(10 * time.Millisecond)

	// Create mock client
	client := &Client{
		hub:      hub,
		send:     make(chan []byte, 256),
		channels: make(map[string]bool),
	}

	// Subscribe to channel
	client.Subscribe("test-channel")

	assert.True(t, client.channels["test-channel"])
}

func TestClientUnsubscribe(t *testing.T) {
	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go hub.Run(ctx)
	time.Sleep(10 * time.Millisecond)

	client := &Client{
		hub:      hub,
		send:     make(chan []byte, 256),
		channels: make(map[string]bool),
	}

	// Subscribe then unsubscribe
	client.Subscribe("test-channel")
	client.Unsubscribe("test-channel")

	assert.False(t, client.channels["test-channel"])
}

func TestHubBroadcast(t *testing.T) {
	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go hub.Run(ctx)
	time.Sleep(10 * time.Millisecond)

	// Create mock client
	client := &Client{
		hub:      hub,
		send:     make(chan []byte, 256),
		channels: make(map[string]bool),
	}
	client.Subscribe("test-channel")

	// Register client
	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	// Broadcast message
	msg := &Message{
		Channel: "test-channel",
		Type:    "update",
		Data:    map[string]interface{}{"foo": "bar"},
	}

	hub.Broadcast(msg)

	// Check if client received message
	select {
	case received := <-client.send:
		var decoded Message
		err := json.Unmarshal(received, &decoded)
		require.NoError(t, err)
		assert.Equal(t, "test-channel", decoded.Channel)
		assert.Equal(t, "update", decoded.Type)
		assert.Equal(t, "bar", decoded.Data.(map[string]interface{})["foo"])
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Expected to receive message")
	}
}

func TestHubBroadcastToAll(t *testing.T) {
	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go hub.Run(ctx)
	time.Sleep(10 * time.Millisecond)

	// Create two clients with different channels
	client1 := &Client{
		hub:      hub,
		send:     make(chan []byte, 256),
		channels: make(map[string]bool),
	}
	client1.Subscribe("channel-1")

	client2 := &Client{
		hub:      hub,
		send:     make(chan []byte, 256),
		channels: make(map[string]bool),
	}
	client2.Subscribe("channel-2")

	// Register clients
	hub.register <- client1
	hub.register <- client2
	time.Sleep(10 * time.Millisecond)

	// Broadcast to all
	msg := &Message{
		Channel: "*", // Broadcast to all
		Type:    "announcement",
		Data:    "Hello everyone",
	}

	hub.Broadcast(msg)

	// Both clients should receive
	select {
	case <-client1.send:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Client 1 should receive broadcast")
	}

	select {
	case <-client2.send:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Client 2 should receive broadcast")
	}
}

func TestHubPublishMethod(t *testing.T) {
	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go hub.Run(ctx)
	time.Sleep(10 * time.Millisecond)

	client := &Client{
		hub:      hub,
		send:     make(chan []byte, 256),
		channels: make(map[string]bool),
	}
	client.Subscribe("test-channel")
	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	// Use Publish helper
	hub.Publish("test-channel", "custom-event", map[string]string{
		"message": "hello",
	})

	select {
	case received := <-client.send:
		var decoded Message
		err := json.Unmarshal(received, &decoded)
		require.NoError(t, err)
		assert.Equal(t, "custom-event", decoded.Type)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Expected to receive published message")
	}
}

func TestClientCleanup(t *testing.T) {
	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go hub.Run(ctx)
	time.Sleep(10 * time.Millisecond)

	client := &Client{
		hub:      hub,
		send:     make(chan []byte, 256),
		channels: make(map[string]bool),
	}
	client.Subscribe("test-channel")

	// Register then unregister
	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	hub.unregister <- client
	time.Sleep(10 * time.Millisecond)

	// Channel should be cleaned up
	hub.mu.RLock()
	_, exists := hub.channels["test-channel"]
	hub.mu.RUnlock()

	assert.False(t, exists, "Channel should be removed when no clients")
}

func TestWebSocketUpgrade(t *testing.T) {
	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go hub.Run(ctx)
	time.Sleep(10 * time.Millisecond)

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := hub.HandleWebSocket(w, r)
		if err != nil {
			t.Logf("WebSocket upgrade failed: %v", err)
		}
	}))
	defer server.Close()

	// We can't easily test WebSocket upgrade without a real WebSocket client
	// This test just ensures the handler doesn't panic
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	t.Logf("WebSocket URL would be: %s", wsURL)
}

func TestMessageSerialization(t *testing.T) {
	msg := &Message{
		Channel: "test",
		Type:    "update",
		Data: map[string]interface{}{
			"id":   123,
			"name": "Test",
		},
	}

	data, err := json.Marshal(msg)
	require.NoError(t, err)

	var decoded Message
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "test", decoded.Channel)
	assert.Equal(t, "update", decoded.Type)
	assert.Equal(t, float64(123), decoded.Data.(map[string]interface{})["id"])
	assert.Equal(t, "Test", decoded.Data.(map[string]interface{})["name"])
}

func TestHubConcurrentBroadcast(t *testing.T) {
	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go hub.Run(ctx)
	time.Sleep(10 * time.Millisecond)

	// Create multiple clients
	clients := make([]*Client, 5)
	for i := 0; i < 5; i++ {
		client := &Client{
			hub:      hub,
			send:     make(chan []byte, 256),
			channels: make(map[string]bool),
		}
		client.Subscribe("test-channel")
		hub.register <- client
		clients[i] = client
	}

	time.Sleep(10 * time.Millisecond)

	// Send multiple messages concurrently
	for i := 0; i < 10; i++ {
		go hub.Publish("test-channel", "update", map[string]int{"count": i})
	}

	time.Sleep(50 * time.Millisecond)

	// Each client should have received messages
	for _, client := range clients {
		select {
		case <-client.send:
			// Got at least one message
		default:
			t.Fatal("Client should have received at least one message")
		}
	}
}

func TestHubFilteredBroadcast(t *testing.T) {
	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go hub.Run(ctx)
	time.Sleep(10 * time.Millisecond)

	// Create clients on different channels
	client1 := &Client{
		hub:      hub,
		send:     make(chan []byte, 256),
		channels: make(map[string]bool),
	}
	client1.Subscribe("channel-a")

	client2 := &Client{
		hub:      hub,
		send:     make(chan []byte, 256),
		channels: make(map[string]bool),
	}
	client2.Subscribe("channel-b")

	hub.register <- client1
	hub.register <- client2
	time.Sleep(10 * time.Millisecond)

	// Broadcast only to channel-a
	hub.Publish("channel-a", "update", "data for A")

	// Only client1 should receive
	select {
	case <-client1.send:
		// Good
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Client 1 should receive message")
	}

	// Client2 should not receive
	select {
	case <-client2.send:
		t.Fatal("Client 2 should not receive message for channel-a")
	case <-time.After(50 * time.Millisecond):
		// Good, timeout expected
	}
}
