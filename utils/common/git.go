package common

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	canonical   = "git@github.com:httpu/packages"
	localPkgDir = ".httpu/packages"
)

// GitCommand represents a new instance of a base git command.
type GitCommand struct {
	cmd *exec.Cmd
}

// NewGitCommand returns a new base git command which is then used to perform
// various git operations such as cloning and pulling.
func NewGitCommand() GitCommand {
	c := exec.Command("git")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return GitCommand{c}
}

// Clone will `git clone` the canonical packages repo to ~/.httpu/packages so
// that the user can access public packages / projects that are stored there.
func (gc *GitCommand) Clone() error {
	hd, err := HomeDir()
	if err != nil {
		return err
	}

	c := gc.cmd
	c.Args = append(
		c.Args, "clone", canonical, fmt.Sprintf("%s/%s", hd, localPkgDir))

	err = c.Start()

	if err != nil {
		return err
	}
	return c.Wait()
}

// Pull will `git pull` the latest updates from the packages repo to keep the
// user's local version up to date.
func (gc *GitCommand) Pull() error {
	hd, err := HomeDir()
	if err != nil {
		return err
	}

	c := gc.cmd
	c.Args = append(
		c.Args, "-C",
		fmt.Sprintf("%s/%s", hd, localPkgDir),
		"pull",
		"origin",
		"master")

	err = c.Start()

	if err != nil {
		return err
	}
	return c.Wait()
}
