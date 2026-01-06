package typegen_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/toutaio/toutago-inertia/pkg/typegen"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  User   `json:"author"`
	Tags    []string `json:"tags"`
}

type Dashboard struct {
	User  User   `json:"user"`
	Posts []Post `json:"posts"`
	Stats struct {
		Total   int `json:"total"`
		Pending int `json:"pending"`
	} `json:"stats"`
}

func TestTypeGenerator_Generate(t *testing.T) {
	t.Run("generates TypeScript interface from struct", func(t *testing.T) {
		gen := typegen.New()
		
		result, err := gen.GenerateInterface("User", User{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Check interface declaration
		if !strings.Contains(result, "export interface User {") {
			t.Error("expected interface declaration")
		}

		// Check fields
		expectedFields := []string{
			"id: number",
			"name: string",
			"email: string",
			"created_at: string",
			"is_active: boolean",
		}

		for _, field := range expectedFields {
			if !strings.Contains(result, field) {
				t.Errorf("expected field %q in output", field)
			}
		}
	})

	t.Run("handles nested structs", func(t *testing.T) {
		gen := typegen.New()
		
		result, err := gen.GenerateInterface("Post", Post{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Check nested User interface
		if !strings.Contains(result, "author: User") {
			t.Error("expected nested User type")
		}

		// Check array type
		if !strings.Contains(result, "tags: string[]") {
			t.Error("expected string array type")
		}
	})

	t.Run("handles inline structs", func(t *testing.T) {
		gen := typegen.New()
		
		result, err := gen.GenerateInterface("Dashboard", Dashboard{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Check inline struct
		if !strings.Contains(result, "stats: {") {
			t.Error("expected inline struct definition")
		}

		if !strings.Contains(result, "total: number") {
			t.Error("expected total field in inline struct")
		}
	})

	t.Run("handles arrays", func(t *testing.T) {
		gen := typegen.New()
		
		result, err := gen.GenerateInterface("Dashboard", Dashboard{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !strings.Contains(result, "posts: Post[]") {
			t.Error("expected Post array type")
		}
	})
}

func TestTypeGenerator_GenerateFile(t *testing.T) {
	t.Run("generates complete TypeScript file", func(t *testing.T) {
		gen := typegen.New()
		gen.Register("User", User{})
		gen.Register("Post", Post{})
		gen.Register("Dashboard", Dashboard{})

		tmpDir := t.TempDir()
		outFile := filepath.Join(tmpDir, "types.ts")

		err := gen.GenerateFile(outFile)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Check file exists
		if _, err := os.Stat(outFile); os.IsNotExist(err) {
			t.Fatal("expected file to be created")
		}

		// Read and verify content
		content, err := os.ReadFile(outFile)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		output := string(content)

		// Check header comment
		if !strings.Contains(output, "// Auto-generated") {
			t.Error("expected auto-generated comment")
		}

		// Check all interfaces are present
		expectedInterfaces := []string{"User", "Post", "Dashboard"}
		for _, name := range expectedInterfaces {
			if !strings.Contains(output, "export interface "+name) {
				t.Errorf("expected interface %s", name)
			}
		}
	})

	t.Run("generates to specified output path", func(t *testing.T) {
		gen := typegen.New()
		gen.Register("User", User{})

		tmpDir := t.TempDir()
		subDir := filepath.Join(tmpDir, "nested", "path")
		outFile := filepath.Join(subDir, "types.ts")

		err := gen.GenerateFile(outFile)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if _, err := os.Stat(outFile); os.IsNotExist(err) {
			t.Fatal("expected file to be created in nested path")
		}
	})
}

func TestTypeGenerator_Options(t *testing.T) {
	t.Run("respects custom indentation", func(t *testing.T) {
		gen := typegen.New(typegen.WithIndent("    ")) // 4 spaces
		
		result, err := gen.GenerateInterface("User", User{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Should use 4 spaces for indentation
		if !strings.Contains(result, "    id: number") {
			t.Error("expected 4-space indentation")
		}
	})

	t.Run("respects custom header", func(t *testing.T) {
		customHeader := "/* Custom Header */"
		gen := typegen.New(typegen.WithHeader(customHeader))
		gen.Register("User", User{})

		tmpDir := t.TempDir()
		outFile := filepath.Join(tmpDir, "types.ts")

		err := gen.GenerateFile(outFile)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		content, err := os.ReadFile(outFile)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		if !strings.Contains(string(content), customHeader) {
			t.Error("expected custom header in output")
		}
	})

	t.Run("handles optional fields", func(t *testing.T) {
		type OptionalFields struct {
			Required string  `json:"required"`
			Optional *string `json:"optional"`
		}

		gen := typegen.New()
		
		result, err := gen.GenerateInterface("OptionalFields", OptionalFields{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Pointer fields should be optional
		if !strings.Contains(result, "optional?: string") {
			t.Error("expected optional field marker for pointer")
		}

		if !strings.Contains(result, "required: string") && strings.Contains(result, "required?:") {
			t.Error("expected required field without optional marker")
		}
	})
}

func TestTypeGenerator_TypeMapping(t *testing.T) {
	t.Run("maps Go types to TypeScript types", func(t *testing.T) {
		type AllTypes struct {
			String    string    `json:"string"`
			Int       int       `json:"int"`
			Int8      int8      `json:"int8"`
			Int16     int16     `json:"int16"`
			Int32     int32     `json:"int32"`
			Int64     int64     `json:"int64"`
			Uint      uint      `json:"uint"`
			Float32   float32   `json:"float32"`
			Float64   float64   `json:"float64"`
			Bool      bool      `json:"bool"`
			Time      time.Time `json:"time"`
			Interface interface{} `json:"interface"`
		}

		gen := typegen.New()
		result, err := gen.GenerateInterface("AllTypes", AllTypes{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expectedMappings := map[string]string{
			"string":    "string: string",
			"int":       "int: number",
			"int8":      "int8: number",
			"int16":     "int16: number",
			"int32":     "int32: number",
			"int64":     "int64: number",
			"uint":      "uint: number",
			"float32":   "float32: number",
			"float64":   "float64: number",
			"bool":      "bool: boolean",
			"time":      "time: string",
			"interface": "interface: any",
		}

		for field, expected := range expectedMappings {
			if !strings.Contains(result, expected) {
				t.Errorf("expected mapping %q for field %s", expected, field)
			}
		}
	})
}

func TestTypeGenerator_EdgeCases(t *testing.T) {
	t.Run("handles empty struct", func(t *testing.T) {
		type Empty struct{}

		gen := typegen.New()
		result, err := gen.GenerateInterface("Empty", Empty{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !strings.Contains(result, "export interface Empty {") {
			t.Error("expected empty interface declaration")
		}
	})

	t.Run("handles struct with no json tags", func(t *testing.T) {
		type NoTags struct {
			Name string
			Age  int
		}

		gen := typegen.New()
		result, err := gen.GenerateInterface("NoTags", NoTags{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Should use field names as-is (or convert to camelCase)
		if !strings.Contains(result, "Name") && !strings.Contains(result, "name") {
			t.Error("expected field Name or name")
		}
	})

	t.Run("handles json tag with omitempty", func(t *testing.T) {
		type WithOmitempty struct {
			Required string `json:"required"`
			Optional string `json:"optional,omitempty"`
		}

		gen := typegen.New()
		result, err := gen.GenerateInterface("WithOmitempty", WithOmitempty{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// omitempty fields should be optional
		if !strings.Contains(result, "optional?: string") {
			t.Error("expected optional marker for omitempty field")
		}
	})

	t.Run("handles json tag with dash (ignored field)", func(t *testing.T) {
		type WithIgnored struct {
			Public  string `json:"public"`
			Ignored string `json:"-"`
		}

		gen := typegen.New()
		result, err := gen.GenerateInterface("WithIgnored", WithIgnored{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !strings.Contains(result, "public: string") {
			t.Error("expected public field")
		}

		// Check that the field itself is not present (not just the type name)
		if strings.Contains(result, "Ignored:") || strings.Contains(result, "ignored:") {
			t.Error("expected ignored field to be excluded")
		}
	})
}
