package template

import (
	"github.com/messiaen/rolles"
	"github.com/spf13/cobra"
)

func NewTemplateCmd() *cobra.Command {
	o := rolles.NewDefaultTemplateOptions()

	ccmd := &cobra.Command{
		Use:   "template",
		Short: "Manage ES templates",
		Long:  ``,
	}

	flags := ccmd.PersistentFlags()

	o.EsOptions.BindFlags(flags)

	f := "temp-dir"
	flags.StringVarP(&o.Dir, f, "d", o.Dir, "root template directory")

	f = "prefix"
	flags.StringVarP(&o.Prefix, f, "p", o.Prefix, "template name prefix")

	f = "name"
	flags.StringVarP(&o.Name, f, "n", o.Name, "template name (all if not specified)")

	ccmd.AddCommand(NewPutCmd(&o))
	ccmd.AddCommand(NewDelCmd(&o))
	return ccmd
}
