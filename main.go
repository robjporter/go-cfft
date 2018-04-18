// https://github.com/juju/errors
// https://github.com/jinzhu/now
// https://github.com/leekchan/accounting
// https://github.com/tidwall/gjson
// https://github.com/tidwall/sjson

package main

import (
	"./application"
)

var (
	App *application.Application
)

func init() {
	App = application.New()
	App.DEBUGOverrideLocalHXServer("http://localhost:5003")
	App.Start()
}

func main() {

}
