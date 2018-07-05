package stash

import (
	"fmt"

	"github.com/hazbo/httpu/utils/varparser"
)

// StashValue represents a value to be stored in memory, that will replace a
// variable that exists within the resource configuration.
//
// e.g. if a Variant.Path is set to ${stash[path]}, if 'path' is stored in
// memory with an associated value, the variable will be replaced with that
// value.
type StashValue struct {
	Name          string   `json:"name"`
	Value         string   `json:"value"`
	JsonPath      []string `json:"jsonPath"`
	Origin        string   `json:"origin"`
	RepeatRequest bool     `json:"repeatRequest"`
}

// store is a map of stash values referenced by name.
type store map[string]StashValue

// Replace is used to replace a variable with a value stored inside of the
// stash.
func (s store) Replace(k string) string {
	return s[k].Value
}

// Store is an instance of the stash store.
var Store store = store{}

// StashValues represents multiple stash values.
type StashValues []StashValue

// Push pushes new stash values to the global stash store.
func (sv *StashValues) Push() {
	for _, s := range *sv {
		Store[s.Name] = s
	}
}

// Set sets a stash value in the map with an associated name.
func Set(key string, value StashValue) {
	Store[key] = value
}

// Get tries to lookup a stash value by name and returns it if it exists.
func Get(key string) (StashValue, error) {
	var (
		v  StashValue
		ok bool
	)
	if v, ok = Store[key]; !ok {
		return StashValue{}, fmt.Errorf("Value not found in stash.")
	}
	return v, nil
}

// Parse uses the built in varparser to find an instance of a variable, and in
// this case replace it with a value that exists with in the stash store.
func Parse(s string) string {
	return varparser.New("stash").Parse(s, Store)
}
