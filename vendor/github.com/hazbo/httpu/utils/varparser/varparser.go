package varparser

import (
	"fmt"
	"strings"
)

type VarReplacer interface {
	Replace(k string) string
}

const (
	tknDelimiter    = '$'
	tknLeftBrace    = '{'
	tknRightBrace   = '}'
	tknLeftbracket  = '['
	tknRightBracket = ']'
)

// VarParser is respionsible for parsing variables passed through into
// configuration files using the following format:
//
// ${kind[key]} where 'kind' is equal to VarParser.kind and 'key' is equal to
// the key that is being searched for within the kind.
type VarParser struct {
	kind string
}

func New(k string) VarParser {
	return VarParser{kind: k}
}

// Parse is a small parser that takes a given variable, in this case something
// like ${stash[name]} and looks inside the stash store to see if name exists.
//
// If it does, it will replace the variable with the value associated with name,
// and if not, it will leave it as it is.
// TODO: This for now will
func (vp VarParser) Parse(s string, vr VarReplacer) string {
	if len(s) == 0 {
		return ""
	}
	for k, _ := range s {
		if k < len(s) {
			switch s[k] {
			case tknDelimiter:

				// vstart (variable start) build up what we're looking for tofind
				// a variable. It'll always start with '$', if that isn't found
				// initially in the switch, we move on. If it's found, we check
				// that there are atleast len(vstart) more characters in the string,
				// before checking to see if it is actually a variable.
				vstart := fmt.Sprintf("%c%c%s%c",
					tknDelimiter,
					tknLeftBrace,
					vp.kind,
					tknLeftbracket)

				// If there are not len(vstart) more characters that follow...
				if k+len(vstart) > len(s) {
					continue
				}

				var va []byte
				if s[k:k+len(vstart)] == vstart {
					// Try to find the string between '[' and ']' and append the
					// byte to the result.
					for i := k + len(vstart); s[i] != tknRightBracket; i++ {
						va = append(va, s[i])
					}
				} else {
					continue
				}

				term := fmt.Sprintf("%s%s%c%c",
					vstart,
					string(va),
					tknRightBracket,
					tknRightBrace)

				stashValue := vr.Replace(string(va))
				if stashValue == "" {
					continue
				}
				s = strings.Replace(s, term, stashValue, 1)
			default:
				continue
			}
		}
	}
	return s
}
