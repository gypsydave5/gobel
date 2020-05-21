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

type readCase struct {
	name    string
	program string
	want    interface{}
}
