# Toutā Inertia Adapter

[![Go Report Card](https://goreportcard.com/badge/github.com/toutaio/toutago-inertia)](https://goreportcard.com/report/github.com/toutaio/toutago-inertia)
[![Go Reference](https://pkg.go.dev/badge/github.com/toutaio/toutago-inertia.svg)](https://pkg.go.dev/github.com/toutaio/toutago-inertia)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Inertia.js adapter for the Toutā framework. Build modern single-page applications using server-side routing without writing a REST API.

## Features

- 🚀 **Complete Inertia.js Protocol** - Full server-side implementation
- 🎨 **Vue 3 Support** - First-class Vue 3 integration with SSR
- ⚡ **Server-Side Rendering** - V8-powered SSR for SEO and performance
- 🔒 **Type Safety** - Auto-generate TypeScript types from Go structs
- 🔄 **Real-time Updates** - WebSocket integration with `useLiveUpdate`
- 📨 **Scéla Integration** - **NEW!** Message bus integration for pub/sub
- 🎯 **HTMX Support** - Full HTMX integration for hypermedia-driven apps
- 🚀 **Lazy Props** - Performance optimization with deferred loading
- ✅ **Validation Helpers** - Built-in error & flash message helpers
- 📦 **No API Needed** - Direct controller → component data flow
- 🧪 **Well Tested** - Comprehensive test suite with >85% coverage

## Installation

```bash
go get github.com/toutaio/toutago-inertia
```

For the Vue 3 client package:

```bash
npm install @toutaio/inertia-vue
```

## Quick Start

### Backend Setup

```go
package main

import (
    "github.com/toutaio/toutago-cosan-router/pkg/cosan"
    "github.com/toutaio/toutago-inertia/pkg/inertia"
)

func main() {
    router := cosan.New()
    
    // Initialize Inertia
    inertiaCfg := inertia.Config{
        RootView: "app.html",
        Version:  "1.0.0",
        SSR:      true,
    }
    
    inertia := inertia.New(inertiaCfg)
    
    // Add middleware
    router.Use(inertia.Middleware())
    
    // Share global data
    inertia.Share("auth", func(ctx cosan.Context) interface{} {
        return ctx.Get("user")
    })
    
    // Routes
    router.GET("/users/:id", ShowUser)
    
    router.Listen(":3000")
}

func ShowUser(ctx cosan.Context) error {
    user := getUserByID(ctx.Param("id"))
    
    return ctx.Inertia("Users/Show", map[string]interface{}{
        "user": user,
    })
}
```

### Frontend Setup (Vue 3)

```javascript
// frontend/app.js
import { createInertiaApp } from '@toutaio/inertia-vue'
import { createSSRApp, h } from 'vue'

createInertiaApp({
    resolve: name => {
        const pages = import.meta.glob('./pages/**/*.vue')
        return pages[`./pages/${name}.vue`]()
    },
    setup({ el, App, props, plugin }) {
        createSSRApp({ render: () => h(App, props) })
            .use(plugin)
            .mount(el)
    },
})
```

```vue
<!-- frontend/pages/Users/Show.vue -->
<script setup>
import { Head } from '@toutaio/inertia-vue'

defineProps({
    user: Object
})
</script>

<template>
    <div>
        <Head :title="user.name" />
        <h1>{{ user.name }}</h1>
        <p>{{ user.email }}</p>
    </div>
</template>
```

## Documentation

- [Getting Started Guide](docs/getting-started.md)
- [HTMX Integration Guide](docs/HTMX.md) - **NEW!** Complete HTMX support
- [Server-Side Rendering](docs/ssr.md)
- [TypeScript Integration](docs/typescript.md)
- [Real-time Updates](docs/realtime.md)
- [API Reference](https://pkg.go.dev/github.com/toutaio/toutago-inertia)

## Examples

- [Blog with SSR](examples/blog-vue/) - Complete blog application
- [Real-time Chat](examples/realtime-chat/) - WebSocket-powered chat
- [Scéla Integration](examples/scela-integration/) - Message bus with filtering

## Requirements

- Go 1.22+
- Node.js 18+ (for frontend builds)
- V8 shared library (for SSR)

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) first.

### Releases

Releases are automated via GitHub Actions. See [Release Process](docs/RELEASING.md) for details on creating new releases.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- Inspired by [Inertia.js](https://inertiajs.com/)
- Part of the [Toutā framework](https://github.com/toutaio/toutago)
