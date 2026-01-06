package inertia_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toutaio/toutago-inertia/pkg/inertia"
)

// TestIsHTMXRequest tests HTMX request detection.
func TestIsHTMXRequest(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		expected bool
	}{
		{
			name:     "standard request",
			headers:  map[string]string{},
			expected: false,
		},
		{
			name:     "HTMX request",
			headers:  map[string]string{"HX-Request": "true"},
			expected: true,
		},
		{
			name:     "HTMX request with other headers",
			headers:  map[string]string{"HX-Request": "true", "HX-Target": "main"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", http.NoBody)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			result := inertia.IsHTMXRequest(req)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestHTMXContext tests HTMX context methods.
func TestHTMXContext(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	require.NoError(t, err)

	t.Run("HTMX redirect", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", http.NoBody)
		req.Header.Set("HX-Request", "true")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		err := ic.HTMXRedirect("/dashboard")
		require.NoError(t, err)

		assert.Equal(t, "/dashboard", w.Header().Get("HX-Redirect"))
	})

	t.Run("HTMX trigger event", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", http.NoBody)
		req.Header.Set("HX-Request", "true")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		err := ic.HTMXTrigger("itemAdded")
		require.NoError(t, err)

		assert.Equal(t, "itemAdded", w.Header().Get("HX-Trigger"))
	})

	t.Run("HTMX trigger with data", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", http.NoBody)
		req.Header.Set("HX-Request", "true")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		data := map[string]interface{}{
			"showMessage": map[string]string{"level": "info", "message": "Item added"},
		}
		err := ic.HTMXTriggerWithData(data)
		require.NoError(t, err)

		header := w.Header().Get("HX-Trigger")
		assert.NotEmpty(t, header)
		assert.Equal(t, '{', rune(header[0]))
	})

	t.Run("HTMX partial render", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", http.NoBody)
		req.Header.Set("HX-Request", "true")
		req.Header.Set("HX-Target", "content")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		html := "<div>Partial content</div>"
		err := ic.HTMXPartial(html)
		require.NoError(t, err)

		assert.Equal(t, html, w.Body.String())
		assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
	})

	t.Run("HTMX reswap", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", http.NoBody)
		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		ic.HTMXReswap("outerHTML")

		assert.Equal(t, "outerHTML", w.Header().Get("HX-Reswap"))
	})

	t.Run("HTMX retarget", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", http.NoBody)
		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		ic.HTMXRetarget("#newTarget")

		assert.Equal(t, "#newTarget", w.Header().Get("HX-Retarget"))
	})

	t.Run("HTMX push URL", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", http.NoBody)
		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		ic.HTMXPushURL("/new-page")

		assert.Equal(t, "/new-page", w.Header().Get("HX-Push-Url"))
	})

	t.Run("HTMX refresh", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", http.NoBody)
		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		err := ic.HTMXRefresh()
		require.NoError(t, err)

		assert.Equal(t, "true", w.Header().Get("HX-Refresh"))
	})
}

// TestGetHTMXHeaders tests extracting HTMX headers from request.
func TestGetHTMXHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/", http.NoBody)
	req.Header.Set("HX-Request", "true")
	req.Header.Set("HX-Target", "main-content")
	req.Header.Set("HX-Trigger", "btn-click")
	req.Header.Set("HX-Trigger-Name", "submitBtn")
	req.Header.Set("HX-Current-URL", "https://example.com/page")

	headers := inertia.GetHTMXHeaders(req)

	assert.True(t, headers.Request)
	assert.Equal(t, "main-content", headers.Target)
	assert.Equal(t, "btn-click", headers.Trigger)
	assert.Equal(t, "submitBtn", headers.TriggerName)
	assert.Equal(t, "https://example.com/page", headers.CurrentURL)
}

// TestHTMXIntegration tests full HTMX integration scenarios.
func TestHTMXIntegration(t *testing.T) {
	config := inertia.Config{
		Version:  "1.0",
		RootView: "app",
	}
	mgr, err := inertia.New(config)
	require.NoError(t, err)

	t.Run("partial update with HTMX", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/update", http.NoBody)
		req.Header.Set("HX-Request", "true")
		req.Header.Set("HX-Target", "user-list")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		// Render partial HTML fragment
		html := `<div id="user-list"><div class="user">John Doe</div></div>`
		err := ic.HTMXPartial(html)
		require.NoError(t, err)

		assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
		assert.Contains(t, w.Body.String(), "John Doe")
	})

	t.Run("HTMX redirect with flash", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/save", http.NoBody)
		req.Header.Set("HX-Request", "true")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		// Save and redirect
		err := ic.HTMXRedirect("/success")
		require.NoError(t, err)

		assert.Equal(t, "/success", w.Header().Get("HX-Redirect"))
	})

	t.Run("trigger client-side event", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/action", http.NoBody)
		req.Header.Set("HX-Request", "true")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		// Trigger event with data
		err := ic.HTMXTriggerWithData(map[string]interface{}{
			"showNotification": map[string]string{
				"level":   "success",
				"message": "Operation completed",
			},
		})
		require.NoError(t, err)

		err = ic.HTMXPartial("<div>Success</div>")
		require.NoError(t, err)

		trigger := w.Header().Get("HX-Trigger")
		assert.Contains(t, trigger, "showNotification")
		assert.Contains(t, trigger, "Operation completed")
	})

	t.Run("out-of-band swap", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/update", http.NoBody)
		req.Header.Set("HX-Request", "true")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		// Update multiple targets
		html := `
			<div id="main-content">Main Updated</div>
			<div id="sidebar" hx-swap-oob="true">Sidebar Updated</div>
		`
		err := ic.HTMXPartial(html)
		require.NoError(t, err)

		assert.Contains(t, w.Body.String(), "Main Updated")
		assert.Contains(t, w.Body.String(), "hx-swap-oob")
	})

	t.Run("chained HTMX operations", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/complex", http.NoBody)
		req.Header.Set("HX-Request", "true")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		// Chain multiple HTMX operations
		err := ic.
			HTMXReswap("outerHTML").
			HTMXRetarget("#result").
			HTMXPushURL("/updated").
			HTMXPartial("<div>Complete</div>")
		require.NoError(t, err)

		// Trigger must be set separately (returns error)
		err = ic.HTMXTrigger("updated")
		require.NoError(t, err)

		assert.Equal(t, "outerHTML", w.Header().Get("HX-Reswap"))
		assert.Equal(t, "#result", w.Header().Get("HX-Retarget"))
		assert.Equal(t, "/updated", w.Header().Get("HX-Push-Url"))
		assert.Equal(t, "updated", w.Header().Get("HX-Trigger"))
	})

	t.Run("hybrid Inertia and HTMX routing", func(t *testing.T) {
		// Regular Inertia request
		inertiaReq := httptest.NewRequest("GET", "/dashboard", http.NoBody)
		inertiaReq.Header.Set("X-Inertia", "true")

		assert.False(t, inertia.IsHTMXRequest(inertiaReq))
		assert.True(t, inertiaReq.Header.Get("X-Inertia") == "true")

		// HTMX request
		htmxReq := httptest.NewRequest("GET", "/partial", http.NoBody)
		htmxReq.Header.Set("HX-Request", "true")

		assert.True(t, inertia.IsHTMXRequest(htmxReq))
		assert.False(t, htmxReq.Header.Get("X-Inertia") == "true")

		// Regular browser request
		browserReq := httptest.NewRequest("GET", "/page", http.NoBody)

		assert.False(t, inertia.IsHTMXRequest(browserReq))
		assert.False(t, browserReq.Header.Get("X-Inertia") == "true")
	})

	t.Run("HTMX with validation errors", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/validate", http.NoBody)
		req.Header.Set("HX-Request", "true")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		// Return validation errors as HTML
		errors := inertia.NewValidationErrors()
		errors.Add("email", "Email is required")
		errors.Add("password", "Password too short")

		html := `<div class="errors">
			<div class="error">Email is required</div>
			<div class="error">Password too short</div>
		</div>`

		err := ic.HTMXPartial(html)
		require.NoError(t, err)

		assert.Contains(t, w.Body.String(), "Email is required")
		assert.Contains(t, w.Body.String(), "Password too short")
	})
}

func TestHTMXReplaceURL(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/profile/update", http.NoBody)
	req.Header.Set("HX-Request", "true")
	w := httptest.NewRecorder()

	ctx := NewMockContext(w, req)
	ic := inertia.NewContext(ctx, mgr)

	err = ic.HTMXReplaceURL("/profile").HTMXPartial("<div>Profile updated</div>")
	require.NoError(t, err)

	assert.Equal(t, "/profile", w.Header().Get("HX-Replace-Url"))
	assert.Contains(t, w.Body.String(), "Profile updated")
}

func TestAlwaysAndAlwaysLazy(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	require.NoError(t, err)

	t.Run("Always prop can be set", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/dashboard", http.NoBody)
		req.Header.Set("X-Inertia", "true")
		w := httptest.NewRecorder()

		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		// Test that Always can be called and chained
		result := ic.Always("alwaysProp", "always here")
		require.NotNil(t, result, "Always should return InertiaContext for chaining")

		err = result.Render("Dashboard", map[string]interface{}{
			"data": "main data",
		})
		require.NoError(t, err)
	})

	t.Run("AlwaysLazy prop can be set", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/dashboard", http.NoBody)
		req.Header.Set("X-Inertia", "true")
		w := httptest.NewRecorder()

		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		called := false
		result := ic.AlwaysLazy("authUser", func() interface{} {
			called = true
			return map[string]string{"name": "Admin"}
		})
		require.NotNil(t, result, "AlwaysLazy should return InertiaContext for chaining")

		err = result.Render("Dashboard", map[string]interface{}{
			"data": "main data",
		})
		require.NoError(t, err)

		assert.True(t, called, "AlwaysLazy should be called")
		assert.Contains(t, w.Body.String(), "authUser")
		assert.Contains(t, w.Body.String(), "Admin")
	})
}
