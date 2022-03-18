package gobel

import (
	"container/list"
	"strconv"
	"strings"
)

type List struct {
	list.Element
}

type Pair struct {
	First interface{}
	Rest  interface{}
}

type Symbol struct {
	Str string
}

// Nil is a more Lispy nil than `nil` - it's a nil *Pair.
var Nil *Pair = nil

// Lexer describes a simple lexer for the Bel language. It can return the current token
// as a string, move to the next token, and flag when the input is at an end.
type Lexer interface {
	Current() string
	Next()
	End() bool
}

func Read(program string) []interface{} {
	var expressions []interface{}
	toks := NewScanLexer(strings.NewReader(program))
	for !toks.End() {
		e := readTokens(toks)
		expressions = append(expressions, e)
	}
	return expressions
}

func readTokens(toks Lexer) interface{} {
	if toks.Current() == "'" {
		toks.Next()
		return &Pair{&Symbol{"quote"}, &Pair{readTokens(toks), Nil}}
	}
	if toks.Current() == "(" {
		toks.Next()
		return readList(toks)
	}
	if strings.HasPrefix(toks.Current(), `"`) {
		s := aString(toks.Current())
		toks.Next()
		return s
	}
	a := atom(toks.Current())
	toks.Next()
	return a
}

func aString(str string) *Pair {
	str = strings.TrimPrefix(str, `"`)
	str = strings.TrimSuffix(str, `"`)
	rs := []rune(str)
	p := Nil
	for i := len(str) - 1; i > -1; i-- {
		p = cons(rs[i], p)
		if rs[i] == '"' || rs[i] == '\\' { // jump past the escape character when reading in
			i--
		}
	}
	return p
}

func readList(toks Lexer) *Pair {
	// () is an alias for Nil
	if toks.Current() == ")" {
		toks.Next()
		return Nil
	}
	head := Pair{}

	head.First = readTokens(toks)
	if toks.Current() == ")" {
		head.Rest = Nil
		toks.Next()
	} else if toks.Current() == "." {
		toks.Next()
		head.Rest = readTokens(toks)
	} else {
		head.Rest = readList(toks)
	}
	return &head
}

func atom(a string) interface{} {
	if a == "nil" {
		return Nil
	}
	i, err := strconv.Atoi(a)
	if err == nil {
		return i
	}
	if a[0] == '\\' {
		return charCodeLookup(a[1:])
	}

	return &Symbol{a}
}

func charCodeLookup(s string) rune {
	if len(s) == 1 {
		// length is 1 so assume this is a direct mapping of rune to rune
		// return the first rune in the string
		for _, r := range s {
			return r
		}
	}

	if s == "bel" {
		return '\a'
	}

	if s == "space" {
		return ' '
	}

	if s == "tab" {
		return '\t'
	}

	return '\000'
}
