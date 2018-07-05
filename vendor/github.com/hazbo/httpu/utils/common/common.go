package common

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
)

// ProjectPath represents the full path for the current project.
// TODO: This shouldn't really be stored here - just for now.
var ProjectPath string

// HomeDir gets the user's home directory.
func HomeDir() (string, error) {
	return homedir.Dir()
}

// MakeHttpuDir creates the ~/.httpu directory if it does not yet exists.
func MakeHttpuDir() error {
	hd, _ := HomeDir()
	pkgDir := fmt.Sprintf("%s/.httpu", hd)
	if _, err := os.Stat(pkgDir); err == nil {
		return nil
	}
	return os.MkdirAll(fmt.Sprintf("%s/.httpu", hd), os.ModePerm)
}

// PackagesExists checks to see if the packages directory within the .httpu
// directory exists.
func PackagesExists() bool {
	hd, _ := HomeDir()
	pkgDir := fmt.Sprintf("%s/.httpu/packages", hd)
	if _, err := os.Stat(pkgDir); err == nil {
		return true
	}
	return false
}
