package app

import (
	"fmt"

	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/ui"
	"github.com/fingcloud/cli/pkg/util"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type CreateOptions struct {
	Name     string
	Platform string
	Region   string
	Plan     string
}

func NewCmdCreate(ctx *cli.Context) *cobra.Command {
	opts := new(CreateOptions)

	cmd := &cobra.Command{
		Use:   "create [app]",
		Short: "create an app",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			opts.Name = args[0]

			util.CheckErr(RunCreate(ctx, opts))
		},
	}

	cmd.Flags().StringVar(&opts.Plan, "plan", "", "your app plan. check 'fing plans'")

	return cmd
}

func RunCreate(ctx *cli.Context, opts *CreateOptions) error {
	var plans []*api.Plan
	if opts.Plan == "" {
		var err error
		plans, err = ctx.Client.PlansList()
		util.CheckErr(err)

		printer := message.NewPrinter(language.English)

		options := funk.Map(plans, func(p *api.Plan) string {
			return printer.Sprintf(
				"[ %s ]\tMemory: %.1fG\tCPU: %.2f\tStorage: %.1fG\tPrice: %d Tomans",
				p.Name, p.Memory, p.CPU, p.Storage, p.MonthlyPrice/10,
			)
		}).([]string)

		var selected string
		err = ui.PromptSelect("select plan", options, &selected)
		util.CheckErr(err)

		fmt.Sscanf(selected, "[ %s ]", &opts.Plan)
	}

	plan, ok := funk.Find(plans, func(p *api.Plan) bool { return p.Name == opts.Plan }).(*api.Plan)
	if !ok {
		return fmt.Errorf("invalid plan %s\nRun 'fing plans' to see available plans", opts.Plan)
	}

	app, err := ctx.Client.AppsCreate(&api.CreateAppOptions{
		Label:    opts.Name,
		Platform: "auto",
		PlanID:   plan.ID,
		Region:   "iran",
	})
	util.CheckErr(err)

	fmt.Printf("app '%s' created\n", app.Name)
	return nil
}
