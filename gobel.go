package gobel

import (
	"container/list"
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

func parse(program string) interface{} {
	toks := NewTokenizer(program)
	return readTokens(toks)
}

func readTokens(toks *Tokenizer) interface{} {
	if toks.Current() == "(" {
		toks.Next()
		return readExpression(toks)
	}
	a := atom(toks.Current())
	toks.Next()
	return a
}

func readExpression(toks *Tokenizer) *Pair {
	if toks.Current() == ")" {
		return Nil
	}
	head := Pair{}
	head.First = readTokens(toks)
	if toks.Current() == ")" {
		head.Rest = Nil
	} else {
		head.Rest = readExpression(toks)
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
	return Symbol{a}
}

func eval(expression interface{}, env map[string]interface{}) interface{} {
	i, _ := expression.(int)
	return i
}
