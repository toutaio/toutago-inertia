package inertia_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toutaio/toutago-inertia/pkg/inertia"
)

func TestMiddleware_DetectInertiaRequest(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	middleware := i.Middleware()

	tests := []struct {
		name        string
		headers     map[string]string
		wantInertia bool
	}{
		{
			name: "with X-Inertia header",
			headers: map[string]string{
				"X-Inertia": "true",
			},
			wantInertia: true,
		},
		{
			name:        "without X-Inertia header",
			headers:     map[string]string{},
			wantInertia: false,
		},
		{
			name: "with X-Inertia false",
			headers: map[string]string{
				"X-Inertia": "false",
			},
			wantInertia: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			w := httptest.NewRecorder()

			called := false
			handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				isInertia := inertia.IsInertiaRequest(r)
				assert.Equal(t, tt.wantInertia, isInertia)
			}))

			handler.ServeHTTP(w, req)
			assert.True(t, called)
		})
	}
}

func TestMiddleware_SetVersionHeader(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	middleware := i.Middleware()

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	assert.Equal(t, "1.0.0", w.Header().Get("X-Inertia-Version"))
}

func TestMiddleware_VersionConflict(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "2.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	middleware := i.Middleware()

	// Client has old version
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Inertia", "true")
	req.Header.Set("X-Inertia-Version", "1.0.0")
	w := httptest.NewRecorder()

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	// Should return 409 Conflict
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestMiddleware_ExternalRedirect(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	middleware := i.Middleware()

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Signal external redirect (don't write status)
		inertia.SetExternalRedirect(r, "https://external.com")
		// Don't call w.WriteHeader - let middleware handle it
	}))

	handler.ServeHTTP(w, req)

	// Should return 409 with X-Inertia-Location header
	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, "https://external.com", w.Header().Get("X-Inertia-Location"))
}

func TestMiddleware_PartialReload(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	middleware := i.Middleware()

	// Request only specific props
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Inertia", "true")
	req.Header.Set("X-Inertia-Partial-Data", "user,posts")
	req.Header.Set("X-Inertia-Partial-Component", "Users/Show")
	w := httptest.NewRecorder()

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get requested only props
		only := inertia.GetPartialOnly(r)
		assert.Equal(t, []string{"user", "posts"}, only)

		component := inertia.GetPartialComponent(r)
		assert.Equal(t, "Users/Show", component)
	}))

	handler.ServeHTTP(w, req)
}

func TestMiddleware_NonInertiaRequest(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	middleware := i.Middleware()

	// Normal browser request (no X-Inertia header)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	called := false
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		// Should still work normally
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
	// Version header should still be set
	assert.Equal(t, "1.0.0", w.Header().Get("X-Inertia-Version"))
}

func TestIsInertiaRequest(t *testing.T) {
	tests := []struct {
		name        string
		headerValue string
		want        bool
	}{
		{"true value", "true", true},
		{"TRUE value", "TRUE", true},
		{"false value", "false", false},
		{"empty value", "", false},
		{"1 value", "1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.headerValue != "" {
				req.Header.Set("X-Inertia", tt.headerValue)
			}

			got := inertia.IsInertiaRequest(req)
			assert.Equal(t, tt.want, got)
		})
	}
}
