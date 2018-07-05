package commands

import (
	"flag"
	"fmt"

	"github.com/hazbo/httpu/meta"
)

var versionFlagSet = flag.NewFlagSet("version", flag.ExitOnError)

func versionValue(args []string) error {
	newFlagSet.Parse(args)

	fmt.Printf("%s v%s\n", meta.Application, meta.Version)

	return nil
}

var versionCmd = &Command{
	Usage: func(arg0 string) {
		fmt.Printf("Usage: %s list [<option>...]\n\nOptions:\n", arg0)
		versionFlagSet.PrintDefaults()
	},
	RunMethod: func(args []string) error {
		return versionValue(args)
	},
}
