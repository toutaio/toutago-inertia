# Inertia with HTTP Router Example

Example showing Inertia.js integration with standard net/http using context wrapper.

## Features Demonstrated

- ✅ Context wrapper pattern
- ✅ Middleware setup
- ✅ Page rendering
- ✅ Form validation
- ✅ Flash messages
- ✅ Error handling

## Run

```bash
go run main.go
```

Visit http://localhost:3000

## Routes

- `GET /` - Home page
- `GET /users` - Users list  
- `GET /users/create` - Create form
- `POST /users/create` - Store user (with validation)

## Key Pattern

```go
// Create simple context adapter
type SimpleContext struct {
    w http.ResponseWriter
    r *http.Request
}

func (c *SimpleContext) Request() *http.Request { return c.r }
func (c *SimpleContext) Response() http.ResponseWriter { return c.w }
// ... implement ContextInterface

// Use in handlers
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := NewSimpleContext(w, r)
    ictx := inertia.NewContext(ctx, mgr)
    ictx.Render("Page", props)
}
```

This pattern works with any router - just implement ContextInterface!
