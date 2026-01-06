# Scéla Integration Example

This example demonstrates how to integrate Scéla message bus with Inertia.js for real-time updates.

## Features

- **Message Bus Integration**: Uses Scéla to publish messages
- **WebSocket Broadcasting**: Automatically forwards messages to WebSocket clients
- **Message Filtering**: Filters which messages are sent to clients
- **Pattern Matching**: Subscribe to message patterns like `notifications.*`
- **Background Jobs**: Example of background tasks publishing updates

## Architecture

```
Backend (Go)                 Scéla Bus                WebSocket Hub              Frontend (Vue)
    │                            │                          │                           │
    ├─ POST /chat/message ──────┼─> Publish("chat.messages")                           │
    │                            │           │               │                           │
    │                            │           └──> Adapter ───┼─> Broadcast               │
    │                            │                           │        │                  │
    │                            │                           │        └──> Client 1 ─────┼─> Vue Component
    │                            │                           │        └──> Client 2 ─────┼─> Vue Component
    │                            │                           │                           │
    ├─ Background Job ───────────┼─> Publish("system.status")                           │
                                 │           │               │                           │
                                 │           └──> Adapter ───┼─> Broadcast ──────────────┼─> All Clients
```

## Running the Example

```bash
# Install dependencies
go mod download

# Run the server
go run main.go

# Open browser to http://localhost:3000/chat
```

## Frontend Setup

On the frontend, subscribe to channels:

```vue
<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const messages = ref([])
let ws = null

onMounted(() => {
  // Connect to WebSocket
  ws = new WebSocket('ws://localhost:3000/ws')
  
  ws.onopen = () => {
    // Subscribe to chat messages
    ws.send(JSON.stringify({
      type: 'subscribe',
      channel: 'chat.messages'
    }))
    
    // Subscribe to all notifications using pattern
    ws.send(JSON.stringify({
      type: 'subscribe',
      channel: 'notifications.*'
    }))
    
    // Subscribe to system status
    ws.send(JSON.stringify({
      type: 'subscribe',
      channel: 'system.status'
    }))
  }
  
  ws.onmessage = (event) => {
    const data = JSON.parse(event.data)
    messages.value.push(data)
  }
})

onUnmounted(() => {
  if (ws) ws.close()
})
</script>
```

## Key Concepts

### Message Filtering

The adapter can filter which messages are forwarded to WebSocket clients:

```go
filter := func(topic string, message interface{}) bool {
    if m, ok := message.(map[string]interface{}); ok {
        return m["public"] != false  // Only forward public messages
    }
    return true
}
adapter := realtime.NewScelaAdapter(bus, hub, realtime.WithFilter(filter))
```

### Pattern Matching

Clients can subscribe to topic patterns:

- `chat.messages` - Exact match
- `notifications.*` - All notifications
- `user.*` - All user events

### Publishing Messages

From anywhere in your application:

```go
// Publish to specific topic
bus.Publish(ctx, "chat.messages", msgData)

// Will be automatically forwarded to WebSocket clients
// subscribed to "chat.messages" or "*"
```

## Benefits

1. **Decoupled Architecture**: Backend services publish to Scéla without knowing about WebSocket
2. **Flexible Routing**: Pattern matching enables flexible subscriptions
3. **Filtering**: Control what reaches clients vs internal-only messages
4. **Scalability**: Scéla handles message distribution efficiently
5. **Testing**: Easy to test by publishing to bus

## Production Considerations

1. **Authentication**: Add authentication to WebSocket endpoint
2. **Authorization**: Filter messages based on user permissions
3. **Rate Limiting**: Prevent clients from overwhelming the server
4. **Message Persistence**: Store messages in database for history
5. **Monitoring**: Track message throughput and WebSocket connections
