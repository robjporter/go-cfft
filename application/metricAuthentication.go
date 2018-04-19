package application

import (
	"github.com/Sirupsen/logrus"
	"github.com/sanity-io/litter"
)

func (a *Application) metricGetToken() bool {
	res,err := a.HX.Authenticate()

	if err != nil {
		litter.Dump(a.HX)
		a.LastError = err
		a.Logger.WithFields(logrus.Fields{"Error": err}).Debug("We were unable to connect to the HX Connect API.")
		return false
	}

	if a.HX.GetResponseOK(res) {
		code := 201
		if override {
			code = 200
		}
		if a.HX.GetResponseCode(res) == code {
			a.Logger.Debug("Querying HX Connect for Authentication.")
			token := a.HX.GetResponseItem(res,"access_token")
			a.HX.SetToken(token.(string))
			return true
		}
		a.Logger.WithFields(logrus.Fields{"ResponseCode": a.HX.GetResponseCode(res)}).Warning("An unexpected response code was received Authentication information.")
	} else {
		a.Logger.WithFields(logrus.Fields{"ResponseOK": false}).Warning("We received a failed attempt at connecting to the Authentication endpoint.")
	}
	return false
}
