package inertia_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toutaio/toutago-inertia/pkg/inertia"
)

// TestValidationHelpers tests validation error helpers.
func TestValidationHelpers(t *testing.T) {
	t.Run("AddError adds single error", func(t *testing.T) {
		errs := make(inertia.ValidationErrors)
		errs.Add("email", "Email is required")

		assert.Len(t, errs["email"], 1)
		assert.Equal(t, "Email is required", errs["email"][0])
	})

	t.Run("AddError appends multiple errors", func(t *testing.T) {
		errs := make(inertia.ValidationErrors)
		errs.Add("email", "Email is required")
		errs.Add("email", "Email must be valid")

		assert.Len(t, errs["email"], 2)
		assert.Equal(t, "Email is required", errs["email"][0])
		assert.Equal(t, "Email must be valid", errs["email"][1])
	})

	t.Run("Has checks for error presence", func(t *testing.T) {
		errs := make(inertia.ValidationErrors)
		errs.Add("email", "Email is required")

		assert.True(t, errs.Has("email"))
		assert.False(t, errs.Has("password"))
	})

	t.Run("First returns first error", func(t *testing.T) {
		errs := make(inertia.ValidationErrors)
		errs.Add("email", "Email is required")
		errs.Add("email", "Email must be valid")

		assert.Equal(t, "Email is required", errs.First("email"))
		assert.Equal(t, "", errs.First("password"))
	})

	t.Run("Any checks if any errors exist", func(t *testing.T) {
		errs := make(inertia.ValidationErrors)
		assert.False(t, errs.Any())

		errs.Add("email", "Email is required")
		assert.True(t, errs.Any())
	})

	t.Run("NewValidationErrors creates new instance", func(t *testing.T) {
		errs := inertia.NewValidationErrors()
		assert.NotNil(t, errs)
		assert.False(t, errs.Any())
	})
}

// TestContextValidationHelpers tests context-level validation helpers.
func TestContextValidationHelpers(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	require.NoError(t, err)

	t.Run("WithError adds single error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users", http.NoBody)
		req.Header.Set("X-Inertia", "true")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		props := map[string]interface{}{
			"user": map[string]string{"name": "John"},
		}

		err := ic.WithError("email", "Email is required").Render("Users/Create", props)
		require.NoError(t, err)

		assert.Contains(t, w.Body.String(), "errors")
		assert.Contains(t, w.Body.String(), "Email is required")
	})

	t.Run("WithError chains multiple errors", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users", http.NoBody)
		req.Header.Set("X-Inertia", "true")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		props := map[string]interface{}{}

		err := ic.
			WithError("email", "Email is required").
			WithError("password", "Password is required").
			Render("Users/Create", props)
		require.NoError(t, err)

		body := w.Body.String()
		assert.Contains(t, body, "Email is required")
		assert.Contains(t, body, "Password is required")
	})
}

// TestFlashHelpers tests flash message helpers.
func TestFlashHelpers(t *testing.T) {
	t.Run("NewFlash creates flash instance", func(t *testing.T) {
		flash := inertia.NewFlash()
		assert.NotNil(t, flash)
	})

	t.Run("Success adds success message", func(t *testing.T) {
		flash := inertia.NewFlash()
		flash.Success("Operation completed")

		assert.Equal(t, "Operation completed", flash["success"])
	})

	t.Run("Error adds error message", func(t *testing.T) {
		flash := inertia.NewFlash()
		flash.Error("Operation failed")

		assert.Equal(t, "Operation failed", flash["error"])
	})

	t.Run("Warning adds warning message", func(t *testing.T) {
		flash := inertia.NewFlash()
		flash.Warning("Please be careful")

		assert.Equal(t, "Please be careful", flash["warning"])
	})

	t.Run("Info adds info message", func(t *testing.T) {
		flash := inertia.NewFlash()
		flash.Info("For your information")

		assert.Equal(t, "For your information", flash["info"])
	})

	t.Run("Custom adds custom flash message", func(t *testing.T) {
		flash := inertia.NewFlash()
		flash.Custom("notification", "Custom message")

		assert.Equal(t, "Custom message", flash["notification"])
	})
}

// TestContextFlashHelpers tests context-level flash helpers.
func TestContextFlashHelpers(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	require.NoError(t, err)

	t.Run("WithSuccess adds success flash", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users", http.NoBody)
		req.Header.Set("X-Inertia", "true")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		props := map[string]interface{}{}

		err := ic.WithSuccess("User created").Render("Users/Index", props)
		require.NoError(t, err)

		assert.Contains(t, w.Body.String(), "User created")
		assert.Contains(t, w.Body.String(), "success")
	})

	t.Run("WithErrorMessage adds error flash", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users", http.NoBody)
		req.Header.Set("X-Inertia", "true")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		props := map[string]interface{}{}

		err := ic.WithErrorMessage("Operation failed").Render("Users/Index", props)
		require.NoError(t, err)

		assert.Contains(t, w.Body.String(), "Operation failed")
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("chain flash methods", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users", http.NoBody)
		req.Header.Set("X-Inertia", "true")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		props := map[string]interface{}{}

		err := ic.
			WithSuccess("Success message").
			WithWarning("Warning message").
			Render("Users/Index", props)
		require.NoError(t, err)

		body := w.Body.String()
		assert.Contains(t, body, "Success message")
		assert.Contains(t, body, "Warning message")
	})
}
