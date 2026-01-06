package ssr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"rogchap.com/v8go"
)

type Config struct {
	PoolSize int
	Timeout  time.Duration
}

type Renderer struct {
	config *Config
	iso    *v8go.Isolate
	bundle string
	pool   chan *v8go.Context
	mu     sync.RWMutex
	closed bool
}

func NewRenderer(cfg ...*Config) (*Renderer, error) {
	config := &Config{
		PoolSize: 10,
		Timeout:  30 * time.Second,
	}
	if len(cfg) > 0 && cfg[0] != nil {
		if cfg[0].PoolSize > 0 {
			config.PoolSize = cfg[0].PoolSize
		}
		if cfg[0].Timeout > 0 {
			config.Timeout = cfg[0].Timeout
		}
	}

	iso := v8go.NewIsolate()
	r := &Renderer{
		config: config,
		iso:    iso,
		pool:   make(chan *v8go.Context, config.PoolSize),
	}

	for i := 0; i < config.PoolSize; i++ {
		ctx := v8go.NewContext(iso)
		r.pool <- ctx
	}

	return r, nil
}

func (r *Renderer) LoadBundle(bundle string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return errors.New("renderer is closed")
	}

	ctx := v8go.NewContext(r.iso)
	defer ctx.Close()

	if _, err := ctx.RunScript("var global = globalThis;", "setup.js"); err != nil {
		return fmt.Errorf("failed to setup global: %w", err)
	}

	_, err := ctx.RunScript(bundle, "bundle.js")
	if err != nil {
		return fmt.Errorf("failed to load bundle: %w", err)
	}

	r.bundle = bundle
	return nil
}

func (r *Renderer) RenderToString(ctx context.Context, pageData map[string]interface{}) (string, error) {
	r.mu.RLock()
	if r.closed {
		r.mu.RUnlock()
		return "", errors.New("renderer is closed")
	}
	r.mu.RUnlock()

	timeout := r.config.Timeout
	if deadline, ok := ctx.Deadline(); ok {
		timeout = time.Until(deadline)
	}

	resultCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		html, err := r.render(pageData)
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- html
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-errCh:
		return "", err
	case html := <-resultCh:
		return html, nil
	case <-time.After(timeout):
		return "", errors.New("render timeout")
	}
}

func (r *Renderer) render(pageData map[string]interface{}) (string, error) {
	var v8ctx *v8go.Context
	select {
	case v8ctx = <-r.pool:
		defer func() { r.pool <- v8ctx }()
	default:
		v8ctx = v8go.NewContext(r.iso)
		defer v8ctx.Close()
	}

	if _, err := v8ctx.RunScript("var global = globalThis;", "setup.js"); err != nil {
		return "", fmt.Errorf("failed to setup global: %w", err)
	}

	if r.bundle != "" {
		if _, err := v8ctx.RunScript(r.bundle, "bundle.js"); err != nil {
			return "", fmt.Errorf("failed to re-run bundle: %w", err)
		}
	}

	pageJSON, err := json.Marshal(pageData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal page data: %w", err)
	}

	script := fmt.Sprintf(`
		(function() {
			var page = %s;
			if (typeof global.render !== 'function') {
				throw new Error('render function not found');
			}
			var result = global.render(page);
			if (typeof result === 'object' && result !== null) {
				return JSON.stringify(result);
			}
			return result;
		})();
	`, string(pageJSON))

	val, err := v8ctx.RunScript(script, "render.js")
	if err != nil {
		return "", fmt.Errorf("render failed: %w", err)
	}

	return val.String(), nil
}

func (r *Renderer) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return nil
	}

	r.closed = true
	close(r.pool)

	for ctx := range r.pool {
		ctx.Close()
	}

	if r.iso != nil {
		r.iso.Dispose()
	}

	return nil
}
