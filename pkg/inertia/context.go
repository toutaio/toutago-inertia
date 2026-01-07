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

	only := GetPartialOnly(req)
	only = ic.appendAlwaysProps(only)

	ic.mergeSharedData(props)
	ic.evaluateLazyProps(props, only)

	page, err := ic.renderPage(component, props, req.URL.Path, only)
	if err != nil {
		return err
	}

	ic.attachPendingData(page)

	res.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(res).Encode(page)
}

// appendAlwaysProps adds "always" props to the only list for partial reloads.
func (ic *InertiaContext) appendAlwaysProps(only []string) []string {
	if len(only) == 0 {
		return only
	}

	only = ic.appendAlwaysRegularProps(only)
	only = ic.appendAlwaysLazyProps(only)
	return only
}

// appendAlwaysRegularProps appends always-included regular props to the only list.
func (ic *InertiaContext) appendAlwaysRegularProps(only []string) []string {
	alwaysPropsInterface := ic.ctx.Get("_inertia_always_props")
	if alwaysPropsInterface == nil {
		return only
	}

	alwaysProps := alwaysPropsInterface.(map[string]interface{})
	for key := range alwaysProps {
		only = append(only, key)
	}
	return only
}

// appendAlwaysLazyProps appends always-included lazy props to the only list.
func (ic *InertiaContext) appendAlwaysLazyProps(only []string) []string {
	lazyPropsInterface := ic.ctx.Get("_inertia_lazy_props")
	if lazyPropsInterface == nil {
		return only
	}

	lazyProps := lazyPropsInterface.(map[string]LazyProp)
	for key, lazyProp := range lazyProps {
		if lazyProp.Group == "always" {
			only = append(only, key)
		}
	}
	return only
}

// mergeSharedData merges context-specific shared data and lazy functions into props.
func (ic *InertiaContext) mergeSharedData(props map[string]interface{}) {
	for key, value := range ic.sharedData {
		if _, exists := props[key]; !exists {
			props[key] = value
		}
	}

	for key, fn := range ic.sharedFuncs {
		if _, exists := props[key]; !exists {
			props[key] = fn()
		}
	}
}

// renderPage renders the page based on whether it's a partial or full reload.
func (ic *InertiaContext) renderPage(
	component string,
	props map[string]interface{},
	path string,
	only []string,
) (*Page, error) {
	if len(only) > 0 {
		return ic.mgr.RenderOnly(component, props, path, only)
	}
	return ic.mgr.Render(component, props, path)
}

// attachPendingData attaches pending errors and flash messages to the page.
func (ic *InertiaContext) attachPendingData(page *Page) {
	if ic.pendingErrors != nil {
		page.WithErrors(ic.pendingErrors)
		ic.pendingErrors = nil
	}

	if ic.pendingFlash != nil {
		page.WithFlash(ic.pendingFlash)
		ic.pendingFlash = nil
	}
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

// WithError adds a single validation error for a field.
func (ic *InertiaContext) WithError(field, message string) *InertiaContext {
	if ic.pendingErrors == nil {
		ic.pendingErrors = NewValidationErrors()
	}
	ic.pendingErrors.Add(field, message)
	return ic
}

// WithSuccess adds a success flash message.
func (ic *InertiaContext) WithSuccess(message string) *InertiaContext {
	if ic.pendingFlash == nil {
		ic.pendingFlash = NewFlash()
	}
	ic.pendingFlash.Success(message)
	return ic
}

// WithErrorMessage adds an error flash message.
func (ic *InertiaContext) WithErrorMessage(message string) *InertiaContext {
	if ic.pendingFlash == nil {
		ic.pendingFlash = NewFlash()
	}
	ic.pendingFlash.Error(message)
	return ic
}

// WithWarning adds a warning flash message.
func (ic *InertiaContext) WithWarning(message string) *InertiaContext {
	if ic.pendingFlash == nil {
		ic.pendingFlash = NewFlash()
	}
	ic.pendingFlash.Warning(message)
	return ic
}

// WithInfo adds an info flash message.
func (ic *InertiaContext) WithInfo(message string) *InertiaContext {
	if ic.pendingFlash == nil {
		ic.pendingFlash = NewFlash()
	}
	ic.pendingFlash.Info(message)
	return ic
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
