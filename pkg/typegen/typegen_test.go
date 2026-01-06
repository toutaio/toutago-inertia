package typegen

import (
"testing"
"time"
)

// Test struct definitions
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
