package reload

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/AlexanderBrese/gomon/pkg/configuration"
	"github.com/AlexanderBrese/gomon/pkg/logging"
	"github.com/AlexanderBrese/gomon/pkg/utils"
)

const checkReloadDelay = 600

func TestReload(t *testing.T) {
	cfg, err := configuration.TestConfiguration()
	cfg.Build.RelSrcDir = "cmd/web"

	if err != nil {
		t.Error(err)
	}
	logger := logging.NewLogger(cfg)
	reloader := NewReload(cfg, logger)

	if err := buildPrepare(cfg); err != nil {
		t.Error(err)
	}

	defer func() {
		if err := runCleanup(reloader); err != nil {
			// TODO: log
			return
		}
	}()

	reloadStart(reloader)
	time.Sleep(checkReloadDelay * time.Millisecond)
	if err := reloadPassed(reloader); err != nil {
		t.Error(err)
	}
}

func reloadStart(reloader *Reload) {
	reloader.Run()
}

func reloadPassed(reloader *Reload) error {
	binary, err := reloader.config.Binary()
	if err != nil {
		return err
	}
	if err := utils.CheckPath(binary); err != nil {
		return fmt.Errorf("error: there was no built binary found at %s", binary)
	}
	return utils.WithLockAndError(&reloader.mu, func() error {
		if !reloader.running {
			return errors.New("error: binary not running")
		}
		return nil
	})
}
