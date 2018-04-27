package application

type MetricAbout struct {
	Build    string
	DVersion string
	Instance string
	Locale   string
	Version  string
}

func (a *Application) metricGetAbout() *MetricAbout {
	res, err := a.HX.About()
	counter := a.Stats.IncreaseCounter("tasks")
	a.addToLogDebug(counter, nil, "Starting to get About metrics.")

	if err != nil {
		a.addToLogDebug(counter, map[string]interface{}{"Error": err}, "We were unable to collect the about information from HX Connect API.")
		a.LastError = err
		return &MetricAbout{}
	}

	if a.HX.GetResponseOK(res) {
		if a.HX.GetResponseCode(res) == 200 {
			a.addToLogDebug(counter, nil, "Querying HX Connect for About information.")

			build := a.HX.GetResponseItemString(res, "build")
			dversion := a.HX.GetResponseItemString(res, "displayVersion")
			instance := a.HX.GetResponseItemString(res, "instanceUuid")
			locale := a.HX.GetResponseItemString(res, "locale")
			version := a.HX.GetResponseItemString(res, "productVersion")

			a.addToLogDebug(counter, map[string]interface{}{"Build": build, "DisplayVersion": dversion, "Locale": locale, "Version": version}, "Querying HX Connect for About information complete.")

			return &MetricAbout{
				Build:    build,
				DVersion: dversion,
				Instance: instance,
				Locale:   locale,
				Version:  version,
			}
		}
		a.addToLogDebug(counter, map[string]interface{}{"ResponseCode": a.HX.GetResponseCode(res)}, "An unexpected response code was received for About information.")
	} else {
		a.addToLogDebug(counter, map[string]interface{}{"ResponseOK": false}, "We received a failed attempt at connecting to the About endpoint.")
	}
	return &MetricAbout{}
}
