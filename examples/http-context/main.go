package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/toutaio/toutago-inertia/pkg/inertia"
)

// SimpleContext implements ContextInterface for http.ResponseWriter
type SimpleContext struct {
	w      http.ResponseWriter
	r      *http.Request
	values map[string]interface{}
}

func NewSimpleContext(w http.ResponseWriter, r *http.Request) *SimpleContext {
	return &SimpleContext{
		w:      w,
		r:      r,
		values: make(map[string]interface{}),
	}
}

func (c *SimpleContext) Request() *http.Request        { return c.r }
func (c *SimpleContext) Response() http.ResponseWriter { return c.w }
func (c *SimpleContext) Set(key string, value interface{}) {
	c.values[key] = value
}
func (c *SimpleContext) Get(key string) interface{} {
	return c.values[key]
}

func main() {
	config := inertia.Config{
		RootView: "templates/app.html",
		Version:  "1.0.0",
	}

	inertiaMgr, err := inertia.New(config)
	if err != nil {
		panic(err)
	}

	inertiaMgr.Share("appName", "HTTP + Inertia")
	inertiaMgr.ShareFunc("timestamp", func() interface{} {
		return time.Now().Unix()
	})

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := NewSimpleContext(w, r)
		ictx := inertia.NewContext(ctx, inertiaMgr)

		err := ictx.Render("Home", map[string]interface{}{
			"greeting": "Welcome to Inertia!",
			"features": []string{
				"No REST API needed",
				"Server-side routing",
				"Modern SPA experience",
			},
		})
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		ctx := NewSimpleContext(w, r)
		ictx := inertia.NewContext(ctx, inertiaMgr)

		users := []map[string]interface{}{
			{"id": 1, "name": "Alice", "email": "alice@example.com"},
			{"id": 2, "name": "Bob", "email": "bob@example.com"},
		}

		err := ictx.Render("Users/Index", map[string]interface{}{
			"users": users,
			"total": len(users),
		})
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	mux.HandleFunc("/users/create", func(w http.ResponseWriter, r *http.Request) {
		ctx := NewSimpleContext(w, r)
		ictx := inertia.NewContext(ctx, inertiaMgr)

		if r.Method == "GET" {
			err := ictx.Render("Users/Create", map[string]interface{}{
				"oldInput": map[string]string{},
			})
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
			return
		}

		// POST - create user
		var input struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			ictx.Error(400, "Invalid input")
			return
		}

		errors := inertia.ValidationErrors{}
		if input.Name == "" {
			errors["name"] = []string{"Name is required"}
		}
		if input.Email == "" {
			errors["email"] = []string{"Email is required"}
		}

		if len(errors) > 0 {
			ictx.WithErrors(errors).Render("Users/Create", map[string]interface{}{
				"oldInput": input,
			})
			return
		}

		flash := inertia.Flash{"success": "User created successfully!"}
		ictx.WithFlash(flash).Redirect("/users")
	})

	// Apply Inertia middleware
	handler := inertiaMgr.Middleware()(mux)

	fmt.Println("Server: http://localhost:3000")
	http.ListenAndServe(":3000", handler)
}
