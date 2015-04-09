package enforcer

import (
	"fmt"
	"os"
	"time"

	"github.com/pivotal-cf-experimental/cf-mysql-quota-enforcer/database"
	"github.com/pivotal-golang/lager"
)

type Enforcer interface {
	EnforceOnce() error
	Run(<-chan os.Signal, chan<- struct{}) error
}

type enforcer struct {
	violatorRepo, reformerRepo database.Repo
	logger                     lager.Logger
}

func NewEnforcer(violatorRepo, reformerRepo database.Repo, logger lager.Logger) Enforcer {
	return &enforcer{
		violatorRepo: violatorRepo,
		reformerRepo: reformerRepo,
		logger:       logger,
	}
}

func (e enforcer) EnforceOnce() error {
	err := e.revokePrivilegesFromViolators()
	if err != nil {
		return err
	}

	err = e.grantPrivilegesToReformed()
	if err != nil {
		return err
	}

	return nil
}

func (e enforcer) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	go func() {
		for {
			err := e.EnforceOnce()
			if err != nil {
				e.logger.Info(fmt.Sprintf("Enforcing Failed: %s", err.Error()))
			}
			time.Sleep(1 * time.Second)
		}
	}()

	close(ready)
	<-signals
	return nil
}

func (e enforcer) revokePrivilegesFromViolators() error {
	e.logger.Info("Looking for violators")

	violators, err := e.violatorRepo.All()
	if err != nil {
		return fmt.Errorf("Finding violators: %s", err.Error())
	}

	for _, db := range violators {
		err = db.RevokePrivileges()
		if err != nil {
			return fmt.Errorf("Revoking privileges: %s", err.Error())
		}

		err = db.KillActiveConnections()
		if err != nil {
			return fmt.Errorf("Resetting active privileges: %s", err.Error())
		}
	}
	return nil
}

func (e enforcer) grantPrivilegesToReformed() error {
	e.logger.Info("Looking for reformers")

	reformers, err := e.reformerRepo.All()
	if err != nil {
		return fmt.Errorf("Finding reformers: %s", err.Error())
	}

	for _, db := range reformers {
		err = db.GrantPrivileges()
		if err != nil {
			return fmt.Errorf("Granting privileges: %s", err.Error())
		}

		err = db.KillActiveConnections()
		if err != nil {
			return fmt.Errorf("Resetting active privileges: %s", err.Error())
		}
	}

	return nil
}
