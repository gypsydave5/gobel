package gobel

import (
	"fmt"
	"strconv"
	"strings"
)

func (p *Pair) String() string {
	if p == Nil {
		return "()"
	}

	var s strings.Builder

	s.WriteString("(")
	s.WriteString(toString(p.First))
	s.WriteString(" . ")
	s.WriteString(toString(p.Rest))
	s.WriteString(")")
	return s.String()
}

func toString(i interface{}) string {
	if v, ok := i.(int); ok {
		return strconv.Itoa(v)
	}
	return i.(fmt.Stringer).String()
}
