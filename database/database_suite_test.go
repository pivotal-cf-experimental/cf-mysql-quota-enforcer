package database_test

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "testing"
)

func TestEnforcer(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Database Suite")
}
