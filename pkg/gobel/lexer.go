package gobel

import (
	"io"
	"text/scanner"
	"unicode"
)

// ScanLexer lexes a Bel program into tokens, represented as strings. Internally
// it wraps the Go `text/scanner.Scanner` and hacks around with it to make it fit.
type ScanLexer struct {
	scanner scanner.Scanner
	tok     rune
	current string
}

func NewScanLexer(r io.Reader) *ScanLexer {
	var s scanner.Scanner
	s.Init(r)
	s.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.ScanInts
	s.IsIdentRune = func(ch rune, i int) bool {
		return ch == '_' ||
			ch == '-' ||
			unicode.IsLetter(ch) ||
			unicode.IsDigit(ch) && i > 0
	}

	l := &ScanLexer{
		scanner: s,
	}
	l.Next()
	return l
}

func (l *ScanLexer) Current() string {
	return l.current
}

func (l *ScanLexer) Next() {
	l.tok = l.scanner.Scan()
	l.current = l.scanner.TokenText()
	if l.tok == '\\' { // small hack to handle Bel characters
		l.scanner.Scan()
		l.current += l.scanner.TokenText()
	}
}

func (l *ScanLexer) End() bool {
	return l.tok == scanner.EOF
}
