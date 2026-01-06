package inertia_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toutaio/toutago-inertia/pkg/inertia"
)

// TestFullRequestCycle tests complete request/response flows.
func TestFullRequestCycle(t *testing.T) {
	config := inertia.Config{
		Version:  "1.0",
		RootView: "app",
	}
	mgr, err := inertia.New(config)
	require.NoError(t, err)

	t.Run("initial page load", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/dashboard", http.NoBody)
		w := httptest.NewRecorder()

		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		err := ic.Render("Dashboard/Index", map[string]interface{}{
			"title": "Dashboard",
			"user": map[string]interface{}{
				"name":  "John Doe",
				"email": "john@example.com",
			},
		})
		require.NoError(t, err)

		// Initial load returns JSON (actual implementation)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
		assert.Contains(t, w.Body.String(), "Dashboard/Index")
		assert.Contains(t, w.Body.String(), "John Doe")
	})

	t.Run("Inertia navigation", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users", http.NoBody)
		req.Header.Set("X-Inertia", "true")
		req.Header.Set("X-Inertia-Version", "1.0")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		err := ic.Render("Users/Index", map[string]interface{}{
			"users": []map[string]interface{}{
				{"name": "Alice", "email": "alice@example.com"},
				{"name": "Bob", "email": "bob@example.com"},
			},
		})
		require.NoError(t, err)

		// Should return JSON
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
		assert.Contains(t, w.Body.String(), "Users/Index")
		assert.Contains(t, w.Body.String(), "Alice")
		assert.Contains(t, w.Body.String(), "Bob")
	})

	t.Run("form submission with validation errors", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/users", http.NoBody)
		req.Header.Set("X-Inertia", "true")
		req.Header.Set("X-Inertia-Version", "1.0")
		req.Header.Set("Referer", "http://example.com/users/new")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		// Validate and return errors
		err := ic.
			WithError("email", "Email is required").
			WithError("password", "Password must be at least 8 characters").
			Back()
		require.NoError(t, err)

		// Back() uses Location() which returns 409 for Inertia requests
		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Equal(t, "http://example.com/users/new", w.Header().Get("X-Inertia-Location"))
	})

	t.Run("successful form submission with flash", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/users", http.NoBody)
		req.Header.Set("X-Inertia", "true")
		req.Header.Set("X-Inertia-Version", "1.0")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		// Create user and redirect
		err := ic.
			WithSuccess("User created successfully").
			Redirect("/users")
		require.NoError(t, err)

		assert.Equal(t, http.StatusSeeOther, w.Code)
		assert.Equal(t, "/users", w.Header().Get("Location"))
	})

	t.Run("lazy props evaluation", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/dashboard", http.NoBody)
		req.Header.Set("X-Inertia", "true")
		req.Header.Set("X-Inertia-Version", "1.0")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		statsEvaluated := false
		analyticsEvaluated := false
		deferredEvaluated := false

		err := ic.
			Lazy("stats", func() interface{} {
				statsEvaluated = true
				return map[string]int{"users": 100}
			}).
			Lazy("analytics", func() interface{} {
				analyticsEvaluated = true
				return map[string]int{"views": 5000}
			}).
			Defer("history", func() interface{} {
				deferredEvaluated = true
				return []string{"action1", "action2"}
			}).
			Render("Dashboard/Index", map[string]interface{}{
				"title": "Dashboard",
			})
		require.NoError(t, err)

		// On full page load: lazy props evaluated, deferred not evaluated
		assert.True(t, statsEvaluated, "lazy props should be evaluated on full load")
		assert.True(t, analyticsEvaluated, "lazy props should be evaluated on full load")
		assert.False(t, deferredEvaluated, "deferred props should not be evaluated")

		// Response should include lazy but not deferred
		assert.Contains(t, w.Body.String(), "stats")
		assert.Contains(t, w.Body.String(), "analytics")
		assert.NotContains(t, w.Body.String(), "history")
	})

	t.Run("external redirect", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/login", http.NoBody)
		req.Header.Set("X-Inertia", "true")
		req.Header.Set("X-Inertia-Version", "1.0")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		// External redirect (e.g., to OAuth provider)
		err := ic.Location("https://oauth.example.com/authorize")
		require.NoError(t, err)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Equal(t, "https://oauth.example.com/authorize", w.Header().Get("X-Inertia-Location"))
	})
}

// TestSharedDataFlow tests shared data across requests.
func TestSharedDataFlow(t *testing.T) {
	config := inertia.Config{
		Version:  "1.0",
		RootView: "app",
	}
	mgr, err := inertia.New(config)
	require.NoError(t, err)

	// Add global shared data
	mgr.Share("appName", "Toutā App")
	mgr.ShareFunc("currentYear", func() interface{} {
		return 2026
	})

	t.Run("shared data included in all responses", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/page", http.NoBody)
		req.Header.Set("X-Inertia", "true")
		req.Header.Set("X-Inertia-Version", "1.0")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		err := ic.Render("Test/Page", map[string]interface{}{
			"title": "Test",
		})
		require.NoError(t, err)

		// Global shared data should be in response
		assert.Contains(t, w.Body.String(), "appName")
		assert.Contains(t, w.Body.String(), "Toutā App")
		assert.Contains(t, w.Body.String(), "currentYear")
		assert.Contains(t, w.Body.String(), "2026")
	})

	t.Run("request-level shared data", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/dashboard", http.NoBody)
		req.Header.Set("X-Inertia", "true")
		req.Header.Set("X-Inertia-Version", "1.0")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		err := ic.
			Share("user", map[string]interface{}{
				"name": "John Doe",
			}).
			Render("Dashboard/Index", map[string]interface{}{
				"title": "Dashboard",
			})
		require.NoError(t, err)

		// Both global and request-level shared data
		assert.Contains(t, w.Body.String(), "appName")
		assert.Contains(t, w.Body.String(), "user")
		assert.Contains(t, w.Body.String(), "John Doe")
	})
}

// TestErrorHandling tests error scenarios.
func TestErrorHandling(t *testing.T) {
	config := inertia.Config{
		Version:  "1.0",
		RootView: "app",
	}
	mgr, err := inertia.New(config)
	require.NoError(t, err)

	t.Run("render 404 error page", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/not-found", http.NoBody)
		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		err := ic.Error(http.StatusNotFound, "Page not found")
		require.NoError(t, err)

		// Error component is just "Error", not "Error/404"
		assert.Contains(t, w.Body.String(), `"component":"Error"`)
		assert.Contains(t, w.Body.String(), "Page not found")
		assert.Contains(t, w.Body.String(), `"status":404`)
	})

	t.Run("render 500 error page", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/error", http.NoBody)
		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		err := ic.Error(http.StatusInternalServerError, "Something went wrong")
		require.NoError(t, err)

		assert.Contains(t, w.Body.String(), `"component":"Error"`)
		assert.Contains(t, w.Body.String(), "Something went wrong")
		assert.Contains(t, w.Body.String(), `"status":500`)
	})
}
