package main

import (
	"log"
	"net/http"
	"time"

	"github.com/toutaio/toutago-cosan-router"
	"github.com/toutaio/toutago-inertia"
	"github.com/toutaio/toutago-inertia/examples/todo-app/handlers"
	"github.com/toutaio/toutago-inertia/examples/todo-app/models"
)

func main() {
	// Initialize Inertia
	inertiaAdapter := inertia.New(inertia.Config{
		RootView:     "app",
		Version:      "1.0.0",
		SSREnabled:   true,
		SSRURL:       "http://localhost:13714",
		AssetURL:     "/build",
		ManifestPath: "public/build/manifest.json",
	})

	// Create router
	router := cosan.New()

	// Apply Inertia middleware
	router.Use(inertia.Middleware(inertiaAdapter))

	// Static files
	router.Static("/build", "./public/build")
	router.Static("/assets", "./public/assets")

	// Routes
	router.GET("/", handlers.HandleHome(inertiaAdapter))
	router.GET("/todos", handlers.HandleTodosList(inertiaAdapter))
	router.POST("/todos", handlers.HandleTodosCreate(inertiaAdapter))
	router.PUT("/todos/:id", handlers.HandleTodosUpdate(inertiaAdapter))
	router.DELETE("/todos/:id", handlers.HandleTodosDelete(inertiaAdapter))
	router.GET("/todos/:id/edit", handlers.HandleTodosEdit(inertiaAdapter))

	// Auth routes
	router.GET("/login", handlers.HandleLoginShow(inertiaAdapter))
	router.POST("/login", handlers.HandleLoginSubmit(inertiaAdapter))
	router.POST("/logout", handlers.HandleLogout(inertiaAdapter))

	// Admin routes (demonstrating nested layouts)
	router.GET("/admin/dashboard", handlers.AdminDashboard(inertiaAdapter))

	// Initialize sample data
	models.InitSampleTodos()

	// Start server
	srv := &http.Server{
		Addr:         ":3000",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Println("Server starting on http://localhost:3000")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
