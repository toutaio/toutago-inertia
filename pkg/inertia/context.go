package inertia

import (
	"encoding/json"
	"net/http"
)

// ContextInterface defines the minimal interface that any router context must implement.
type ContextInterface interface {
	Request() *http.Request
	Response() http.ResponseWriter
	Set(key string, value interface{})
	Get(key string) interface{}
}

// InertiaContext wraps a router context and provides Inertia-specific methods.
//
//nolint:revive // InertiaContext name is intentional for clarity in Inertia-specific context.
type InertiaContext struct {
	ctx           ContextInterface
	mgr           *Inertia
	sharedData    map[string]interface{}
	sharedFuncs   map[string]SharedDataFunc
	pendingErrors ValidationErrors
	pendingFlash  Flash
}

// NewContext creates a new Inertia context wrapper.
func NewContext(ctx ContextInterface, mgr *Inertia) *InertiaContext {
	return &InertiaContext{
		ctx:         ctx,
		mgr:         mgr,
		sharedData:  make(map[string]interface{}),
		sharedFuncs: make(map[string]SharedDataFunc),
	}
}

// Share adds context-specific shared data.
func (ic *InertiaContext) Share(key string, value interface{}) *InertiaContext {
	ic.sharedData[key] = value
	return ic
}

// ShareFunc adds context-specific lazy shared data.
func (ic *InertiaContext) ShareFunc(key string, fn SharedDataFunc) *InertiaContext {
	ic.sharedFuncs[key] = fn
	return ic
}

// WithErrors adds validation errors to the next render.
func (ic *InertiaContext) WithErrors(errors ValidationErrors) *InertiaContext {
	ic.pendingErrors = errors
	return ic
}

// WithFlash adds flash messages to the next render.
func (ic *InertiaContext) WithFlash(flash Flash) *InertiaContext {
	ic.pendingFlash = flash
	return ic
}

// Render renders an Inertia page with context-specific data.
func (ic *InertiaContext) Render(component string, props map[string]interface{}) error {
	req := ic.ctx.Request()
	res := ic.ctx.Response()

	// Get partial reload info
	only := GetPartialOnly(req)

	// Merge context-specific shared data into props first
	// (before filtering for partial reloads)
	for key, value := range ic.sharedData {
		if _, exists := props[key]; !exists {
			props[key] = value
		}
	}

	// Add context-specific lazy shared data
	for key, fn := range ic.sharedFuncs {
		if _, exists := props[key]; !exists {
			props[key] = fn()
		}
	}

	var page *Page
	var err error

	if len(only) > 0 {
		// Partial reload
		page, err = ic.mgr.RenderOnly(component, props, req.URL.Path, only)
	} else {
		// Full page load
		page, err = ic.mgr.Render(component, props, req.URL.Path)
	}

	if err != nil {
		return err
	}

	// Add pending errors
	if ic.pendingErrors != nil {
		page.WithErrors(ic.pendingErrors)
		ic.pendingErrors = nil
	}

	// Add pending flash messages
	if ic.pendingFlash != nil {
		page.WithFlash(ic.pendingFlash)
		ic.pendingFlash = nil
	}

	// Send JSON response
	res.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(res).Encode(page)
}

// Redirect performs an internal redirect.
func (ic *InertiaContext) Redirect(url string) error {
	return ic.mgr.Redirect(ic.ctx.Response(), ic.ctx.Request(), url)
}

// Location performs an external redirect.
func (ic *InertiaContext) Location(url string) error {
	return ic.mgr.Location(ic.ctx.Response(), ic.ctx.Request(), url)
}

// Back redirects to the previous page.
func (ic *InertiaContext) Back() error {
	return ic.mgr.Back(ic.ctx.Response(), ic.ctx.Request())
}

// Error renders an error page.
func (ic *InertiaContext) Error(status int, message string) error {
	page, err := ic.mgr.Error(status, message, ic.ctx.Request().URL.Path, ic.ctx.Request())
	if err != nil {
		return err
	}

	res := ic.ctx.Response()
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	return json.NewEncoder(res).Encode(page)
}
