# TypeScript Type Generation Watch Mode Example

This example demonstrates how to use the TypeScript type generator in watch mode to automatically regenerate types when Go files change.

## Usage

```bash
go run main.go
```

The watcher will:
1. Generate initial TypeScript types in `frontend/types.ts`
2. Watch `models/user.go` and `models/post.go` for changes
3. Automatically regenerate types when files are modified
4. Run until you press Ctrl+C

## Try it out

1. Start the watcher:
   ```bash
   go run main.go
   ```

2. In another terminal, modify one of the model files:
   ```bash
   echo '// Updated' >> models/user.go
   ```

3. Watch the watcher automatically regenerate the TypeScript file

## Features

- **File watching**: Monitors specific Go files or entire directories
- **Debouncing**: Prevents excessive regeneration during rapid changes
- **Error handling**: Continues watching even if generation fails
- **Graceful shutdown**: Clean exit on Ctrl+C
