package template

import (
	"context"

	"github.com/messiaen/rolles"
	"github.com/spf13/cobra"
)

func NewDelCmd(o *rolles.TemplateOptions) *cobra.Command {
	ccmd := &cobra.Command{
		Use:   "del",
		Short: "Delete template",
		Long:  ``,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			gatheredTemplates, err := o.GatherTemplates()
			if err != nil {
				return err
			}
			es, err := o.EsOptions.NewClient()
			if err != nil {
				return err
			}
			ctx := context.Background()
			for name := range gatheredTemplates {
				if err := o.DelTemplate(es, ctx, name); err != nil {
					return err
				}
			}
			return nil
		},
	}
	return ccmd
}
