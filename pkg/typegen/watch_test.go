package typegen

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func TestWatcher_Watch(t *testing.T) {
	// Create temporary directory for test files
	tmpDir := t.TempDir()

	// Create test Go file
	goFile := filepath.Join(tmpDir, "test.go")
	initialContent := `package test
type User struct {
	Name string
}`
	if err := os.WriteFile(goFile, []byte(initialContent), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create output file path
	outFile := filepath.Join(tmpDir, "types.ts")

	// Create watcher
	watcher := NewWatcher()
	if watcher == nil {
		t.Fatal("NewWatcher returned nil")
	}

	// Add file to watch
	if err := watcher.AddFile(goFile); err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	// Set output path
	watcher.SetOutput(outFile)

	// Set generator function
	var generated atomic.Int32
	watcher.SetGenerator(func() error {
		generated.Add(1)
		content := "export interface User { name: string; }"
		return os.WriteFile(outFile, []byte(content), 0600)
	})

	// Start watching in background
	errCh := make(chan error, 1)
	go func() {
		errCh <- watcher.Watch()
	}()

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Initial generation should happen
	if generated.Load() != 1 {
		t.Errorf("Expected 1 initial generation, got %d", generated.Load())
	}

	// Modify the file
	updatedContent := `package test
type User struct {
	Name string
	Email string
}`
	if err := os.WriteFile(goFile, []byte(updatedContent), 0600); err != nil {
		t.Fatalf("Failed to update test file: %v", err)
	}

	// Wait for regeneration (with timeout)
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if generated.Load() >= 2 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Should have regenerated
	if generated.Load() < 2 {
		t.Errorf("Expected at least 2 generations after file change, got %d", generated.Load())
	}

	// Stop watcher
	watcher.Stop()

	// Wait for watcher to stop
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Watcher returned error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Error("Watcher did not stop in time")
	}
}

func TestWatcher_Debounce(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "test.go")

	// Create initial file
	if err := os.WriteFile(goFile, []byte("package test"), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create watcher with short debounce
	watcher := NewWatcher()
	watcher.SetDebounce(200 * time.Millisecond)

	if err := watcher.AddFile(goFile); err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	// Track generations
	var generated atomic.Int32
	watcher.SetGenerator(func() error {
		generated.Add(1)
		return nil
	})

	// Start watching
	go watcher.Watch()
	defer watcher.Stop()

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)
	initialGen := generated.Load()

	// Make multiple rapid changes
	for i := range 5 {
		content := "package test\n// Change " + string(rune(i))
		if err := os.WriteFile(goFile, []byte(content), 0600); err != nil {
			t.Fatalf("Failed to update file: %v", err)
		}
		time.Sleep(20 * time.Millisecond) // Rapid changes
	}

	// Wait for debounce to settle
	time.Sleep(400 * time.Millisecond)

	// Should have debounced (not 5 regenerations)
	totalGen := generated.Load() - initialGen
	if totalGen >= 5 {
		t.Errorf("Debouncing failed: expected < 5 generations, got %d", totalGen)
	}
	if totalGen == 0 {
		t.Error("No regeneration after file changes")
	}
}

func TestWatcher_MultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple Go files
	file1 := filepath.Join(tmpDir, "user.go")
	file2 := filepath.Join(tmpDir, "post.go")

	if err := os.WriteFile(file1, []byte("package test"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("package test"), 0600); err != nil {
		t.Fatal(err)
	}

	watcher := NewWatcher()
	watcher.SetDebounce(100 * time.Millisecond)

	// Add both files
	if err := watcher.AddFile(file1); err != nil {
		t.Fatalf("AddFile failed for file1: %v", err)
	}
	if err := watcher.AddFile(file2); err != nil {
		t.Fatalf("AddFile failed for file2: %v", err)
	}

	var generated atomic.Int32
	watcher.SetGenerator(func() error {
		generated.Add(1)
		return nil
	})

	go watcher.Watch()
	defer watcher.Stop()

	time.Sleep(200 * time.Millisecond)
	initialGen := generated.Load()

	// Modify file1
	if err := os.WriteFile(file1, []byte("package test\n// Modified"), 0600); err != nil {
		t.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)

	if generated.Load() <= initialGen {
		t.Error("Expected regeneration after file1 change")
	}

	file1Gen := generated.Load()

	// Modify file2
	if err := os.WriteFile(file2, []byte("package test\n// Modified"), 0600); err != nil {
		t.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)

	if generated.Load() <= file1Gen {
		t.Error("Expected regeneration after file2 change")
	}
}

func TestWatcher_ErrorHandling(t *testing.T) {
	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "test.go")

	if err := os.WriteFile(goFile, []byte("package test"), 0600); err != nil {
		t.Fatal(err)
	}

	watcher := NewWatcher()
	watcher.SetDebounce(100 * time.Millisecond)
	if err := watcher.AddFile(goFile); err != nil {
		t.Fatal(err)
	}

	// Set generator that returns error
	var errorCount atomic.Int32
	watcher.SetGenerator(func() error {
		errorCount.Add(1)
		return os.ErrPermission
	})

	// Set error handler
	var lastError atomic.Value
	watcher.SetErrorHandler(func(err error) {
		lastError.Store(err)
	})

	go watcher.Watch()
	defer watcher.Stop()

	time.Sleep(300 * time.Millisecond)

	// Should have attempted generation
	if errorCount.Load() == 0 {
		t.Error("Generator was not called")
	}

	// Error handler should have been called
	if lastError.Load() == nil {
		t.Error("Error handler was not called")
	}

	// Modify file to trigger another error
	if err := os.WriteFile(goFile, []byte("package test\n// Modified"), 0600); err != nil {
		t.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)

	// Watcher should continue after errors
	if errorCount.Load() < 2 {
		t.Error("Watcher did not continue after error")
	}
}

func TestWatcher_AddDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Create some Go files in directory
	file1 := filepath.Join(tmpDir, "user.go")
	file2 := filepath.Join(tmpDir, "post.go")
	file3 := filepath.Join(tmpDir, "readme.md") // Non-Go file

	if err := os.WriteFile(file1, []byte("package test"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("package test"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file3, []byte("# README"), 0600); err != nil {
		t.Fatal(err)
	}

	watcher := NewWatcher()
	watcher.SetDebounce(100 * time.Millisecond)

	// Add directory (should watch all .go files)
	if err := watcher.AddDirectory(tmpDir); err != nil {
		t.Fatalf("AddDirectory failed: %v", err)
	}

	var generated atomic.Int32
	watcher.SetGenerator(func() error {
		generated.Add(1)
		return nil
	})

	go watcher.Watch()
	defer watcher.Stop()

	time.Sleep(300 * time.Millisecond)
	initialGen := generated.Load()

	// Modify Go file
	if err := os.WriteFile(file1, []byte("package test\n// Modified"), 0600); err != nil {
		t.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)

	if generated.Load() <= initialGen {
		t.Error("Expected regeneration after Go file change")
	}

	goFileGen := generated.Load()

	// Modify non-Go file (should not trigger regeneration)
	if err := os.WriteFile(file3, []byte("# UPDATED"), 0600); err != nil {
		t.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)

	if generated.Load() != goFileGen {
		t.Error("Non-Go file change triggered regeneration")
	}
}

func TestWatcher_Stop(t *testing.T) {
	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "test.go")

	if err := os.WriteFile(goFile, []byte("package test"), 0600); err != nil {
		t.Fatal(err)
	}

	watcher := NewWatcher()
	if err := watcher.AddFile(goFile); err != nil {
		t.Fatal(err)
	}

	var generated atomic.Int32
	watcher.SetGenerator(func() error {
		generated.Add(1)
		return nil
	})

	errCh := make(chan error, 1)
	go func() {
		errCh <- watcher.Watch()
	}()

	time.Sleep(100 * time.Millisecond)
	beforeStop := generated.Load()

	// Stop watcher
	watcher.Stop()

	// Wait for Watch to return
	select {
	case <-errCh:
		// Good, watcher stopped
	case <-time.After(1 * time.Second):
		t.Fatal("Watcher did not stop in time")
	}

	// Modify file after stop
	if err := os.WriteFile(goFile, []byte("package test\n// After stop"), 0600); err != nil {
		t.Fatal(err)
	}

	time.Sleep(300 * time.Millisecond)

	// Should not regenerate after stop
	if generated.Load() != beforeStop {
		t.Error("Watcher regenerated after Stop()")
	}
}

func TestWatcher_NoGenerator(t *testing.T) {
	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "test.go")

	if err := os.WriteFile(goFile, []byte("package test"), 0600); err != nil {
		t.Fatal(err)
	}

	watcher := NewWatcher()
	if err := watcher.AddFile(goFile); err != nil {
		t.Fatal(err)
	}

	// Don't set generator - should handle gracefully
	errCh := make(chan error, 1)
	go func() {
		errCh <- watcher.Watch()
	}()

	time.Sleep(200 * time.Millisecond)
	watcher.Stop()

	select {
	case err := <-errCh:
		// Should not panic, but may return error
		_ = err
	case <-time.After(1 * time.Second):
		t.Fatal("Watcher did not stop")
	}
}

func TestWatcher_InvalidPath(t *testing.T) {
	watcher := NewWatcher()

	// Try to add non-existent file
	err := watcher.AddFile("/non/existent/file.go")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	// Try to add non-existent directory
	err = watcher.AddDirectory("/non/existent/directory")
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
}
