package assetfs

import (
	"fmt"
	"os"
	"time"

	"github.com/lorenzogood/blog/internal/watcher"
	"go.uber.org/zap"
)

func NewWatched(baseDir string) (*AssetFS, error) {
	logger := zap.L().Named("assetfs_watched")

	var a AssetFS
	rebuildFunc := func() {
		logger.Info("building")
		as, err := New(os.DirFS(baseDir), ".")
		if err != nil {
			logger.Error("failed to rebuild assetfs", zap.Error(err))
			return
		}
		a = *as
		logger.Debug("built")
	}

	rebuildFunc()

	if err := watcher.NewWatcher(baseDir, 100*time.Millisecond, rebuildFunc); err != nil {
		return nil, fmt.Errorf("error watching asset dir: %w", err)
	}

	return &a, nil
}

func NewWatchedPermanent(baseDir, webBase string) (*PermanentFS, error) {
	logger := zap.L().Named("assetfs_watched")

	var a PermanentFS
	rebuildFunc := func() {
		logger.Info("building")
		as, err := NewPermanent(os.DirFS(baseDir), ".", webBase)
		if err != nil {
			logger.Error("failed to rebuild assetfs", zap.Error(err))
			return
		}
		a = *as
		logger.Debug("built")
	}

	rebuildFunc()

	if err := watcher.NewWatcher(baseDir, 100*time.Millisecond, rebuildFunc); err != nil {
		return nil, fmt.Errorf("error watching asset dir: %w", err)
	}

	return &a, nil
}
