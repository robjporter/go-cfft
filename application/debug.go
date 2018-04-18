package application


func (a *Application) DEBUGOverrideLocalHXServer(server string) {
	a.HX.Credentials.Url = server
	override = true
}
