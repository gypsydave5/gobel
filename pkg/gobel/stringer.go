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

	for {
		s.WriteString(toString(p.First))
		if p.Rest == Nil || p.Rest == nil {
			s.WriteString(")")
			return s.String()
		}
		if next, ok := p.Rest.(*Pair); ok {
			s.WriteString(" ")
			p = next
			continue
		}
		s.WriteString(" . ")
		s.WriteString(toString(p.Rest))
		s.WriteString(")")
		break
	}

	return s.String()
}

func toString(i interface{}) string {
	if v, ok := i.(int); ok {
		return strconv.Itoa(v)
	}
	return i.(fmt.Stringer).String()
}
