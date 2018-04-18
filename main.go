// https://github.com/juju/errors
// https://github.com/jinzhu/now
// https://github.com/leekchan/accounting
// https://github.com/tidwall/gjson
// https://github.com/tidwall/sjson
// https://github.com/inconshreveable/go-update

// git tag -a v0.2.3 -m "First release"
// git tag -fa v0.2.3
// git push origin v0.2.3
// goreleaser

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
