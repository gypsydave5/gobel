package gobel

import (
	"reflect"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	t.Run("reader macros", func(t *testing.T) {
		t.Run("quote", func(t *testing.T) {
			cases := []readCase{
				{"quote symbol", "'a", &Pair{&Symbol{"quote"}, &Pair{&Symbol{"a"}, Nil}}},
				{"quote list", "'(1)", &Pair{&Symbol{"quote"}, &Pair{&Pair{1, Nil}, Nil}}},
			}
			testReadCases(cases, t)
		})
	})

	t.Run("numbers", func(t *testing.T) {
		cases := []readCase{
			{"integer", "1", 1},
		}
		testReadCases(cases, t)
	})

	t.Run("symbols", func(t *testing.T) {
		cases := []readCase{
			{"symbol", "symbol", &Symbol{"symbol"}},
		}
		testReadCases(cases, t)
	})

	t.Run("characters", func(t *testing.T) {
		cases := []readCase{
			{"a", `\a`, 'a'},
			{"bel", `\bel`, '\a'},
			{"space", `\space`, ' '},
			{"tab", `\tab`, '\t'},
		}

		t.Run("alphanumeric", func(t *testing.T) {
			for i := 33; i <= 126; i++ {
				r := rune(i)

				t.Run(string(r), func(t *testing.T) {
					var s strings.Builder
					s.WriteRune('\\')
					s.WriteRune(r)
					got := Read(s.String())[0]
					if !reflect.DeepEqual(r, got) {
						t.Errorf("Expected %#v when reading '%s' but got %#v", r, s.String(), got)
					}
				})
			}
		})
		testReadCases(cases, t)
	})

	t.Run("lists", func(t *testing.T) {
		cases := []readCase{
			{"empty list", "()", Nil},
			{"nil", "nil", Nil},
			{"one item list", "(1)", &Pair{1, Nil}},
			{"two item list", "(1 2)", &Pair{1, &Pair{2, Nil}}},
			{"pair", "(1 . 2)", &Pair{1, 2}},
			{"three item list", "(1 2 3)", &Pair{1, &Pair{2, &Pair{3, Nil}}}},
			{"dotted list", "(1 2 . 3)", &Pair{1, &Pair{2, 3}}},
			{"nested list", "((1))", &Pair{&Pair{1, Nil}, Nil}},
			{"simplest two lists", "( () () )", &Pair{Nil, &Pair{Nil, Nil}}},
			{"simplest three lists", "( () () () )", &Pair{Nil, &Pair{Nil, &Pair{Nil, Nil}}}},
			{"list in second place", "( 1 () )", &Pair{1, &Pair{Nil, Nil}}},
			{"simple two lists", "( (1) (1) )", &Pair{&Pair{1, Nil}, &Pair{&Pair{1, Nil}, Nil}}},
			{"two lists", "(+ (+ 1 2) (+ 3 4))",
				&Pair{&Symbol{"+"}, &Pair{Read("(+ 1 2)")[0], &Pair{Read("(+ 3 4)")[0], Nil}}}},
		}

		testReadCases(cases, t)
	})

	t.Run("strings", func(t *testing.T) {
		cases := []readCase{
			{"simple string", `"hello"`, &Pair{'h', &Pair{'e', &Pair{'l', &Pair{'l', &Pair{'o', Nil}}}}}},
			{"string with space", `"h o"`, &Pair{'h', &Pair{' ', &Pair{'o', Nil}}}},
			{"string with quote", `"\""`, &Pair{'"', Nil}},
		}
		testReadCases(cases, t)
	})

	t.Run("misc", func(t *testing.T) {
		cases := []readCase{
			{"mixed", "(if nil 1 2)", &Pair{&Symbol{"if"}, &Pair{Nil, &Pair{1, &Pair{2, Nil}}}}},
		}
		testReadCases(cases, t)
	})
}

func testReadCases(cases []readCase, t *testing.T) {
	for i := range cases {
		c := cases[i]
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := Read(c.program)[0]
			if !reflect.DeepEqual(c.want, got) {
				t.Errorf("Expected %v but got %v", c.want, got)
			}
		})
	}
}

func TestEval(t *testing.T) {
	emptyEnv := NewEnv(nil)
	oneEnv := NewEnv(nil)
	oneEnv.set("one", 1)

	t.Run("types", func(t *testing.T) {
		cases := []evalCase{
			{"integer", []interface{}{1}, emptyEnv, 1},
			{"symbol", []interface{}{&Symbol{"one"}}, oneEnv, 1},
			{"multiple expressions", Read("1 2 3"), GlobalEnv(), 3},
		}

		testEvalCases(cases, t)
	})

	t.Run("addition", func(t *testing.T) {
		cases := []evalCase{
			{"addition", []interface{}{&Pair{&Symbol{"+"}, &Pair{1, &Pair{2, Nil}}}}, GlobalEnv(), 3},
			{"more addition", Read("(+ 1 2 3 4 5)"), GlobalEnv(), 15},
			{"empty addition", Read("(+)"), GlobalEnv(), 0},
			{"nested addition", Read("(+ (+ 2 2) (+ 3 3))"), GlobalEnv(), 10},
		}

		testEvalCases(cases, t)
	})

	t.Run("subtraction", func(t *testing.T) {
		cases := []evalCase{
			{"subtract", Read("(-)"), GlobalEnv(), 0},
			{"subtract", Read("(- 1)"), GlobalEnv(), -1},
			{"subtract", Read("(- 6 4)"), GlobalEnv(), 2},
			{"subtract", Read("(- 20 2 2 2)"), GlobalEnv(), 14},
			{"subtract", Read("(- 20 (+ 2 2 2) (- 10))"), GlobalEnv(), 24},
		}
		testEvalCases(cases, t)
	})

	t.Run("if", func(t *testing.T) {
		cases := []evalCase{
			{"if true", Read("(if 1 6 7)"), GlobalEnv(), 6},
			{"if nil", Read("(if nil 6 7)"), GlobalEnv(), 7},
			{"do not eval third if true", Read("(if 1 6 garbage)"), GlobalEnv(), 6},
			{"do not eval second if false", Read("(if nil rubbish 7)"), GlobalEnv(), 7},
			{"bel if", Read("(if nil rubbish nil more-rubbish 7 )"), GlobalEnv(), 7},
			{"bel if shortened", Read("(if nil rubbish)"), GlobalEnv(), Nil},
			{"bel if bit longer", Read("(if nil rubbish nil balls nil crap)"), GlobalEnv(), Nil},
		}
		testEvalCases(cases, t)
	})

	t.Run("quote", func(t *testing.T) {
		cases := []evalCase{
			{"quote", Read("(quote a)"), GlobalEnv(), &Symbol{"a"}},
		}
		testEvalCases(cases, t)
	})

	t.Run("set", func(t *testing.T) {
		cases := []evalCase{
			{"simple set", Read("(set (quote x) 1) x"), GlobalEnv(), 1},
			{"fancy quote set", Read("(set 'x 55) x"), GlobalEnv(), 55},
		}
		testEvalCases(cases, t)
	})
}

func testEvalCases(cases []evalCase, t *testing.T) {
	t.Helper()
	for i := range cases {
		c := cases[i]
		t.Run(c.name, func(t *testing.T) {
			got := Eval(c.expression, c.env)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Expected %v to evaluate to %v but got %v", c.expression[0], c.want, got)
			}
		})
	}
}

type readCase struct {
	name    string
	program string
	want    interface{}
}

type evalCase struct {
	name       string
	expression []interface{}
	env        *Env
	want       interface{}
}
