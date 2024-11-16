package watcher

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/bep/debounce"
	"github.com/fsnotify/fsnotify"
)

func NewWatcher(baseDir string, timeout time.Duration, rebuildFunc func()) error {
	baseDir = filepath.Clean(baseDir)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating fsnotify watcher: %w", err)
	}

	err = filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walkdir error: %w", err)
		}

		if !d.IsDir() {
			return nil
		}

		if err := watcher.Add(path); err != nil {
			return fmt.Errorf("error adding %s to watcher: %w", path, err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	debouncer := debounce.New(timeout)

	go func() {
		for range watcher.Events {
			debouncer(rebuildFunc)
		}
	}()

	return nil
}
