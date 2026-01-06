# Full-Stack Inertia Example

This example demonstrates a complete application using Toutago Inertia with Vue 3, including:

- ✅ Server-side rendering (SSR)
- ✅ TypeScript type generation from Go structs
- ✅ Form handling with validation
- ✅ Navigation with Inertia Link
- ✅ Shared layout components
- ✅ Asset versioning
- ✅ Flash messages
- ✅ Error handling

## Running the Example

### 1. Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install NPM dependencies
cd frontend
npm install
```

### 2. Build Frontend

```bash
cd frontend
npm run build
```

### 3. Run the Server

```bash
# From the fullstack directory
go run backend/main.go
```

### 4. Visit the Application

Open http://localhost:3000 in your browser.

## Development Mode

For development with hot-reload:

```bash
# Terminal 1: Run frontend dev server
cd frontend
npm run dev

# Terminal 2: Run Go server
go run backend/main.go
```

## Project Structure

```
fullstack/
├── backend/
│   ├── main.go           # Server entry point
│   ├── handlers/         # HTTP handlers
│   └── models/           # Data models
├── frontend/
│   ├── src/
│   │   ├── components/   # Vue components
│   │   ├── pages/        # Page components
│   │   ├── types/        # Generated TypeScript types
│   │   ├── app.ts        # Frontend entry point
│   │   └── ssr.ts        # SSR entry point
│   ├── public/           # Static assets
│   └── package.json
└── go.mod
```

## Features Demonstrated

### Type Safety

Go structs are automatically converted to TypeScript types:

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

Generates TypeScript types that can be imported in Vue components.

### Form Handling

Uses Inertia's form helper for easy form submission with validation:

```typescript
const form = useForm({
  name: '',
  email: ''
})

form.post('/users')
```

### Navigation

Client-side navigation without full page reloads:

```vue
<Link href="/dashboard">Dashboard</Link>
```

### Shared Layouts

Pages can share common layouts:

```vue
<template>
  <Layout>
    <h1>Page Content</h1>
  </Layout>
</template>
```
