package inertia_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/toutaio/toutago-inertia/pkg/inertia"
)

// BenchmarkRender benchmarks the Render function.
func BenchmarkRender(b *testing.B) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	if err != nil {
		b.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/users", http.NoBody)
	req.Header.Set("X-Inertia", "true")

	props := map[string]interface{}{
		"users": []map[string]string{
			{"name": "John", "email": "john@example.com"},
			{"name": "Jane", "email": "jane@example.com"},
		},
		"total": 2,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		if err := ic.Render("Users/Index", props); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderWithSharedData benchmarks rendering with shared data.
func BenchmarkRenderWithSharedData(b *testing.B) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	if err != nil {
		b.Fatal(err)
	}

	mgr.Share("app", map[string]interface{}{
		"name":    "MyApp",
		"version": "1.0.0",
	})

	req := httptest.NewRequest("GET", "/users", http.NoBody)
	req.Header.Set("X-Inertia", "true")

	props := map[string]interface{}{
		"users": []map[string]string{
			{"name": "John", "email": "john@example.com"},
			{"name": "Jane", "email": "jane@example.com"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		if err := ic.Render("Users/Index", props); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderWithLazyProps benchmarks rendering with lazy props.
func BenchmarkRenderWithLazyProps(b *testing.B) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	if err != nil {
		b.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/dashboard", http.NoBody)
	req.Header.Set("X-Inertia", "true")

	props := map[string]interface{}{
		"stats": map[string]int{"visits": 100},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		ic.Lazy("expensive", func() interface{} {
			return map[string]interface{}{
				"data": "expensive computation",
			}
		})

		if err := ic.Render("Dashboard", props); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPartialReload benchmarks partial reload performance.
func BenchmarkPartialReload(b *testing.B) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	if err != nil {
		b.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/users", http.NoBody)
	req.Header.Set("X-Inertia", "true")
	req.Header.Set("X-Inertia-Partial-Data", "users")
	req.Header.Set("X-Inertia-Partial-Component", "Users/Index")

	props := map[string]interface{}{
		"users": []map[string]string{
			{"name": "John", "email": "john@example.com"},
			{"name": "Jane", "email": "jane@example.com"},
		},
		"stats": map[string]int{"total": 100},
		"meta":  map[string]string{"title": "Users"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()

		// Run through middleware to set context values
		middleware := mgr.Middleware()
		var capturedReq *http.Request
		handler := middleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			capturedReq = r
		}))
		handler.ServeHTTP(w, req)

		// Now render with captured request
		w = httptest.NewRecorder()
		ctx := NewMockContext(w, capturedReq)
		ic := inertia.NewContext(ctx, mgr)

		if err := ic.Render("Users/Index", props); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkHTMXPartial benchmarks HTMX partial rendering.
func BenchmarkHTMXPartial(b *testing.B) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	mgr, err := inertia.New(config)
	if err != nil {
		b.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/users", http.NoBody)
	req.Header.Set("HX-Request", "true")

	html := `<div id="user-list">
		<div class="user">John</div>
		<div class="user">Jane</div>
	</div>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		ctx := NewMockContext(w, req)
		ic := inertia.NewContext(ctx, mgr)

		if err := ic.HTMXPartial(html); err != nil {
			b.Fatal(err)
		}
	}
}
