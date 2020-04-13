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

var Nil *Pair = nil

type Tokenizer struct {
	tokens  []string
	current int
}

func NewTokenizer(program string) *Tokenizer {
	return &Tokenizer{tokenize(program), 0}
}

func (t *Tokenizer) Current() string {
	return t.tokens[t.current]
}

func (t *Tokenizer) Next() {
	t.current += 1
}

func (t *Tokenizer) End() bool {
	return t.current == len(t.tokens)
}

func Read(program string) []interface{} {
	var expressions []interface{}
	toks := NewTokenizer(program)
	for !toks.End() {
		e := readTokens(toks)
		expressions = append(expressions, e)
	}
	return expressions
}

func readTokens(toks *Tokenizer) interface{} {
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
	str = strings.Trim(str, `"`)
	rs := []rune(str)
	p := Nil
	for i := len(str) - 1; i > -1; i-- {
		p = cons(rs[i], p)
	}
	return p
}

func cons(i interface{}, p *Pair) *Pair {
	return &Pair{i, p}
}

func readList(toks *Tokenizer) *Pair {
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
	if expression == Nil {
		return Nil
	}

	if expression == nil {
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

	m["if"] = func(l *Pair, env Env) interface{} {
		condition := eval(l.First, env)
		if condition != Nil {
			return eval(l.Rest.(*Pair).First, env)
		}
		return eval(l.Rest.(*Pair).Rest.(*Pair).First, env)
	}

	return m
}
