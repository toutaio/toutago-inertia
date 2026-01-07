package realtime

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/toutaio/toutago-scela-bus/pkg/scela"
)

// ScelaAdapter bridges Scéla message bus to WebSocket hub.
type ScelaAdapter struct {
	bus          scela.Bus
	hub          *Hub
	subscription scela.Subscription
	filter       MessageFilter
	mu           sync.RWMutex
	closed       bool
}

// MessageFilter determines if a message should be forwarded to WebSocket
type MessageFilter func(topic string, message interface{}) bool

// ScelaOption configures the Scéla adapter
type ScelaOption func(*ScelaAdapter)

// WithFilter sets a message filter
func WithFilter(filter MessageFilter) ScelaOption {
	return func(a *ScelaAdapter) {
		a.filter = filter
	}
}

// NewScelaAdapter creates a new Scéla-to-WebSocket adapter
func NewScelaAdapter(bus scela.Bus, hub *Hub, opts ...ScelaOption) *ScelaAdapter {
	adapter := &ScelaAdapter{
		bus: bus,
		hub: hub,
	}

	for _, opt := range opts {
		opt(adapter)
	}

	// Subscribe to all topics with wildcard using HandlerFunc
	subscription, err := bus.Subscribe("*", scela.HandlerFunc(adapter.handleMessage))
	if err != nil {
		// Log error but continue - subscription might still work
		return adapter
	}
	adapter.subscription = subscription

	return adapter
}

// handleMessage is called by Scéla when a message is published
func (a *ScelaAdapter) handleMessage(ctx context.Context, msg scela.Message) error {
	a.mu.RLock()
	if a.closed {
		a.mu.RUnlock()
		return nil
	}
	a.mu.RUnlock()

	// Apply filter if set
	if a.filter != nil && !a.filter(msg.Topic(), msg.Payload()) {
		return nil
	}

	// Serialize message to JSON
	data, err := json.Marshal(msg.Payload())
	if err != nil {
		return err
	}

	// Get the topic as the channel
	channel := msg.Topic()

	// Broadcast to all clients on matching channels
	a.hub.mu.RLock()
	defer a.hub.mu.RUnlock()

	for client := range a.hub.clients {
		// Check if client is subscribed to any matching channel
		client.mu.RLock()
		matched := false
		for clientChannel := range client.channels {
			if matchesPattern(clientChannel, channel) {
				matched = true
				break
			}
		}
		client.mu.RUnlock()

		if matched {
			select {
			case client.send <- data:
			default:
				// Client buffer full, skip
			}
		}
	}

	return nil
}

// matchesPattern checks if a channel pattern matches a topic
func matchesPattern(pattern, topic string) bool {
	// Exact match
	if pattern == topic {
		return true
	}

	// Wildcard match (simple implementation)
	if pattern == "*" {
		return true
	}

	// Pattern matching (e.g., "user.*" matches "user.created")
	if len(pattern) > 2 && pattern[len(pattern)-2:] == ".*" {
		prefix := pattern[:len(pattern)-2]
		if len(topic) > len(prefix) && topic[:len(prefix)] == prefix && topic[len(prefix)] == '.' {
			return true
		}
	}

	// Prefix wildcard (e.g., "*.created" matches "user.created")
	if len(pattern) > 2 && pattern[:2] == "*." {
		suffix := pattern[2:]
		if len(topic) > len(suffix) && topic[len(topic)-len(suffix):] == suffix {
			// Check there's a dot before suffix
			dotPos := len(topic) - len(suffix) - 1
			if dotPos >= 0 && topic[dotPos] == '.' {
				return true
			}
		}
	}

	return false
}

// Close stops the adapter and unsubscribes from Scéla
func (a *ScelaAdapter) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return nil
	}

	a.closed = true

	if a.subscription != nil {
		a.subscription.Unsubscribe()
	}

	return nil
}
