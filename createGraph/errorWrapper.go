package main

import (
	"github.com/arangodb/go-driver"
	"github.com/pkg/errors"
)

func init() {
	driver.WithStack = errors.WithStack
	driver.Cause = errors.Cause
}
