package inertia

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Response represents an Inertia.js page response.
type Response struct {
	Component string                 `json:"component"`
	Props     map[string]interface{} `json:"props"`
	URL       string                 `json:"url"`
	Version   string                 `json:"version"`
}

// MarshalJSON implements json.Marshaler.
func (r Response) MarshalJSON() ([]byte, error) {
	type Alias Response
	return json.Marshal(Alias(r))
}

// Page represents an Inertia page with all data.
type Page struct {
	Component string                 `json:"component"`
	Props     map[string]interface{} `json:"props"`
	URL       string                 `json:"url"`
	Version   string                 `json:"version"`
}

// NewPage creates a new Inertia page.
func NewPage(component string, props map[string]interface{}, url, version string) *Page {
	if props == nil {
		props = make(map[string]interface{})
	}
	return &Page{
		Component: component,
		Props:     props,
		URL:       url,
		Version:   version,
	}
}

// MergeSharedData merges shared data into the page props.
func (p *Page) MergeSharedData(shared map[string]interface{}) {
	for key, value := range shared {
		// Don't overwrite existing props
		if _, exists := p.Props[key]; !exists {
			p.Props[key] = value
		}
	}
}

// Config holds Inertia configuration.
type Config struct {
	RootView string // Path to root template
	Version  string // Asset version
	SSR      bool   // Enable server-side rendering
	AssetURL string // Base URL for assets
}

// Validate checks if the config is valid.
func (c Config) Validate() error {
	if c.RootView == "" {
		return errors.New("inertia: RootView is required")
	}
	return nil
}

// SharedDataFunc is a function that returns shared data.
type SharedDataFunc func() interface{}

// Inertia is the main Inertia instance.
type Inertia struct {
	config     Config
	version    string
	sharedData map[string]interface{}
	sharedFunc map[string]SharedDataFunc
}

// New creates a new Inertia instance.
func New(config Config) (*Inertia, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	version := config.Version
	if version == "" {
		version = "1" // Default version
	}

	return &Inertia{
		config:     config,
		version:    version,
		sharedData: make(map[string]interface{}),
		sharedFunc: make(map[string]SharedDataFunc),
	}, nil
}

// Share adds a static shared value
func (i *Inertia) Share(key string, value interface{}) {
	i.sharedData[key] = value
}

// ShareFunc adds a function that provides shared data
func (i *Inertia) ShareFunc(key string, fn SharedDataFunc) {
	i.sharedFunc[key] = fn
}

// GetSharedData returns all shared data (static + evaluated functions)
func (i *Inertia) GetSharedData() map[string]interface{} {
	result := make(map[string]interface{})

	// Add static shared data
	for key, value := range i.sharedData {
		result[key] = value
	}

	// Evaluate and add function-based shared data
	for key, fn := range i.sharedFunc {
		result[key] = fn()
	}

	return result
}

// Version returns the current asset version
func (i *Inertia) Version() string {
	return i.version
}

// SetVersion updates the asset version
func (i *Inertia) SetVersion(version string) {
	i.version = version
}

// Render creates an Inertia response
func (i *Inertia) Render(component string, props map[string]interface{}, url string) (*Page, error) {
	if component == "" {
		return nil, fmt.Errorf("inertia: component name is required")
	}

	if url == "" {
		return nil, fmt.Errorf("inertia: URL is required")
	}

	if props == nil {
		props = make(map[string]interface{})
	}

	page := NewPage(component, props, url, i.version)
	page.MergeSharedData(i.GetSharedData())

	return page, nil
}

// RenderOnly creates an Inertia response with only specified props
func (i *Inertia) RenderOnly(component string, props map[string]interface{}, url string, only []string) (*Page, error) {
	if component == "" {
		return nil, fmt.Errorf("inertia: component name is required")
	}

	if url == "" {
		return nil, fmt.Errorf("inertia: URL is required")
	}

	if props == nil {
		props = make(map[string]interface{})
	}

	// Filter props to only include requested ones
	filteredProps := make(map[string]interface{})
	for _, key := range only {
		if val, ok := props[key]; ok {
			filteredProps[key] = val
		}
	}

	page := NewPage(component, filteredProps, url, i.version)
	// Shared data is always included
	page.MergeSharedData(i.GetSharedData())

	return page, nil
}
