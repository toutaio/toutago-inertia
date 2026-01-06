# Advanced Usage Guide

## Table of Contents

1. [Validation & Flash Helpers](#validation--flash-helpers)
2. [Lazy Props & Performance](#lazy-props--performance)
3. [Server-Side Rendering (SSR)](#server-side-rendering-ssr)
4. [TypeScript Type Generation](#typescript-type-generation)
5. [Advanced Form Handling](#advanced-form-handling)
6. [Asset Versioning](#asset-versioning)
7. [Partial Reloads](#partial-reloads)
8. [Custom Context Wrappers](#custom-context-wrappers)
9. [Error Handling](#error-handling)
10. [Testing](#testing)
11. [Performance Optimization](#performance-optimization)

## Validation & Flash Helpers

### Validation Helpers

Simplify error handling with built-in validation helpers.

```go
func CreateUser(ctx *YourContext) error {
    ic := inertia.NewContext(ctx, mgr)
    
    // Simple chaining
    if email == "" || !isValid(email) {
        return ic.
            WithError("email", "Email is required").
            WithError("email", "Email must be valid").
            Back()
    }
    
    // Or build manually
    errors := inertia.NewValidationErrors()
    errors.Add("name", "Name is required")
    errors.Add("email", "Email must be valid")
    
    if errors.Any() {
        return ic.WithErrors(errors).Back()
    }
    
    // Success
    return ic.WithSuccess("User created!").Redirect("/users")
}
```

**Validation Methods:**
- `NewValidationErrors()` - Create new instance
- `Add(field, message)` - Add error to field
- `Has(field)` - Check if field has errors
- `First(field)` - Get first error for field
- `Any()` - Check if any errors exist

### Flash Message Helpers

Show temporary messages to users.

```go
func UpdateProfile(ctx *YourContext) error {
    ic := inertia.NewContext(ctx, mgr)
    
    // Different message types
    return ic.
        WithSuccess("Profile updated successfully").
        WithWarning("Subscription expires in 7 days").
        WithInfo("Check your email for confirmation").
        Redirect("/profile")
}

// Or build flash manually
flash := inertia.NewFlash()
flash.Success("Operation completed")
flash.Error("Something went wrong")
flash.Warning("Please review changes")
flash.Info("For your information")
flash.Custom("customKey", "Custom message")

return ic.WithFlash(flash).Back()
```

**Flash Methods:**
- `NewFlash()` - Create new instance
- `Success(msg)` - Add success message
- `Error(msg)` - Add error message
- `Warning(msg)` - Add warning message
- `Info(msg)` - Add info message
- `Custom(key, msg)` - Add custom message

### Client-Side Usage

```vue
<script setup>
import { useForm } from '@toutaio/inertia-vue'

const form = useForm({
    email: '',
    password: '',
})
</script>

<template>
    <!-- Validation errors -->
    <div v-if="form.errors.email" class="error">
        {{ form.errors.email[0] }}
    </div>
    
    <!-- Flash messages -->
    <div v-if="$page.props.success" class="alert-success">
        {{ $page.props.success }}
    </div>
    <div v-if="$page.props.error" class="alert-danger">
        {{ $page.props.error }}
    </div>
</template>
```

## Lazy Props & Performance

Optimize performance by deferring expensive computations.

### Types of Lazy Props

**1. Lazy() - Excluded from partial reloads**
```go
ic.Lazy("analytics", func() interface{} {
    return calculateAnalytics() // Only on full page load
})
```

**2. AlwaysLazy() - Always included, but lazily evaluated**
```go
ic.AlwaysLazy("auth", func() interface{} {
    return getCurrentUser() // Always included, computed once
})
```

**3. Defer() - Only when explicitly requested**
```go
ic.Defer("comments", func() interface{} {
    return loadComments() // Only when client requests
})
```

### Complete Example

```go
func ShowDashboard(ctx *YourContext) error {
    ic := inertia.NewContext(ctx, mgr)
    
    return ic.
        // Always included
        AlwaysLazy("auth", func() interface{} {
            return map[string]interface{}{
                "user":        getCurrentUser(ctx),
                "permissions": loadPermissions(ctx),
            }
        }).
        // Skipped on navigation
        Lazy("stats", func() interface{} {
            return calculateDashboardStats()
        }).
        // Only when requested
        Defer("auditLog", func() interface{} {
            return loadFullAuditLog()
        }).
        Render("Dashboard/Index", map[string]interface{}{
            "title": "Dashboard",
        })
}
```

### Client-Side Deferred Loading

```javascript
import { router } from '@inertiajs/vue3'

// Load deferred prop on demand
function loadComments() {
    router.reload({
        only: ['comments'], // Request specific props
    })
}
```

### Performance Benefits

```go
// ❌ Bad - always runs expensive query
func ShowPost(ctx *YourContext) error {
    relatedPosts := findRelatedPosts() // Runs every time!
    
    return ic.Render("Posts/Show", map[string]interface{}{
        "related": relatedPosts,
    })
}

// ✅ Good - only runs when needed
func ShowPost(ctx *YourContext) error {
    return ic.
        Lazy("related", func() interface{} {
            return findRelatedPosts() // Only on full load
        }).
        Render("Posts/Show", props)
}
```

## Server-Side Rendering (SSR)

### Why SSR?

- Better SEO (search engines see complete HTML)
- Faster initial page load
- Better perceived performance
- Social media preview cards work correctly

### Setup SSR

**1. Create SSR entry point:**

```typescript
// frontend/ssr.ts
import { createSSRApp } from 'vue'
import { createInertiaSSRApp, createSSRPage } from '@toutaio/inertia-vue'

export async function render(page: any) {
  const { html } = await createInertiaSSRApp({
    page,
    resolve: (name) => {
      const pages = import.meta.glob('./pages/**/*.vue', { eager: true })
      return pages[`./pages/${name}.vue`]
    },
  })
  
  return createSSRPage({
    html,
    page,
    title: 'My App',
    head: '<link rel="stylesheet" href="/dist/app.css">',
  })
}
```

**2. Build SSR bundle:**

```javascript
// build-ssr.js
const esbuild = require('esbuild')

esbuild.build({
  entryPoints: ['frontend/ssr.ts'],
  bundle: true,
  platform: 'node',
  format: 'cjs',
  outfile: 'dist/ssr.cjs',
  external: ['vue'],
})
```

**3. Use in Go:**

```go
import (
    "encoding/json"
    "os/exec"
)

func renderWithSSR(component string, props map[string]interface{}) (string, error) {
    page := map[string]interface{}{
        "component": component,
        "props":     props,
        "url":       "/current-url",
        "version":   "1.0",
    }
    
    pageJSON, _ := json.Marshal(page)
    
    cmd := exec.Command("node", "-e", fmt.Sprintf(`
        const { render } = require('./dist/ssr.cjs');
        render(%s).then(html => console.log(html));
    `, pageJSON))
    
    output, err := cmd.Output()
    return string(output), err
}
```

### SSR with Streaming

For even better performance, stream the response:

```go
func streamSSR(w http.ResponseWriter, component string, props map[string]interface{}) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    
    // Write initial HTML
    w.Write([]byte("<!DOCTYPE html><html><head>"))
    w.Write([]byte("<title>Loading...</title>"))
    w.Write([]byte("</head><body>"))
    
    if f, ok := w.(http.Flusher); ok {
        f.Flush()
    }
    
    // Render component
    html, _ := renderWithSSR(component, props)
    w.Write([]byte(html))
    
    w.Write([]byte("</body></html>"))
}
```

## TypeScript Type Generation

### Automatic Generation

Generate TypeScript interfaces from Go structs automatically:

```go
package models

//go:generate go run github.com/toutaio/toutago-inertia/cmd/typegen -struct=User,Post,Comment -output=../frontend/types/models.ts

type User struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"createdAt"`
}

type Post struct {
    ID       int64  `json:"id"`
    Title    string `json:"title"`
    Content  string `json:"content"`
    AuthorID int64  `json:"authorId"`
    Author   *User  `json:"author,omitempty"`
}
```

Run:
```bash
go generate ./models
```

Generates:
```typescript
// frontend/types/models.ts
export interface User {
  id: number
  name: string
  email: string
  createdAt: string
}

export interface Post {
  id: number
  title: string
  content: string
  authorId: number
  author?: User
}
```

### Custom Type Mappings

Use `ts:""` tag for custom TypeScript types:

```go
type Config struct {
    // JSON number, but TypeScript should treat as specific type
    Status int `json:"status" ts:"'active' | 'inactive'"`
    
    // Custom union type
    Role string `json:"role" ts:"'admin' | 'user' | 'guest'"`
    
    // Metadata as any
    Meta interface{} `json:"meta" ts:"Record<string, any>"`
}
```

Generates:
```typescript
export interface Config {
  status: 'active' | 'inactive'
  role: 'admin' | 'user' | 'guest'
  meta: Record<string, any>
}
```

### Watch Mode

For development, watch for changes:

```bash
# Add to package.json
{
  "scripts": {
    "types:watch": "nodemon --exec 'go generate ./models' --watch '**/*.go'"
  }
}
```

## Advanced Form Handling

### File Uploads

```vue
<script setup lang="ts">
import { useForm } from '@toutaio/inertia-vue'

const form = useForm({
  title: '',
  image: null as File | null,
})

const submit = () => {
  form.post('/posts', {
    onSuccess: () => form.reset(),
  })
}
</script>

<template>
  <form @submit.prevent="submit">
    <input v-model="form.title">
    <input type="file" @input="form.image = $event.target.files[0]">
    
    <div v-if="form.progress">
      Upload progress: {{ form.progress.percentage }}%
    </div>
    
    <button :disabled="form.processing">Upload</button>
  </form>
</template>
```

Backend:
```go
func createPost(w http.ResponseWriter, r *http.Request) {
    r.ParseMultipartForm(10 << 20) // 10 MB max
    
    file, handler, err := r.FormFile("image")
    if err != nil {
        // Handle error
    }
    defer file.Close()
    
    // Save file...
    
    ctx := inertia.WrapContext(w, r)
    ctx.WithFlash("success", "Post created!").Redirect("/posts")
}
```

### Form Validation

```vue
<script setup lang="ts">
import { useForm } from '@toutaio/inertia-vue'

const form = useForm({
  email: '',
  password: '',
})

const submit = () => {
  form.post('/login', {
    onError: () => {
      // Form automatically populates errors
      // form.errors.email, form.errors.password available
    },
  })
}
</script>

<template>
  <form @submit.prevent="submit">
    <div>
      <input v-model="form.email" type="email">
      <span v-if="form.errors.email" class="error">
        {{ form.errors.email }}
      </span>
    </div>
    
    <div>
      <input v-model="form.password" type="password">
      <span v-if="form.errors.password" class="error">
        {{ form.errors.password }}
      </span>
    </div>
    
    <button :disabled="form.processing">Login</button>
  </form>
</template>
```

Backend:
```go
func login(w http.ResponseWriter, r *http.Request) {
    ctx := inertia.WrapContext(w, r)
    
    email := r.FormValue("email")
    password := r.FormValue("password")
    
    errors := make(map[string]string)
    
    if email == "" {
        errors["email"] = "Email is required"
    }
    if len(password) < 8 {
        errors["password"] = "Password must be at least 8 characters"
    }
    
    if len(errors) > 0 {
        ctx.WithErrors(errors).Back()
        return
    }
    
    // Authenticate user...
    ctx.Redirect("/dashboard")
}
```

### Dirty State Tracking

Prevent accidental navigation away from unsaved forms:

```vue
<script setup lang="ts">
import { useForm } from '@toutaio/inertia-vue'
import { onBeforeUnmount } from 'vue'

const form = useForm({
  title: '',
  content: '',
})

// Warn before leaving if form has changes
onBeforeUnmount(() => {
  if (form.isDirty) {
    return confirm('You have unsaved changes. Are you sure?')
  }
})
</script>
```

## Lazy Data Evaluation

Only compute expensive data when specifically requested:

```go
func main() {
    app := inertia.New(&inertia.Config{
        Version: "1.0",
        RootView: "app.html",
    })
    
    // Expensive query - only run when needed
    app.ShareFunc("stats", func(r *http.Request) interface{} {
        return map[string]interface{}{
            "totalUsers":  countUsers(),      // Expensive!
            "totalPosts":  countPosts(),      // Expensive!
            "totalViews":  countPageViews(),  // Expensive!
        }
    })
}
```

Request specific data:

```typescript
import { router } from '@toutaio/inertia-vue'

// Only load 'stats' data
router.reload({ only: ['stats'] })
```

Or exclude data:

```typescript
// Load everything except 'stats'
router.reload({ except: ['stats'] })
```

## Asset Versioning

Automatically detect and force reload when assets change:

```go
func main() {
    // Read asset hash from build manifest
    manifest, _ := os.ReadFile("dist/manifest.json")
    var data map[string]interface{}
    json.Unmarshal(manifest, &data)
    
    app := inertia.New(&inertia.Config{
        Version:  data["hash"].(string),
        RootView: "app.html",
    })
    
    // When version changes, Inertia forces full reload
}
```

Generate manifest during build:

```javascript
// build.js
const crypto = require('crypto')
const fs = require('fs')

esbuild.build({
  // ...build config
}).then(() => {
  const files = fs.readdirSync('dist')
  const content = files.join('')
  const hash = crypto.createHash('md5').update(content).digest('hex')
  
  fs.writeFileSync('dist/manifest.json', JSON.stringify({ hash }))
})
```

## Partial Reloads

Only reload specific parts of the page:

```go
func showPost(w http.ResponseWriter, r *http.Request) {
    ctx := inertia.WrapContext(w, r)
    
    post := getPost(id)
    comments := getComments(id) // Expensive query
    
    // Only re-render these props when needed
    ctx.RenderOnly("Posts/Show", map[string]interface{}{
        "post":     post,
        "comments": comments,
    }, []string{"comments"}) // Only reload comments if requested
}
```

Request partial reload:

```typescript
import { router } from '@toutaio/inertia-vue'

// Only reload comments
router.reload({ only: ['comments'] })
```

Useful for:
- Polling for updates
- Infinite scroll pagination
- Refreshing specific sections

## Custom Context Wrappers

Create router-specific wrappers:

```go
// For chi router
type ChiContext struct {
    *inertia.InertiaContext
    chi chi.Context
}

func WrapChi(w http.ResponseWriter, r *http.Request) *ChiContext {
    return &ChiContext{
        InertiaContext: inertia.WrapContext(w, r),
        chi:            chi.RouteContext(r.Context()),
    }
}

func (c *ChiContext) URLParam(key string) string {
    return c.chi.URLParam(key)
}

// Usage
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := WrapChi(w, r)
    id := ctx.URLParam("id")
    
    ctx.Render("Posts/Show", map[string]interface{}{
        "post": getPost(id),
    })
}
```

## Error Handling

### Custom Error Pages

```go
func errorHandler(w http.ResponseWriter, r *http.Request, err error, status int) {
    ctx := inertia.WrapContext(w, r)
    
    ctx.Error("Error", map[string]interface{}{
        "message": err.Error(),
        "status":  status,
    }, status)
}

// Usage
func handler(w http.ResponseWriter, r *http.Request) {
    post, err := getPost(id)
    if err != nil {
        errorHandler(w, r, err, http.StatusNotFound)
        return
    }
    
    ctx := inertia.WrapContext(w, r)
    ctx.Render("Posts/Show", map[string]interface{}{
        "post": post,
    })
}
```

Error page component:

```vue
<!-- pages/Error.vue -->
<script setup lang="ts">
defineProps<{
  message: string
  status: number
}>()
</script>

<template>
  <div class="error-page">
    <h1>{{ status }}</h1>
    <p>{{ message }}</p>
    <Link href="/">Go Home</Link>
  </div>
</template>
```

### Global Error Handler

```typescript
// frontend/app.ts
import { router } from '@toutaio/inertia-vue'

router.on('error', (event) => {
  console.error('Inertia error:', event.detail.errors)
  
  // Show toast notification
  showToast('An error occurred', 'error')
})
```

## Testing

### Testing Vue Components

```typescript
// tests/components/PostCard.spec.ts
import { mount } from '@vue/test-utils'
import { Link } from '@toutaio/inertia-vue'
import PostCard from '../components/PostCard.vue'

describe('PostCard', () => {
  it('renders post title', () => {
    const wrapper = mount(PostCard, {
      props: {
        post: {
          id: 1,
          title: 'Test Post',
        },
      },
      global: {
        components: { Link },
      },
    })
    
    expect(wrapper.text()).toContain('Test Post')
  })
})
```

### Testing Inertia Handlers

```go
// handlers_test.go
func TestShowPost(t *testing.T) {
    req := httptest.NewRequest("GET", "/posts/1", nil)
    req.Header.Set("X-Inertia", "true")
    w := httptest.NewRecorder()
    
    showPost(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    var page inertia.Page
    json.Unmarshal(w.Body.Bytes(), &page)
    
    assert.Equal(t, "Posts/Show", page.Component)
    assert.NotNil(t, page.Props["post"])
}
```

## Performance Optimization

### Shared Data Caching

Cache expensive shared data:

```go
var userCache sync.Map

app.ShareFunc("auth", func(r *http.Request) interface{} {
    userID := getSessionUserID(r)
    
    if cached, ok := userCache.Load(userID); ok {
        return cached
    }
    
    user := getUser(userID)
    userCache.Store(userID, user)
    
    return user
})
```

### Lazy Loading Components

Load heavy components only when needed:

```typescript
// frontend/app.ts
import { defineAsyncComponent } from 'vue'

createInertiaApp({
  resolve: (name) => {
    // Lazy load heavy components
    if (name === 'Dashboard') {
      return defineAsyncComponent(() => import('./pages/Dashboard.vue'))
    }
    
    const pages = import.meta.glob('./pages/**/*.vue', { eager: true })
    return pages[`./pages/${name}.vue`]
  },
})
```

### Prefetching

Prefetch pages before navigation:

```typescript
import { router } from '@toutaio/inertia-vue'

// On hover, prefetch page
const prefetch = () => {
  router.prefetch('/posts/123')
}
```

```vue
<template>
  <Link href="/posts/123" @mouseenter="prefetch">
    View Post
  </Link>
</template>
```

### Response Compression

```go
import "github.com/gorilla/handlers"

func main() {
    handler := inertiaApp.Middleware()(http.HandlerFunc(routes))
    compressed := handlers.CompressHandler(handler)
    
    http.ListenAndServe(":3000", compressed)
}
```

## Next Steps

- Explore [examples](../examples/) for complete working applications
- Check out the [API reference](../README.md)
- Learn about [testing strategies](../npm/inertia-vue/README.md)
- Read about [deployment best practices](MIGRATION.md)
