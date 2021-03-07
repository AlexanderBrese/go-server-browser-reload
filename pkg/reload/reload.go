package reload

import (
	"sync"

	"github.com/AlexanderBrese/GOATmon/pkg/configuration"
)

type Reload struct {
	config *configuration.Configuration

	running         bool
	startBuilding   chan bool
	stop            chan bool
	stopRunning     chan bool
	FinishedRunning chan bool
	mu              sync.Mutex
}

func NewReload(cfg *configuration.Configuration) *Reload {
	return &Reload{
		config:          cfg,
		running:         false,
		startBuilding:   make(chan bool, 1),
		stop:            make(chan bool, 1),
		stopRunning:     make(chan bool),
		FinishedRunning: make(chan bool, 1),
	}
}

func (r *Reload) Cleanup() {
	r.BuildCleanup()
	r.RunCleanup()
}

func (r *Reload) Reload() {
	r.Cleanup()
	go r.start()
}

func (r *Reload) start() error {
	r.startBuilding <- true
	defer func() {
		<-r.startBuilding
	}()

	select {
	case <-r.stop:
		return nil
	default:
	}
	if err := r.build(); err != nil {
		return err
	}

	select {
	case <-r.stop:
		return nil
	default:
	}

	return r.run()
}
