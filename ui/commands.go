package ui

import (
	"fmt"
	"os/exec"

	"github.com/hazbo/httpu/stash"
	"github.com/jroimartin/gocui"
)

// Command requires each member to have an Execute func to allow each and any
// command to be executable while in command mode.
type Command interface {
	Execute(g *gocui.Gui, cmd string, args []string) error
}

// EchoCommand represents the command that echos text out into the request
// view screen.
//
// Usage: echo Hello, World!
type EchoCommand struct {
}

// Execute will print each echo argument to the request view screen.
func (ec EchoCommand) Execute(g *gocui.Gui, cmd string, args []string) error {
	defer cmdBarRefresh(g)
	RequestView.Clear()
	for _, a := range args {
		fmt.Fprintf(RequestView, "%s ", a)
	}
	return nil
}

// ClearCommand represents the command that will clear all contents from the
// request view screen.
//
// Usage: clear
type ClearCommand struct {
}

// Execute will clear the contents of the request view screen.
func (cc ClearCommand) Execute(g *gocui.Gui, cmd string, args []string) error {
	defer cmdBarRefresh(g)
	RequestView.Clear()
	return nil
}

// ListCommandsCommand represents the command that will list all available
// commands that can be ran while in command mode.
//
// Usage: list-commands
type ListCommandsCommand struct {
}

// Execute will list all available commands, in the request view screen.
func (lcc ListCommandsCommand) Execute(g *gocui.Gui, cmd string, args []string) error {
	defer cmdBarRefresh(g)
	if len(args) > 0 {
		return fmt.Errorf("list-commands expects 0 arguments, %d passed.", len(args))
	}
	RequestView.Clear()
	// TODO: These should be listed in order, not at random
	for c, _ := range Commands {
		fmt.Fprintf(RequestView, "%s\n", c)
	}
	return nil
}

type StashCommand struct {
}

func (sc StashCommand) Execute(g *gocui.Gui, cmd string, args []string) error {
	defer cmdBarRefresh(g)
	RequestView.Clear()
	for n, s := range stash.Store {
		fmt.Fprintf(RequestView, "%s: %s\n", n, s.Value)
	}
	return nil
}

type WelcomeCommand struct {
}

func (wc WelcomeCommand) Execute(g *gocui.Gui, cmd string, args []string) error {
	defer cmdBarRefresh(g)
	RequestView.Clear()
	requestViewSetup(g)
	return nil
}

type ShellCommand struct {
}

func (sc ShellCommand) Execute(g *gocui.Gui, cmd string, args []string) error {
	defer cmdBarRefresh(g)
	RequestView.Clear()

	command := exec.Command(args[0])

	for _, v := range args[1:] {
		command.Args = append(command.Args, v)
	}

	output, _ := command.CombinedOutput()

	fmt.Fprintf(RequestView, "%s", output)

	return nil
}

// Commands is a map of all available commands to be used while in command mode.
var Commands map[string]Command = map[string]Command{
	"clear":         ClearCommand{},
	"echo":          EchoCommand{},
	"list-commands": ListCommandsCommand{},
	"stash":         StashCommand{},
	"welcome":       WelcomeCommand{},
	"!":             ShellCommand{},
}
