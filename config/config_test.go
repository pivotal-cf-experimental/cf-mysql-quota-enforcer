package config_test

import (
	. "github.com/pivotal-cf-experimental/cf-mysql-quota-enforcer/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("Validate", func() {
		var config Config

		BeforeEach(func() {
			config = Config{
				Host:         "fake-host",
				Port:         9999,
				User:         "fake-user",
				Password:     "fake-password",
				ReadOnlyUser: "fake-read-only-user",
				DBName:       "fake-db-name",
			}
		})

		It("validates a valid config file", func() {
			err := config.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when Host is not specified", func() {
			BeforeEach(func() {
				config.Host = ""
			})

			It("returns a validation error", func() {
				err := config.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Host"))
			})
		})

		Context("when Port is not specified", func() {
			BeforeEach(func() {
				config.Port = 0
			})

			It("returns a validation error", func() {
				err := config.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Port"))
			})
		})

		Context("when User is not specified", func() {
			BeforeEach(func() {
				config.User = ""
			})

			It("returns a validation error", func() {
				err := config.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("User"))
			})
		})

		Context("when ReadOnlyUser is not specified", func() {
			BeforeEach(func() {
				config.ReadOnlyUser = ""
			})

			It("returns a validation error", func() {
				err := config.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ReadOnlyUser"))
			})
		})

		Context("when Password is not specified", func() {
			BeforeEach(func() {
				config.Password = ""
			})

			It("returns a validation error", func() {
				err := config.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Password"))
			})
		})

		Context("when DBName is not specified", func() {
			BeforeEach(func() {
				config.DBName = ""
			})

			It("returns a validation error", func() {
				err := config.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("DBName"))
			})
		})
	})
})
