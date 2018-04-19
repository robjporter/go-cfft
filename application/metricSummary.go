package application

import (
	"strings"

	"github.com/Sirupsen/logrus"
)

type MetricSummary struct {
	ReplicationCompliance string
	SpaceStatus           string
	TotalCapacity         int64
	UsedCapacity          int64
	FreeCapacity          int64
	CompressionSavings    float64
	DeduplicationSavings  float64
	TotalSavings          float64
	ZoneType              int
	NumZones              int
	HealthState           string
	HealthMessage         string
}

func (a *Application) metricGetSummary() *MetricSummary {
	res, err := a.HX.ClusterSummary()
	taskCounterNumber := a.Stats.GetCounter("tasks")

	if err != nil {
		a.Logger.WithFields(logrus.Fields{"Task Number": taskCounterNumber}).Debug("We were unable to collect the Summary information from HX Connect API.")
		a.LastError = err
		return &MetricSummary{}
	}

	if a.HX.GetResponseOK(res) {
		if a.HX.GetResponseCode(res) == 200 {
			var metric MetricSummary
			a.Logger.WithFields(logrus.Fields{"Task Number": taskCounterNumber}).Debug("Querying HX Connect for Summary information.")

			metric.ReplicationCompliance = strings.Title(a.HX.GetResponseItemString(res, "dataReplicationCompliance"))
			metric.SpaceStatus = a.HX.GetResponseItemString(res, "spaceStatus")
			metric.TotalCapacity = a.HX.GetResponseItemInt64(res, "totalCapacity")
			metric.UsedCapacity = a.HX.GetResponseItemInt64(res, "usedCapacity")
			metric.FreeCapacity = a.HX.GetResponseItemInt64(res, "freeCapacity")
			metric.CompressionSavings = a.HX.GetResponseItemFloat(res, "compressionSavings")
			metric.DeduplicationSavings = a.HX.GetResponseItemFloat(res, "deduplicationSavings")
			metric.TotalSavings = a.HX.GetResponseItemFloat(res, "totalSavings")
			metric.ZoneType = a.HX.GetResponseItemInt(res, "zoneType")
			metric.NumZones = a.HX.GetResponseItemInt(res, "numZones")
			metric.HealthState = strings.Title(a.HX.GetResponseItemString(res, "resiliencyInfo.state"))
			metric.HealthMessage = strings.TrimSpace(a.HX.GetResponseItemString(res, "resiliencyInfo.messages.0"))

			a.Logger.WithFields(logrus.Fields{"Task Number": taskCounterNumber}).Debug("Querying HX Connect for Summary information complete.")

			return &metric
		}
		a.Logger.WithFields(logrus.Fields{"Task Number": taskCounterNumber, "ResponseCode": a.HX.GetResponseCode(res)}).Warning("An unexpected response code was received for Summary information.")
	} else {
		a.Logger.WithFields(logrus.Fields{"Task Number": taskCounterNumber, "ResponseOK": false}).Warning("We received a failed attempt at connecting to the Summary endpoint.")
	}

	return &MetricSummary{}
}
