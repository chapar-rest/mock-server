package utils

import (
	"os"
	"testing"
)

var (
	name    = "mock-server"
	version = "snapshot"
	env     = "local"
)

// Name returns the name of the app.
func Name() string { return name }

// Version returns the version of the app.
func Version() string { return version }

// SetVersion records the build version, stamped into main by ko's ldflags.
func SetVersion(v string) { version = v }

// Env returns the environment of the app.
func Env() string { return env }

func init() {
	if env = os.Getenv("MOCK_ENV"); env == "" && !testing.Testing() {
		panic("MOCK_ENV is not set")
	}
}
