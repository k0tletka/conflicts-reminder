//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	cmdLocation      = "cmd/reminder.go"
	artifactLocation = "build/reminder"
)

var (
	buildEnv = map[string]string{
		"CGO_ENABLED": "0",
	}
)

func Build() error {
	mg.Deps(Tidy)

	return sh.RunWithV(
		buildEnv,
		"go", "build", "-o", artifactLocation, cmdLocation,
	)
}

func Tidy() error {
	return sh.RunV("go", "mod", "tidy")
}
