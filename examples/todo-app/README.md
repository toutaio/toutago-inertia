# Todo App Example

A complete example demonstrating Toutago Inertia integration with Vue 3.

## Features

- Full-stack Go + Vue 3 application
- Server-side rendering (SSR)
- Type-safe props with TypeScript
- Form handling with validation
- Authentication example
- Nested layouts demonstration
- Asset bundling with esbuild

## Running the Example

### Prerequisites

- Go 1.21+
- Node.js 18+

### Setup

1. Install Go dependencies:
```bash
go mod download
```

2. Install Node dependencies:
```bash
npm install
```

3. Generate TypeScript types:
```bash
go run cmd/typegen/main.go
```

4. Build frontend assets:
```bash
npm run build
```

5. Run the server:
```bash
go run main.go
```

6. Visit http://localhost:3000

### Development Mode

Run with hot reloading:

```bash
# Terminal 1: Start Go server
go run main.go

# Terminal 2: Watch frontend changes
npm run dev
```

## Project Structure

```
todo-app/
├── cmd/
│   └── typegen/        # TypeScript type generation
├── handlers/           # HTTP handlers
├── models/            # Data models
├── views/             # Vue components
│   ├── layouts/       # Layout components
│   └── pages/         # Page components
├── public/            # Static assets
├── main.go            # Application entry
└── ssr.js             # SSR entry point
```

## Key Concepts Demonstrated

### 1. Basic Rendering

```go
func HandleHome(ctx *cosan.Context) error {
    return ctx.Inertia("Home", inertia.Props{
        "greeting": "Hello, World!",
    })
}
```

### 2. Type-Safe Props

```typescript
// Auto-generated from Go structs
interface TodoPageProps {
    todos: Todo[];
    filter: string;
}
```

### 3. Form Handling

```vue
<script setup lang="ts">
import { useForm } from '@toutaio/inertia-vue'

const form = useForm({
    title: '',
    description: ''
})

const submit = () => {
    form.post('/todos')
}
</script>
```

### 4. Server-Side Rendering

SSR is automatically handled by the middleware. Pages are rendered on the server for initial requests and client-side for subsequent navigation.

### 5. Flash Messages

```go
func HandleCreate(ctx *cosan.Context) error {
    // ... create todo ...
    
    ctx.Session().Flash("success", "Todo created!")
    return ctx.InertiaRedirect("/todos")
}
```

### 6. Nested Layouts

The app demonstrates nested layouts with the Admin section:

```typescript
// In views/app.ts
import { withLayout, resolvePageLayout } from '@toutaio/inertia-vue'

createInertiaApp({
  resolve: async (name) => {
    const page = await pages[`./pages/${name}.vue`]()
    
    // Check if page has custom layout
    const pageLayout = resolvePageLayout(page.default)
    if (pageLayout) {
      page.default = withLayout(page.default, pageLayout)
    }
    
    // Wrap with app layout
    page.default = withLayout(page.default, AppLayout)
    return page
  }
})
```

Visit `/admin/dashboard` to see nested layouts in action (App Layout → Admin Layout → Page).

## Learn More

- [Toutago Documentation](https://github.com/toutaio/toutago)
- [Inertia.js](https://inertiajs.com/)
- [Vue 3](https://vuejs.org/)
