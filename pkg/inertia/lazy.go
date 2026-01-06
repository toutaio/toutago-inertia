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
	// Get lazy props from context
	lazyPropsInterface := ic.ctx.Get("_inertia_lazy_props")
	if lazyPropsInterface == nil {
		return
	}
	lazyProps := lazyPropsInterface.(map[string]LazyProp)

	// Get always props from context
	alwaysPropsInterface := ic.ctx.Get("_inertia_always_props")
	if alwaysPropsInterface != nil {
		alwaysProps := alwaysPropsInterface.(map[string]interface{})
		for key, value := range alwaysProps {
			if _, exists := props[key]; !exists {
				props[key] = value
			}
		}
	}

	// Determine if this is a partial reload
	isPartial := len(only) > 0

	for key, lazyProp := range lazyProps {
		shouldEvaluate := false

		switch lazyProp.Group {
		case "always":
			// Always evaluate, regardless of partial reload
			shouldEvaluate = true

		case "lazy":
			if !isPartial {
				// Evaluate on full page load
				shouldEvaluate = true
			} else {
				// Only evaluate if explicitly requested in partial reload
				for _, requestedKey := range only {
					if requestedKey == key {
						shouldEvaluate = true
						break
					}
				}
			}

		case "defer":
			// Only evaluate if explicitly requested (never on full load)
			if isPartial {
				for _, requestedKey := range only {
					if requestedKey == key {
						shouldEvaluate = true
						break
					}
				}
			}
		}

		if shouldEvaluate {
			if _, exists := props[key]; !exists {
				props[key] = lazyProp.Evaluator()
			}
		}
	}
}
