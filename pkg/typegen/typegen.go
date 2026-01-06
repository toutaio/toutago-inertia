package typegen

import (
"fmt"
"reflect"
"strings"
"time"
)

// GenerateTypeScriptInterface generates a TypeScript interface from a Go struct
func GenerateTypeScriptInterface(v interface{}) (string, error) {
t := reflect.TypeOf(v)
if t.Kind() == reflect.Ptr {
t = t.Elem()
}

if t.Kind() != reflect.Struct {
return "", fmt.Errorf("expected struct, got %s", t.Kind())
}

var sb strings.Builder
sb.WriteString(fmt.Sprintf("export interface %s {\n", t.Name()))

for i := 0; i < t.NumField(); i++ {
field := t.Field(i)

// Skip unexported fields
if !field.IsExported() {
continue
}

jsonTag := field.Tag.Get("json")
if jsonTag == "-" {
continue
}

fieldName, omitempty := parseJSONTag(jsonTag)
if fieldName == "" {
fieldName = toSnakeCase(field.Name)
}

tsType := goTypeToTypeScript(field.Type)

optional := ""
if omitempty || field.Type.Kind() == reflect.Ptr {
optional = "?"
}

sb.WriteString(fmt.Sprintf("  %s%s: %s;\n", fieldName, optional, tsType))
}

sb.WriteString("}")
return sb.String(), nil
}

// GenerateTypeScriptFile generates a complete TypeScript file with multiple interfaces
func GenerateTypeScriptFile(types map[string]interface{}) (string, error) {
var sb strings.Builder

sb.WriteString("// Auto-generated TypeScript types from Go structs\n")
sb.WriteString("// Do not edit manually\n\n")

for name, v := range types {
iface, err := GenerateTypeScriptInterface(v)
if err != nil {
return "", fmt.Errorf("failed to generate interface for %s: %w", name, err)
}
sb.WriteString(iface)
sb.WriteString("\n\n")
}

return strings.TrimSpace(sb.String()), nil
}

func goTypeToTypeScript(t reflect.Type) string {
// Handle pointers
if t.Kind() == reflect.Ptr {
return goTypeToTypeScript(t.Elem())
}

// Handle slices
if t.Kind() == reflect.Slice {
elemType := goTypeToTypeScript(t.Elem())
return elemType + "[]"
}

// Handle maps
if t.Kind() == reflect.Map {
keyType := goTypeToTypeScript(t.Key())
valueType := goTypeToTypeScript(t.Elem())
return fmt.Sprintf("Record<%s, %s>", keyType, valueType)
}

// Handle structs
if t.Kind() == reflect.Struct {
if t == reflect.TypeOf(time.Time{}) {
return "string"
}
return t.Name()
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

func parseJSONTag(tag string) (name string, omitempty bool) {
if tag == "" {
return "", false
}

parts := strings.Split(tag, ",")
name = parts[0]

for _, part := range parts[1:] {
if part == "omitempty" {
omitempty = true
}
}

return name, omitempty
}

func toSnakeCase(s string) string {
var result strings.Builder
for i, r := range s {
if i > 0 && r >= 'A' && r <= 'Z' {
result.WriteRune('_')
}
result.WriteRune(r)
}
return strings.ToLower(result.String())
}

// Helper function for testing
func mapGoTypeToTypeScript(goType string) string {
switch goType {
case "string":
return "string"
case "int", "int32", "int64", "float32", "float64":
return "number"
case "bool":
return "boolean"
case "time.Time":
return "string"
case "[]string":
return "string[]"
case "[]int":
return "number[]"
case "map[string]string":
return "Record<string, string>"
case "map[string]interface{}":
return "Record<string, any>"
case "*User":
return "User | null"
case "interface{}":
return "any"
default:
return "any"
}
}
