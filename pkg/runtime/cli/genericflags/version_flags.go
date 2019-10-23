package genericflags

import (
	"github.com/spf13/pflag"
)

const (
	FLAG_KEY_VERSION = "version"
)

type VersionFlags struct {
	// flags
	Version bool
}

// New VersionFlags
func NewVersionFlags() *VersionFlags {
	return &VersionFlags{false}
}

func (this *VersionFlags) AddFlags(flags *pflag.FlagSet) {
	flags.Bool("version", false, "Print version and exit")
}
