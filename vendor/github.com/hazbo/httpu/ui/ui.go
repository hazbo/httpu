package ui

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/hazbo/httpu"
	"github.com/jroimartin/gocui"
	"github.com/mitchellh/go-wordwrap"
)

// Ui is a wrapper for the gocui.Gui which is created and started at in main.
type Ui struct {
	Gui *gocui.Gui
}

const (
	cmdBar           = "cmd_bar"
	requestView      = "request_view"
	responseView     = "response_view"
	statusCodeView   = "status_code_view"
	requestTimeView  = "request_time_view"
	defaultPromptMsg = "(httpu) "
	commandPromptMsg = "(httpu) :"
	welcomeMessage   = "Welcome to httpu!"
)

// Mode is the mode of the command bar, it being either in default mode, or
// toggled into command mode, in which additional commands within httpu can be
// used.
//
// Default is always set to defaultMode, in which commands are not used.
type Mode int

const (
	DefaultMode Mode = iota
	CommandMode
)

var (
	// HttpuMode is always initially set to default.
	HttpuMode       = DefaultMode
	RequestView     *gocui.View
	ResponseView    *gocui.View
	CmdBarView      *gocui.View
	StatusCodeView  *gocui.View
	RequestTimeView *gocui.View

	readmeMsg = ""
)

// Toggle changes the mode from either default to command or the other way
// around.
func (m *Mode) Toggle() {
	switch HttpuMode {
	case DefaultMode:
		*m = CommandMode
	case CommandMode:
		*m = DefaultMode
	}
}

// New returns a new instance of UI with a pre-configured layout.
func New() Ui {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}

	g.SetManagerFunc(layout)

	return Ui{Gui: g}
}

// Start starts the main loop for the UI
func (u Ui) Start() {
	defer u.Gui.Close()

	u.setBindings()

	cmdBarSetup(u.Gui)
	requestViewSetup(u.Gui)

	if err := u.Gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (u Ui) Exit() error {
	return gocui.ErrQuit
}

// layout generates the default layout for httpu which is then passed into
// SetManagerFunc.
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	g.Cursor = true

	uis := NewUiSpec(maxX, maxY)

	createRequestView(g, uis.RequestViewSpec)
	err := createResponseView(g, uis.ResponseViewSpec)
	if err != nil {
		log.Fatal(err)
	}
	createCmdBar(g, uis.CmdBarViewSpec)

	err = createStatusCodeView(g, uis.StatusCodeViewSpec)
	if err != nil {
		log.Fatal(err)
	}

	err = createRequestTimeView(g, uis.RequestTimeViewSpec)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// createCmdBar creates the command bar at the bottom of the program where
// commands can be inputted.
//
// ┌──────────────────────────┐
// │                          │
// │        CommandBar        │
// │            |             │
// │            ▼             │
// │                          │
// │┌────────────────────────┐│
// │└────────────────────────┘│
// └──────────────────────────┘
func createCmdBar(g *gocui.Gui, vs ViewSpec) error {
	d := vs.Dimensions()
	v, err := g.SetView(cmdBar, d[0], d[1], d[2], d[3])
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	v.Editable = true
	v.Editor = gocui.EditorFunc(cmdBarEditor)
	v.Frame = true

	CmdBarView = v
	return nil
}

func cmdBarEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		cx, _ := v.Cursor()
		if cx > len(defaultPromptMsg) && HttpuMode == DefaultMode {
			v.EditDelete(true)
		}
		if cx > len(commandPromptMsg) && HttpuMode == CommandMode {
			v.EditDelete(true)
		}
	}
	defaultKeyPress(v)
}

// cmdBarRefresh clears the buffer of the cmdBar and replaces it with the
// default prompt message based on the current mode of httpu.
func cmdBarRefresh(g *gocui.Gui) {
	v, err := g.View(cmdBar)
	if err != nil {
		log.Fatal(err)
	}

	v.Clear()

	switch HttpuMode {
	case DefaultMode:
		fmt.Fprintf(v, "%s", defaultPromptMsg)
		v.SetCursor(len(defaultPromptMsg), 0)
		return
	case CommandMode:
		fmt.Fprintf(v, "%s", commandPromptMsg)
		v.SetCursor(len(commandPromptMsg), 0)
		return
	}
}

func cmdBarSetup(g *gocui.Gui) {
	g.Update(func(g *gocui.Gui) error {
		_, err := g.View(cmdBar)
		if err != nil {
			return err
		}
		cmdBarRefresh(g)
		_, err = g.SetCurrentView(cmdBar)
		return err
	})
}

func cmdBarBuffer() string {
	var msg string
	switch HttpuMode {
	case DefaultMode:
		msg = defaultPromptMsg
	case CommandMode:
		msg = commandPromptMsg
	}
	b := CmdBarView.Buffer()
	return strings.Trim(
		strings.Replace(b, msg, "", 1), "\n")
}

// createRequestView creates the left-hand view where request data can be seen
// and modified.
//
// ┌──────────────────────────┐
// │┌─────────┐               │
// ││         │               │
// ││         │               │
// ││         │  RequestView  │
// ││         │     ◄---      │
// ││         │               │
// │└─────────┘               │
// └──────────────────────────┘
func createRequestView(g *gocui.Gui, vs ViewSpec) error {
	d := vs.Dimensions()
	v, err := g.SetView(requestView, d[0], d[1], d[2], d[3])
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	v.Wrap = vs.Wrap
	v.Editable = true
	RequestView = v
	return nil
}

func requestViewSetup(g *gocui.Gui) {
	g.Update(func(g *gocui.Gui) error {
		rdmeFile := fmt.Sprintf("%s/README", httpu.Session().ProjectPath)
		if _, err := os.Stat(rdmeFile); err == nil {
			rdme, _ := ioutil.ReadFile(rdmeFile)
			readmeMsg = string(rdme)
			x, _ := RequestView.Size()

			wrapped := wordwrap.WrapString(readmeMsg, uint(x))

			fmt.Fprintln(RequestView, wrapped)
		} else {
			fmt.Fprintln(RequestView, welcomeMessage)
		}
		return nil
	})
}

// createResponseView creates the right-hand view where response data can be
// viewedafter a request has been made.
//
// ┌──────────────────────────┐
// │               ┌─────────┐│
// │               │         ││
// │               │         ││
// │ ResponseView  │         ││
// │     ---►      │         ││
// │               │         ││
// │               └─────────┘│
// └──────────────────────────┘
func createResponseView(g *gocui.Gui, vs ViewSpec) error {
	d := vs.Dimensions()
	v, err := g.SetView(responseView, d[0], d[1], d[2], d[3])
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	v.Wrap = vs.Wrap
	v.Editable = true
	ResponseView = v
	return nil
}

// createStatusCodeView creates the view in which the status code from an
// executed request will appear.
func createStatusCodeView(g *gocui.Gui, vs ViewSpec) error {
	d := vs.Dimensions()
	v, err := g.SetView(statusCodeView, d[0], d[1], d[2], d[3])
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	v.Wrap = vs.Wrap
	v.Title = vs.Title
	StatusCodeView = v
	return nil
}

// createRequestTimeView creates the view in which the time it has taken to make
// a request is displayed.
func createRequestTimeView(g *gocui.Gui, vs ViewSpec) error {
	d := vs.Dimensions()
	v, err := g.SetView(requestTimeView, d[0], d[1], d[2], d[3])
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	v.Title = vs.Title
	RequestTimeView = v
	return nil
}
