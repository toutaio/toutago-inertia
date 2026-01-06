package inertia_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toutaio/toutago-inertia/pkg/inertia"
)

// MockContext simulates Cosan context for testing.
type MockContext struct {
	req    *http.Request
	res    http.ResponseWriter
	params map[string]string
	values map[string]interface{}
}

func NewMockContext(w http.ResponseWriter, r *http.Request) *MockContext {
	return &MockContext{
		req:    r,
		res:    w,
		params: make(map[string]string),
		values: make(map[string]interface{}),
	}
}

func (c *MockContext) Request() *http.Request        { return c.req }
func (c *MockContext) Response() http.ResponseWriter { return c.res }
func (c *MockContext) Set(key string, value interface{}) {
	c.values[key] = value
}
func (c *MockContext) Get(key string) interface{} {
	return c.values[key]
}

func TestInertiaContext_Render(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/users", http.NoBody)
	req.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()
	ctx := NewMockContext(w, req)

	// Create context wrapper
	ictx := inertia.NewContext(ctx, mgr)

	// Render using context
	err = ictx.Render("Users/Index", map[string]interface{}{
		"users": []string{"Alice", "Bob"},
	})
	require.NoError(t, err)

	// Verify response
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "Users/Index")
	assert.Contains(t, w.Body.String(), "Alice")
}

func TestInertiaContext_Redirect(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/users", http.NoBody)
	req.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()
	ctx := NewMockContext(w, req)

	ictx := inertia.NewContext(ctx, mgr)

	// Redirect
	err = ictx.Redirect("/users/1")
	require.NoError(t, err)

	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "/users/1", w.Header().Get("Location"))
}

func TestInertiaContext_Location(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/logout", http.NoBody)
	req.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()
	ctx := NewMockContext(w, req)

	ictx := inertia.NewContext(ctx, mgr)

	// External redirect
	err = ictx.Location("https://example.com")
	require.NoError(t, err)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, "https://example.com", w.Header().Get("X-Inertia-Location"))
}

func TestInertiaContext_Back(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/cancel", http.NoBody)
	req.Header.Set("X-Inertia", "true")
	req.Header.Set("Referer", "/users")
	w := httptest.NewRecorder()
	ctx := NewMockContext(w, req)

	ictx := inertia.NewContext(ctx, mgr)

	err = ictx.Back()
	require.NoError(t, err)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, "/users", w.Header().Get("X-Inertia-Location"))
}

func TestInertiaContext_WithErrors(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/users", http.NoBody)
	req.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()
	ctx := NewMockContext(w, req)

	ictx := inertia.NewContext(ctx, mgr)

	errors := inertia.ValidationErrors{
		"email": []string{"Email is required"},
		"name":  []string{"Name is required"},
	}

	err = ictx.WithErrors(errors).Render("Users/Create", map[string]interface{}{
		"oldInput": map[string]string{"email": ""},
	})
	require.NoError(t, err)

	assert.Contains(t, w.Body.String(), "errors")
	assert.Contains(t, w.Body.String(), "Email is required")
}

func TestInertiaContext_WithFlash(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/users", http.NoBody)
	req.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()
	ctx := NewMockContext(w, req)

	ictx := inertia.NewContext(ctx, mgr)

	flash := inertia.Flash{
		"success": "User created successfully",
	}

	err = ictx.WithFlash(flash).Render("Users/Index", map[string]interface{}{
		"users": []string{},
	})
	require.NoError(t, err)

	assert.Contains(t, w.Body.String(), "success")
	assert.Contains(t, w.Body.String(), "User created successfully")
}

func TestInertiaContext_Share(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	require.NoError(t, err)

	// Global shared data
	mgr.Share("appName", "Test App")

	req := httptest.NewRequest("GET", "/users", http.NoBody)
	req.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()
	ctx := NewMockContext(w, req)

	ictx := inertia.NewContext(ctx, mgr)

	// Context-specific shared data
	ictx.Share("user", map[string]string{"name": "Alice"})

	err = ictx.Render("Users/Index", map[string]interface{}{})
	require.NoError(t, err)

	// Should have both global and context shared data
	assert.Contains(t, w.Body.String(), "Test App")
	assert.Contains(t, w.Body.String(), "Alice")
}

func TestInertiaContext_RenderOnly(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/users", http.NoBody)
	req.Header.Set("X-Inertia", "true")
	req.Header.Set("X-Inertia-Partial-Data", "users")
	req.Header.Set("X-Inertia-Partial-Component", "Users/Index")
	w := httptest.NewRecorder()

	// Pass through middleware to set up context
	middleware := mgr.Middleware()
	var capturedReq *http.Request
	handler := middleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		capturedReq = r
	}))
	handler.ServeHTTP(w, req)

	// Now use the request with middleware-set context
	w = httptest.NewRecorder()
	ctx := NewMockContext(w, capturedReq)
	ictx := inertia.NewContext(ctx, mgr)

	err = ictx.Render("Users/Index", map[string]interface{}{
		"users":  []string{"Alice"},
		"stats":  map[string]int{"total": 1},
		"recent": []string{},
	})
	require.NoError(t, err)

	// Should only include requested prop
	assert.Contains(t, w.Body.String(), "users")
	assert.NotContains(t, w.Body.String(), "stats")
	assert.NotContains(t, w.Body.String(), "recent")
}

func TestInertiaContext_ShareFunc(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/dashboard", http.NoBody)
	req.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()

	ctx := NewMockContext(w, req)
	ic := inertia.NewContext(ctx, mgr)

	// Add lazy shared data function
	called := false
	ic.ShareFunc("currentUser", func() interface{} {
		called = true
		return map[string]interface{}{
			"name":  "John Doe",
			"email": "john@example.com",
		}
	})

	err = ic.Render("Dashboard", map[string]interface{}{
		"stats": map[string]int{"visits": 100},
	})
	require.NoError(t, err)

	// ShareFunc should have been called
	assert.True(t, called, "ShareFunc should be called during render")
	assert.Contains(t, w.Body.String(), "currentUser")
	assert.Contains(t, w.Body.String(), "John Doe")
}

func TestInertiaContext_WithInfo(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/settings", http.NoBody)
	req.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()

	ctx := NewMockContext(w, req)
	ic := inertia.NewContext(ctx, mgr)

	err = ic.WithInfo("Settings saved successfully").Render("Settings/Index", map[string]interface{}{
		"settings": map[string]string{"theme": "dark"},
	})
	require.NoError(t, err)

	assert.Contains(t, w.Body.String(), "info")
	assert.Contains(t, w.Body.String(), "Settings saved successfully")
}
