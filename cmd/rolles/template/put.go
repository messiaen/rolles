package template

import (
	"context"

	"github.com/messiaen/rolles"
	"github.com/spf13/cobra"
)

func NewPutCmd(o *rolles.TemplateOptions) *cobra.Command {
	ccmd := &cobra.Command{
		Use:   "put",
		Short: "Put template",
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
			for name, fn := range gatheredTemplates {
				if err := o.PutTemplateFromFilename(es, ctx, name, fn); err != nil {
					return err
				}
			}
			return nil
		},
	}
	return ccmd
}
