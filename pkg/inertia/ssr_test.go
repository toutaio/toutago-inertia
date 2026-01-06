package inertia

import (
	"context"
	"strings"
	"testing"

	"github.com/toutaio/toutago-inertia/pkg/ssr"
)

func TestSSRIntegration(t *testing.T) {
	t.Run("can attach SSR renderer", func(t *testing.T) {
		renderer, err := ssr.NewRenderer()
		if err != nil {
			t.Fatalf("failed to create renderer: %v", err)
		}
		defer renderer.Close()

		i, _ := New(Config{RootView: "app"})
		i.SetSSRRenderer(renderer)

		if i.ssrRenderer == nil {
			t.Error("expected SSR renderer to be set")
		}
	})

	t.Run("SSR renderer can render page data", func(t *testing.T) {
		renderer, _ := ssr.NewRenderer()
		defer renderer.Close()

		bundle := `
			global.render = function(page) {
				return '<div id="app"><h1>' + page.component + '</h1><p>' + (page.props.message || '') + '</p></div>';
			};
		`
		renderer.LoadBundle(bundle)

		i, _ := New(Config{RootView: "app"})
		i.SetSSRRenderer(renderer)

		page := NewPage("Home", map[string]interface{}{
			"message": "Hello SSR",
		}, "/", "1")

		html, err := i.RenderSSR(context.Background(), page)
		if err != nil {
			t.Fatalf("SSR render failed: %v", err)
		}

		if !strings.Contains(html, "<h1>Home</h1>") {
			t.Error("expected component name in SSR HTML")
		}
		if !strings.Contains(html, "<p>Hello SSR</p>") {
			t.Error("expected props in SSR HTML")
		}
	})

	t.Run("SSR returns empty when no renderer set", func(t *testing.T) {
		i, _ := New(Config{RootView: "app"})
		page := NewPage("Home", map[string]interface{}{}, "/", "1")

		html, err := i.RenderSSR(context.Background(), page)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if html != "" {
			t.Error("expected empty HTML when no SSR renderer")
		}
	})

	t.Run("SSR handles render errors gracefully", func(t *testing.T) {
		renderer, _ := ssr.NewRenderer()
		defer renderer.Close()

		errorBundle := `
			global.render = function(page) {
				throw new Error('Render failed');
			};
		`
		renderer.LoadBundle(errorBundle)

		i, _ := New(Config{RootView: "app"})
		i.SetSSRRenderer(renderer)

		page := NewPage("Home", map[string]interface{}{}, "/", "1")
		_, err := i.RenderSSR(context.Background(), page)
		if err == nil {
			t.Error("expected error from SSR render")
		}
	})
}

func TestSSRWithComplexData(t *testing.T) {
	renderer, _ := ssr.NewRenderer()
	defer renderer.Close()

	bundle := `
		global.render = function(page) {
			var html = '<div id="app">';
			html += '<h1>' + page.component + '</h1>';
			if (page.props.user) {
				html += '<p>User: ' + page.props.user.name + '</p>';
			}
			if (page.props.items && page.props.items.length) {
				html += '<ul>';
				page.props.items.forEach(function(item) {
					html += '<li>' + item + '</li>';
				});
				html += '</ul>';
			}
			html += '</div>';
			return html;
		};
	`
	renderer.LoadBundle(bundle)

	i, _ := New(Config{RootView: "app"})
	i.SetSSRRenderer(renderer)

	t.Run("renders nested objects", func(t *testing.T) {
		page := NewPage("Profile", map[string]interface{}{
			"user": map[string]interface{}{
				"name": "John Doe",
			},
		}, "/profile", "1")

		html, err := i.RenderSSR(context.Background(), page)
		if err != nil {
			t.Fatalf("SSR render failed: %v", err)
		}

		if !strings.Contains(html, "John Doe") {
			t.Error("expected nested user data in HTML")
		}
	})

	t.Run("renders arrays", func(t *testing.T) {
		page := NewPage("List", map[string]interface{}{
			"items": []string{"Item 1", "Item 2", "Item 3"},
		}, "/list", "1")

		html, err := i.RenderSSR(context.Background(), page)
		if err != nil {
			t.Fatalf("SSR render failed: %v", err)
		}

		if !strings.Contains(html, "Item 1") {
			t.Error("expected array items in HTML")
		}
		if !strings.Contains(html, "<ul>") {
			t.Error("expected list markup")
		}
	})
}
