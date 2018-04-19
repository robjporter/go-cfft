package application

import "github.com/Sirupsen/logrus"

type MetricAbout struct {
	Build    string
	DVersion string
	Instance string
	Locale   string
	Version  string
}

func (a *Application) metricGetAbout() *MetricAbout {
	res,err := a.HX.About()

	if err != nil {
		a.Logger.Debug("We were unable to collect the about information from HX Connect API.")
		a.LastError = err
		return &MetricAbout{}
	}

	if a.HX.GetResponseOK(res) {
		if a.HX.GetResponseCode(res) == 200 {
			a.Logger.Debug("Querying HX Connect for About information.")

			build := a.HX.GetResponseItemString(res,"build")
			dversion := a.HX.GetResponseItemString(res,"displayVersion")
			instance := a.HX.GetResponseItemString(res,"instanceUuid")
			locale := a.HX.GetResponseItemString(res,"locale")
			version := a.HX.GetResponseItemString(res,"productVersion")

			a.Logger.WithFields(logrus.Fields{"Build": build, "DisplayVersion": dversion, "Locale": locale, "Version": version}).Debug("Querying HX Connect for About information complete.")

			return &MetricAbout{
				Build:    build,
				DVersion: dversion,
				Instance: instance,
				Locale:   locale,
				Version:  version,
			}
		}
		a.Logger.WithFields(logrus.Fields{"ResponseCode": a.HX.GetResponseCode(res)}).Warning("An unexpected response code was received for About information.")
	} else {
		a.Logger.WithFields(logrus.Fields{"ResponseOK": false}).Warning("We received a failed attempt at connecting to the About endpoint.")
	}
	return &MetricAbout{}
}
