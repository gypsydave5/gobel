package gobel

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"
)

func tokenize(s string) []string {
	s = strings.Replace(s, "(", " ( ", -1)
	s = strings.Replace(s, ")", " ) ", -1)
	ss := strings.Fields(s)
	return ss
}

type List struct {
	list.Element
}

type Pair struct {
	First interface{}
	Rest  *Pair
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

func Parse(program string) []interface{} {
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
	a := atom(toks.Current())
	toks.Next()
	return a
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
	} else {
		head.Rest = readList(toks)
	}
	return &head
}

func atom(a string) interface{} {
	if a == "nil" {
		return nil
	}
	i, err := strconv.Atoi(a)
	if err == nil {
		return i
	}
	return &Symbol{a}
}

func Eval(expressions []interface{}, env map[string]interface{}) interface{} {
	var r interface{}
	for i := range expressions {
		r = eval(expressions[i], env)
	}
	return r
}

func eval(expression interface{}, env map[string]interface{}) interface{} {
	if expression == nil {
		return nil
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
		ff := f.(func(l *Pair, env map[string]interface{}) interface{})
		return ff(p.Rest, env)
	}

	return "WTF???"
}

func DefaultEnv() map[string]interface{} {
	m := make(map[string]interface{})

	m["+"] = func(l *Pair, env map[string]interface{}) interface{} {
		result := 0
		next := l
		for next != nil {
			result += eval(next.First, env).(int)
			next = next.Rest
		}
		return result
	}

	m["-"] = func(l *Pair, env map[string]interface{}) interface{} {
		result := 0
		next := l
		if next == nil {
			return 0
		}
		if next.Rest == nil {
			return -eval(next.First, env).(int)
		}
		result = eval(next.First, env).(int)
		next = next.Rest
		for next != nil {
			result -= eval(next.First, env).(int)
			next = next.Rest
		}
		return result
	}

	m["if"] = func(l *Pair, env map[string]interface{}) interface{} {
		condition := eval(l.First, env)
		fmt.Println(condition)
		if condition != nil {
			return eval(l.Rest.First, env)
		}
		return eval(l.Rest.Rest.First, env)
	}

	return m
}
