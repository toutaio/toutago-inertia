package inertia

import (
	"net/http"
)

// ValidationErrors represents form validation errors
type ValidationErrors map[string][]string

// Flash represents flash messages
type Flash map[string]string

// Location performs an external redirect (409 for Inertia, 302 for browsers)
func (i *Inertia) Location(w http.ResponseWriter, r *http.Request, url string) error {
	if IsInertiaRequest(r) {
		w.Header().Set("X-Inertia-Location", url)
		w.WriteHeader(http.StatusConflict)
		return nil
	}

	// Regular browser redirect
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

// Back redirects back to the previous page (using Referer header)
func (i *Inertia) Back(w http.ResponseWriter, r *http.Request) error {
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}

	return i.Location(w, r, referer)
}

// Redirect performs an internal redirect
func (i *Inertia) Redirect(w http.ResponseWriter, r *http.Request, url string) error {
	if IsInertiaRequest(r) {
		// For Inertia requests, always use 303 See Other to change method to GET
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusSeeOther)
		return nil
	}

	// Regular browser redirect
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

// Error creates an error page response
func (i *Inertia) Error(status int, message string, url string, r *http.Request) (*Page, error) {
	props := map[string]interface{}{
		"status":  status,
		"message": message,
	}

	page := NewPage("Error", props, url, i.version)
	page.MergeSharedData(i.GetSharedData())

	return page, nil
}

// WithErrors adds validation errors to the page props
func (p *Page) WithErrors(errors ValidationErrors) *Page {
	p.Props["errors"] = errors
	return p
}

// WithFlash adds flash messages to the page props
func (p *Page) WithFlash(flash Flash) *Page {
	for key, value := range flash {
		p.Props[key] = value
	}
	return p
}
