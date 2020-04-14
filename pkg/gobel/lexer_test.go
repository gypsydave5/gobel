package gobel

import (
	"reflect"
	"strings"
	"testing"
)

func TestScanLexer(t *testing.T) {
	cases := []struct {
		name    string
		program string
		want    []string
	}{
		{"simple", "(+ 1 1 )", []string{"(", "+", "1", "1", ")"}},
		{"character", `\bel`, []string{`\bel`}},
		{"character quote", `\"`, []string{`\"`}},
		{"string with embedded quote", `"\""`, []string{`"\""`}},
		{"pair of nils", "(() ())", []string{"(", "(", ")", "(", ")", ")"}},
		{"pair of ones", "((1) (1))", []string{"(", "(", "1", ")", "(", "1", ")", ")"}},
		{"snake_case", "snake_case", []string{"snake_case"}},
		{"kebab-case-OK", "kebab-case", []string{"kebab-case"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := NewScanLexer(strings.NewReader(c.program))
			var tokens []string
			for ; !l.End(); l.Next() {
				tokens = append(tokens, l.Current())
			}
			if !reflect.DeepEqual(tokens, c.want) {
				t.Errorf("wanted %v but got %v", c.want, tokens)
			}
		})
	}
}
