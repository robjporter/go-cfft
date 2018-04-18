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
	err := a.HX.ClusterSummary()

	if err != nil {
		a.Logger.Debug("We were unable to collect the Summary information from HX Connect API.")
		a.LastError = err
		return &MetricSummary{}
	}

	if a.HX.GetResponseOK() {
		if a.HX.GetResponseCode() == 200 {
			var metric MetricSummary
			a.Logger.Debug("Querying HX Connect for Summary information.")

			metric.ReplicationCompliance = strings.Title(a.HX.GetResponseItemString("dataReplicationCompliance"))
			metric.SpaceStatus = a.HX.GetResponseItemString("spaceStatus")
			metric.TotalCapacity = a.HX.GetResponseItemInt64("totalCapacity")
			metric.UsedCapacity = a.HX.GetResponseItemInt64("usedCapacity")
			metric.FreeCapacity = a.HX.GetResponseItemInt64("freeCapacity")
			metric.CompressionSavings = a.HX.GetResponseItemFloat("compressionSavings")
			metric.DeduplicationSavings = a.HX.GetResponseItemFloat("deduplicationSavings")
			metric.TotalSavings = a.HX.GetResponseItemFloat("totalSavings")
			metric.ZoneType = a.HX.GetResponseItemInt("zoneType")
			metric.NumZones = a.HX.GetResponseItemInt("numZones")
			metric.HealthState = strings.Title(a.HX.GetResponseItemString("resiliencyInfo.state"))
			metric.HealthMessage = strings.TrimSpace(a.HX.GetResponseItemString("resiliencyInfo.messages.0"))

			a.Logger.WithFields(logrus.Fields{}).Debug("Querying HX Connect for Summary information complete.")

			return &metric
		}
		a.Logger.WithFields(logrus.Fields{"ResponseCode": a.HX.GetResponseCode()}).Warning("An unexpected response code was received for Summary information.")
	} else {
		a.Logger.WithFields(logrus.Fields{"ResponseOK": false}).Warning("We received a failed attempt at connecting to the Summary endpoint.")
	}

	return &MetricSummary{}
}
