package application

import (
	"time"
	
	"github.com/labstack/echo"
)

func (a *Application) applicationBanner() string {
	tmp := "  _____       ______    ______  _____" + "\n"
	tmp += " /  __ \\      |  ___|   |  ___||_   _|" + "\n"
	tmp += " | /  \\/      | |_      | |_     | |" + "\n"
	tmp += " | |          |  _|     |  _|    | |" + "\n"
	tmp += " | \\__/\\      | |       | |      | |" + "\n"
	tmp += "  \\____/apital\\_|lexible\\_|inance\\_/ool" + "\n"
	tmp += "===========================================" + "\n"
	tmp += "App Version: " + a.Versions.appVersion + " | Server Version: " + echo.Version + "\n"
	tmp += "===========================================" + "\n"
	tmp += time.Now().Format("Monday, 02-Jan-06 15:04:05")
	tmp += a.Server.Server.Addr + "\n"
	// TODO: FIX Server port not being displayed.
	return tmp
}
