package inertia_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toutaio/toutago-inertia/pkg/inertia"
)

func TestLocation_ExternalRedirect(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()

	// Use Location method
	err = i.Location(w, req, "https://external.com")
	require.NoError(t, err)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, "https://external.com", w.Header().Get("X-Inertia-Location"))
}

func TestLocation_NonInertiaRequest(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	// Regular browser request
	req := httptest.NewRequest("GET", "/test", http.NoBody)
	w := httptest.NewRecorder()

	err = i.Location(w, req, "https://external.com")
	require.NoError(t, err)

	// Should return normal 302 redirect
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "https://external.com", w.Header().Get("Location"))
}

func TestBack_InertiaRequest(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("X-Inertia", "true")
	req.Header.Set("Referer", "/previous-page")
	w := httptest.NewRecorder()

	err = i.Back(w, req)
	require.NoError(t, err)

	// Should redirect to referer
	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, "/previous-page", w.Header().Get("X-Inertia-Location"))
}

func TestBack_NoReferer(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("X-Inertia", "true")
	// No referer header
	w := httptest.NewRecorder()

	err = i.Back(w, req)
	require.NoError(t, err)

	// Should redirect to root
	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, "/", w.Header().Get("X-Inertia-Location"))
}

func TestRedirect_InertiaRequest_GET(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()

	err = i.Redirect(w, req, "/dashboard")
	require.NoError(t, err)

	// For GET requests, use 303 See Other
	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "/dashboard", w.Header().Get("Location"))
}

func TestRedirect_InertiaRequest_POST(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	// POST request
	req := httptest.NewRequest("POST", "/test", http.NoBody)
	req.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()

	err = i.Redirect(w, req, "/dashboard")
	require.NoError(t, err)

	// For POST requests, use 303 See Other to change to GET
	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "/dashboard", w.Header().Get("Location"))
}

func TestRedirect_InertiaRequest_PUT(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	// PUT request
	req := httptest.NewRequest("PUT", "/test", http.NoBody)
	req.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()

	err = i.Redirect(w, req, "/dashboard")
	require.NoError(t, err)

	// For PUT/PATCH/DELETE, use 303 See Other
	assert.Equal(t, http.StatusSeeOther, w.Code)
}

func TestRedirect_NonInertiaRequest(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	// No X-Inertia header
	w := httptest.NewRecorder()

	err = i.Redirect(w, req, "/dashboard")
	require.NoError(t, err)

	// Regular browser redirect
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/dashboard", w.Header().Get("Location"))
}

func TestError_Response(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	req.Header.Set("X-Inertia", "true")

	// Create error response
	errorPage, err := i.Error(404, "Page not found", "/404", req)
	require.NoError(t, err)

	assert.Equal(t, "Error", errorPage.Component)
	assert.Equal(t, 404, errorPage.Props["status"])
	assert.Equal(t, "Page not found", errorPage.Props["message"])
}

func TestValidationErrors(t *testing.T) {
	errors := inertia.ValidationErrors{
		"email":    []string{"Email is required", "Email must be valid"},
		"password": []string{"Password is too short"},
	}

	assert.Equal(t, 2, len(errors))
	assert.Contains(t, errors, "email")
	assert.Contains(t, errors, "password")
	assert.Equal(t, 2, len(errors["email"]))
	assert.Equal(t, 1, len(errors["password"]))
}

func TestFlashData(t *testing.T) {
	flash := inertia.Flash{
		"success": "User created successfully",
		"error":   "Failed to save",
		"info":    "Please check your email",
	}

	assert.Equal(t, 3, len(flash))
	assert.Equal(t, "User created successfully", flash["success"])
}
