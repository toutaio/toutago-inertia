# TypeScript Type Generation Example

This example demonstrates how to generate TypeScript type definitions from Go structs.

## Usage

```bash
go run main.go
```

This will generate `examples/types/generated.ts` with TypeScript interfaces for:
- User
- Post  
- DashboardData

## Watch Mode

For development, use the watch mode example to automatically regenerate types when Go files change:

```bash
cd ../typegen-watch
go run main.go
```

See [typegen-watch/README.md](../typegen-watch/README.md) for details on automatic type regeneration.

## Integration with Build Process

You can integrate type generation into your build process:

```bash
# In your Makefile or build script
go run examples/typegen/main.go
```

Or use `go generate`:

```go
//go:generate go run examples/typegen/main.go
```

Then run:

```bash
go generate ./...
```

## Frontend Integration

Once generated, import the types in your Vue components:

```typescript
import type { User, Post, DashboardData } from './types/generated'

// Use with Inertia page props
const props = defineProps<{
  dashboard: DashboardData
}>()
```
