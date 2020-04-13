package gobel_test

import (
	g "github.com/gypsydave5/gobel/pkg/gobel"
	"testing"
)

func TestStringer(t *testing.T) {
	t.Run("pairs", func(t *testing.T) {
		cases := []struct {
			name     string
			list     *g.Pair
			stringed string
		}{
			{"nil", nil, "()"},
			{"Nil", g.Nil, "()"},
			{"simple pair", &g.Pair{1, 2}, "(1 . 2)"},
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
