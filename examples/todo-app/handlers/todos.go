package handlers

import (
	"github.com/toutaio/toutago-cosan-router"
	"github.com/toutaio/toutago-inertia"
	"github.com/toutaio/toutago-inertia/examples/todo-app/models"
)

// HomePageProps defines props for the home page
type HomePageProps struct {
	Greeting string `json:"greeting"`
	User     *User  `json:"user,omitempty"`
}

// User represents the authenticated user
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// HandleHome handles the home page
func HandleHome(adapter *inertia.Inertia) cosan.HandlerFunc {
	return func(ctx *cosan.Context) error {
		return ctx.Inertia("Home", inertia.Props{
			"greeting": "Welcome to Toutago + Inertia!",
			"user":     getCurrentUser(ctx),
		})
	}
}

// TodosListPageProps defines props for the todos list page
type TodosListPageProps struct {
	Todos  []*models.Todo      `json:"todos"`
	Filter models.TodosFilter  `json:"filter"`
	Flash  map[string]string   `json:"flash,omitempty"`
}

// HandleTodosList handles the todos list page
func HandleTodosList(adapter *inertia.Inertia) cosan.HandlerFunc {
	return func(ctx *cosan.Context) error {
		filter := models.TodosFilter{
			Status: ctx.Query("status", "all"),
			Search: ctx.Query("search", ""),
		}

		todos := models.GetAll(filter)

		return ctx.Inertia("Todos/Index", inertia.Props{
			"todos":  todos,
			"filter": filter,
			"flash":  ctx.Session().Flash(),
		})
	}
}

// TodosCreateInput represents input for creating a todo
type TodosCreateInput struct {
	Title       string `json:"title" validate:"required,min=3"`
	Description string `json:"description"`
}

// HandleTodosCreate handles creating a new todo
func HandleTodosCreate(adapter *inertia.Inertia) cosan.HandlerFunc {
	return func(ctx *cosan.Context) error {
		var input TodosCreateInput
		if err := ctx.BindJSON(&input); err != nil {
			return ctx.InertiaValidationErrors(map[string]string{
				"title": "Invalid input",
			})
		}

		// Validate
		if len(input.Title) < 3 {
			return ctx.InertiaValidationErrors(map[string]string{
				"title": "Title must be at least 3 characters",
			})
		}

		// Create todo
		todo := models.Create(&models.Todo{
			Title:       input.Title,
			Description: input.Description,
			Completed:   false,
		})

		if todo == nil {
			return ctx.InertiaError(500, "Failed to create todo")
		}

		ctx.Session().Flash("success", "Todo created successfully!")
		return ctx.InertiaRedirect("/todos")
	}
}

// TodosUpdateInput represents input for updating a todo
type TodosUpdateInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// HandleTodosUpdate handles updating a todo
func HandleTodosUpdate(adapter *inertia.Inertia) cosan.HandlerFunc {
	return func(ctx *cosan.Context) error {
		id := ctx.ParamInt("id")
		
		var input TodosUpdateInput
		if err := ctx.BindJSON(&input); err != nil {
			return ctx.InertiaValidationErrors(map[string]string{
				"title": "Invalid input",
			})
		}

		todo := models.Update(id, &models.Todo{
			Title:       input.Title,
			Description: input.Description,
			Completed:   input.Completed,
		})

		if todo == nil {
			return ctx.InertiaError(404, "Todo not found")
		}

		ctx.Session().Flash("success", "Todo updated successfully!")
		return ctx.InertiaRedirect("/todos")
	}
}

// HandleTodosDelete handles deleting a todo
func HandleTodosDelete(adapter *inertia.Inertia) cosan.HandlerFunc {
	return func(ctx *cosan.Context) error {
		id := ctx.ParamInt("id")

		if !models.Delete(id) {
			return ctx.InertiaError(404, "Todo not found")
		}

		ctx.Session().Flash("success", "Todo deleted successfully!")
		return ctx.InertiaRedirect("/todos")
	}
}

// TodosEditPageProps defines props for the todo edit page
type TodosEditPageProps struct {
	Todo  *models.Todo      `json:"todo"`
	Flash map[string]string `json:"flash,omitempty"`
}

// HandleTodosEdit handles the todo edit page
func HandleTodosEdit(adapter *inertia.Inertia) cosan.HandlerFunc {
	return func(ctx *cosan.Context) error {
		id := ctx.ParamInt("id")
		todo := models.GetByID(id)

		if todo == nil {
			ctx.Session().Flash("error", "Todo not found")
			return ctx.InertiaRedirect("/todos")
		}

		return ctx.Inertia("Todos/Edit", inertia.Props{
			"todo":  todo,
			"flash": ctx.Session().Flash(),
		})
	}
}

// Helper to get current user (mock implementation)
func getCurrentUser(ctx *cosan.Context) *User {
	// In real app, get from session
	userID := ctx.Session().Get("user_id")
	if userID == nil {
		return nil
	}
	
	return &User{
		ID:    userID.(int),
		Name:  "John Doe",
		Email: "john@example.com",
	}
}
