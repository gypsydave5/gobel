package gobel

import "strings"

func tokenize(s string) []string {
	s = strings.Replace(s, "(", " ( ", -1)
	s = strings.Replace(s, ")", " ) ", -1)
	ss := strings.Fields(s)
	return ss
}
