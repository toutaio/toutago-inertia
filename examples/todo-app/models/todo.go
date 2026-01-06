package models

import (
	"sync"
	"time"
)

// Todo represents a todo item
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TodosFilter represents filter options for todos
type TodosFilter struct {
	Status string `json:"status"` // all, active, completed
	Search string `json:"search"`
}

// In-memory storage (for demo purposes)
var (
	todos      = make(map[int]*Todo)
	todosmu    sync.RWMutex
	nextID     = 1
)

// InitSampleTodos initializes some sample todos
func InitSampleTodos() {
	Create(&Todo{
		Title:       "Learn Toutago",
		Description: "Explore the Toutago framework features",
		Completed:   false,
	})
	Create(&Todo{
		Title:       "Build an app",
		Description: "Create a full-stack application with Toutago and Vue",
		Completed:   false,
	})
	Create(&Todo{
		Title:       "Deploy to production",
		Description: "Deploy the application to a production server",
		Completed:   false,
	})
}

// GetAll returns all todos matching the filter
func GetAll(filter TodosFilter) []*Todo {
	todosmu.RLock()
	defer todosmu.RUnlock()

	var result []*Todo
	for _, todo := range todos {
		// Apply status filter
		if filter.Status == "active" && todo.Completed {
			continue
		}
		if filter.Status == "completed" && !todo.Completed {
			continue
		}

		// Apply search filter
		if filter.Search != "" {
			// Simple case-insensitive search
			// In production, use proper search library
			continue
		}

		result = append(result, todo)
	}

	return result
}

// GetByID returns a todo by ID
func GetByID(id int) *Todo {
	todosmu.RLock()
	defer todosmu.RUnlock()
	return todos[id]
}

// Create creates a new todo
func Create(todo *Todo) *Todo {
	todosmu.Lock()
	defer todosmu.Unlock()

	todo.ID = nextID
	nextID++
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	todos[todo.ID] = todo

	return todo
}

// Update updates an existing todo
func Update(id int, updates *Todo) *Todo {
	todosmu.Lock()
	defer todosmu.Unlock()

	todo, exists := todos[id]
	if !exists {
		return nil
	}

	if updates.Title != "" {
		todo.Title = updates.Title
	}
	if updates.Description != "" {
		todo.Description = updates.Description
	}
	todo.Completed = updates.Completed
	todo.UpdatedAt = time.Now()

	return todo
}

// Delete deletes a todo
func Delete(id int) bool {
	todosmu.Lock()
	defer todosmu.Unlock()

	_, exists := todos[id]
	if !exists {
		return false
	}

	delete(todos, id)
	return true
}
