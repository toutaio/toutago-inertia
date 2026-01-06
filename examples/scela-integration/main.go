package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/toutaio/toutago-cosan-router/pkg/cosan"
	"github.com/toutaio/toutago-inertia/pkg/inertia"
	"github.com/toutaio/toutago-inertia/pkg/realtime"
	"github.com/toutaio/toutago-scela-bus/pkg/scela"
)

func main() {
	// Create Scéla message bus
	bus := scela.New()
	defer bus.Close()

	// Create Inertia instance
	irt := inertia.New()
	irt.SetVersion("1.0.0")
	irt.SetRootTemplate("app.html")

	// Create WebSocket hub for real-time updates
	hub := realtime.NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	// Create Scéla adapter with filtering
	filter := func(topic string, message interface{}) bool {
		// Only forward public messages to WebSocket
		if m, ok := message.(map[string]interface{}); ok {
			return m["public"] != false
		}
		return true
	}
	adapter := realtime.NewScelaAdapter(bus, hub, realtime.WithFilter(filter))
	defer adapter.Close()

	// Setup router
	router := cosan.New()

	// Add Inertia middleware
	router.Use(inertia.Middleware(irt))

	// WebSocket endpoint
	router.Get("/ws", func(ctx cosan.Context) error {
		return hub.HandleWebSocket(ctx.Response(), ctx.Request())
	})

	// Chat room endpoint
	router.Get("/chat", func(ctx cosan.Context) error {
		ictx := inertia.GetContext(ctx)
		return ictx.Render("Chat", inertia.Props{
			"messages": getRecentMessages(),
		})
	})

	// Post message endpoint
	router.Post("/chat/message", func(ctx cosan.Context) error {
		// Get message from form
		message := ctx.Request().FormValue("message")
		user := ctx.Request().FormValue("user")

		// Create message data
		msgData := map[string]interface{}{
			"id":        generateID(),
			"user":      user,
			"message":   message,
			"timestamp": time.Now(),
			"public":    true, // Will pass filter
		}

		// Publish to Scéla bus - will be forwarded to WebSocket clients
		err := bus.Publish(context.Background(), "chat.messages", msgData)
		if err != nil {
			return err
		}

		// Save to database
		saveMessage(msgData)

		// Return to chat page
		ictx := inertia.GetContext(ctx)
		return ictx.Back()
	})

	// User activity endpoint (internal only)
	router.Post("/user/activity", func(ctx cosan.Context) error {
		activity := map[string]interface{}{
			"user":   ctx.Request().FormValue("user"),
			"action": ctx.Request().FormValue("action"),
			"public": false, // Won't pass filter
		}

		// Publish to Scéla - won't be forwarded to WebSocket
		bus.Publish(context.Background(), "user.activity", activity)

		// But still logged internally
		logActivity(activity)

		return ctx.JSON(200, map[string]string{"status": "ok"})
	})

	// Notification endpoint with pattern matching
	router.Post("/notify/:channel", func(ctx cosan.Context) error {
		channel := ctx.Param("channel")

		notification := map[string]interface{}{
			"channel": channel,
			"title":   ctx.Request().FormValue("title"),
			"body":    ctx.Request().FormValue("body"),
			"public":  true,
		}

		// Publish to specific channel
		// Clients subscribed to "notifications.*" will receive this
		topic := fmt.Sprintf("notifications.%s", channel)
		bus.Publish(context.Background(), topic, notification)

		return ctx.JSON(200, map[string]string{"status": "sent"})
	})

	// Example: Background job publishing updates
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Publish system status update
				status := map[string]interface{}{
					"status":     "healthy",
					"users":      getActiveUserCount(),
					"memory":     getMemoryUsage(),
					"timestamp":  time.Now(),
					"public":     true,
				}
				bus.Publish(context.Background(), "system.status", status)

			case <-ctx.Done():
				return
			}
		}
	}()

	// Start server
	log.Println("Server starting on :3000")
	log.Println("WebSocket endpoint: ws://localhost:3000/ws")
	log.Println("Chat: http://localhost:3000/chat")
	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatal(err)
	}
}

// Helper functions
func getRecentMessages() []map[string]interface{} {
	// In real app, fetch from database
	return []map[string]interface{}{
		{"id": "1", "user": "Alice", "message": "Hello!", "timestamp": time.Now()},
		{"id": "2", "user": "Bob", "message": "Hi there!", "timestamp": time.Now()},
	}
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func saveMessage(msg map[string]interface{}) {
	// Save to database
	log.Printf("Saving message: %v", msg)
}

func logActivity(activity map[string]interface{}) {
	// Log internal activity
	log.Printf("Activity: %v", activity)
}

func getActiveUserCount() int {
	return 42 // Example
}

func getMemoryUsage() string {
	return "150MB" // Example
}
