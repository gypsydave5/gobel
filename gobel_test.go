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
	expression := 1
	want := 1
	env := make(map[string]interface{})
	got := eval(expression, env)
	if got != want {
		t.Errorf("Expected %#v to evaluate to %#v but got %#v", expression, want, got)
	}
}
