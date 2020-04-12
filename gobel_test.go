package gobel

import (
	"reflect"
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

func TestParse(t *testing.T) {
	cases := []struct {
		name    string
		program string
		want    interface{}
	}{
		{"integer", "1", 1},
		{"symbol", "symbol", Symbol{"symbol"}},
		{"empty list", "()", Nil},
		{"one item list", "(1)", &Pair{1, nil}},
		{"two item list", "(1 2)", &Pair{1, &Pair{2, nil}}},
		{"three item list", "(1 2 3)", &Pair{1, &Pair{2, &Pair{3, nil}}}},
		{"nested list", "((1))", &Pair{&Pair{1, nil}, nil}},
		{"funky expression", "(if nil 1 2)", &Pair{Symbol{"if"}, &Pair{nil, &Pair{1, &Pair{2, nil}}}}},
	}

	for i := range cases {
		c := cases[i]
		t.Run(c.name, func(t *testing.T) {
			got := parse(c.program)
			if !reflect.DeepEqual(c.want, got) {
				t.Errorf("Expected %#v but got %#v", c.want, got)
			}
		})
	}
}

func TestEval(t *testing.T) {
	emptyEnv := make(map[string]interface{})
	oneEnv := make(map[string]interface{})
	oneEnv["one"] = 1

	cases := []struct {
		expression interface{}
		env        map[string]interface{}
		want       interface{}
	}{
		{1, emptyEnv, 1},
		{&Symbol{"one"}, oneEnv, 1},
		{&Pair{Symbol{"+"}, &Pair{1, &Pair{2, nil}}}, defaultEnv(), 3},
		{parse("(+ 1 2 3 4 5"), defaultEnv(), 15},
		{parse("(+)"), defaultEnv(), 0},
	}
	for _, c := range cases {
		got := eval(c.expression, c.env)
		if got != c.want {
			t.Errorf("Expected %#v to evaluate to %#v but got %#v", c.expression, c.want, got)
		}
	}
}
