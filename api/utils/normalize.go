package utils

import (
	"golang.org/x/text/unicode/norm"
	"strings"
)

var arabicToPersian = map[rune]rune{
	'ك': 'ک',
	'ي': 'ی',
	'ۀ': 'ه',
	'ۂ': 'ه',
	'ۃ': 'ه',
	'ة': 'ه',
}

func NormalizeProfessor(name string) string {
	var builder strings.Builder
	builder.Grow(len(name)) // Optimize memory allocation

	for _, r := range name {
		if p, ok := arabicToPersian[r]; ok {
			builder.WriteRune(p)
		} else {
			builder.WriteRune(r)
		}
	}

	normalized := norm.NFC.String(builder.String())

	return strings.Join(strings.Fields(normalized), " ")
}
