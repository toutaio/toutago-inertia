# Ritual Grove Integration Guide

This guide explains how to integrate Inertia.js support into Toutā Ritual Grove rituals.

## Overview

The `toutago-inertia` package can be used in Ritual Grove rituals to scaffold full-stack applications with modern SPA capabilities while maintaining server-side rendering and type safety.

## Existing Inertia Ritual

### fullstack-inertia-vue

A complete ritual for creating full-stack applications with Inertia.js and Vue 3:

**Location**: `rituals/fullstack-inertia-vue/`

**Features**:
- Vue 3 frontend with Inertia.js adapter
- TypeScript support
- Vite build tool
- Server-side rendering (optional)
- Authentication scaffolding (optional)
- Database integration (PostgreSQL or MySQL)

**Usage**:
```bash
ritual run fullstack-inertia-vue
```

**Prompts**:
- `project_name`: Project name (kebab-case)
- `module_path`: Go module path (e.g., github.com/user/project)
- `db_driver`: Database driver (postgres or mysql)
- `port`: Server port (default: 8080)
- `use_auth`: Include authentication (default: true)
- `use_ssr`: Enable Server-Side Rendering (default: false)

## Adding Inertia to Existing Rituals

### Step 1: Add Dependency

Add `toutago-inertia` to your ritual's dependencies:

```yaml
dependencies:
  packages:
    - github.com/toutaio/toutago-inertia
```

### Step 2: Add Frontend Choice Question

Add a question to let users choose their frontend stack:

```yaml
questions:
  - name: frontend
    type: select
    prompt: "Choose your frontend stack"
    options:
      - inertia-vue    # Modern SPA with Vue 3
      - htmx           # Hypermedia-driven
      - traditional    # Server-side templates
    default: inertia-vue
    required: true
```

### Step 3: Add Conditional Hooks

Use conditional hooks to scaffold the chosen frontend:

```yaml
hooks:
  post_generation:
    # Common setup
    - task: go-mod-init
      config:
        module: "{{.module_path}}"
    
    # Inertia-specific setup
    - task: directory-create
      condition: "{{eq .frontend \"inertia-vue\"}}"
      config:
        path: "frontend"
    
    - task: npm-init
      condition: "{{eq .frontend \"inertia-vue\"}}"
      config:
        working_dir: "."
    
    - task: npm-install
      condition: "{{eq .frontend \"inertia-vue\"}}"
      config:
        packages:
          - "@toutaio/inertia-vue@^0.2.0"
          - "vue@^3.4.0"
          - "vite@^5.0.0"
        dev: true
```

### Step 4: Create Conditional Templates

Create template variants for each frontend choice:

```
templates/
  ├── main.go.tmpl                    # Common backend
  ├── inertia/                        # Inertia templates
  │   ├── frontend/
  │   │   ├── app.js.tmpl
  │   │   ├── pages/
  │   │   │   └── Home.vue.tmpl
  │   │   └── layouts/
  │   │       └── Default.vue.tmpl
  │   └── vite.config.js.tmpl
  ├── htmx/                           # HTMX templates
  │   └── views/
  └── traditional/                    # Traditional templates
      └── views/
```

Use conditional rendering in hooks:

```yaml
- task: template-render
  condition: "{{eq .frontend \"inertia-vue\"}}"
  config:
    source: "inertia/frontend/app.js.tmpl"
    destination: "frontend/app.js"

- task: template-render
  condition: "{{eq .frontend \"htmx\"}}"
  config:
    source: "htmx/views/index.html.tmpl"
    destination: "views/index.html"
```

## Template Examples

### Backend Handler (main.go.tmpl)

```go
package main

import (
    "github.com/toutaio/toutago-cosan-router"
    "github.com/toutaio/toutago-inertia"
)

func main() {
    router := cosan.New()
    inertiaHandler := inertia.New()
    
    {{if eq .frontend "inertia-vue"}}
    // Inertia middleware
    router.Use(inertia.Middleware(inertiaHandler))
    
    // Shared data
    inertiaHandler.Share("app_name", "{{.app_name}}")
    
    // Routes
    router.GET("/", func(ctx *cosan.Context) error {
        return ctx.Inertia().Render("Home", map[string]interface{}{
            "message": "Welcome to {{.app_name}}",
        })
    })
    {{else if eq .frontend "htmx"}}
    // HTMX routes
    router.GET("/", homeHandler)
    {{else}}
    // Traditional routes
    router.GET("/", traditionalHomeHandler)
    {{end}}
    
    router.Listen(":{{.port}}")
}
```

### Frontend Entry (frontend/app.js.tmpl)

```javascript
import { createApp, h } from 'vue'
import { createInertiaApp } from '@toutaio/inertia-vue'

createInertiaApp({
  resolve: name => {
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

### Vue Page (frontend/pages/Home.vue.tmpl)

```vue
<template>
  <Default>
    <div class="container">
      <h1>{{.app_name}}</h1>
      <p>{{ message }}</p>
    </div>
  </Default>
</template>

<script setup>
import Default from '../layouts/Default.vue'

defineProps({
  message: String
})
</script>
```

## Best Practices

### 1. Progressive Enhancement

Start with a simple ritual and add Inertia support gradually:

```yaml
# Phase 1: Basic ritual
questions:
  - name: project_name
  - name: module_path

# Phase 2: Add frontend choice
questions:
  - name: frontend
    type: select
    options: [inertia-vue, traditional]
```

### 2. Sensible Defaults

Choose defaults based on the ritual's purpose:

- **Admin panels**: Default to `inertia-vue` (rich interactivity)
- **Marketing sites**: Default to `traditional` (SEO, simplicity)
- **APIs**: Default to `none` (no frontend)

### 3. Conditional Dependencies

Only install what's needed:

```yaml
- task: go-get
  condition: "{{eq .frontend \"inertia-vue\"}}"
  config:
    packages:
      - github.com/toutaio/toutago-inertia

- task: npm-install
  condition: "{{eq .frontend \"inertia-vue\"}}"
  config:
    packages:
      - "@toutaio/inertia-vue"
```

### 4. Documentation

Include frontend-specific documentation:

```yaml
- task: template-render
  condition: "{{eq .frontend \"inertia-vue\"}}"
  config:
    source: "docs/INERTIA_SETUP.md.tmpl"
    destination: "docs/SETUP.md"
```

## Testing Ritual Integration

### 1. Create Test Ritual

```bash
cd toutago-ritual-grove
ritual create test-inertia-ritual
```

### 2. Add Templates

Follow the structure in `fullstack-inertia-vue` ritual.

### 3. Test Generation

```bash
ritual run test-inertia-ritual --dry-run
```

### 4. Verify Output

Check that:
- Dependencies are installed
- Templates are rendered correctly
- Build scripts work
- Development server starts

## Migration from Traditional Templates

### Step 1: Identify Dynamic Pages

Convert pages with high interactivity to Inertia:

- Forms with validation
- Real-time updates
- Complex state management

### Step 2: Create Inertia Handlers

Replace template rendering with Inertia:

```go
// Before
func handler(ctx *cosan.Context) error {
    return ctx.Render("view.html", data)
}

// After
func handler(ctx *cosan.Context) error {
    return ctx.Inertia().Render("View", map[string]interface{}{
        "data": data,
    })
}
```

### Step 3: Create Vue Components

Convert HTML templates to Vue components:

```html
<!-- Before: view.html -->
<h1>{{.title}}</h1>

<!-- After: View.vue -->
<template>
  <h1>{{ title }}</h1>
</template>

<script setup>
defineProps({ title: String })
</script>
```

### Step 4: Add Build Process

Add Vite configuration and build scripts.

## Example: Blog Ritual with Inertia

See `rituals/blog/` for a complete example of a ritual supporting multiple frontend choices including Inertia.js.

## Resources

- [Inertia.js Documentation](https://inertiajs.com)
- [toutago-inertia API Documentation](./docs/API.md)
- [Example Applications](./examples/)
- [Ritual Grove Documentation](../toutago-ritual-grove/README.md)

## Support

For questions or issues:
- GitHub Issues: https://github.com/toutaio/toutago-inertia/issues
- Ritual Grove Issues: https://github.com/toutaio/toutago-ritual-grove/issues
