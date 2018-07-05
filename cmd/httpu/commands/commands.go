package commands

// Command defines the functionality and useage text for any given subcommand.
// These may be from either built-in commands, inline or external commands.
type Command struct {
	Usage     func(string)
	RunMethod func([]string) error
}

// Run will run the given command with the arguments passed
func (cmd *Command) Run(args []string) error {
	return cmd.RunMethod(args)
}

// CommandMap defines the string key for a given command, which is then used
// when invoking the cnf commands.
type CommandMap map[string]*Command

// Commands is the list of commands within a map
var Commands = CommandMap{
	"new":     newCmd,
	"pull":    pullCmd,
	"version": versionCmd,
}
