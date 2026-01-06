# Real-time Chat Example

This example demonstrates WebSocket-based real-time updates using `useLiveUpdate` composable.

## Features

- Real-time message broadcasting to all connected clients
- WebSocket connection status indicator
- Auto-reconnection on disconnect
- Simple chat interface with Vue 3

## Running the Example

1. Install dependencies:
```bash
cd examples/chat
npm install
```

2. Build the frontend:
```bash
npm run build
```

3. Run the server:
```bash
go run main.go
```

4. Open http://localhost:8080 in multiple browser windows to test real-time communication

## Architecture

- **Backend**: Go HTTP server with WebSocket support
- **WebSocket Hub**: Manages client connections and message broadcasting
- **Frontend**: Vue 3 with `useLiveUpdate` composable
- **Real-time**: New messages are broadcast to all connected clients via WebSocket

## Code Highlights

### Backend (main.go)

```go
// Initialize WebSocket hub
hub := realtime.NewHub()
ctx := context.Background()
go hub.Run(ctx)

// Broadcast new messages to all clients
hub.Publish("chat", "message", msg)
```

### Frontend (Chat.vue)

```typescript
const { connected, on } = useLiveUpdate(wsUrl);

// Listen for new messages
on('chat', (message: Message) => {
  liveMessages.value.push(message);
});
```

## useLiveUpdate API

- `connected`: Reactive ref indicating WebSocket connection status
- `on(channel, handler)`: Subscribe to channel and handle messages
- `off(channel, handler)`: Unsubscribe handler from channel
- `disconnect()`: Manually close WebSocket connection

## Options

```typescript
useLiveUpdate(url, {
  reconnect: true,              // Auto-reconnect on disconnect
  reconnectDelay: 1000,         // Delay between reconnect attempts (ms)
  maxReconnectAttempts: 10,     // Maximum reconnection attempts
});
```
