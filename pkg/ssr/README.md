# Server-Side Rendering (SSR) with V8

This package provides server-side rendering support for Inertia.js using V8 (via v8go).

## Features

- Server-side rendering with V8 JavaScript engine
- Context pooling for performance
- Timeout protection
- Error handling with graceful fallback
- Support for complex data structures

## Installation

The SSR package requires `v8go`:

```bash
go get rogchap.com/v8go
```

## Basic Usage

```go
package main

import (
    "github.com/toutaio/toutago-inertia/pkg/inertia"
    "github.com/toutaio/toutago-inertia/pkg/ssr"
)

func main() {
    // Create SSR renderer
    renderer, err := ssr.NewRenderer()
    if err != nil {
        panic(err)
    }
    defer renderer.Close()

    // Load your Vue SSR bundle
    bundle := `
        global.render = function(page) {
            return '<div id="app"><h1>' + page.component + '</h1></div>';
        };
    `
    if err := renderer.LoadBundle(bundle); err != nil {
        panic(err)
    }

    // Create Inertia instance with SSR
    i, _ := inertia.New(inertia.Config{RootView: "app"})
    i.SetSSRRenderer(renderer)

    // Now your pages can be server-side rendered
}
```

## Configuration

```go
cfg := &ssr.Config{
    PoolSize: 10,              // Number of V8 contexts in pool (default: 10)
    Timeout:  30 * time.Second, // Render timeout (default: 30s)
}
renderer, err := ssr.NewRenderer(cfg)
```

## Advanced Usage

### With Vue SSR

```javascript
// ssr-bundle.js - Built with esbuild or webpack
const { createSSRApp } = require('vue')
const { renderToString } = require('@vue/server-renderer')

global.render = async function(page) {
    const app = createSSRApp({
        data: () => page.props,
        template: `<div id="app">
            <h1>{{ title }}</h1>
            <p>{{ message }}</p>
        </div>`
    })
    
    const html = await renderToString(app)
    return html
}
```

### With Head Management

```javascript
global.render = function(page) {
    return {
        html: '<div id="app">...</div>',
        head: '<title>' + page.props.title + '</title>'
    }
}
```

### Error Handling

```go
page := inertia.NewPage("Home", props, "/", "1")
html, err := i.RenderSSR(context.Background(), page)
if err != nil {
    // SSR failed - fall back to client-side rendering
    log.Printf("SSR error: %v", err)
    // Return client-side template
}
```

## Performance

The SSR renderer uses context pooling to avoid creating new V8 contexts for each request:

- Pool size: 10 contexts by default
- Context reuse for better performance
- Automatic scaling when pool is exhausted

Benchmark results:
```
BenchmarkSSRRender-8    5000    250 μs/op
```

## Limitations

- V8 contexts are not goroutine-safe (handled internally with pooling)
- Timeout applies per render (default 30s)
- Bundle must define `global.render` function
- No DOM APIs (server-side only)

## See Also

- [Inertia.js SSR Documentation](https://inertiajs.com/server-side-rendering)
- [v8go Documentation](https://pkg.go.dev/rogchap.com/v8go)
