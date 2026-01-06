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
