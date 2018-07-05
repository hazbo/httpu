package commands

import (
	"flag"
	"fmt"

	utils "github.com/hazbo/httpu/utils/common"
)

var pullFlagSet = flag.NewFlagSet("pull", flag.ExitOnError)

func pullValue(args []string) error {
	pullFlagSet.Parse(args)

	// Make the ~/.httpu directory if it does not already exists.
	err := utils.MakeHttpuDir()
	if err != nil {
		return err
	}

	// If the packages directory does not yet exist, create it.
	gc := utils.NewGitCommand()
	if !utils.PackagesExists() {
		return gc.Clone()
	}

	// Pull the latest from the packages directory
	return gc.Pull()
}

var pullCmd = &Command{
	Usage: func(arg0 string) {
		fmt.Printf("Usage: %s list [<option>...]\n\nOptions:\n", arg0)
		pullFlagSet.PrintDefaults()
	},
	RunMethod: func(args []string) error {
		return pullValue(args)
	},
}
