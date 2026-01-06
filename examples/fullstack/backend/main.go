package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/toutaio/toutago-cosan-router/pkg/cosan"
	"github.com/toutaio/toutago-inertia/pkg/inertia"
)

//go:embed frontend/dist/*
var assets embed.FS

//go:embed frontend/dist/ssr/ssr.js
var ssrBundle []byte

// User represents a user in our system
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

// Post represents a blog post
type Post struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Author  User      `json:"author"`
	Created time.Time `json:"created"`
}

// In-memory data store
var (
	users = []User{
		{ID: 1, Name: "Alice Johnson", Email: "alice@example.com", CreatedAt: time.Now().AddDate(0, -6, 0)},
		{ID: 2, Name: "Bob Smith", Email: "bob@example.com", CreatedAt: time.Now().AddDate(0, -3, 0)},
		{ID: 3, Name: "Carol White", Email: "carol@example.com", CreatedAt: time.Now().AddDate(0, -1, 0)},
	}

	posts = []Post{
		{
			ID:      1,
			Title:   "Getting Started with Toutago",
			Content: "Toutago is a modern Go framework...",
			Author:  users[0],
			Created: time.Now().AddDate(0, 0, -7),
		},
		{
			ID:      2,
			Title:   "Understanding Inertia.js",
			Content: "Inertia allows you to build SPAs...",
			Author:  users[1],
			Created: time.Now().AddDate(0, 0, -3),
		},
	}

	nextUserID = 4
	nextPostID = 3
)

func main() {
	// Initialize Inertia
	inertiaManager := inertia.New(inertia.Config{
		RootTemplate: "app",
		SSR: inertia.SSRConfig{
			Enabled: true,
			Bundle:  ssrBundle,
		},
		Version: "1.0",
	})

	// Create router
	router := cosan.New()

	// Serve static assets
	staticFS := http.FS(assets)
	router.ServeFiles("/dist/*filepath", http.FileServer(staticFS))

	// Apply Inertia middleware
	router.Use(inertiaManager.Middleware())

	// Routes
	router.GET("/", handleHome(inertiaManager))
	router.GET("/users", handleUsersList(inertiaManager))
	router.GET("/users/:id", handleUserDetail(inertiaManager))
	router.GET("/posts", handlePostsList(inertiaManager))
	router.GET("/posts/create", handlePostCreate(inertiaManager))
	router.POST("/posts", handlePostStore(inertiaManager))
	router.GET("/about", handleAbout(inertiaManager))

	// Start server
	fmt.Println("Server starting on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func handleHome(i *inertia.Inertia) cosan.HandlerFunc {
	return func(ctx *cosan.Context) error {
		return i.Render(ctx.Writer, ctx.Request, "Home", inertia.Props{
			"message": "Welcome to Toutago Inertia!",
			"stats": map[string]int{
				"users": len(users),
				"posts": len(posts),
			},
		})
	}
}

func handleUsersList(i *inertia.Inertia) cosan.HandlerFunc {
	return func(ctx *cosan.Context) error {
		return i.Render(ctx.Writer, ctx.Request, "Users/Index", inertia.Props{
			"users": users,
		})
	}
}

func handleUserDetail(i *inertia.Inertia) cosan.HandlerFunc {
	return func(ctx *cosan.Context) error {
		id := ctx.Param("id")

		var user *User
		for _, u := range users {
			if fmt.Sprintf("%d", u.ID) == id {
				user = &u
				break
			}
		}

		if user == nil {
			ctx.Writer.WriteHeader(http.StatusNotFound)
			return i.Render(ctx.Writer, ctx.Request, "Error", inertia.Props{
				"status":  404,
				"message": "User not found",
			})
		}

		// Find user's posts
		var userPosts []Post
		for _, p := range posts {
			if p.Author.ID == user.ID {
				userPosts = append(userPosts, p)
			}
		}

		return i.Render(ctx.Writer, ctx.Request, "Users/Show", inertia.Props{
			"user":  user,
			"posts": userPosts,
		})
	}
}

func handlePostsList(i *inertia.Inertia) cosan.HandlerFunc {
	return func(ctx *cosan.Context) error {
		return i.Render(ctx.Writer, ctx.Request, "Posts/Index", inertia.Props{
			"posts": posts,
		})
	}
}

func handlePostCreate(i *inertia.Inertia) cosan.HandlerFunc {
	return func(ctx *cosan.Context) error {
		return i.Render(ctx.Writer, ctx.Request, "Posts/Create", inertia.Props{
			"users": users,
		})
	}
}

func handlePostStore(i *inertia.Inertia) cosan.HandlerFunc {
	return func(ctx *cosan.Context) error {
		// Parse form data
		if err := ctx.Request.ParseForm(); err != nil {
			return err
		}

		title := ctx.Request.FormValue("title")
		content := ctx.Request.FormValue("content")

		// Validation
		errors := make(map[string]string)
		if title == "" {
			errors["title"] = "Title is required"
		}
		if content == "" {
			errors["content"] = "Content is required"
		}

		if len(errors) > 0 {
			return i.Render(ctx.Writer, ctx.Request, "Posts/Create", inertia.Props{
				"users":  users,
				"errors": errors,
				"old": map[string]string{
					"title":   title,
					"content": content,
				},
			})
		}

		// Create post
		post := Post{
			ID:      nextPostID,
			Title:   title,
			Content: content,
			Author:  users[0], // Default to first user
			Created: time.Now(),
		}
		posts = append(posts, post)
		nextPostID++

		// Redirect with success message
		i.ShareFlash(ctx.Writer, "success", "Post created successfully!")
		i.Location(ctx.Writer, "/posts")
		return nil
	}
}

func handleAbout(i *inertia.Inertia) cosan.HandlerFunc {
	return func(ctx *cosan.Context) error {
		return i.Render(ctx.Writer, ctx.Request, "About", inertia.Props{
			"framework": map[string]string{
				"name":    "Toutago",
				"version": "1.0.0",
			},
		})
	}
}
