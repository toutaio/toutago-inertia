package typegen

import (
	"os"
	"testing"
	"time"
)

// Test struct definitions.
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	AuthorID int    `json:"author_id"`
	Author   *User  `json:"author,omitempty"`
}

type PageProps struct {
	User  User   `json:"user"`
	Posts []Post `json:"posts"`
	Count int    `json:"count"`
}

func TestGenerateTypeScriptInterface(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:  "simple struct",
			input: User{},
			expected: `export interface User {
  id: number;
  name: string;
  email: string;
  active: boolean;
  created_at: string;
}`,
		},
		{
			name:  "nested struct",
			input: Post{},
			expected: `export interface Post {
  id: number;
  title: string;
  content: string;
  author_id: number;
  author?: User;
}`,
		},
		{
			name:  "struct with array",
			input: PageProps{},
			expected: `export interface PageProps {
  user: User;
  posts: Post[];
  count: number;
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateTypeScriptInterface(tt.input)
			if err != nil {
				t.Fatalf("GenerateTypeScriptInterface() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("GenerateTypeScriptInterface() =\n%v\n\nwant:\n%v", result, tt.expected)
			}
		})
	}
}

func TestGenerateTypeScriptFile(t *testing.T) {
	types := map[string]interface{}{
		"User":      User{},
		"Post":      Post{},
		"PageProps": PageProps{},
	}

	result, err := GenerateTypeScriptFile(types)
	if err != nil {
		t.Fatalf("GenerateTypeScriptFile() error = %v", err)
	}

	// Check that all interfaces are present
	expectedInterfaces := []string{
		"export interface User {",
		"export interface Post {",
		"export interface PageProps {",
	}

	for _, expected := range expectedInterfaces {
		if !contains(result, expected) {
			t.Errorf("GenerateTypeScriptFile() missing interface: %s", expected)
		}
	}
}

func TestMapGoTypeToTypeScript(t *testing.T) {
	tests := []struct {
		goType string
		want   string
	}{
		{"string", "string"},
		{"int", "number"},
		{"int32", "number"},
		{"int64", "number"},
		{"float32", "number"},
		{"float64", "number"},
		{"bool", "boolean"},
		{"time.Time", "string"},
		{"[]string", "string[]"},
		{"[]int", "number[]"},
		{"map[string]string", "Record<string, string>"},
		{"map[string]interface{}", "Record<string, any>"},
		{"*User", "User | null"},
		{"interface{}", "any"},
	}

	for _, tt := range tests {
		t.Run(tt.goType, func(t *testing.T) {
			got := mapGoTypeToTypeScript(tt.goType)
			if got != tt.want {
				t.Errorf("mapGoTypeToTypeScript(%q) = %q, want %q", tt.goType, got, tt.want)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr || contains(s[1:], substr)))
}

func TestNew(t *testing.T) {
	gen := New()
	if gen == nil {
		t.Fatal("New() returned nil")
	}
	if gen.types == nil {
		t.Error("New() created generator with nil types map")
	}
}

func TestRegister(t *testing.T) {
	gen := New()
	gen.Register("User", User{})
	gen.Register("Post", Post{})

	if len(gen.types) != 2 {
		t.Errorf("Register() added %d types, want 2", len(gen.types))
	}

	if _, ok := gen.types["User"]; !ok {
		t.Error("Register() did not add User type")
	}
	if _, ok := gen.types["Post"]; !ok {
		t.Error("Register() did not add Post type")
	}
}

func TestGenerateFile(t *testing.T) {
	gen := New()
	gen.Register("User", User{})
	gen.Register("Post", Post{})

	tmpDir := t.TempDir()
	outputPath := tmpDir + "/types.ts"

	err := gen.GenerateFile(outputPath)
	if err != nil {
		t.Fatalf("GenerateFile() error = %v", err)
	}

	// Verify file exists
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	// Verify content contains both types
	contentStr := string(content)
	if !contains(contentStr, "export interface User") {
		t.Error("Generated file missing User interface")
	}
	if !contains(contentStr, "export interface Post") {
		t.Error("Generated file missing Post interface")
	}
}

func TestGenerateFileNestedDirectory(t *testing.T) {
	gen := New()
	gen.Register("User", User{})

	tmpDir := t.TempDir()
	outputPath := tmpDir + "/nested/deep/types.ts"

	err := gen.GenerateFile(outputPath)
	if err != nil {
		t.Fatalf("GenerateFile() with nested directory error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("GenerateFile() did not create nested directories: %v", err)
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"UserID", "user_i_d"},
		{"HTTPResponse", "h_t_t_p_response"},
		{"SimpleString", "simple_string"},
		{"ID", "i_d"},
		{"", ""},
		{"A", "a"},
		{"UserName", "user_name"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toSnakeCase(tt.input)
			if got != tt.want {
				t.Errorf("toSnakeCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
