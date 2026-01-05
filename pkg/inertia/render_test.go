package inertia_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toutaio/toutago-inertia/pkg/inertia"
)

func TestRender_WithPartialData(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	// Add shared data
	i.Share("app_name", "Test App")
	i.Share("user", map[string]string{"name": "John"})

	// Render with props
	props := map[string]interface{}{
		"posts": []string{"Post 1", "Post 2"},
		"count": 42,
	}

	page, err := i.Render("Posts/Index", props, "/posts")
	require.NoError(t, err)

	// Should have shared data + props
	assert.Equal(t, "Test App", page.Props["app_name"])
	assert.NotNil(t, page.Props["user"])
	assert.NotNil(t, page.Props["posts"])
	assert.Equal(t, 42, page.Props["count"])
}

func TestRender_OnlyPartialProps(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	i.Share("app_name", "Test App")

	props := map[string]interface{}{
		"posts":  []string{"Post 1"},
		"users":  []string{"User 1"},
		"count":  10,
		"active": true,
	}

	// Render with only filter
	only := []string{"posts", "count"}
	page, err := i.RenderOnly("Posts/Index", props, "/posts", only)
	require.NoError(t, err)

	// Should only have requested props (+ shared data always included)
	assert.Contains(t, page.Props, "app_name") // Shared data always included
	assert.Contains(t, page.Props, "posts")
	assert.Contains(t, page.Props, "count")
	assert.NotContains(t, page.Props, "users")
	assert.NotContains(t, page.Props, "active")
}

func TestPage_ToJSON(t *testing.T) {
	page := inertia.Page{
		Component: "Users/Show",
		Props: map[string]interface{}{
			"user": map[string]string{
				"name":  "Alice",
				"email": "alice@example.com",
			},
		},
		URL:     "/users/1",
		Version: "1.0.0",
	}

	data, err := json.Marshal(page)
	require.NoError(t, err)

	// Should be valid JSON with expected structure
	var decoded map[string]interface{}
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "Users/Show", decoded["component"])
	assert.Equal(t, "/users/1", decoded["url"])
	assert.Equal(t, "1.0.0", decoded["version"])
	assert.NotNil(t, decoded["props"])
}

func TestRender_ValidationErrors(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	tests := []struct {
		name      string
		component string
		url       string
		wantErr   bool
	}{
		{
			name:      "valid",
			component: "Home/Index",
			url:       "/",
			wantErr:   false,
		},
		{
			name:      "missing component",
			component: "",
			url:       "/",
			wantErr:   true,
		},
		{
			name:      "missing url",
			component: "Home/Index",
			url:       "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := i.Render(tt.component, nil, tt.url)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPage_MergeSharedData_NoOverwrite(t *testing.T) {
	shared := map[string]interface{}{
		"app_name": "Shared App",
		"version":  "1.0",
	}

	props := map[string]interface{}{
		"app_name": "Override App", // Should NOT be overwritten
		"title":    "Home",
	}

	page := inertia.NewPage("Home/Index", props, "/", "1.0.0")
	page.MergeSharedData(shared)

	// Props should take precedence
	assert.Equal(t, "Override App", page.Props["app_name"])
	assert.Equal(t, "Home", page.Props["title"])
	assert.Equal(t, "1.0", page.Props["version"]) // From shared
}

func TestInertia_LazySharedData(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	callCount := 0
	i.ShareFunc("dynamic", func() interface{} {
		callCount++
		return map[string]int{"count": callCount}
	})

	// Each render should call the function
	_, err = i.Render("Test/One", nil, "/test1")
	require.NoError(t, err)
	assert.Equal(t, 1, callCount)

	_, err = i.Render("Test/Two", nil, "/test2")
	require.NoError(t, err)
	assert.Equal(t, 2, callCount)
}

func TestRender_EmptyProps(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	// nil props should work
	page, err := i.Render("Home/Index", nil, "/")
	require.NoError(t, err)
	assert.NotNil(t, page.Props)
	assert.Equal(t, 0, len(page.Props))

	// Empty map should work
	page, err = i.Render("Home/Index", map[string]interface{}{}, "/")
	require.NoError(t, err)
	assert.NotNil(t, page.Props)
}
