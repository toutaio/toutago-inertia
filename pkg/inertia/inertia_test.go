package inertia_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toutaio/toutago-inertia/pkg/inertia"
)

func TestResponse_Creation(t *testing.T) {
	response := inertia.Response{
		Component: "Users/Show",
		Props: map[string]interface{}{
			"user": map[string]string{
				"name":  "John Doe",
				"email": "john@example.com",
			},
		},
		URL:     "/users/1",
		Version: "1.0.0",
	}

	assert.Equal(t, "Users/Show", response.Component)
	assert.Equal(t, "/users/1", response.URL)
	assert.Equal(t, "1.0.0", response.Version)
	assert.NotNil(t, response.Props)
}

func TestResponse_MarshalJSON(t *testing.T) {
	response := inertia.Response{
		Component: "Users/Index",
		Props: map[string]interface{}{
			"users": []string{"Alice", "Bob"},
		},
		URL:     "/users",
		Version: "1.0.0",
	}

	// Response should be JSON serializable
	data, err := response.MarshalJSON()
	require.NoError(t, err)
	assert.Contains(t, string(data), "Users/Index")
	assert.Contains(t, string(data), "Alice")
}

func TestConfig_Validation(t *testing.T) {
	tests := []struct {
		name    string
		config  inertia.Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: inertia.Config{
				RootView: "app.html",
				Version:  "1.0.0",
			},
			wantErr: false,
		},
		{
			name: "missing root view",
			config: inertia.Config{
				Version: "1.0.0",
			},
			wantErr: true,
		},
		{
			name: "empty version is ok (will auto-generate)",
			config: inertia.Config{
				RootView: "app.html",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNew(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)
	assert.NotNil(t, i)
}

func TestNew_InvalidConfig(t *testing.T) {
	config := inertia.Config{
		// Missing RootView
		Version: "1.0.0",
	}

	_, err := inertia.New(config)
	assert.Error(t, err)
}

func TestInertia_Share(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	// Share a simple value
	i.Share("app_name", "My App")

	// Share a function
	i.ShareFunc("user", func() interface{} {
		return map[string]string{"name": "Test User"}
	})

	// Verify shared data is stored
	shared := i.GetSharedData()
	assert.Contains(t, shared, "app_name")
	assert.Contains(t, shared, "user")
}

func TestInertia_Version(t *testing.T) {
	config := inertia.Config{
		RootView: "app.html",
		Version:  "1.0.0",
	}

	i, err := inertia.New(config)
	require.NoError(t, err)

	assert.Equal(t, "1.0.0", i.Version())

	// Update version
	i.SetVersion("2.0.0")
	assert.Equal(t, "2.0.0", i.Version())
}

func TestPage_Creation(t *testing.T) {
	page := inertia.Page{
		Component: "Dashboard/Index",
		Props: map[string]interface{}{
			"title": "Dashboard",
		},
		URL:     "/dashboard",
		Version: "1.0.0",
	}

	assert.Equal(t, "Dashboard/Index", page.Component)
	assert.Equal(t, "/dashboard", page.URL)
}

func TestPage_WithSharedData(t *testing.T) {
	sharedData := map[string]interface{}{
		"auth": map[string]string{
			"user": "admin",
		},
		"flash": map[string]string{
			"success": "Saved!",
		},
	}

	pageProps := map[string]interface{}{
		"title": "Home",
	}

	page := inertia.NewPage("Home/Index", pageProps, "/home", "1.0.0")
	page.MergeSharedData(sharedData)

	// Shared data should be merged with props
	assert.Contains(t, page.Props, "auth")
	assert.Contains(t, page.Props, "flash")
	assert.Contains(t, page.Props, "title")
}
