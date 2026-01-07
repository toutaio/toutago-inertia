package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/toutaio/toutago-inertia/pkg/typegen"
)

// Example structs to generate TypeScript types for
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Post struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	AuthorID int    `json:"authorId"`
}

func main() {
	// Create generator
	gen := typegen.New()
	gen.Register("User", User{})
	gen.Register("Post", Post{})

	outputPath := "frontend/types.ts"

	// Create watcher
	watcher := typegen.NewWatcher()

	// Add files to watch (or use AddDirectory)
	if err := watcher.AddFile("models/user.go"); err != nil {
		log.Printf("Warning: %v", err)
	}
	if err := watcher.AddFile("models/post.go"); err != nil {
		log.Printf("Warning: %v", err)
	}

	// Set output path
	watcher.SetOutput(outputPath)

	// Set generator function
	watcher.SetGenerator(func() error {
		fmt.Println("Regenerating TypeScript types...")
		if err := gen.GenerateFile(outputPath); err != nil {
			return err
		}
		fmt.Printf("✓ Generated %s\n", outputPath)
		return nil
	})

	// Set error handler
	watcher.SetErrorHandler(func(err error) {
		log.Printf("Error: %v", err)
	})

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\nStopping watcher...")
		watcher.Stop()
	}()

	// Start watching
	fmt.Println("Watching for changes... (Press Ctrl+C to stop)")
	if err := watcher.Watch(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Watcher stopped")
}
