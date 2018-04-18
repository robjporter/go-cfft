package application

import (
    "fmt"
	"time"
	"net/http"
    "github.com/labstack/echo"
)

func (a *Application) routesHomePing(c echo.Context) error {
    return c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
}

func (a *Application) routesHomeIndex(c echo.Context) error {
  return c.Render(200,"index.html",map[string]interface{}{"appname":"APPNAME","title":"TITLE"})
}

func (a *Application) customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	errorPage := fmt.Sprintf("%d.html", code)
	//c.Logger().Error(err)
	c.Render(code,errorPage,map[string]interface{}{"appname":"APPNAME","title":"TITLE"})
}