package cmd

import (
	"github.com/spf13/cobra"
	initcmd "github.com/thecxx/go-driver/pkg/builder/cmd/init"
	"github.com/thecxx/go-driver/pkg/runtime/cli/genericflags"
	"github.com/thecxx/go-driver/pkg/runtime/unexpected"
)

// Builder command
func NewBuilderCommand() *cobra.Command {
	cmds := &cobra.Command{
		Use:     "builder",
		Version: "1.0.0",
		Short:   "",
		Long:    "",
		// EntryPoint
		Run: func(cmd *cobra.Command, args []string) {
			unexpected.CheckError(cmd.Help())
		},
	}

	// flags
	if flags := cmds.Flags(); flags != nil {
		genericflags.NewHelpFlags().AddFlags(flags)
		genericflags.NewVersionFlags().AddFlags(flags)
	}

	// subcommands
	cmds.AddCommand(initcmd.NewInitCommand())

	return cmds
}
