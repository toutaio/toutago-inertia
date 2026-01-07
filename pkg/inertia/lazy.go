package inertia

// LazyProp represents a lazily-evaluated property.
type LazyProp struct {
	Evaluator func() interface{}
	Group     string // "lazy", "always", or "defer"
}

// Lazy adds a lazily-evaluated prop that is excluded from partial reloads
// unless explicitly requested.
func (ic *InertiaContext) Lazy(key string, fn func() interface{}) *InertiaContext {
	if ic.ctx.Get("_inertia_lazy_props") == nil {
		ic.ctx.Set("_inertia_lazy_props", make(map[string]LazyProp))
	}
	lazyProps := ic.ctx.Get("_inertia_lazy_props").(map[string]LazyProp)
	lazyProps[key] = LazyProp{
		Evaluator: fn,
		Group:     "lazy",
	}
	return ic
}

// Always adds a prop that is always included, even in partial reloads.
func (ic *InertiaContext) Always(key string, value interface{}) *InertiaContext {
	if ic.ctx.Get("_inertia_always_props") == nil {
		ic.ctx.Set("_inertia_always_props", make(map[string]interface{}))
	}
	alwaysProps := ic.ctx.Get("_inertia_always_props").(map[string]interface{})
	alwaysProps[key] = value
	return ic
}

// AlwaysLazy adds a lazily-evaluated prop that is always included.
func (ic *InertiaContext) AlwaysLazy(key string, fn func() interface{}) *InertiaContext {
	if ic.ctx.Get("_inertia_lazy_props") == nil {
		ic.ctx.Set("_inertia_lazy_props", make(map[string]LazyProp))
	}
	lazyProps := ic.ctx.Get("_inertia_lazy_props").(map[string]LazyProp)
	lazyProps[key] = LazyProp{
		Evaluator: fn,
		Group:     "always",
	}
	return ic
}

// Defer adds a prop that is never included unless explicitly requested.
// Useful for expensive computations that should only load on demand.
func (ic *InertiaContext) Defer(key string, fn func() interface{}) *InertiaContext {
	if ic.ctx.Get("_inertia_lazy_props") == nil {
		ic.ctx.Set("_inertia_lazy_props", make(map[string]LazyProp))
	}
	lazyProps := ic.ctx.Get("_inertia_lazy_props").(map[string]LazyProp)
	lazyProps[key] = LazyProp{
		Evaluator: fn,
		Group:     "defer",
	}
	return ic
}

// evaluateLazyProps evaluates lazy props based on the request type.
func (ic *InertiaContext) evaluateLazyProps(props map[string]interface{}, only []string) {
	ic.mergeAlwaysProps(props)

	lazyProps := ic.getLazyPropsFromContext()
	if lazyProps == nil {
		return
	}

	isPartial := len(only) > 0
	for key, lazyProp := range lazyProps {
		if ic.shouldEvaluateLazyProp(key, lazyProp, isPartial, only) {
			ic.evaluatePropIfNotExists(props, key, lazyProp)
		}
	}
}

// getLazyPropsFromContext retrieves lazy props from the context.
func (ic *InertiaContext) getLazyPropsFromContext() map[string]LazyProp {
	lazyPropsInterface := ic.ctx.Get("_inertia_lazy_props")
	if lazyPropsInterface == nil {
		return nil
	}
	return lazyPropsInterface.(map[string]LazyProp)
}

// mergeAlwaysProps merges always props into the props map.
func (ic *InertiaContext) mergeAlwaysProps(props map[string]interface{}) {
	alwaysPropsInterface := ic.ctx.Get("_inertia_always_props")
	if alwaysPropsInterface == nil {
		return
	}

	alwaysProps := alwaysPropsInterface.(map[string]interface{})
	for key, value := range alwaysProps {
		if _, exists := props[key]; !exists {
			props[key] = value
		}
	}
}

// shouldEvaluateLazyProp determines if a lazy prop should be evaluated.
func (ic *InertiaContext) shouldEvaluateLazyProp(
	key string,
	lazyProp LazyProp,
	isPartial bool,
	only []string,
) bool {
	switch lazyProp.Group {
	case "always":
		return true
	case "lazy":
		return ic.shouldEvaluateLazyGroup(key, isPartial, only)
	case "defer":
		return ic.shouldEvaluateDeferGroup(key, isPartial, only)
	default:
		return false
	}
}

// shouldEvaluateLazyGroup determines if a "lazy" group prop should be evaluated.
func (ic *InertiaContext) shouldEvaluateLazyGroup(key string, isPartial bool, only []string) bool {
	if !isPartial {
		return true
	}
	return ic.isKeyRequested(key, only)
}

// shouldEvaluateDeferGroup determines if a "defer" group prop should be evaluated.
func (ic *InertiaContext) shouldEvaluateDeferGroup(key string, isPartial bool, only []string) bool {
	if !isPartial {
		return false
	}
	return ic.isKeyRequested(key, only)
}

// isKeyRequested checks if a key is in the requested keys list.
func (ic *InertiaContext) isKeyRequested(key string, only []string) bool {
	for _, requestedKey := range only {
		if requestedKey == key {
			return true
		}
	}
	return false
}

// evaluatePropIfNotExists evaluates a lazy prop if it doesn't already exist.
func (ic *InertiaContext) evaluatePropIfNotExists(
	props map[string]interface{},
	key string,
	lazyProp LazyProp,
) {
	if _, exists := props[key]; !exists {
		props[key] = lazyProp.Evaluator()
	}
}
