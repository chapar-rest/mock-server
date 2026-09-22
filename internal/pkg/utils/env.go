package utils

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// MustProcess parses environment variables into cfg using envconfig and
// panics on any parse or validation error. A malformed or missing required
// variable must stop the process at boot, never leave zero values behind.
func MustProcess(prefix string, cfg any) {
	if err := envconfig.Process(prefix, cfg); err != nil {
		panic(fmt.Sprintf("envconfig.Process(prefix=%q): %v", prefix, err))
	}
}
