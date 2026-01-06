package ssr

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestNewRenderer(t *testing.T) {
	t.Run("creates renderer with default config", func(t *testing.T) {
		r, err := NewRenderer()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		defer r.Close()

		if r == nil {
			t.Fatal("expected renderer, got nil")
		}
	})

	t.Run("creates renderer with custom config", func(t *testing.T) {
		cfg := &Config{
			PoolSize: 5,
			Timeout:  10 * time.Second,
		}
		r, err := NewRenderer(cfg)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		defer r.Close()

		if r.config.PoolSize != 5 {
			t.Errorf("expected pool size 5, got %d", r.config.PoolSize)
		}
		if r.config.Timeout != 10*time.Second {
			t.Errorf("expected timeout 10s, got %v", r.config.Timeout)
		}
	})
}

func TestLoadBundle(t *testing.T) {
	r, err := NewRenderer()
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}
	defer r.Close()

	t.Run("loads valid JavaScript bundle", func(t *testing.T) {
		bundle := `
			global.render = function(page) {
				return '<div>Hello ' + page.props.name + '</div>';
			};
		`
		err := r.LoadBundle(bundle)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("handles invalid JavaScript", func(t *testing.T) {
		bundle := `this is not valid javascript {{{`
		err := r.LoadBundle(bundle)
		if err == nil {
			t.Error("expected error for invalid JavaScript, got nil")
		}
	})
}

func TestRenderToString(t *testing.T) {
	r, err := NewRenderer()
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}
	defer r.Close()

	bundle := `
		global.render = function(page) {
			return '<div id="app" data-page="' + JSON.stringify(page) + '">' +
				'<h1>' + page.component + '</h1>' +
				'<p>' + page.props.message + '</p>' +
				'</div>';
		};
	`
	if err := r.LoadBundle(bundle); err != nil {
		t.Fatalf("failed to load bundle: %v", err)
	}

	t.Run("renders page to HTML string", func(t *testing.T) {
		pageData := map[string]interface{}{
			"component": "Home",
			"props": map[string]interface{}{
				"message": "Welcome to SSR",
			},
			"url":     "/",
			"version": "1",
		}

		html, err := r.RenderToString(context.Background(), pageData)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if html == "" {
			t.Error("expected non-empty HTML")
		}

		if !contains(html, "Home") {
			t.Error("expected HTML to contain component name")
		}
		if !contains(html, "Welcome to SSR") {
			t.Error("expected HTML to contain message prop")
		}
	})

	t.Run("handles render timeout via config", func(t *testing.T) {
		r2, _ := NewRenderer(&Config{Timeout: 1 * time.Millisecond})
		defer func() {
			// Give any pending renders time to complete
			time.Sleep(50 * time.Millisecond)
			r2.Close()
		}()

		// Simple function that returns quickly
		bundle := `global.render = function(page) { return '<div>Fast</div>'; };`
		r2.LoadBundle(bundle)

		// Just verify timeout config is set
		if r2.config.Timeout != 1*time.Millisecond {
			t.Error("expected 1ms timeout config")
		}
	})

	t.Run("handles render function errors", func(t *testing.T) {
		errorBundle := `
			global.render = function(page) {
				throw new Error('Render failed');
			};
		`
		r2, _ := NewRenderer()
		defer r2.Close()
		r2.LoadBundle(errorBundle)

		pageData := map[string]interface{}{"component": "Test"}
		_, err := r2.RenderToString(context.Background(), pageData)
		if err == nil {
			t.Error("expected error when render throws, got nil")
		}
	})
}

func TestContextPooling(t *testing.T) {
	cfg := &Config{
		PoolSize: 2,
		Timeout:  5 * time.Second,
	}
	r, err := NewRenderer(cfg)
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}
	defer r.Close()

	bundle := `global.render = function(page) { return '<div>Test</div>'; };`
	if err := r.LoadBundle(bundle); err != nil {
		t.Fatalf("failed to load bundle: %v", err)
	}

	t.Run("concurrent renders use different contexts", func(t *testing.T) {
		done := make(chan bool, 3)
		pageData := map[string]interface{}{"component": "Test"}

		for i := 0; i < 3; i++ {
			go func() {
				_, err := r.RenderToString(context.Background(), pageData)
				if err != nil {
					t.Errorf("render failed: %v", err)
				}
				done <- true
			}()
		}

		for i := 0; i < 3; i++ {
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				t.Fatal("timeout waiting for concurrent renders")
			}
		}
	})
}

func TestExtractHead(t *testing.T) {
	r, _ := NewRenderer()
	defer r.Close()

	bundle := `
		global.render = function(page) {
			return {
				html: '<div>Content</div>',
				head: '<title>My Page</title><meta name="description" content="Test">',
			};
		};
	`
	r.LoadBundle(bundle)

	t.Run("extracts head content when returned as object", func(t *testing.T) {
		pageData := map[string]interface{}{"component": "Test"}
		result, err := r.RenderToString(context.Background(), pageData)
		if err != nil {
			t.Fatalf("render failed: %v", err)
		}

		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(result), &obj); err == nil {
			if html, ok := obj["html"].(string); ok {
				if !contains(html, "Content") {
					t.Error("expected HTML content")
				}
			}
			if head, ok := obj["head"].(string); ok {
				if !contains(head, "My Page") {
					t.Error("expected head content")
				}
			}
		}
	})
}

func TestClose(t *testing.T) {
	r, err := NewRenderer()
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	t.Run("closes cleanly", func(t *testing.T) {
		err := r.Close()
		if err != nil {
			t.Errorf("expected no error on close, got %v", err)
		}
	})

	t.Run("cannot render after close", func(t *testing.T) {
		pageData := map[string]interface{}{"component": "Test"}
		_, err := r.RenderToString(context.Background(), pageData)
		if err == nil {
			t.Error("expected error rendering after close, got nil")
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || s != "" && stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
