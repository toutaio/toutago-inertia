package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/toutaio/toutago-inertia/pkg/inertia"
	"github.com/toutaio/toutago-inertia/pkg/realtime"
)

type Message struct {
	ID        int       `json:"id"`
	User      string    `json:"user"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	messages   []Message
	messagesMu sync.RWMutex
	messageID  = 0
)

func main() {
	// Initialize Inertia
	config := inertia.Config{
		RootView: "templates/app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	if err != nil {
		panic(err)
	}

	mgr.Share("appName", "Chat Example")

	// Initialize WebSocket hub
	hub := realtime.NewHub()
	ctx := context.Background()
	go hub.Run(ctx)

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		messagesMu.RLock()
		msgs := messages
		messagesMu.RUnlock()

		page, _ := mgr.Render("Chat", map[string]interface{}{
			"messages": msgs,
		}, r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(page)
	})

	mux.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			User string `json:"user"`
			Text string `json:"text"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		messagesMu.Lock()
		messageID++
		msg := Message{
			ID:        messageID,
			User:      input.User,
			Text:      input.Text,
			Timestamp: time.Now(),
		}
		messages = append(messages, msg)
		messagesMu.Unlock()

		// Broadcast to all connected clients
		hub.Publish("chat", "message", msg)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(msg)
	})

	// WebSocket endpoint
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if err := hub.HandleWebSocket(w, r); err != nil {
			log.Printf("WebSocket error: %v", err)
		}
	})

	// Serve static files
	mux.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("dist"))))

	handler := mgr.Middleware()(mux)

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
