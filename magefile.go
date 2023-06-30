//go:build mage

package main

import "github.com/magefile/mage/sh"

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
	return sh.RunWithV(
		buildEnv,
		"go", "build", "--ldflags=\"-s -w\"", "-o", artifactLocation, cmdLocation,
	)
}
