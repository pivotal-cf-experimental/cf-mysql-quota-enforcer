package enforcer

import (
	"os"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf-experimental/cf-mysql-quota-enforcer/clock"
	"github.com/tedsuo/ifrit"
)

type runner struct {
	enforcer Enforcer
	clock    clock.Clock
	pause    time.Duration
	logger   lager.Logger
}

func NewRunner(enforcer Enforcer, clock clock.Clock, pause time.Duration,
	logger lager.Logger) ifrit.Runner {
	return &runner{
		enforcer: enforcer,
		clock:    clock,
		pause:    pause,
		logger:   logger,
	}
}

func (r runner) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	close(ready)
	for {
		err := r.enforcer.EnforceOnce()
		if err != nil {
			r.logger.Error("Enforcing Failed", err)
		}
		select {
		case <-signals:
			return nil
		case <-r.clock.After(r.pause):
		}
	}
}
