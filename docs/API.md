# API Reference

Complete API reference for toutago-inertia.

## Table of Contents

- [Core API](#core-api)
- [Context Methods](#context-methods)
- [Middleware](#middleware)
- [TypeScript Code Generation](#typescript-code-generation)
- [HTMX Support](#htmx-support)
- [Vue Components](#vue-components)
- [Vue Composables](#vue-composables)

## Core API

### New()

Creates a new Inertia instance.

```go
func New(rootTemplate string, config ...Config) *Inertia
```

**Parameters:**
- `rootTemplate` (string): Path to the root HTML template
- `config` (...Config): Optional configuration (version, shared data)

**Returns:** `*Inertia`

**Example:**
```go
i := inertia.New("views/app.html", inertia.Config{
    Version: "v1.0.0",
})
```

### Render()

Renders an Inertia page with props.

```go
func (i *Inertia) Render(w http.ResponseWriter, r *http.Request, component string, props Props) error
```

**Parameters:**
- `w` (http.ResponseWriter): HTTP response writer
- `r` (*http.Request): HTTP request
- `component` (string): Vue component name
- `props` (Props): Page props (map[string]interface{})

**Returns:** `error`

**Example:**
```go
err := i.Render(w, r, "Home", inertia.Props{
    "user": user,
    "posts": posts,
})
```

### RenderOnly()

Renders only specified props (partial reload).

```go
func (i *Inertia) RenderOnly(w http.ResponseWriter, r *http.Request, component string, props Props, only []string) error
```

**Parameters:**
- `w` (http.ResponseWriter): HTTP response writer
- `r` (*http.Request): HTTP request
- `component` (string): Vue component name
- `props` (Props): All available props
- `only` ([]string): Props to include in response

**Returns:** `error`

**Example:**
```go
err := i.RenderOnly(w, r, "Dashboard", inertia.Props{
    "stats": stats,
    "notifications": notifications,
}, []string{"stats"})
```

### Share()

Adds shared data available to all pages.

```go
func (i *Inertia) Share(key string, value interface{})
```

**Parameters:**
- `key` (string): Shared data key
- `value` (interface{}): Value or lazy function

**Example:**
```go
i.Share("appName", "My App")
i.Share("user", func(r *http.Request) interface{} {
    return getUserFromSession(r)
})
```

### ShareFunc()

Adds lazy-evaluated shared data.

```go
func (i *Inertia) ShareFunc(key string, fn SharedDataFunc)
```

**Parameters:**
- `key` (string): Shared data key
- `fn` (SharedDataFunc): Function returning value

**Example:**
```go
i.ShareFunc("csrf", func(r *http.Request) interface{} {
    return csrf.Token(r)
})
```

### SetVersion()

Sets the asset version for cache busting.

```go
func (i *Inertia) SetVersion(version string)
```

**Parameters:**
- `version` (string): Asset version

**Example:**
```go
i.SetVersion("v2.0.0")
```

## Context Methods

### Render()

Renders an Inertia page from a context.

```go
func (c *InertiaContext) Render(component string, props Props) error
```

**Example:**
```go
return c.Render("Users/Index", inertia.Props{
    "users": users,
})
```

### RenderOnly()

Partial reload from context.

```go
func (c *InertiaContext) RenderOnly(component string, props Props, only []string) error
```

### Location()

External redirect (full page reload).

```go
func (c *InertiaContext) Location(url string) error
```

**Example:**
```go
return c.Location("https://example.com")
```

### Redirect()

Internal redirect (Inertia navigation).

```go
func (c *InertiaContext) Redirect(url string) error
```

**Example:**
```go
return c.Redirect("/dashboard")
```

### Back()

Redirect back to previous page.

```go
func (c *InertiaContext) Back() error
```

### Share()

Share request-specific data.

```go
func (c *InertiaContext) Share(key string, value interface{})
```

**Example:**
```go
c.Share("flash", "Record saved!")
```

### WithErrors()

Redirect with validation errors.

```go
func (c *InertiaContext) WithErrors(errors map[string]string) *InertiaContext
```

**Example:**
```go
return c.WithErrors(map[string]string{
    "email": "Email is required",
}).Back()
```

### WithFlash()

Redirect with flash messages.

```go
func (c *InertiaContext) WithFlash(key string, value interface{}) *InertiaContext
```

**Example:**
```go
return c.WithFlash("success", "User created!").Redirect("/users")
```

### WithSuccess()

Flash success message.

```go
func (c *InertiaContext) WithSuccess(message string) *InertiaContext
```

### WithError()

Flash error message.

```go
func (c *InertiaContext) WithError(message string) *InertiaContext
```

### WithWarning()

Flash warning message.

```go
func (c *InertiaContext) WithWarning(message string) *InertiaContext
```

### WithInfo()

Flash info message.

```go
func (c *InertiaContext) WithInfo(message string) *InertiaContext
```

### Always()

Props always included (never lazy).

```go
func (c *InertiaContext) Always(props Props) *InertiaContext
```

**Example:**
```go
c.Always(inertia.Props{
    "auth": user,
})
```

### AlwaysLazy()

Lazy props always evaluated.

```go
func (c *InertiaContext) AlwaysLazy(key string, fn SharedDataFunc) *InertiaContext
```

## Middleware

### Middleware()

Creates Inertia middleware for handling requests.

```go
func (i *Inertia) Middleware(next http.Handler) http.Handler
```

**Features:**
- Detects Inertia requests
- Handles version conflicts (409)
- Manages partial reloads
- Sets appropriate headers

**Example:**
```go
mux := http.NewServeMux()
mux.Handle("/", i.Middleware(handler))
```

## TypeScript Code Generation

### typegen.New()

Creates a TypeScript generator.

```go
func New() *TypeGen
```

**Example:**
```go
gen := typegen.New()
```

### Register()

Registers a Go struct for type generation.

```go
func (g *TypeGen) Register(name string, v interface{}, opts ...RegisterOption)
```

**Parameters:**
- `name` (string): TypeScript interface name
- `v` (interface{}): Go struct instance
- `opts` (...RegisterOption): Options (WithHeader, WithExport)

**Example:**
```go
gen.Register("User", User{}, typegen.WithExport())
```

### GenerateFile()

Generates TypeScript file.

```go
func (g *TypeGen) GenerateFile(filename string) error
```

**Example:**
```go
err := gen.GenerateFile("resources/types/generated.ts")
```

### WithHeader()

Adds a header comment to generated file.

```go
func WithHeader(header string) RegisterOption
```

### WithExport()

Makes interface exported.

```go
func WithExport() RegisterOption
```

## HTMX Support

### IsHTMXRequest()

Checks if request is from HTMX.

```go
func (c *InertiaContext) IsHTMXRequest() bool
```

### GetHTMXHeaders()

Gets HTMX request headers.

```go
func (c *InertiaContext) GetHTMXHeaders() HTMXHeaders
```

**Returns:**
```go
type HTMXHeaders struct {
    Request      bool
    Trigger      string
    TriggerName  string
    Target       string
    CurrentURL   string
    Boosted      bool
    HistoryRestore bool
}
```

### HTMXPartial()

Renders HTML partial for HTMX.

```go
func (c *InertiaContext) HTMXPartial(html string) error
```

**Example:**
```go
return c.HTMXPartial("<div>Updated content</div>")
```

### HTMXRedirect()

HTMX redirect via HX-Redirect header.

```go
func (c *InertiaContext) HTMXRedirect(url string) error
```

### HTMXRefresh()

Triggers full page refresh.

```go
func (c *InertiaContext) HTMXRefresh() error
```

### HTMXTrigger()

Triggers HTMX events.

```go
func (c *InertiaContext) HTMXTrigger(events map[string]interface{}) *InertiaContext
```

**Example:**
```go
c.HTMXTrigger(map[string]interface{}{
    "showNotification": map[string]string{
        "message": "Saved!",
        "level": "success",
    },
})
```

### HTMXReswap()

Changes swap behavior.

```go
func (c *InertiaContext) HTMXReswap(strategy string) *InertiaContext
```

**Strategies:** innerHTML, outerHTML, beforebegin, afterbegin, beforeend, afterend, delete, none

### HTMXRetarget()

Changes target element.

```go
func (c *InertiaContext) HTMXRetarget(selector string) *InertiaContext
```

### HTMXPushURL()

Pushes URL to history.

```go
func (c *InertiaContext) HTMXPushURL(url string) *InertiaContext
```

### HTMXReplaceURL()

Replaces current URL.

```go
func (c *InertiaContext) HTMXReplaceURL(url string) *InertiaContext
```

## Vue Components

### Link

Navigation component for Inertia links.

```vue
<Link href="/users" method="get" :data="{ search: query }">
  View Users
</Link>
```

**Props:**
- `href` (string): Target URL
- `method` (string): HTTP method (default: "get")
- `data` (object): Request data
- `headers` (object): Additional headers
- `replace` (boolean): Replace history
- `preserveScroll` (boolean): Maintain scroll position
- `preserveState` (boolean): Preserve component state
- `only` (string[]): Partial reload props

### Head

Manages document head.

```vue
<Head title="Dashboard">
  <meta name="description" content="User dashboard">
</Head>
```

**Props:**
- `title` (string): Page title

## Vue Composables

### usePage()

Access current page data.

```typescript
const page = usePage<PageProps>()
```

**Returns:**
```typescript
{
  component: string
  props: PageProps
  url: string
  version: string | null
}
```

### useForm()

Form helper with validation.

```typescript
const form = useForm({
  email: '',
  password: ''
})

form.post('/login')
```

**Methods:**
- `get(url, options)` - GET request
- `post(url, options)` - POST request
- `put(url, options)` - PUT request
- `patch(url, options)` - PATCH request
- `delete(url, options)` - DELETE request
- `submit(method, url, options)` - Generic submit
- `reset()` - Reset form
- `clearErrors()` - Clear errors
- `setError(field, message)` - Set field error

**Properties:**
- `data` - Form data (reactive)
- `errors` - Validation errors
- `processing` - Submission state
- `progress` - Upload progress
- `wasSuccessful` - Success state
- `recentlySuccessful` - Recent success (2s)
- `hasErrors` - Has validation errors
- `isDirty` - Form has changes

### useRemember()

Remember component state.

```typescript
const filters = useRemember({ search: '' }, 'filters')
```

**Parameters:**
- `initialValue` (T): Initial value
- `key` (string): Storage key
- `storage` ('local' | 'session'): Storage type (default: 'session')

**Returns:** `Ref<T>`

## Types

### Props

```go
type Props map[string]interface{}
```

### SharedDataFunc

```go
type SharedDataFunc func(*http.Request) interface{}
```

### Config

```go
type Config struct {
    Version    string
    SharedData map[string]interface{}
}
```

### Page

```go
type Page struct {
    Component string                 `json:"component"`
    Props     map[string]interface{} `json:"props"`
    URL       string                 `json:"url"`
    Version   string                 `json:"version"`
}
```

## Error Handling

All methods return `error` which should be checked:

```go
if err := c.Render("Users/Index", props); err != nil {
    log.Printf("Render error: %v", err)
    http.Error(w, "Internal Server Error", 500)
    return
}
```

Common errors:
- Template parsing errors
- JSON encoding errors
- Response write errors
- Invalid configuration

## Best Practices

1. **Always check errors** returned by Render methods
2. **Use SharedDataFunc** for expensive computations
3. **Leverage partial reloads** for better UX
4. **Type your page props** with TypeScript generation
5. **Use flash messages** for user feedback
6. **Implement proper validation** with WithErrors
7. **Use HTMX** for simple interactions
8. **Remember component state** for better UX

## See Also

- [Getting Started Guide](getting-started.md)
- [Advanced Features](ADVANCED.md)
- [HTMX Integration](HTMX.md)
- [Migration Guide](MIGRATION.md)
