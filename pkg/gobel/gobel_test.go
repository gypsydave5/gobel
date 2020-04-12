package gobel

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	cases := []struct {
		program string
		want    []string
	}{
		{"(+ 1 1)", []string{"(", "+", "1", "1", ")"}},
		{"(() ())", []string{"(", "(", ")", "(", ")", ")"}},
		{"((1) (1))", []string{"(", "(", "1", ")", "(", "1", ")", ")"}},
	}

	for _, c := range cases {
		got := tokenize(c.program)

		if !reflect.DeepEqual(c.want, got) {
			t.Errorf("Expected %#v but got %#v", c.want, got)
		}
	}
}

func TestParse(t *testing.T) {
	cases := []struct {
		name    string
		program string
		want    interface{}
	}{
		{"integer", "1", 1},
		{"symbol", "symbol", &Symbol{"symbol"}},
		{"empty list", "()", Nil},
		{"nil", "nil", Nil},
		{"one item list", "(1)", &Pair{1, Nil}},
		{"two item list", "(1 2)", &Pair{1, &Pair{2, Nil}}},
		{"three item list", "(1 2 3)", &Pair{1, &Pair{2, &Pair{3, Nil}}}},
		{"nested list", "((1))", &Pair{&Pair{1, Nil}, Nil}},
		{"simplest two lists", "( () () )", &Pair{Nil, &Pair{Nil, Nil}}},
		{"simplest three lists", "( () () () )", &Pair{Nil, &Pair{Nil, &Pair{Nil, Nil}}}},
		{"list in second place", "( 1 () )", &Pair{1, &Pair{Nil, Nil}}},
		{"simple two lists", "( (1) (1) )", &Pair{&Pair{1, Nil}, &Pair{&Pair{1, Nil}, Nil}}},
		{"two lists", "(+ (+ 1 2) (+ 3 4))",
			&Pair{&Symbol{"+"}, &Pair{Parse("(+ 1 2)")[0], &Pair{Parse("(+ 3 4)")[0], Nil}}}},
		{"funky expression", "(if nil 1 2)", &Pair{&Symbol{"if"}, &Pair{Nil, &Pair{1, &Pair{2, Nil}}}}},
	}

	for i := range cases {
		c := cases[i]
		t.Run(c.name, func(t *testing.T) {
			got := Parse(c.program)[0]
			if !reflect.DeepEqual(c.want, got) {
				t.Errorf("Expected %#v but got %#v", c.want, got)
			}
		})
	}
}

type evalCase struct {
	name       string
	expression []interface{}
	env        map[string]interface{}
	want       interface{}
}

func TestEval(t *testing.T) {
	emptyEnv := make(map[string]interface{})
	oneEnv := make(map[string]interface{})
	oneEnv["one"] = 1

	cases := []evalCase{
		{"integer", []interface{}{1}, emptyEnv, 1},
		{"symbol", []interface{}{&Symbol{"one"}}, oneEnv, 1},
		{"addition", []interface{}{&Pair{&Symbol{"+"}, &Pair{1, &Pair{2, Nil}}}}, DefaultEnv(), 3},
		{"more addition", Parse("(+ 1 2 3 4 5)"), DefaultEnv(), 15},
		{"empty addition", Parse("(+)"), DefaultEnv(), 0},
		{"nested addition", Parse("(+ (+ 2 2) (+ 3 3))"), DefaultEnv(), 10},
		{"multiple expressions", Parse("1 2 3"), DefaultEnv(), 3},
	}

	testEvalCases(cases, t)

	t.Run("subtraction", func(t *testing.T) {
		cases := []evalCase{
			{"subtract", Parse("(-)"), DefaultEnv(), 0},
			{"subtract", Parse("(- 1)"), DefaultEnv(), -1},
			{"subtract", Parse("(- 6 4)"), DefaultEnv(), 2},
			{"subtract", Parse("(- 20 2 2 2)"), DefaultEnv(), 14},
			{"subtract", Parse("(- 20 (+ 2 2 2) (- 10))"), DefaultEnv(), 24},
		}

		testEvalCases(cases, t)
	})

	t.Run("if", func(t *testing.T) {
		cases := []evalCase{
			{"if true", Parse("(if 1 6 7)"), DefaultEnv(), 6},
			{"if nil", Parse("(if nil 6 7)"), DefaultEnv(), 7},
		}
		testEvalCases(cases, t)
	})
}

func testEvalCases(cases []evalCase, t *testing.T) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Eval(c.expression, c.env)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Expected %#v to evaluate to %#v but got %#v", c.expression, c.want, got)
			}
		})
	}
}
