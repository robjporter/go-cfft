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
	err := a.HX.ClusterPlatform()

	if err != nil {
		a.Logger.Debug("We were unable to collect the Platform information from HX Connect API.")
		a.LastError = err
		return &MetricPlatform{}
	}

	if a.HX.GetResponseOK() {
		if a.HX.GetResponseCode() == 200 {
			var metric MetricPlatform
			a.Logger.Debug("Querying HX Connect for Platform information.")

			metric.EffectiveCores = a.HX.GetResponseItemInt64("numCpuCores")
			metric.EffectiveCPU = a.HX.GetResponseItemInt64("effectiveCpu")
			metric.EffectiveMemory = a.HX.GetResponseItemInt64("effectiveMemory")
			metric.TotalCPU = a.HX.GetResponseItemInt64("totalCpu")
			metric.TotalMemory = a.HX.GetResponseItemInt64("totalMemory")
			metric.OverallStatus = a.HX.GetResponseItemString("overallStatus")

			a.Logger.WithFields(logrus.Fields{"Cores": metric.EffectiveCores, "CPU": metric.TotalCPU, "Memory": metric.TotalMemory, "Status": metric.OverallStatus}).Debug("Querying HX Connect for Platform information complete.")

			return &metric
		}
		a.Logger.WithFields(logrus.Fields{"ResponseCode": a.HX.GetResponseCode()}).Warning("An unexpected response code was received for Platform information.")
	} else {
		a.Logger.WithFields(logrus.Fields{"ResponseOK": false}).Warning("We received a failed attempt at connecting to the Platform endpoint.")
	}
	return &MetricPlatform{}
}
