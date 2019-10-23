package init

import (
	"github.com/spf13/cobra"
	"github.com/thecxx/go-driver/pkg/runtime/cli/genericflags"
	"github.com/thecxx/go-driver/pkg/runtime/unexpected"
)

//
type Options struct {
}

// Validate options
func (o *Options) Validate() error {

	return nil
}

//
type Command struct {
	opts *Options
}

// New Command
func NewCommand(opts *Options) *Command {
	return &Command{opts}
}

// Validate environment
func (c *Command) Validate(cmd *cobra.Command, args []string) error {

	// options

	// args

	return nil
}

// Execute "run" command
func (c *Command) Run(cmd *cobra.Command) error {
	return nil
}

// New init command
func NewInitCommand() *cobra.Command {
	opts := &Options{}
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init package",
		Long:  "",
		// EntryPoint
		Run: func(cmd *cobra.Command, args []string) {
			if c := NewCommand(opts); c != nil {
				unexpected.CheckError(c.Validate(cmd, args))
				unexpected.CheckError(c.Run(cmd))
			}
		},
	}

	// flags
	if flags := cmd.Flags(); flags != nil {
		genericflags.NewHelpFlags().AddFlags(flags)
		// options
	}

	return cmd
}
