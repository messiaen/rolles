package alias

import (
	"context"

	"github.com/messiaen/rolles"
	"github.com/spf13/cobra"
)

func NewRolloverCmd(o *rolles.AliasOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "rollover",
		Short: "Use Es rollover api to rollover alias according to its config",
		Long:  ``,
		RunE: func(ccmd *cobra.Command, args []string) error {
			cfgs, err := o.GatherAliasCfgsFromFile()
			if err != nil {
				return err
			}

			es, err := o.EsOptions.NewClient()
			if err != nil {
				return err
			}
			ctx := context.Background()

			for _, c := range cfgs {
				err := o.Rollover(es, ctx, c)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}
