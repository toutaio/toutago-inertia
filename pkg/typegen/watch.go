package typegen

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher watches Go files and regenerates TypeScript types on changes.
type Watcher struct {
	watcher      *fsnotify.Watcher
	files        map[string]bool
	outputPath   string
	generator    func() error
	errorHandler func(error)
	debounce     time.Duration
	stopCh       chan struct{}
	mu           sync.Mutex
	timer        *time.Timer
}

// NewWatcher creates a new file watcher.
func NewWatcher() *Watcher {
	return &Watcher{
		files:    make(map[string]bool),
		debounce: 300 * time.Millisecond,
		stopCh:   make(chan struct{}),
	}
}

// AddFile adds a Go file to watch.
func (w *Watcher) AddFile(path string) error {
	// Check if file exists
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("file does not exist: %w", err)
	}

	w.mu.Lock()
	w.files[path] = true
	w.mu.Unlock()

	return nil
}

// AddDirectory adds all Go files in a directory to watch.
func (w *Watcher) AddDirectory(dir string) error {
	// Check if directory exists
	if info, err := os.Stat(dir); err != nil {
		return fmt.Errorf("directory does not exist: %w", err)
	} else if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", dir)
	}

	// Find all .go files
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	w.mu.Lock()
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".go" {
			path := filepath.Join(dir, entry.Name())
			w.files[path] = true
		}
	}
	w.mu.Unlock()

	return nil
}

// SetOutput sets the output path for generated TypeScript files.
func (w *Watcher) SetOutput(path string) {
	w.mu.Lock()
	w.outputPath = path
	w.mu.Unlock()
}

// SetGenerator sets the function to call when regenerating types.
func (w *Watcher) SetGenerator(fn func() error) {
	w.mu.Lock()
	w.generator = fn
	w.mu.Unlock()
}

// SetErrorHandler sets the function to call when errors occur.
func (w *Watcher) SetErrorHandler(fn func(error)) {
	w.mu.Lock()
	w.errorHandler = fn
	w.mu.Unlock()
}

// SetDebounce sets the debounce duration (minimum time between regenerations).
func (w *Watcher) SetDebounce(d time.Duration) {
	w.mu.Lock()
	w.debounce = d
	w.mu.Unlock()
}

// Watch starts watching files and regenerating on changes.
func (w *Watcher) Watch() error {
	var err error
	w.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer w.watcher.Close()

	// Add all files to watcher
	w.mu.Lock()
	for file := range w.files {
		if err := w.watcher.Add(file); err != nil {
			w.mu.Unlock()
			return fmt.Errorf("failed to watch file %s: %w", file, err)
		}
	}
	w.mu.Unlock()

	// Initial generation
	w.generate()

	// Watch for changes
	for {
		select {
		case <-w.stopCh:
			return nil

		case event, ok := <-w.watcher.Events:
			if !ok {
				return nil
			}

			// Only care about write and create events for Go files
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				if filepath.Ext(event.Name) == ".go" {
					w.debounceGenerate()
				}
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return nil
			}
			w.handleError(fmt.Errorf("watcher error: %w", err))
		}
	}
}

// Stop stops the watcher.
func (w *Watcher) Stop() {
	close(w.stopCh)
}

// debounceGenerate schedules a generation after debounce period.
func (w *Watcher) debounceGenerate() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.timer != nil {
		w.timer.Stop()
	}

	w.timer = time.AfterFunc(w.debounce, func() {
		w.generate()
	})
}

// generate calls the generator function.
func (w *Watcher) generate() {
	w.mu.Lock()
	gen := w.generator
	w.mu.Unlock()

	if gen == nil {
		return
	}

	if err := gen(); err != nil {
		w.handleError(err)
	}
}

// handleError calls the error handler if set.
func (w *Watcher) handleError(err error) {
	w.mu.Lock()
	handler := w.errorHandler
	w.mu.Unlock()

	if handler != nil {
		handler(err)
	}
}
