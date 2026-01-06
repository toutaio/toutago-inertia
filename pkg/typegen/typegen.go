package typegen

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// Generator generates TypeScript type definitions from Go structs
type Generator struct {
	types  map[string]interface{}
	indent string
	header string
}

// Option is a functional option for Generator
type Option func(*Generator)

// WithIndent sets the indentation string
func WithIndent(indent string) Option {
	return func(g *Generator) {
		g.indent = indent
	}
}

// WithHeader sets a custom header comment
func WithHeader(header string) Option {
	return func(g *Generator) {
		g.header = header
	}
}

// New creates a new TypeScript type generator
func New(opts ...Option) *Generator {
	g := &Generator{
		types:  make(map[string]interface{}),
		indent: "  ",
		header: "// Auto-generated TypeScript types. DO NOT EDIT.",
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

// Register registers a type to be generated
func (g *Generator) Register(name string, v interface{}) {
	g.types[name] = v
}

// GenerateInterface generates a TypeScript interface for a single struct
func (g *Generator) GenerateInterface(name string, v interface{}) (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("export interface %s {\n", name))

	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("expected struct, got %s", t.Kind())
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get JSON tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue // Skip fields with json:"-"
		}

		// Parse JSON tag
		fieldName, optional := parseJSONTag(jsonTag, field.Name)

		// Check if pointer (makes it optional)
		if field.Type.Kind() == reflect.Ptr {
			optional = true
		}

		// Generate TypeScript type
		tsType := g.goTypeToTS(field.Type)

		// Build field declaration
		optMarker := ""
		if optional {
			optMarker = "?"
		}

		sb.WriteString(fmt.Sprintf("%s%s%s: %s;\n", g.indent, fieldName, optMarker, tsType))
	}

	sb.WriteString("}")

	return sb.String(), nil
}

// GenerateFile generates a TypeScript file with all registered types
func (g *Generator) GenerateFile(outPath string) error {
	var sb strings.Builder

	// Write header
	sb.WriteString(g.header)
	sb.WriteString("\n\n")

	// Generate all interfaces
	for name, v := range g.types {
		iface, err := g.GenerateInterface(name, v)
		if err != nil {
			return fmt.Errorf("failed to generate interface %s: %w", name, err)
		}

		sb.WriteString(iface)
		sb.WriteString("\n\n")
	}

	// Ensure directory exists
	dir := filepath.Dir(outPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(outPath, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// goTypeToTS converts a Go type to a TypeScript type
func (g *Generator) goTypeToTS(t reflect.Type) string {
	// Handle pointers
	if t.Kind() == reflect.Ptr {
		return g.goTypeToTS(t.Elem())
	}

	// Handle arrays/slices
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		elemType := g.goTypeToTS(t.Elem())
		return elemType + "[]"
	}

	// Handle maps
	if t.Kind() == reflect.Map {
		keyType := g.goTypeToTS(t.Key())
		valueType := g.goTypeToTS(t.Elem())
		return fmt.Sprintf("Record<%s, %s>", keyType, valueType)
	}

	// Handle structs
	if t.Kind() == reflect.Struct {
		// Special case for time.Time
		if t == reflect.TypeOf(time.Time{}) {
			return "string"
		}

		// Check if it's a named type
		if t.Name() != "" {
			return t.Name()
		}

		// Inline struct
		return g.generateInlineStruct(t)
	}

	// Handle basic types
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Interface:
		return "any"
	default:
		return "any"
	}
}

// generateInlineStruct generates an inline TypeScript type for anonymous structs
func (g *Generator) generateInlineStruct(t reflect.Type) string {
	var sb strings.Builder

	sb.WriteString("{\n")

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		fieldName, optional := parseJSONTag(jsonTag, field.Name)

		if field.Type.Kind() == reflect.Ptr {
			optional = true
		}

		tsType := g.goTypeToTS(field.Type)

		optMarker := ""
		if optional {
			optMarker = "?"
		}

		sb.WriteString(fmt.Sprintf("%s%s%s%s: %s;\n", g.indent, g.indent, fieldName, optMarker, tsType))
	}

	sb.WriteString(g.indent)
	sb.WriteString("}")

	return sb.String()
}

// parseJSONTag parses a JSON tag and returns the field name and whether it's optional
func parseJSONTag(tag, defaultName string) (string, bool) {
	if tag == "" {
		return defaultName, false
	}

	parts := strings.Split(tag, ",")
	name := parts[0]

	if name == "" {
		name = defaultName
	}

	optional := false
	for _, part := range parts[1:] {
		if part == "omitempty" {
			optional = true
		}
	}

	return name, optional
}
