package main

import (
	"os"
	"strings"

	"github.com/messiaen/rolles/cmd/rolles/alias"
	"github.com/messiaen/rolles/cmd/rolles/template"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRollesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rolles",
		Short: "Manage Elasticsearch temporal indices with rollover strategy",
	}

	cmd.AddCommand(template.NewTemplateCmd())
	cmd.AddCommand(alias.NewAliasCmd())
	return cmd
}

func Execute() error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("ROLLES")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	ccmd := NewRollesCommand()

	return ccmd.Execute()
}

func main() {
	if err := Execute(); err != nil {
		os.Exit(1)
	}
}
