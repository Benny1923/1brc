package utils

import (
	"bytes"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func Split(text []byte) (name, temp []byte) {
	idx := bytes.IndexByte(text, ';')
	name = text[0:idx]
	temp = text[idx+1:]
	return
}

func ParseFloat(text []byte) float64 {
	var result float64
	var decimal float64
	length := len(text)

	// Check for negative sign
	start := 0
	minus := text[0] == '-'
	if minus {
		start = 1
	}

	// Process digits before decimal
	for i := start; i < length; i++ {
		b := text[i]
		if b == '.' {
			decimal = 1
			continue
		}
		if decimal == 0 {
			result = result*10 + float64(b-'0')
		} else {
			decimal *= 10
			result += float64(b-'0') / decimal
		}
	}

	if minus {
		return -result
	}
	return result
}

func RemoveDiacritics(s string) string {
	result := make([]byte, len(s))
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	nDst, _, err := t.Transform(result, []byte(s), true)
	if err != nil {
		panic(err)
	}
	return string(result[:nDst])
}
