package gobel

import (
	"container/list"
	"strconv"
	"strings"
)

func tokenize(s string) []string {
	// TODO: better tokenization that this old tut
	s = strings.Replace(s, "\\(", "\000", -1)
	s = strings.Replace(s, "(", " ( ", -1)
	s = strings.Replace(s, "\000", "\\(", -1)
	s = strings.Replace(s, "\\)", "\000", -1)
	s = strings.Replace(s, ")", " ) ", -1)
	s = strings.Replace(s, "\000", "\\)", -1)
	ss := strings.Fields(s)
	return ss
}

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

func cons(i interface{}, p *Pair) *Pair {
	return &Pair{i, p}
}

func readList(toks Lexer) *Pair {
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

func Eval(expressions []interface{}, env Env) interface{} {
	var r interface{}
	for i := range expressions {
		r = eval(expressions[i], env)
	}
	return r
}

func eval(expression interface{}, env Env) interface{} {
	if isNil(expression) {
		return Nil
	}

	i, ok := expression.(int)
	if ok {
		return i
	}

	s, ok := expression.(*Symbol)
	if ok {
		return env[s.Str]
	}

	p, ok := expression.(*Pair)
	if ok {
		f := eval(p.First, env)
		ff := f.(func(l *Pair, env Env) interface{})
		return ff(p.Rest.(*Pair), env)
	}

	return "WTF???"
}

type Env = map[string]interface{}

func DefaultEnv() Env {
	m := make(Env)

	m["+"] = func(l *Pair, env Env) interface{} {
		result := 0
		next := l
		for next != nil {
			result += eval(next.First, env).(int)
			next = next.Rest.(*Pair)
		}
		return result
	}

	m["-"] = func(l *Pair, env Env) interface{} {
		result := 0
		next := l
		if next == Nil {
			return 0
		}
		if next.Rest == Nil {
			return -eval(next.First, env).(int)
		}
		result = eval(next.First, env).(int)
		next = next.Rest.(*Pair)
		for next != Nil {
			result -= eval(next.First, env).(int)
			next = next.Rest.(*Pair)
		}
		return result
	}

	m["if"] = belIf

	return m
}

func belIf(l *Pair, env Env) interface{} {
	condition := eval(l.First, env)
	if !isNil(condition) {
		return eval(car(cdr(l).(*Pair)), env)
	}

	if v, ok := l.Rest.(*Pair).Rest.(*Pair); ok && isNil(v) {
		return Nil
	}

	if v, ok := l.Rest.(*Pair).Rest.(*Pair).Rest.(*Pair); ok && isNil(v) {
		return l.Rest.(*Pair).Rest.(*Pair).First
	}

	return belIf(l.Rest.(*Pair).Rest.(*Pair), env)
}

func id(a, b interface{}) bool {
	if ap, aok := a.(*Pair); aok {
		if bp, bok := b.(*Pair); bok {
			return bp == ap
		}
	}

	return false
}

func isNil(i interface{}) bool {
	return id(i, Nil)
}

func car(p *Pair) interface{} {
	return p.First
}

func cdr(p *Pair) interface{} {
	return p.Rest
}
