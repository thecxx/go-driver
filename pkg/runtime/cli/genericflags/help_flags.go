package genericflags

import (
	"github.com/spf13/pflag"
)

const (
	FLAG_KEY_HELP = "help"
)

type HelpFlags struct {
	// flags
	Help bool
}

// New HelpFlags
func NewHelpFlags() *HelpFlags {
	return &HelpFlags{false}
}

func (this *HelpFlags) AddFlags(flags *pflag.FlagSet) {
	flags.BoolVar(&this.Help, FLAG_KEY_HELP, this.Help, "Print usage and this help message and exit")
}
