package cmd

import (
	"os"
	"strings"
)

type CmdOptions interface {
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}

func FlagName(prefix, name string) string {
	return strings.Join([]string{prefix, name}, "-")
}
