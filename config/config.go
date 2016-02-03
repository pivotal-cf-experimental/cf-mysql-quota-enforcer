package config

import (
	"errors"
	"fmt"

	"gopkg.in/validator.v2"
)

type Config struct {
	Host         string `yaml:"Host" validate:"nonzero"`
	Port         int    `yaml:"Port" validate:"nonzero"`
	User         string `yaml:"User" validate:"nonzero"`
	Password     string `yaml:"Password"` //blank Password is allowed
	ReadOnlyUser string `yaml:"ReadOnlyUser" validate:"nonzero"`
	DBName       string `yaml:"DBName"` //blank DBName is allowed
}

func (c Config) Validate() error {
	err := validator.Validate(c)
	var errString string
	if err != nil {
		errString = formatErrorString(err)
	}

	if len(errString) > 0 {
		return errors.New(fmt.Sprintf("Validation errors: %s\n", errString))
	}
	return nil
}

func formatErrorString(err error) string {
	errs := err.(validator.ErrorMap)
	var errsString string
	for fieldName, validationMessage := range errs {
		errsString += fmt.Sprintf("%s : %s\n", fieldName, validationMessage)
	}
	return errsString
}
