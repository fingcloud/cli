package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Exact(num int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == num {
			return nil
		}
		return fmt.Errorf("%q requires exactly %d args", cmd.CommandPath(), num)
	}
}

func NoArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}

	if cmd.HasSubCommands() {
		return cmd.Help()
	}

	return fmt.Errorf("%q accepts no args", cmd.CommandPath())
}
