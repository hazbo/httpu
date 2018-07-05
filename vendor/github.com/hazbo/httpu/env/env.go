package env

import (
	"os"

	"github.com/hazbo/httpu/utils/varparser"
)

// store is a map of stash values referenced by name.
type store string

// Replace is used to replace a variable with a value stored inside of the
// stash.
func (s store) Replace(k string) string {
	return os.Getenv(k)
}

// Store is an instance of the stash store.
var Store store

// Parse uses the built in varparser to find an instance of a variable, and in
// this case replace it with a value that exists with in the stash store.
func Parse(s string) string {
	return varparser.New("env").Parse(s, Store)
}
