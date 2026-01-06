# HTMX Integration Guide

## Overview

toutago-inertia provides comprehensive HTMX support, allowing you to build hypermedia-driven applications alongside or instead of traditional SPAs. This guide shows you how to use HTMX features effectively.

## Quick Start

### Detecting HTMX Requests

```go
import "github.com/toutaio/toutago-inertia/pkg/inertia"

func handler(w http.ResponseWriter, r *http.Request) {
    if inertia.IsHTMXRequest(r) {
        // This is an HTMX request
        // Render partial HTML instead of full page
    } else {
        // Regular request - render full page
    }
}
```

### Getting HTMX Headers

```go
headers := inertia.GetHTMXHeaders(r)

fmt.Println("Target:", headers.Target)         // HX-Target
fmt.Println("Trigger:", headers.Trigger)       // HX-Trigger
fmt.Println("Current URL:", headers.CurrentURL) // HX-Current-URL
fmt.Println("Boosted:", headers.Boosted)       // HX-Boosted
```

## Response Helpers

### HTMXRedirect - Client-Side Redirects

Redirect the browser without a full page reload:

```go
func updateHandler(ctx *YourContext) error {
    ic := inertia.NewContext(ctx, inertiaManager)
    
    // Save data...
    
    // Redirect to dashboard
    return ic.HTMXRedirect("/dashboard")
}
```

### HTMXTrigger - Client-Side Events

Trigger client-side events that your JavaScript can listen for:

```go
func deleteHandler(ctx *YourContext) error {
    ic := inertia.NewContext(ctx, inertiaManager)
    
    // Delete item...
    
    // Trigger event to update UI
    return ic.HTMXTrigger("itemDeleted")
}
```

Client-side:

```javascript
document.body.addEventListener("itemDeleted", function(evt) {
    showNotification("Item deleted successfully");
});
```

### HTMXTriggerWithData - Events with Payload

Send structured data with your events:

```go
func saveHandler(ctx *YourContext) error {
    ic := inertia.NewContext(ctx, inertiaManager)
    
    // Save item...
    
    data := map[string]interface{}{
        "showMessage": map[string]string{
            "level":   "success",
            "message": "Item saved successfully",
        },
        "updateCount": map[string]interface{}{
            "count": newCount,
        },
    }
    
    return ic.HTMXTriggerWithData(data)
}
```

Client-side:

```javascript
document.body.addEventListener("showMessage", function(evt) {
    const detail = evt.detail;
    showNotification(detail.level, detail.message);
});

document.body.addEventListener("updateCount", function(evt) {
    document.getElementById("count").textContent = evt.detail.count;
});
```

### HTMXPartial - Render HTML Fragments

Return HTML fragments for HTMX to swap into the page:

```go
func listItemsHandler(ctx *YourContext) error {
    ic := inertia.NewContext(ctx, inertiaManager)
    
    items := getItems()
    
    // Render partial HTML
    html := renderItemList(items) // "<ul><li>Item 1</li>...</ul>"
    
    return ic.HTMXPartial(html)
}
```

HTML:

```html
<div id="items" hx-get="/items" hx-trigger="load">
    Loading...
</div>
```

## Advanced Features

### HTMXReswap - Change Swap Strategy

Control how HTMX swaps content:

```go
func updateHandler(ctx *YourContext) error {
    ic := inertia.NewContext(ctx, inertiaManager)
    
    html := "<div>New content</div>"
    
    // Use outerHTML to replace the entire element
    ic.HTMXReswap("outerHTML")
    
    return ic.HTMXPartial(html)
}
```

Swap strategies:
- `innerHTML` - Replace inner content (default)
- `outerHTML` - Replace entire element
- `beforebegin` - Insert before element
- `afterbegin` - Insert at start of element
- `beforeend` - Insert at end of element
- `afterend` - Insert after element
- `delete` - Delete element
- `none` - Don't swap

### HTMXRetarget - Dynamic Target Selection

Change which element receives the response:

```go
func errorHandler(ctx *YourContext) error {
    ic := inertia.NewContext(ctx, inertiaManager)
    
    // Retarget to error container
    ic.HTMXRetarget("#error-messages")
    
    return ic.HTMXPartial("<div class='error'>Something went wrong</div>")
}
```

### HTMXPushURL - Browser History

Update the browser URL and history:

```go
func filterHandler(ctx *YourContext) error {
    ic := inertia.NewContext(ctx, inertiaManager)
    
    // Push new URL to history
    ic.HTMXPushURL("/items?filter=active")
    
    html := renderFilteredItems()
    return ic.HTMXPartial(html)
}
```

### HTMXReplaceURL - Update Without History

Replace the current URL without adding to history:

```go
func sortHandler(ctx *YourContext) error {
    ic := inertia.NewContext(ctx, inertiaManager)
    
    // Update URL without history entry
    ic.HTMXReplaceURL("/items?sort=date")
    
    html := renderSortedItems()
    return ic.HTMXPartial(html)
}
```

### HTMXRefresh - Force Page Reload

Trigger a full page refresh:

```go
func resetHandler(ctx *YourContext) error {
    ic := inertia.NewContext(ctx, inertiaManager)
    
    // Clear all state...
    
    // Force page refresh
    return ic.HTMXRefresh()
}
```

## Chaining Helpers

Many HTMX helpers return `*InertiaContext`, allowing you to chain them:

```go
func updateHandler(ctx *YourContext) error {
    ic := inertia.NewContext(ctx, inertiaManager)
    
    html := "<div>Updated content</div>"
    
    return ic.
        HTMXReswap("outerHTML").
        HTMXRetarget("#content").
        HTMXPushURL("/updated").
        HTMXPartial(html)
}
```

## Complete Example

Here's a complete example combining multiple HTMX features:

```go
package main

import (
    "net/http"
    "github.com/toutaio/toutago-inertia/pkg/inertia"
)

type Item struct {
    ID   int
    Name string
}

func main() {
    config := inertia.Config{
        RootView: "app.html",
        Version:  "1.0.0",
    }
    
    mgr, _ := inertia.New(config)
    
    http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
        ctx := &MyContext{req: r, res: w}
        ic := inertia.NewContext(ctx, mgr)
        
        if inertia.IsHTMXRequest(r) {
            // HTMX request - return partial
            items := []Item{{1, "Item 1"}, {2, "Item 2"}}
            html := renderItemsPartial(items)
            ic.HTMXPartial(html)
        } else {
            // Regular request - return full Inertia page
            items := []Item{{1, "Item 1"}, {2, "Item 2"}}
            ic.Render("Items/Index", map[string]interface{}{
                "items": items,
            })
        }
    })
    
    http.HandleFunc("/items/create", func(w http.ResponseWriter, r *http.Request) {
        ctx := &MyContext{req: r, res: w}
        ic := inertia.NewContext(ctx, mgr)
        
        // Create item...
        newItem := Item{ID: 3, Name: "New Item"}
        
        // Trigger event and redirect
        ic.HTMXTrigger("itemCreated")
        ic.HTMXRedirect("/items")
    })
    
    http.ListenAndServe(":3000", nil)
}

func renderItemsPartial(items []Item) string {
    html := "<ul>"
    for _, item := range items {
        html += fmt.Sprintf("<li>%s</li>", item.Name)
    }
    html += "</ul>"
    return html
}
```

HTML Template:

```html
<!DOCTYPE html>
<html>
<head>
    <title>HTMX + Inertia</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
    <h1>Items</h1>
    
    <!-- HTMX-powered list -->
    <div id="items" hx-get="/items" hx-trigger="load">
        Loading...
    </div>
    
    <!-- Create form -->
    <form hx-post="/items/create" hx-target="#items">
        <input type="text" name="name" />
        <button type="submit">Add Item</button>
    </form>
    
    <script>
        // Listen for custom events
        document.body.addEventListener("itemCreated", function() {
            console.log("Item created!");
        });
    </script>
</body>
</html>
```

## Mixing Inertia.js and HTMX

You can use both Inertia.js and HTMX in the same application:

```go
func handler(ctx *YourContext) error {
    ic := inertia.NewContext(ctx, mgr)
    
    // Check request type
    if inertia.IsHTMXRequest(ctx.Request()) {
        // Simple interaction - use HTMX
        return ic.HTMXPartial("<div>Quick update</div>")
    } else if ctx.Request().Header.Get("X-Inertia") != "" {
        // Complex SPA view - use Inertia
        return ic.Render("Dashboard/Index", props)
    } else {
        // First page load - render Inertia root
        return ic.Render("Dashboard/Index", props)
    }
}
```

**Use HTMX for:**
- Simple interactions
- Partial page updates
- Progressive enhancement
- Forms and validation

**Use Inertia.js for:**
- Full SPA experiences
- Complex client-side routing
- Rich interactive UIs
- Vue/React components

## Best Practices

1. **Progressive Enhancement**: Start with HTMX for simple interactions, upgrade to Inertia for complex features
2. **Consistent Responses**: Check request type and return appropriate format
3. **Error Handling**: Use `HTMXRetarget` to show errors in specific containers
4. **Events**: Use `HTMXTrigger` for loose coupling between components
5. **URL Management**: Use `HTMXPushURL` for meaningful URL changes, `HTMXReplaceURL` for filtering/sorting

## Resources

- [HTMX Documentation](https://htmx.org/)
- [Inertia.js Documentation](https://inertiajs.com/)
- [toutago-inertia GitHub](https://github.com/toutaio/toutago-inertia)
