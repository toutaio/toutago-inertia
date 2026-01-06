package handlers

import (
	"github.com/toutaio/toutago-cosan-router"
	"github.com/toutaio/toutago-inertia"
)

// LoginPageProps defines props for the login page
type LoginPageProps struct {
	Flash map[string]string `json:"flash,omitempty"`
}

// HandleLoginShow shows the login page
func HandleLoginShow(adapter *inertia.Inertia) cosan.HandlerFunc {
	return func(ctx *cosan.Context) error {
		return ctx.Inertia("Auth/Login", inertia.Props{
			"flash": ctx.Session().Flash(),
		})
	}
}

// LoginInput represents login form input
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// HandleLoginSubmit handles login form submission
func HandleLoginSubmit(adapter *inertia.Inertia) cosan.HandlerFunc {
	return func(ctx *cosan.Context) error {
		var input LoginInput
		if err := ctx.BindJSON(&input); err != nil {
			return ctx.InertiaValidationErrors(map[string]string{
				"email": "Invalid input",
			})
		}

		// Validate email
		if input.Email == "" {
			return ctx.InertiaValidationErrors(map[string]string{
				"email": "Email is required",
			})
		}

		// Validate password
		if len(input.Password) < 6 {
			return ctx.InertiaValidationErrors(map[string]string{
				"password": "Password must be at least 6 characters",
			})
		}

		// Mock authentication (in real app, check against database)
		if input.Email != "demo@example.com" || input.Password != "password" {
			return ctx.InertiaValidationErrors(map[string]string{
				"email": "Invalid credentials",
			})
		}

		// Set session
		ctx.Session().Set("user_id", 1)
		ctx.Session().Flash("success", "Logged in successfully!")

		return ctx.InertiaRedirect("/")
	}
}

// HandleLogout handles logout
func HandleLogout(adapter *inertia.Inertia) cosan.HandlerFunc {
	return func(ctx *cosan.Context) error {
		ctx.Session().Delete("user_id")
		ctx.Session().Flash("success", "Logged out successfully!")
		return ctx.InertiaRedirect("/login")
	}
}
