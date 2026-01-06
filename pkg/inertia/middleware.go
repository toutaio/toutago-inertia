package inertia

import (
	"context"
	"net/http"
	"strings"
)

// Context keys for storing request data.
type contextKey string

const (
	contextKeyInertia          contextKey = "inertia"
	contextKeyPartialOnly      contextKey = "partial_only"
	contextKeyPartialComponent contextKey = "partial_component"
	contextKeyExternalRedirect contextKey = "external_redirect"
)

// Middleware returns an HTTP middleware that handles Inertia requests.
func (i *Inertia) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Always set version header
			w.Header().Set("X-Inertia-Version", i.version)

			// Check if this is an Inertia request
			isInertia := IsInertiaRequest(r)

			if isInertia {
				// Store Inertia flag in context
				ctx := context.WithValue(r.Context(), contextKeyInertia, true)

				// Check version match
				clientVersion := r.Header.Get("X-Inertia-Version")
				if clientVersion != "" && clientVersion != i.version {
					// Version mismatch - force reload
					w.WriteHeader(http.StatusConflict)
					return
				}

				// Handle partial reloads
				if partialData := r.Header.Get("X-Inertia-Partial-Data"); partialData != "" {
					only := strings.Split(partialData, ",")
					for i := range only {
						only[i] = strings.TrimSpace(only[i])
					}
					ctx = context.WithValue(ctx, contextKeyPartialOnly, only)
				}

				if partialComponent := r.Header.Get("X-Inertia-Partial-Component"); partialComponent != "" {
					ctx = context.WithValue(ctx, contextKeyPartialComponent, partialComponent)
				}

				r = r.WithContext(ctx)
			}

			// Wrap response writer to intercept status code
			wrapped := &responseWriter{ResponseWriter: w, request: r}

			// Call next handler
			next.ServeHTTP(wrapped, r)

			// Check for external redirect after handler (but before response is written)
			if isInertia && !wrapped.written {
				if location := GetExternalRedirect(r); location != "" {
					w.Header().Set("X-Inertia-Location", location)
					w.WriteHeader(http.StatusConflict)
					return
				}
			}
		})
	}
}

// responseWriter wraps http.ResponseWriter to track if response was written.
type responseWriter struct {
	http.ResponseWriter
	request *http.Request
	written bool
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.written = true
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(b)
}

// IsInertiaRequest checks if the request is an Inertia request.
func IsInertiaRequest(r *http.Request) bool {
	value := r.Header.Get("X-Inertia")
	return strings.EqualFold(value, "true")
}

// GetPartialOnly returns the list of props to include in partial reload.
func GetPartialOnly(r *http.Request) []string {
	if only, ok := r.Context().Value(contextKeyPartialOnly).([]string); ok {
		return only
	}
	return nil
}

// GetPartialComponent returns the component name for partial reload.
func GetPartialComponent(r *http.Request) string {
	if component, ok := r.Context().Value(contextKeyPartialComponent).(string); ok {
		return component
	}
	return ""
}

// SetExternalRedirect marks the request for external redirect.
func SetExternalRedirect(r *http.Request, url string) {
	ctx := context.WithValue(r.Context(), contextKeyExternalRedirect, url)
	*r = *r.WithContext(ctx)
}

// GetExternalRedirect gets the external redirect URL if set.
func GetExternalRedirect(r *http.Request) string {
	if url, ok := r.Context().Value(contextKeyExternalRedirect).(string); ok {
		return url
	}
	return ""
}
