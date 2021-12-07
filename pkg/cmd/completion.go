package cmd

import (
	"os"

	"github.com/fingcloud/cli/pkg/cli"
	"github.com/spf13/cobra"
)

func NewCmdCompletion(ctx *cli.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

  $ source <(fing completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ fing completion bash > /etc/bash_completion.d/fing
  # macOS:
  $ fing completion bash > /usr/local/etc/bash_completion.d/fing

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ fing completion zsh > "${fpath[1]}/_fing"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ fing completion fish | source

  # To load completions for each session, execute once:
  $ fing completion fish > ~/.config/fish/completions/fing.fish

PowerShell:

  PS> fing completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> fing completion powershell > fing.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}

	return cmd
}
