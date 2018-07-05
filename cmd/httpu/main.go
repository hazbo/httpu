package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/hazbo/httpu/cmd/httpu/commands"
)

const usageMessageTemplate = `usage: %s <command>
Where <command> is one of:
  %s
For individual command usage, run:
  %s help <command>
`

func usage() {
	command := os.Args[0]
	var subcommands []string
	for subcommand := range commands.Commands {
		subcommands = append(subcommands, subcommand)
	}
	sort.Strings(subcommands)
	fmt.Printf(usageMessageTemplate, command, strings.Join(subcommands, "\n  "), command)
}

func help() {
	if len(os.Args) < 3 {
		usage()
		return
	}
	cmdmap := commands.Commands
	subcommand, ok := cmdmap[os.Args[2]]
	if !ok {
		fmt.Printf("Unknown command %q\n", os.Args[2])
		usage()
		return
	}
	subcommand.Usage(os.Args[0])
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "help" {
		help()
		return
	}
	if len(os.Args) < 2 {
		help()
		return
	}
	cmdmap := commands.Commands
	subcommand, ok := cmdmap[os.Args[1]]
	if !ok {
		fmt.Printf("Unknown command: %q\n", os.Args[1])
		usage()
		return
	}
	if err := subcommand.Run(os.Args[2:]); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
