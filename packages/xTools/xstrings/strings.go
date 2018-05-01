package xstrings

import (
	"bytes"
	"unicode/utf8"
)

func Center(str, padding string, width uint) string {
	text := []rune(str)
	if len(text) >= int(width) {
		return str
	}
	padrunes := []rune(padding)

	out := make([]rune, int(width))
	pos := 0

	padwidth := int(width)/2 - 2 - len(text)/2
	if len(str)%2 == 0 {
		padwidth++
	}
	for i := 0; i < padwidth; i++ {
		out[pos] = padrunes[i%len(padrunes)]
		pos++
	}
	out[pos] = ' '
	pos++
	for i := 0; i < len(text); i++ {
		out[pos] = text[i]
		pos++
	}
	out[pos] = ' '
	pos++

	if len(str)%2 == 1 {
		padwidth++
	}
	for i := 0; i < padwidth; i++ {
		out[pos] = padrunes[i%len(padrunes)]
		pos++
	}
	return string(out)
}

func Len(str string) int {
	return utf8.RuneCountInString(str)
}

func writePadString(output *bytes.Buffer, pad string, padLen, remains int) {
	var r rune
	var size int

	repeats := remains / padLen

	for i := 0; i < repeats; i++ {
		output.WriteString(pad)
	}

	remains = remains % padLen

	if remains != 0 {
		for i := 0; i < remains; i++ {
			r, size = utf8.DecodeRuneInString(pad)
			output.WriteRune(r)
			pad = pad[size:]
		}
	}
}
