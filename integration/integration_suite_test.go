package enforcer_test

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/fraenkel/candiedyaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"github.com/pivotal-cf-experimental/cf-mysql-quota-enforcer/config"
	"github.com/pivotal-cf-experimental/cf-mysql-quota-enforcer/database"

	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"os"
	"os/exec"
)

var brokerDBName string
var rootConfig config.Config
var binaryPath string

var configFile string

func TestEnforcer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Enforcer Suite")
}

func newRootDatabaseConfig(dbName string) config.Config {

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	Expect(err).ToNot(HaveOccurred())

	dbConfig := config.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     dbPort,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   dbName,
	}
	err = dbConfig.Validate()
	errString := "Error generating config. Provide all config properties as environment variable prefixed with 'DB_' (e.g. DB_HOST=localhost)"
	Expect(err).ToNot(HaveOccurred(), errString)

	return dbConfig
}

var _ = BeforeSuite(func() {
	initConfig := newRootDatabaseConfig("")

	brokerDBName = fmt.Sprintf("quota_enforcer_integration_enforcer_test_%d", GinkgoParallelNode())
	rootConfig = newRootDatabaseConfig(brokerDBName)

	initDB, err := database.NewConnection(initConfig)
	Expect(err).ToNot(HaveOccurred())
	defer initDB.Close()

	_, err = initDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", brokerDBName))
	Expect(err).ToNot(HaveOccurred())

	db, err := database.NewConnection(rootConfig)
	Expect(err).ToNot(HaveOccurred())
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS service_instances (
    id int(11) NOT NULL AUTO_INCREMENT,
    guid varchar(255),
    plan_guid varchar(255),
    max_storage_mb int(11) NOT NULL DEFAULT '0',
    db_name varchar(255),
    PRIMARY KEY (id)
	)`)
	Expect(err).ToNot(HaveOccurred())

	binaryPath, err = gexec.Build("github.com/pivotal-cf-experimental/cf-mysql-quota-enforcer")
	Expect(err).ToNot(HaveOccurred())

	_, err = os.Stat(binaryPath)
	if err != nil {
		Expect(os.IsExist(err)).To(BeTrue())
	}

	tempDir, err := ioutil.TempDir(os.TempDir(), "quota-enforcer-integration-test")
	Expect(err).NotTo(HaveOccurred())

	configFile = filepath.Join(tempDir, "quotaEnforcerConfig.yml")
	writeConfig()
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()

	_, err := os.Stat(binaryPath)
	if err != nil {
		Expect(os.IsExist(err)).To(BeFalse())
	}

	var emptyConfig config.Config
	if rootConfig != emptyConfig {
		db, err := database.NewConnection(rootConfig)
		Expect(err).ToNot(HaveOccurred())
		defer db.Close()

		_, err = db.Exec("DROP TABLE IF EXISTS service_instances")
		Expect(err).ToNot(HaveOccurred())

		_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", brokerDBName))
		Expect(err).ToNot(HaveOccurred())

	}
})

func startEnforcerWithFlags(flags ...string) *gexec.Session {

	flags = append(
		flags,
		fmt.Sprintf("-configFile=%s", configFile),
		"-logLevel=debug",
	)

	command := exec.Command(
		binaryPath,
		flags...,
	)

	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).ShouldNot(HaveOccurred())

	return session
}

func runEnforcerContinuously() *gexec.Session {
	session := startEnforcerWithFlags()
	Eventually(session.Out).Should(gbytes.Say("Running continuously"))
	return session
}

func runEnforcerOnce() {
	session := startEnforcerWithFlags("-runOnce")

	Eventually(session.Out).Should(gbytes.Say("Running once"))
	// Wait for the process to finish naturally.
	// This should not take a long time
	session.Wait(5 * time.Second)
	Expect(session.ExitCode()).To(Equal(0), string(session.Err.Contents()))
}

func writeConfig() {
	fileToWrite, err := os.Create(configFile)
	Expect(err).ShouldNot(HaveOccurred())

	encoder := candiedyaml.NewEncoder(fileToWrite)
	err = encoder.Encode(rootConfig)
	Expect(err).ShouldNot(HaveOccurred())
}

func getDirOfCurrentFile() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
