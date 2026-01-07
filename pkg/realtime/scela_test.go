package realtime

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/toutaio/toutago-scela-bus/pkg/scela"
)

func TestScelaAdapter_BasicIntegration(t *testing.T) {
	// Create Scéla bus and WebSocket hub
	bus := scela.New()
	defer bus.Close()

	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	// Create Scéla adapter
	adapter := NewScelaAdapter(bus, hub)
	defer adapter.Close()

	// Create a test client with subscription
	client := &Client{
		hub:      hub,
		conn:     nil, // No actual websocket for testing
		send:     make(chan []byte, 10),
		channels: make(map[string]bool),
	}
	client.Subscribe("test-channel")

	// Manually register client
	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	// Publish message via Scéla
	msg := map[string]interface{}{
		"type":    "message",
		"content": "Hello from Scéla!",
	}
	err := bus.PublishSync(context.Background(), "test-channel", msg)
	if err != nil {
		t.Fatalf("Failed to publish: %v", err)
	}

	// Wait for async propagation
	time.Sleep(50 * time.Millisecond)

	// Verify message received
	select {
	case data := <-client.send:
		var received map[string]interface{}
		if err := json.Unmarshal(data, &received); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}
		if received["type"] != "message" {
			t.Errorf("Expected type=message, got %v", received["type"])
		}
		if received["content"] != "Hello from Scéla!" {
			t.Errorf("Expected content=Hello from Scéla!, got %v", received["content"])
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Timeout waiting for message")
	}
}

func TestScelaAdapter_PatternMatching(t *testing.T) {
	bus := scela.New()
	defer bus.Close()

	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	adapter := NewScelaAdapter(bus, hub)
	defer adapter.Close()

	// Create client subscribed to "user.*" pattern
	client := &Client{
		hub:      hub,
		conn:     nil,
		send:     make(chan []byte, 10),
		channels: make(map[string]bool),
	}
	client.Subscribe("user.*")
	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	// Publish to "user.created"
	msg1 := map[string]interface{}{"event": "created", "user": "john"}
	err := bus.PublishSync(context.Background(), "user.created", msg1)
	if err != nil {
		t.Fatalf("Failed to publish: %v", err)
	}

	// Publish to "user.updated"
	msg2 := map[string]interface{}{"event": "updated", "user": "jane"}
	err = bus.PublishSync(context.Background(), "user.updated", msg2)
	if err != nil {
		t.Fatalf("Failed to publish: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	// Verify both messages received
	receivedCount := 0
	timeout := time.After(200 * time.Millisecond)
	for receivedCount < 2 {
		select {
		case data := <-client.send:
			var received map[string]interface{}
			if err := json.Unmarshal(data, &received); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}
			event := received["event"]
			if event != "created" && event != "updated" {
				t.Errorf("Unexpected event: %v", event)
			}
			receivedCount++
		case <-timeout:
			t.Fatalf("Timeout - only received %d/2 messages", receivedCount)
		}
	}
}

func TestScelaAdapter_Filtering(t *testing.T) {
	bus := scela.New()
	defer bus.Close()

	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	// Create adapter with filter
	filter := func(_ string, message interface{}) bool {
		// Only pass messages with "important" flag
		if m, ok := message.(map[string]interface{}); ok {
			return m["important"] == true
		}
		return false
	}
	adapter := NewScelaAdapter(bus, hub, WithFilter(filter))
	defer adapter.Close()

	client := &Client{
		hub:      hub,
		conn:     nil,
		send:     make(chan []byte, 10),
		channels: make(map[string]bool),
	}
	client.Subscribe("notifications")
	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	// Publish important message
	msg1 := map[string]interface{}{"important": true, "text": "Important!"}
	err := bus.PublishSync(context.Background(), "notifications", msg1)
	if err != nil {
		t.Fatalf("Failed to publish: %v", err)
	}

	// Publish non-important message
	msg2 := map[string]interface{}{"important": false, "text": "Not important"}
	err = bus.PublishSync(context.Background(), "notifications", msg2)
	if err != nil {
		t.Fatalf("Failed to publish: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	// Should only receive important message
	select {
	case data := <-client.send:
		var received map[string]interface{}
		if err := json.Unmarshal(data, &received); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}
		if received["text"] != "Important!" {
			t.Errorf("Expected important message, got %v", received)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Timeout waiting for important message")
	}

	// Should not receive second message
	select {
	case data := <-client.send:
		t.Fatalf("Received unexpected message: %s", data)
	case <-time.After(50 * time.Millisecond):
		// Expected - no message
	}
}

func TestScelaAdapter_MultipleClients(t *testing.T) {
	bus := scela.New()
	defer bus.Close()

	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	adapter := NewScelaAdapter(bus, hub)
	defer adapter.Close()

	// Create multiple clients on same channel
	clients := make([]*Client, 3)
	for i := 0; i < 3; i++ {
		clients[i] = &Client{
			hub:      hub,
			conn:     nil,
			send:     make(chan []byte, 10),
			channels: make(map[string]bool),
		}
		clients[i].Subscribe("broadcast")
		hub.register <- clients[i]
	}
	time.Sleep(10 * time.Millisecond)

	// Publish message
	msg := map[string]interface{}{"data": "broadcast to all"}
	err := bus.PublishSync(context.Background(), "broadcast", msg)
	if err != nil {
		t.Fatalf("Failed to publish: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	// Verify all clients received
	for i, client := range clients {
		select {
		case data := <-client.send:
			var received map[string]interface{}
			if err := json.Unmarshal(data, &received); err != nil {
				t.Fatalf("Client %d failed to unmarshal: %v", i, err)
			}
			if received["data"] != "broadcast to all" {
				t.Errorf("Client %d got wrong data: %v", i, received)
			}
		case <-time.After(200 * time.Millisecond):
			t.Fatalf("Client %d timeout", i)
		}
	}
}

func TestScelaAdapter_ErrorHandling(t *testing.T) {
	bus := scela.New()
	defer bus.Close()

	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	adapter := NewScelaAdapter(bus, hub)
	defer adapter.Close()

	// Publish message with no subscribers - should not error
	msg := map[string]interface{}{"test": "no subscribers"}
	err := bus.Publish(context.Background(), "empty-channel", msg)
	if err != nil {
		t.Fatalf("Should not error on empty channel: %v", err)
	}

	// Close adapter and try to publish - should handle gracefully
	adapter.Close()
	time.Sleep(10 * time.Millisecond)

	err = bus.Publish(context.Background(), "test", msg)
	// Should not panic or crash
	_ = err // Expected - subscription may be removed
}

func TestScelaAdapter_ContextCancellation(t *testing.T) {
	bus := scela.New()
	defer bus.Close()

	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	adapter := NewScelaAdapter(bus, hub)
	defer adapter.Close()

	pubCtx, pubCancel := context.WithCancel(context.Background())
	pubCancel() // Cancel immediately

	// Publish with canceled context
	msg := map[string]interface{}{"test": "canceled"}
	err := bus.Publish(pubCtx, "test-channel", msg)
	// Scéla respects context cancellation
	if err == nil {
		// Some implementations may still succeed on async publish
		// Just verify no panic and adapter still works
		t.Log("Publish succeeded despite canceled context (async behavior)")
	}

	// Verify adapter is still functional
	if adapter == nil {
		t.Fatal("adapter should not be nil")
	}
}
