package alias

import (
	"context"
	"time"

	"github.com/messiaen/rolles"
	"github.com/spf13/cobra"
)

func NewPutCmd(o *rolles.AliasOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "put",
		Short: "Put index by alias",
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

			t := time.Now()
			for _, c := range cfgs {
				err := o.PutAlias(es, ctx, c, t)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}
