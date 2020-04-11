package gobel

import (
	"container/list"
	"reflect"
	"strconv"
	"testing"
)

func TestTokenize(t *testing.T) {
	program := "(+ 1 1)"
	want := []string{"(", "+", "1", "1", ")"}
	got := tokenize(program)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Expected %#v but got %#v", want, got)
	}
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

func TestParse(t *testing.T) {
	t.Run("integer", func(t *testing.T) {
		program := "1"
		want := 1
		got := parse(program)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Expected %#v but got %#v", want, got)
		}
	})
	t.Run("symbol", func(t *testing.T) {
		program := "symbol"
		want := Symbol{"symbol"}
		got := parse(program)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Expected %#v but got %#v", want, got)
		}
	})
}

func parse(program string) interface{} {
	toks := tokenize(program)
	return readTokens(toks)

}

func readTokens(toks []string) interface{} {
	return atom(toks[0])
}

func atom(a string) interface{} {
	i, err := strconv.Atoi(a)
	if err == nil {
		return i
	}
	return Symbol{a}
}
