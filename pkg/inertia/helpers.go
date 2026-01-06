package inertia

// Add adds a validation error for a field.
func (v ValidationErrors) Add(field, message string) {
	v[field] = append(v[field], message)
}

// Has checks if a field has any errors.
func (v ValidationErrors) Has(field string) bool {
	_, exists := v[field]
	return exists
}

// First returns the first error message for a field, or empty string if none.
func (v ValidationErrors) First(field string) string {
	if errors, exists := v[field]; exists && len(errors) > 0 {
		return errors[0]
	}
	return ""
}

// Any returns true if there are any validation errors.
func (v ValidationErrors) Any() bool {
	return len(v) > 0
}

// NewValidationErrors creates a new ValidationErrors instance.
func NewValidationErrors() ValidationErrors {
	return make(ValidationErrors)
}

// Success adds a success flash message.
func (f Flash) Success(message string) {
	f["success"] = message
}

// Error adds an error flash message.
func (f Flash) Error(message string) {
	f["error"] = message
}

// Warning adds a warning flash message.
func (f Flash) Warning(message string) {
	f["warning"] = message
}

// Info adds an info flash message.
func (f Flash) Info(message string) {
	f["info"] = message
}

// Custom adds a custom flash message with the given key.
func (f Flash) Custom(key, message string) {
	f[key] = message
}

// NewFlash creates a new Flash instance.
func NewFlash() Flash {
	return make(Flash)
}
