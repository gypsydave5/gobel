package gobel

import (
	"fmt"
	"strconv"
	"strings"
)

func (p *Pair) String() string {
	if p == Nil {
		return "nil"
	}

	var s strings.Builder
	var actualString strings.Builder
	isActualString := true

	s.WriteString("(")
	actualString.WriteRune('"')

	for {
		s.WriteString(toString(p.First))
		if isActualString {
			if r, ok := p.First.(rune); ok {
				if r == '"' {
					actualString.WriteRune('\\')
				}
				actualString.WriteRune(r)
			} else {
				isActualString = false
			}
		}

		if p.Rest == Nil || p.Rest == nil {
			if isActualString {
				actualString.WriteRune('"')
				return actualString.String()
			}
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

	if isActualString {
		actualString.WriteRune('"')
		return actualString.String()
	}
	return s.String()
}

func (s Symbol) String() string {
	return s.Str
}

func toString(i interface{}) string {
	if v, ok := i.(int); ok {
		return strconv.Itoa(v)
	}

	if v, ok := i.(rune); ok {
		return fmt.Sprintf("\\%s", string(v))
	}

	return i.(fmt.Stringer).String()
}
