package printer

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/scanner"
)

// StringPrinter defines the behaviour of an object which must be able to
// 'PrintString', by a passed string, through to an instance of io.Writer.
type StringPrinter interface {
	PrintString(w io.Writer, a string)
}

// JSONPrinter is a StringPrinter which supports indentation and syntax
// highlighting for printing JSON to a terminal window using ANSI codes to
// define the syntax colours.
type JSONPrinter struct {
	depth  int
	spaces int
	b      bytes.Buffer
}

// NewJSONPrinter returns the default instance of JSONPrinter with a default
// indentation setting of 2 spaces.
func NewJSONPrinter() *JSONPrinter {
	return &JSONPrinter{
		depth:  0,
		spaces: 2,
	}
}

// PrintString prints a JSON string to an io.Writer, after scanning and syntax
// highlighting via ANSI codes has taken place.
func (jp *JSONPrinter) PrintString(w io.Writer, a string) {
	jp.printString(w, a)
}

func (jp *JSONPrinter) printString(w io.Writer, a string) {
	var s scanner.Scanner

	s.Init(strings.NewReader(a))

	for t := s.Scan(); t != scanner.EOF; t = s.Scan() {
		ct := s.TokenText()
		switch true {
		case jp.isStringTkn(ct):
			jp.b.WriteRune('"')
			jp.b.WriteString(Color(ct[1:len(ct)-1], ColorRed))
			jp.b.WriteRune('"')
		case jp.isBoolTkn(ct):
			jp.b.WriteString(Color(ct, ColorBlue))
		case ct == "{" || ct == "[":
			jp.b.WriteString(ct)
			jp.newline()
			jp.depth++
			jp.indent()
		case ct == "}" || ct == "]":
			jp.newline()
			jp.depth--
			jp.indent()
			jp.b.WriteString(ct)
		case ct == ",":
			jp.b.WriteString(ct)
			jp.newline()
			jp.indent()
		case ct == ":":
			jp.b.WriteString(ct)
		default:
			jp.b.WriteString(Color(ct, ColorGreen))
		}
	}
	fmt.Fprint(w, jp.b.String())
}

// isStringTkn checks to see if the given token is a string.
func (jp *JSONPrinter) isStringTkn(s string) bool {
	if len(s) != 0 && s[0] == '"' && s[len(s)-1] == '"' {
		return true
	}
	return false
}

// isBoolTkn checks to see if the given token is a boolean.
func (jp *JSONPrinter) isBoolTkn(s string) bool {
	return s == "true" || s == "false"
}

// indent writes empty space runes to the buffer depending on what the depth is
// and how many spaces are configured to be written as per the JSONPrinter.
func (jp *JSONPrinter) indent() {
	for i := 0; i < jp.depth*jp.spaces; i++ {
		jp.b.WriteRune(' ')
	}
}

// newline writes a newline rune to the buffer.
func (jp *JSONPrinter) newline() {
	jp.b.WriteRune('\n')
}
