package ui

import (
	"fmt"
	"log"
	"strings"

	"github.com/hazbo/httpu/resource"
	"github.com/jroimartin/gocui"
)

// setBindings applies all keybindings for the user interface.
func (u *Ui) setBindings() {
	var err error

	err = u.Gui.SetKeybinding(
		"", gocui.KeyCtrlC, gocui.ModNone, quit)
	if err != nil {
		log.Panicln(err)
	}

	err = u.Gui.SetKeybinding(
		"", gocui.KeyCtrlS, gocui.ModNone, switchTopView)
	if err != nil {
		log.Panicln(err)
	}

	err = u.Gui.SetKeybinding(
		"", gocui.KeyCtrlW, gocui.ModNone, switchCmdView)
	if err != nil {
		log.Panicln(err)
	}

	err = u.Gui.SetKeybinding(
		cmdBar, gocui.KeyArrowUp, gocui.ModNone, switchModeCmd)
	if err != nil {
		log.Panicln(err)
	}

	err = u.Gui.SetKeybinding(
		cmdBar, gocui.KeyArrowDown, gocui.ModNone, switchModeDefault)
	if err != nil {
		log.Panicln(err)
	}

	err = u.Gui.SetKeybinding(
		cmdBar, gocui.KeyArrowLeft, gocui.ModNone, focusRequest)
	if err != nil {
		log.Panicln(err)
	}

	err = u.Gui.SetKeybinding(
		cmdBar, gocui.KeyArrowRight, gocui.ModNone, focusResponse)
	if err != nil {
		log.Panicln(err)
	}

	err = u.Gui.SetKeybinding(
		cmdBar, gocui.KeyEnter, gocui.ModNone, cmdEnter)
	if err != nil {
		log.Panicln(err)
	}

	err = u.Gui.SetKeybinding(
		cmdBar, gocui.KeyEnter, gocui.ModNone, defaultEnter)
	if err != nil {
		log.Panicln(err)
	}
}

// quit quits the program.
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func switchTopView(g *gocui.Gui, v *gocui.View) error {
	var err error
	switch g.CurrentView().Name() {
	case requestView:
		_, err = g.SetCurrentView(responseView)
	case responseView:
		_, err = g.SetCurrentView(requestView)
	}
	return err
}

func switchCmdView(g *gocui.Gui, v *gocui.View) error {
	var err error
	switch g.CurrentView().Name() {
	case requestView, responseView:
		_, err = g.SetCurrentView(cmdBar)
	}
	return err
}

// switchModeCmd switches the current mode to CommandMode.
func switchModeCmd(g *gocui.Gui, v *gocui.View) error {
	cmdBarRefresh(g)
	if HttpuMode == CommandMode {
		return nil
	}
	HttpuMode = CommandMode
	cmdBarRefresh(g)
	return nil
}

// switchModeDefault switches the current mode to DefaultMode.
func switchModeDefault(g *gocui.Gui, v *gocui.View) error {
	cmdBarRefresh(g)
	if HttpuMode == DefaultMode {
		return nil
	}
	HttpuMode = DefaultMode
	cmdBarRefresh(g)
	return nil
}

func defaultEnter(g *gocui.Gui, v *gocui.View) error {
	if HttpuMode == CommandMode {
		return nil
	}

	rp := strings.Split(cmdBarBuffer(), ".")

	if len(rp) == 2 {
		req, ok := resource.Requests[rp[0]]
		if !ok {
			// There was no request find, we will fail sliently at this point.
			return nil
		}
		v, err := req.Variant(rp[1])
		if err != nil {
			return err
		}
		resp, stat, err := makeRequestWithVariant(&req, &v)
		if err != nil {
			return err
		}

		writeRequestDataVariant(&req, &v)
		writeResponseData(resp, stat)
	}

	if len(rp) == 1 {
		req, ok := resource.Requests[rp[0]]
		if !ok {
			// There was no request find, we will fail sliently at this point
			// also.
			return nil
		}

		resp, stat, err := makeRequest(&req)
		if err != nil {
			return err
		}

		writeRequestData(&req)
		writeResponseData(resp, stat)
	}

	return nil
}

func cmdEnter(g *gocui.Gui, v *gocui.View) error {
	if HttpuMode == DefaultMode {
		return nil
	}
	cmdStr := strings.Trim(
		strings.Replace(v.Buffer(), commandPromptMsg, "", 1), "\n")

	parts := strings.Split(cmdStr, " ")

	cmdName, args := parts[0], parts[1:]

	// TODO: Report an error here if the command is not found
	if cmd, ok := Commands[cmdName]; ok {
		err := cmd.Execute(g, cmdName, args)
		if err != nil {
			rv, _ := g.View(requestView)
			fmt.Fprintf(rv, "%s\n", err)
		}
	}
	return nil
}

func focusRequest(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetCurrentView(requestView)
	return err
}

func focusResponse(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetCurrentView(responseView)
	return err
}
