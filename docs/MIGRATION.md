# Migration Guide

## Migrating from Traditional Server-Side Rendered Apps

This guide helps you migrate from traditional Go server-side rendered applications to Toutā with Inertia.js.

## Overview

**Traditional Approach:**
- Full page reloads on every navigation
- Server renders complete HTML for each request
- Difficult to add interactivity
- Page state lost on navigation

**Inertia Approach:**
- SPA-like navigation without page reloads
- Server sends JSON data, client renders Vue components
- Rich interactivity with Vue ecosystem
- Persistent layouts and shared state

## Step-by-Step Migration

### 1. Install Dependencies

```bash
# Go dependencies
go get github.com/toutaio/toutago-inertia

# NPM dependencies
npm install @toutaio/inertia-vue vue@^3
npm install -D esbuild @vue/compiler-sfc
```

### 2. Setup Inertia Middleware

**Before (traditional):**
```go
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", homeHandler)
    http.ListenAndServe(":3000", mux)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "home.html", data)
}
```

**After (with Inertia):**
```go
import "github.com/toutaio/toutago-inertia"

func main() {
    inertiaApp := inertia.New(&inertia.Config{
        Version: "1.0",
        RootView: "app.html",
    })
    
    handler := inertiaApp.Middleware()(http.HandlerFunc(routes))
    http.ListenAndServe(":3000", handler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    ctx := inertia.WrapContext(w, r)
    ctx.Render("Home", map[string]interface{}{
        "user": getCurrentUser(r),
        "posts": getPosts(),
    })
}
```

### 3. Convert Templates to Vue Components

**Before (HTML template):**
```html
<!-- home.html -->
<html>
<body>
    <h1>Welcome {{.User.Name}}</h1>
    <ul>
        {{range .Posts}}
        <li><a href="/posts/{{.ID}}">{{.Title}}</a></li>
        {{end}}
    </ul>
</body>
</html>
```

**After (Vue component):**
```vue
<!-- frontend/pages/Home.vue -->
<script setup lang="ts">
import { Link } from '@toutaio/inertia-vue'
import type { User, Post } from '../types'

defineProps<{
  user: User
  posts: Post[]
}>()
</script>

<template>
  <div>
    <h1>Welcome {{ user.name }}</h1>
    <ul>
      <li v-for="post in posts" :key="post.id">
        <Link :href="`/posts/${post.id}`">{{ post.title }}</Link>
      </li>
    </ul>
  </div>
</template>
```

### 4. Create Root Template

Create a minimal HTML file that loads your Vue app:

```html
<!-- app.html -->
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>My App</title>
</head>
<body>
    <div id="app" data-page="{{ .Page }}"></div>
    <script src="/dist/app.js"></script>
</body>
</html>
```

### 5. Setup Frontend Entry Point

```typescript
// frontend/app.ts
import { createApp, h } from 'vue'
import { createInertiaApp } from '@toutaio/inertia-vue'

createInertiaApp({
  resolve: (name) => {
    const pages = import.meta.glob('./pages/**/*.vue', { eager: true })
    return pages[`./pages/${name}.vue`]
  },
  setup({ el, App, props, plugin }) {
    createApp({ render: () => h(App, props) })
      .use(plugin)
      .mount(el)
  },
})
```

### 6. Update Forms

**Before:**
```html
<form method="POST" action="/posts">
    <input name="title" value="{{.Title}}">
    <button type="submit">Save</button>
</form>
```

**After:**
```vue
<script setup lang="ts">
import { useForm } from '@toutaio/inertia-vue'

const form = useForm({
  title: '',
  content: '',
})

const submit = () => {
  form.post('/posts', {
    onSuccess: () => form.reset(),
  })
}
</script>

<template>
  <form @submit.prevent="submit">
    <input v-model="form.title" :disabled="form.processing">
    <div v-if="form.errors.title">{{ form.errors.title }}</div>
    <button type="submit" :disabled="form.processing">Save</button>
  </form>
</template>
```

### 7. Handle Redirects

**Before:**
```go
http.Redirect(w, r, "/posts", http.StatusSeeOther)
```

**After:**
```go
ctx := inertia.WrapContext(w, r)
ctx.Redirect("/posts")
```

### 8. Flash Messages

**Before:**
```go
// Set in session
session.Flash("success", "Post created!")

// Read in template
{{.Flash.success}}
```

**After:**
```go
// Backend
ctx.WithFlash("success", "Post created!").Redirect("/posts")

// Frontend
<script setup>
import { usePage } from '@toutaio/inertia-vue'

const page = usePage()
</script>

<template>
  <div v-if="page.props.flash?.success" class="alert">
    {{ page.props.flash.success }}
  </div>
</template>
```

## Common Patterns

### Shared Layouts

**Before:** Duplicate header/footer in each template

**After:**
```vue
<!-- frontend/layouts/AppLayout.vue -->
<script setup lang="ts">
import { Link } from '@toutaio/inertia-vue'
</script>

<template>
  <div>
    <nav>
      <Link href="/">Home</Link>
      <Link href="/about">About</Link>
    </nav>
    <main>
      <slot /> <!-- Page content here -->
    </main>
  </div>
</template>
```

```vue
<!-- frontend/pages/Home.vue -->
<script setup lang="ts">
import AppLayout from '../layouts/AppLayout.vue'
</script>

<template>
  <AppLayout>
    <h1>Home Page</h1>
  </AppLayout>
</template>
```

### Pagination

**Before:**
```go
tmpl.ExecuteTemplate(w, "posts.html", map[string]interface{}{
    "posts": posts,
    "nextPage": page + 1,
    "prevPage": page - 1,
})
```

**After:**
```go
ctx.Render("Posts/Index", map[string]interface{}{
    "posts": inertia.Pagination{
        Data: posts,
        CurrentPage: page,
        PerPage: 20,
        Total: totalCount,
    },
})
```

```vue
<template>
  <div>
    <div v-for="post in posts.data" :key="post.id">
      {{ post.title }}
    </div>
    <Link :href="`/posts?page=${posts.currentPage + 1}`">Next</Link>
  </div>
</template>
```

### Authentication

Share user data across all pages:

```go
func main() {
    inertiaApp := inertia.New(&inertia.Config{
        Version: "1.0",
        RootView: "app.html",
    })
    
    // Share auth user on every request
    inertiaApp.ShareFunc("auth", func(r *http.Request) interface{} {
        return map[string]interface{}{
            "user": getCurrentUser(r),
        }
    })
    
    // ...
}
```

Access in any component:

```vue
<script setup>
import { usePage } from '@toutaio/inertia-vue'

const page = usePage()
const user = computed(() => page.props.auth?.user)
</script>

<template>
  <div v-if="user">
    Welcome {{ user.name }}
  </div>
</template>
```

## TypeScript Integration

Generate TypeScript types from Go structs:

```go
// backend/models/user.go
//go:generate go run github.com/toutaio/toutago-inertia/cmd/typegen -struct User -output ../frontend/types/models.ts

type User struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

Run:
```bash
go generate ./...
```

Use in Vue:
```typescript
import type { User } from './types/models'

const props = defineProps<{
  user: User
}>()
```

## Migration Checklist

- [ ] Install Inertia dependencies (Go + NPM)
- [ ] Setup Inertia middleware
- [ ] Create root view template (app.html)
- [ ] Setup frontend build (esbuild/Vite)
- [ ] Convert one page to Vue component
- [ ] Test navigation works
- [ ] Migrate remaining pages gradually
- [ ] Update forms to use useForm
- [ ] Setup shared data (auth, flash)
- [ ] Add TypeScript types
- [ ] Update redirects to use ctx.Redirect()
- [ ] Test error handling
- [ ] Setup SSR (optional)

## Troubleshooting

### Page doesn't load

- Check `data-page` attribute in root template
- Verify component name matches `ctx.Render("ComponentName")`
- Check browser console for errors

### Forms don't submit

- Use `@submit.prevent` to prevent default form submission
- Use `form.post()` instead of regular form action
- Check backend is receiving JSON correctly

### Flash messages don't show

- Use `ctx.WithFlash()` before redirect
- Access via `usePage().props.flash`
- Remember flash data only available after redirect

### Navigation doesn't work

- Use `<Link>` component from `@toutaio/inertia-vue`
- Don't use regular `<a>` tags for internal links
- External links work automatically

## Next Steps

1. Read the [Advanced Usage Guide](ADVANCED.md)
2. Check out the [examples](../examples/)
3. Learn about [SSR setup](../examples/todo-app/README.md#ssr)
4. Explore [form handling patterns](../npm/inertia-vue/README.md#forms)
