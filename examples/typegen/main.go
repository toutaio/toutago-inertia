package main

import (
	"log"
	"time"

	"github.com/toutaio/toutago-inertia/pkg/typegen"
)

// Example domain models
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Avatar    *string   `json:"avatar,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

type Post struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Author      User       `json:"author"`
	Tags        []string   `json:"tags"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

type DashboardData struct {
	User  User   `json:"user"`
	Posts []Post `json:"posts"`
	Stats struct {
		TotalPosts   int `json:"total_posts"`
		PendingPosts int `json:"pending_posts"`
		DraftPosts   int `json:"draft_posts"`
	} `json:"stats"`
}

func main() {
	// Create a new type generator
	gen := typegen.New()

	// Register all types that should be exported
	gen.Register("User", User{})
	gen.Register("Post", Post{})
	gen.Register("DashboardData", DashboardData{})

	// Generate TypeScript file
	outPath := "examples/types/generated.ts"
	if err := gen.GenerateFile(outPath); err != nil {
		log.Fatalf("Failed to generate types: %v", err)
	}

	log.Printf("Successfully generated TypeScript types at %s", outPath)
}
