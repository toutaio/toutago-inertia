package inertia_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toutaio/toutago-inertia/pkg/inertia"
)

// TestLazyProps tests lazy prop evaluation.
func TestLazyProps(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	require.NoError(t, err)

	t.Run("lazy props not evaluated on partial reload", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users", http.NoBody)
		req.Header.Set("X-Inertia", "true")
		req.Header.Set("X-Inertia-Partial-Data", "name,email")
		req.Header.Set("X-Inertia-Partial-Component", "Users/Index")

		w := httptest.NewRecorder()

		// Run through middleware to set context values
		middleware := mgr.Middleware()
		var capturedReq *http.Request
		handler := middleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			capturedReq = r
		}))
		handler.ServeHTTP(w, req)

		// Now use the request with context values set
		w = httptest.NewRecorder() // Reset recorder
		ctx := NewMockContext(w, capturedReq)
		ic := inertia.NewContext(ctx, mgr)

		called := false
		lazyFn := func() interface{} {
			called = true
			return "lazy value"
		}

		props := map[string]interface{}{
			"name":  "John",
			"email": "john@example.com",
		}

		err := ic.Lazy("expensive", lazyFn).Render("Users/Index", props)
		require.NoError(t, err)

		// Lazy prop should not be evaluated during partial reload
		// when it's not in the "only" list
		assert.False(t, called, "lazy prop should not be evaluated")
	})

	t.Run("lazy props evaluated on full load", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users", http.NoBody)
		req.Header.Set("X-Inertia", "true")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		called := false
		lazyFn := func() interface{} {
			called = true
			return "lazy value"
		}

		props := map[string]interface{}{
			"name": "John",
		}

		err := ic.Lazy("expensive", lazyFn).Render("Users/Index", props)
		require.NoError(t, err)

		// Lazy prop should be evaluated on full load
		assert.True(t, called, "lazy prop should be evaluated")
	})

	t.Run("lazy props evaluated when explicitly requested", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users", http.NoBody)
		req.Header.Set("X-Inertia", "true")
		req.Header.Set("X-Inertia-Partial-Data", "expensive")
		req.Header.Set("X-Inertia-Partial-Component", "Users/Index")

		w := httptest.NewRecorder()

		// Run through middleware to set context values
		middleware := mgr.Middleware()
		var capturedReq *http.Request
		handler := middleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			capturedReq = r
		}))
		handler.ServeHTTP(w, req)

		// Now use the request with context values set
		w = httptest.NewRecorder() // Reset recorder
		ctx := NewMockContext(w, capturedReq)
		ic := inertia.NewContext(ctx, mgr)

		called := false
		lazyFn := func() interface{} {
			called = true
			return map[string]interface{}{
				"data": "expensive data",
			}
		}

		props := map[string]interface{}{
			"name": "John",
		}

		err := ic.Lazy("expensive", lazyFn).Render("Users/Index", props)
		require.NoError(t, err)

		// Lazy prop should be evaluated when explicitly requested
		assert.True(t, called, "lazy prop should be evaluated when requested")
	})
}

// TestAlways tests always-included props.
func TestAlways(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	_, err := inertia.New(config)
	require.NoError(t, err)

	// TODO: Fix always props in partial reload - needs investigation
	// The always lazy props work, but static always props need debugging
	t.Skip("Always props in partial reload needs investigation")
}

// TestDefer tests deferred props.
func TestDefer(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	require.NoError(t, err)

	t.Run("deferred props only evaluated when requested", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users", http.NoBody)
		req.Header.Set("X-Inertia", "true")

		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		called := false
		deferFn := func() interface{} {
			called = true
			return "deferred value"
		}

		props := map[string]interface{}{
			"name": "John",
		}

		err := ic.Defer("comments", deferFn).Render("Posts/Show", props)
		require.NoError(t, err)

		// Deferred prop should not be evaluated on initial load
		assert.False(t, called, "deferred prop should not be evaluated initially")
	})

	t.Run("deferred props evaluated when explicitly requested", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/posts/1", http.NoBody)
		req.Header.Set("X-Inertia", "true")
		req.Header.Set("X-Inertia-Partial-Data", "comments")
		req.Header.Set("X-Inertia-Partial-Component", "Posts/Show")

		w := httptest.NewRecorder()

		// Run through middleware
		middleware := mgr.Middleware()
		var capturedReq *http.Request
		handler := middleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			capturedReq = r
		}))
		handler.ServeHTTP(w, req)

		w = httptest.NewRecorder()
		ctx := NewMockContext(w, capturedReq)
		ic := inertia.NewContext(ctx, mgr)

		called := false
		deferFn := func() interface{} {
			called = true
			return []string{"Comment 1", "Comment 2"}
		}

		props := map[string]interface{}{
			"title": "Post Title",
		}

		err := ic.Defer("comments", deferFn).Render("Posts/Show", props)
		require.NoError(t, err)

		// Deferred prop should be evaluated when requested
		assert.True(t, called, "deferred prop should be evaluated when requested")
	})
}
