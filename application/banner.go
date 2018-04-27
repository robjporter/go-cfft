package application

import (
	"fmt"
	"time"

	"../packages/ts"
	"../packages/xTools/xstrings"

	"github.com/labstack/echo"
)

func (a *Application) applicationDisplayBanner() {
	size, _ := ts.GetSize()
	width := size.Col()
	tmp := CLEARSCREEN
	tmp += xstrings.Center("  _____       ______    ______  _____", width, " ") + "\n"
	tmp += xstrings.Center(" /  __ \\      |  ___|   |  ___||_   _|", width, " ") + "\n"
	tmp += xstrings.Center(" | /  \\/      | |_      | |_     | |", width, " ") + "\n"
	tmp += xstrings.Center(" | |          |  _|     |  _|    | |", width, " ") + "\n"
	tmp += xstrings.Center(" | \\__/\\      | |       | |      | |", width, " ") + "\n"
	tmp += xstrings.Center("  \\____/apital\\_|lexible\\_|inance\\_/ool", width, " ") + "\n"
	tmp += xstrings.Center("===========================================", width, " ") + "\n"
	tmp += xstrings.Center("App Version: "+a.Versions.appVersion+" | Server Version: "+echo.Version, width, " ") + "\n"
	tmp += xstrings.Center("===========================================", width, " ") + "\n"
	tmp += xstrings.Center(time.Now().Format("Monday, 02-Jan-06 15:04:05"), width, " ") + "\n"
	tmp += a.Server.Server.Addr + "\n"
	// TODO: FIX Server port not being displayed.
	fmt.Println(tmp)
}
