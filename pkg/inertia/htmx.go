package inertia

import (
	"encoding/json"
	"net/http"
)

const htmxTrueValue = "true"

// HTMXHeaders contains request headers sent by HTMX.
type HTMXHeaders struct {
	Request        bool   // HX-Request
	Target         string // HX-Target
	Trigger        string // HX-Trigger
	TriggerName    string // HX-Trigger-Name
	CurrentURL     string // HX-Current-URL
	Boosted        bool   // HX-Boosted
	HistoryRestore bool   // HX-History-Restore-Request
	Prompt         string // HX-Prompt
}

// IsHTMXRequest checks if the request is from HTMX.
func IsHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == htmxTrueValue
}

// GetHTMXHeaders extracts HTMX-specific headers from the request.
func GetHTMXHeaders(r *http.Request) HTMXHeaders {
	return HTMXHeaders{
		Request:        r.Header.Get("HX-Request") == htmxTrueValue,
		Target:         r.Header.Get("HX-Target"),
		Trigger:        r.Header.Get("HX-Trigger"),
		TriggerName:    r.Header.Get("HX-Trigger-Name"),
		CurrentURL:     r.Header.Get("HX-Current-URL"),
		Boosted:        r.Header.Get("HX-Boosted") == htmxTrueValue,
		HistoryRestore: r.Header.Get("HX-History-Restore-Request") == htmxTrueValue,
		Prompt:         r.Header.Get("HX-Prompt"),
	}
}

// HTMXRedirect sends an HTMX redirect response.
func (ic *InertiaContext) HTMXRedirect(url string) error {
	res := ic.ctx.Response()
	res.Header().Set("HX-Redirect", url)
	res.WriteHeader(http.StatusOK)
	return nil
}

// HTMXTrigger triggers a client-side event.
func (ic *InertiaContext) HTMXTrigger(event string) error {
	res := ic.ctx.Response()
	res.Header().Set("HX-Trigger", event)
	return nil
}

// HTMXTriggerWithData triggers a client-side event with data.
func (ic *InertiaContext) HTMXTriggerWithData(data map[string]interface{}) error {
	res := ic.ctx.Response()
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	res.Header().Set("HX-Trigger", string(jsonData))
	return nil
}

// HTMXPartial renders an HTML partial for HTMX.
func (ic *InertiaContext) HTMXPartial(html string) error {
	res := ic.ctx.Response()
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	_, err := res.Write([]byte(html))
	return err
}

// HTMXReswap changes the swap strategy.
func (ic *InertiaContext) HTMXReswap(strategy string) *InertiaContext {
	res := ic.ctx.Response()
	res.Header().Set("HX-Reswap", strategy)
	return ic
}

// HTMXRetarget changes the target element.
func (ic *InertiaContext) HTMXRetarget(target string) *InertiaContext {
	res := ic.ctx.Response()
	res.Header().Set("HX-Retarget", target)
	return ic
}

// HTMXPushURL pushes a new URL to the browser history.
func (ic *InertiaContext) HTMXPushURL(url string) *InertiaContext {
	res := ic.ctx.Response()
	res.Header().Set("HX-Push-Url", url)
	return ic
}

// HTMXReplaceURL replaces the current URL in browser history.
func (ic *InertiaContext) HTMXReplaceURL(url string) *InertiaContext {
	res := ic.ctx.Response()
	res.Header().Set("HX-Replace-Url", url)
	return ic
}

// HTMXRefresh triggers a client-side page refresh.
func (ic *InertiaContext) HTMXRefresh() error {
	res := ic.ctx.Response()
	res.Header().Set("HX-Refresh", htmxTrueValue)
	res.WriteHeader(http.StatusOK)
	return nil
}
