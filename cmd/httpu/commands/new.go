package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/hazbo/httpu"
	"github.com/hazbo/httpu/ui"
	"github.com/joho/godotenv"
)

var newFlagSet = flag.NewFlagSet("new", flag.ExitOnError)

var (
	newEnvFile = newFlagSet.String(
		"e", "", "Loads .env file to start httpu with environment variables")
)

func newValue(args []string) error {
	newFlagSet.Parse(args)

	// Get the first argument for the `new` command
	p := args[0]

	if *newEnvFile != "" {
		err := godotenv.Load(*newEnvFile)
		if err != nil {
			return fmt.Errorf("Could not find .env file: %s", *newEnvFile)
		}

		// If an env file has been passed, this becomes argument 2, over 0.
		p = args[2]
	}

	if len(args) == 0 {
		fmt.Printf("Expecting 1 argument, 0 passed")
		os.Exit(1)
	}

	err := httpu.ConfigureFromFile(p)
	if err != nil {
		return err
	}

	// Start the terminal user interface!
	ui.New().Start()

	return nil
}

var newCmd = &Command{
	Usage: func(arg0 string) {
		fmt.Printf("Usage: %s new <package_name> [<options>...]\n\nOptions:\n", arg0)
		newFlagSet.PrintDefaults()
	},
	RunMethod: func(args []string) error {
		return newValue(args)
	},
}
