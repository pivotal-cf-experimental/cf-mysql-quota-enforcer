package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tedsuo/ifrit"

	"github.com/cloudfoundry-incubator/cf-lager"
	"github.com/pivotal-cf-experimental/cf-mysql-quota-enforcer/config"
	"github.com/pivotal-cf-experimental/cf-mysql-quota-enforcer/database"
	"github.com/pivotal-cf-experimental/cf-mysql-quota-enforcer/enforcer"
	"github.com/pivotal-golang/lager"
)

type Config struct {
	Host         string
	Port         int
	User         string
	Password     string
	BrokerDBName string
}

func main() {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	runOnce := flags.Bool("runOnce", false, "Run only once instead of continuously")
	configFile := flags.String("configFile", "", "Location of config file")
	cf_lager.AddFlags(flags)
	flags.Parse(os.Args[1:])
	logger, _ := cf_lager.New("Quota Enforcer")

	logger.Info("Config file", lager.Data{"configFile": configFile})
	config, err := config.Load(*configFile)
	if err != nil {
		panic(err.Error())
	}

	brokerDBName := config.DBName
	if brokerDBName == "" {
		logger.Fatal("Must specify DBName in the config file", nil)
	}

	db, err := database.NewConnection(*config)
	if db != nil {
		defer db.Close()
	}

	if err != nil {
		panic(err.Error())
	}

	logger.Info(
		"Database connection established.",
		lager.Data{
			"Host":         config.Host,
			"Port":         config.Port,
			"User":         config.User,
			"DatabaseName": brokerDBName,
		})

	violatorRepo := database.NewViolatorRepo(brokerDBName, db, logger)
	reformerRepo := database.NewReformerRepo(brokerDBName, db, logger)

	e := enforcer.NewEnforcer(violatorRepo, reformerRepo, logger)
	r := enforcer.NewRunner(e, logger)

	if *runOnce {
		logger.Info("Running once")

		err := e.EnforceOnce()
		if err != nil {
			logger.Info(fmt.Sprintf("Quota Enforcing Failed: %s", err.Error()))
		}
	} else {
		process := ifrit.Invoke(r)
		logger.Info("Running continuously")

		err := <-process.Wait()
		if err != nil {
			logger.Fatal("Quota Enforcing Failed", err)
		}
	}
}

func enforce(e enforcer.Enforcer, logger lager.Logger) {
	logger.Info("Enforcing")
	err := e.EnforceOnce()
	if err != nil {
		logger.Info(fmt.Sprintf("Enforcing Failed: %s", err.Error()))
	}
}
