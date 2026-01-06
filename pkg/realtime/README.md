# Real-time WebSocket Package

The `pkg/realtime` package provides WebSocket-based real-time communication for Inertia.js applications.

## Features

- **WebSocket Hub**: Centralized connection management
- **Channel-based Broadcasting**: Publish messages to specific channels
- **Auto-cleanup**: Automatic client disconnection handling
- **Concurrent-safe**: Thread-safe client and channel management
- **Ping/Pong**: Built-in keepalive mechanism
- **Flexible Filtering**: Broadcast to all or specific channels

## Basic Usage

### Server Setup

```go
package main

import (
    "context"
    "net/http"
    
    "github.com/toutaio/toutago-inertia/pkg/realtime"
)

func main() {
    // Create a new hub
    hub := realtime.NewHub()
    
    // Run the hub (in background)
    ctx := context.Background()
    go hub.Run(ctx)
    
    // WebSocket endpoint
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        if err := hub.HandleWebSocket(w, r); err != nil {
            log.Printf("WebSocket error: %v", err)
        }
    })
    
    // Broadcast messages
    hub.Publish("notifications", "update", map[string]interface{}{
        "message": "Hello, world!",
        "timestamp": time.Now(),
    })
    
    http.ListenAndServe(":8080", nil)
}
```

### Client Setup (Vue 3)

```vue
<script setup>
import { useLiveUpdate } from '@toutaio/inertia-vue';

const { connected, on } = useLiveUpdate('ws://localhost:8080/ws');

// Listen for messages on a channel
on('notifications', (data) => {
  console.log('Received:', data);
});
</script>

<template>
  <div>
    <span v-if="connected">🟢 Connected</span>
    <span v-else>🔴 Disconnected</span>
  </div>
</template>
```

## API Reference

### Hub

#### `NewHub() *Hub`

Creates a new WebSocket hub instance.

```go
hub := realtime.NewHub()
```

#### `Hub.Run(ctx context.Context)`

Starts the hub's message processing loop. Should be run in a goroutine.

```go
ctx := context.Background()
go hub.Run(ctx)
```

#### `Hub.Publish(channel, msgType string, data interface{})`

Publishes a message to a specific channel.

```go
hub.Publish("chat", "message", map[string]string{
    "user": "Alice",
    "text": "Hello!",
})
```

#### `Hub.Broadcast(msg *Message)`

Broadcasts a message with full control over the message structure.

```go
hub.Broadcast(&realtime.Message{
    Channel: "notifications",
    Type:    "alert",
    Data:    "System maintenance in 5 minutes",
})
```

#### `Hub.HandleWebSocket(w http.ResponseWriter, r *http.Request) error`

Upgrades an HTTP connection to WebSocket and registers the client.

```go
http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
    hub.HandleWebSocket(w, r)
})
```

### Message Structure

```go
type Message struct {
    Channel string      `json:"channel"` // Target channel or "*" for broadcast
    Type    string      `json:"type"`    // Message type
    Data    interface{} `json:"data"`    // Message payload
}
```

### Client

Clients are managed automatically by the hub. When a client connects, they can subscribe to channels by sending subscription messages:

```json
{
  "type": "subscribe",
  "channel": "chat"
}
```

To unsubscribe:

```json
{
  "type": "unsubscribe",
  "channel": "chat"
}
```

## Broadcasting Strategies

### Broadcast to Specific Channel

```go
hub.Publish("user:123", "notification", data)
```

Only clients subscribed to `user:123` receive the message.

### Broadcast to All Clients

```go
hub.Publish("*", "announcement", "Server maintenance scheduled")
```

All connected clients receive the message regardless of subscriptions.

## Configuration

### Connection Settings

Default settings can be modified in `realtime.go`:

```go
const (
    writeWait      = 10 * time.Second  // Write timeout
    pongWait       = 60 * time.Second  // Pong timeout
    pingPeriod     = 54 * time.Second  // Ping interval
    maxMessageSize = 512 * 1024        // Max message size
)
```

### Upgrader Settings

```go
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(_ *http.Request) bool {
        return true // Configure for production
    },
}
```

⚠️ **Security Note**: In production, implement proper origin checking in `CheckOrigin`.

## Error Handling

The hub handles common WebSocket errors gracefully:

- **Connection Drops**: Clients are automatically unregistered
- **Buffer Overflow**: Clients with full buffers are disconnected
- **Parse Errors**: Invalid messages are silently ignored
- **Cleanup**: Channels with no clients are removed automatically

## Performance

- **Concurrent Broadcasting**: Messages are sent to all clients concurrently
- **Buffered Channels**: 256-message buffer per client
- **Automatic Cleanup**: No memory leaks from disconnected clients
- **Ping/Pong**: Detects dead connections within ~60 seconds

## Testing

Run tests with:

```bash
go test ./pkg/realtime/... -v
```

Current test coverage: **35.2%** (12 tests passing)

## Examples

See [examples/chat](../../examples/chat) for a complete real-time chat application demonstrating:

- WebSocket hub setup
- Message broadcasting
- Vue client with `useLiveUpdate`
- Connection status handling
- Auto-reconnection

## Dependencies

- `github.com/gorilla/websocket` v1.5.3

## Thread Safety

All hub operations are thread-safe and can be called from multiple goroutines concurrently:

- `Publish()` - Safe to call from multiple handlers
- `Broadcast()` - Thread-safe message distribution
- Client registration/unregistration - Handled via channels

## Future Enhancements

Planned for future versions:

- [ ] Scéla bus integration
- [ ] Message persistence/replay
- [ ] Redis-backed distributed pub/sub
- [ ] Message filtering predicates
- [ ] Client authentication/authorization
- [ ] Rate limiting
- [ ] Metrics and monitoring hooks
