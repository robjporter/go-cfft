package application

import (
	"time"
	"bytes"
	"unicode/utf8"
	
	"../packages/ts"

	"github.com/labstack/echo"
)

func (a *Application) applicationBanner() string {
	size, _ := ts.GetSize()
	width := size.Col()
	tmp := "\033[H\033[2J"
	tmp += Center("  _____       ______    ______  _____",width," ") + "\n"
	tmp += Center(" /  __ \\      |  ___|   |  ___||_   _|",width," ") + "\n"
	tmp += Center(" | /  \\/      | |_      | |_     | |",width," ") + "\n"
	tmp += Center(" | |          |  _|     |  _|    | |",width," ") + "\n"
	tmp += Center(" | \\__/\\      | |       | |      | |",width," ") + "\n"
	tmp += Center("  \\____/apital\\_|lexible\\_|inance\\_/ool",width," ") + "\n"
	tmp += Center("===========================================",width," ") + "\n"
	tmp += Center("App Version: " + a.Versions.appVersion + " | Server Version: " + echo.Version,width," ") + "\n"
	tmp += Center("===========================================",width," ") + "\n"
	tmp += Center(time.Now().Format("Monday, 02-Jan-06 15:04:05"),width," ") + "\n"
	tmp += a.Server.Server.Addr + "\n"
	// TODO: FIX Server port not being displayed.
	return tmp
}

func Center(str string, length int, pad string) string {
	l := Len(str)

	if l >= length || pad == "" {
		return str
	}

	remains := length - l
	padLen := Len(pad)

	output := &bytes.Buffer{}
	output.Grow(len(str) + (remains/padLen+1)*len(pad))
	writePadString(output, pad, padLen, remains/2)
	output.WriteString(str)
	writePadString(output, pad, padLen, (remains+1)/2)
	return output.String()
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