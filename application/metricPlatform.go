package application

import "github.com/Sirupsen/logrus"

type MetricPlatform struct {
	EffectiveCPU    int64
	EffectiveMemory int64
	TotalCPU        int64
	TotalMemory     int64
	OverallStatus   string
	EffectiveCores  int64
}

func (a *Application) metricGetPlatform() *MetricPlatform {
	res, err := a.HX.ClusterPlatform()

	if err != nil {
		a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("We were unable to collect the Platform information from HX Connect API.")
		a.LastError = err
		return &MetricPlatform{}
	}

	if a.HX.GetResponseOK(res) {
		if a.HX.GetResponseCode(res) == 200 {
			var metric MetricPlatform
			a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Querying HX Connect for Platform information.")

			metric.EffectiveCores = a.HX.GetResponseItemInt64(res, "numCpuCores")
			metric.EffectiveCPU = a.HX.GetResponseItemInt64(res, "effectiveCpu")
			metric.EffectiveMemory = a.HX.GetResponseItemInt64(res, "effectiveMemory")
			metric.TotalCPU = a.HX.GetResponseItemInt64(res, "totalCpu")
			metric.TotalMemory = a.HX.GetResponseItemInt64(res, "totalMemory")
			metric.OverallStatus = a.HX.GetResponseItemString(res, "overallStatus")

			a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks"), "Cores": metric.EffectiveCores, "CPU": metric.TotalCPU, "Memory": metric.TotalMemory, "Status": metric.OverallStatus}).Debug("Querying HX Connect for Platform information complete.")

			return &metric
		}
		a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks"), "ResponseCode": a.HX.GetResponseCode(res)}).Warning("An unexpected response code was received for Platform information.")
	} else {
		a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks"), "ResponseOK": false}).Warning("We received a failed attempt at connecting to the Platform endpoint.")
	}
	return &MetricPlatform{}
}
