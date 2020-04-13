package gobel_test

import (
	g "github.com/gypsydave5/gobel/pkg/gobel"
	"testing"
)

func TestStringer(t *testing.T) {
	t.Run("symbols", func(t *testing.T) {
		t.Parallel()
		s := &g.Symbol{"symbol"}
		if s.String() != s.Str {
			t.Errorf("Expected %q but got %q", s.String(), s.Str)
		}
	})
	t.Run("pairs", func(t *testing.T) {
		cases := []struct {
			name     string
			list     *g.Pair
			stringed string
		}{
			{"nil", nil, "()"},
			{"Nil", g.Nil, "()"},
			{"simple pair", &g.Pair{1, 2}, "(1 . 2)"},
			{"simple proper list", &g.Pair{1, g.Nil}, "(1)"},
			{"two item proper list", &g.Pair{1, &g.Pair{2, g.Nil}}, "(1 2)"},
			{"three item proper list", g.Read("(1 2 3)")[0].(*g.Pair), "(1 2 3)"},
			{"nested lists", g.Read("((1) (2 (3)))")[0].(*g.Pair), "((1) (2 (3)))"},
			{"dotted list", &g.Pair{1, &g.Pair{2, 3}}, "(1 2 . 3)"},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				if c.list.String() != c.stringed {
					t.Errorf("Expected %q but got %q", c.stringed, c.list.String())
				}
			})
		}
	})
}
