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
		fmt.Printf("Usage: %s pull [<options>...]\n\nOptions:\n", arg0)
		pullFlagSet.PrintDefaults()
		if flag.NFlag() < 1 {
			fmt.Println("\tNo options for this command")
		}
	},
	RunMethod: func(args []string) error {
		return pullValue(args)
	},
}
