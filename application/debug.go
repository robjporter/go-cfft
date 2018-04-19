package application

import "github.com/Sirupsen/logrus"

func (a *Application) DEBUGOverrideLocalHXServer(server string) {
	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Overriding DB loaded information.")
	a.HX.Credentials.Url = server
	override = true
}
