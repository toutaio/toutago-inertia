# Getting Started with Toutā Inertia

Build modern single-page applications without the complexity of a REST API.

## Installation

```bash
go get github.com/toutaio/toutago-inertia
```

## Quick Start

### 1. Initialize Inertia

```go
config := inertia.Config{
    RootView: "templates/app.html",
    Version:  "1.0.0",
}

inertiaMgr, err := inertia.New(config)
if err != nil {
    panic(err)
}
```

### 2. Add Middleware

```go
mux := http.NewServeMux()
handler := inertiaMgr.Middleware()(mux)
```

### 3. Create Handlers

```go
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    page, _ := inertiaMgr.Render("Home/Index", map[string]interface{}{
        "message": "Hello Inertia!",
    }, r.URL.Path)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(page)
}
```

## Core Concepts

### Shared Data

```go
inertiaMgr.Share("appName", "My App")
inertiaMgr.ShareFunc("user", getCurrentUser)
```

### Redirects

```go
inertiaMgr.Redirect(w, r, "/dashboard")  // Internal
inertiaMgr.Location(w, r, "https://...")  // External
inertiaMgr.Back(w, r)                     // Go back
```

### Validation Errors

```go
errors := inertia.ValidationErrors{
    "email": []string{"Email is required"},
}
page.WithErrors(errors)
```

### Flash Messages

```go
flash := inertia.Flash{
    "success": "Saved successfully!",
}
page.WithFlash(flash)
```

## Next Steps

- See [examples](../examples/)
- Read [API reference](https://pkg.go.dev/github.com/toutaio/toutago-inertia)
