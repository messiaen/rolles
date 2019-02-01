package alias

import (
	"github.com/messiaen/rolles"
	"github.com/spf13/cobra"
)

func NewAliasCmd() *cobra.Command {
	o := rolles.NewDefaultAliasOptions()

	ccmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage ES indices via their aliases",
		Long:  ``,
	}

	// defaultOpts := NewDefaultAliasOptions()
	// defaultOpts.BindFlags(ccmd.PersistentFlags())

	flags := ccmd.PersistentFlags()

	o.EsOptions.BindFlags(flags)

	f := "config"
	flags.StringVarP(&o.Config, f, "c", o.Config, "indices configuration file")

	f = "prefix"
	flags.StringVarP(&o.Prefix, f, "p", o.Prefix, "alias name prefix")

	f = "name"
	flags.StringVarP(&o.Name, f, "n", o.Name, "alias name (all if not specified)")

	ccmd.AddCommand(NewPutCmd(&o))
	ccmd.AddCommand(NewDelCmd(&o))
	return ccmd
}
