package gobel

import (
	"reflect"
	"testing"
)

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
			{"simple set", Read("(set x 1) x"), GlobalEnv(), 1},
			{"fancy quote set", Read("(set x 55) x"), GlobalEnv(), 55},
		}
		testEvalCases(cases, t)
	})

	t.Run("a simple procedure", func(t *testing.T) {
		cases := []evalCase{
			{"test-procedure", Read("(test-procedure 1 1)"), GlobalEnv(), 2},
		}
		testEvalCases(cases, t)
	})

	t.Run("lambda", func(t *testing.T) {
		cases := []evalCase{
			{"lambda", Read("((lambda (x) x) 1)"), GlobalEnv(), 1},
			{"lambda lambda", Read("((lambda (x) (+ x x)) 1)"), GlobalEnv(), 2},
			{"lambda the ultimate", Read("((lambda (x y) (+ x y)) 3 4)"), GlobalEnv(), 7},
		}
		testEvalCases(cases, t)
	})

	t.Run("define", func(t *testing.T) {
		cases := []evalCase{
			{"define double", Read("(define double (x) (+ x x)) (double 4)"), GlobalEnv(), 8},
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
				t.Fatalf("Expected %v to evaluate to %v but got %+v", c.expression[0], c.want, got)
			}
		})
	}
}

type evalCase struct {
	name       string
	expression []interface{}
	env        *Env
	want       interface{}
}
