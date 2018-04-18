package application

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/timshannon/bolthold"
)

type MetricData struct {
	CollectionTime           time.Time
	CollectionDuration       string
	Submitted                bool
	SubmittedOn              time.Time
	SubmittedTransactionCode string
	About                    *MetricAbout
	Info                     *[]MetricInfo
	Savings                  *MetricSavings
	Platform                 *MetricPlatform
	Datastores               *MetricDatastores
	Nodes                    *MetricNodes
	VM                       *MetricVirtualMachines
	Summary                  *MetricSummary
}

func (a *Application) gatherAndRecordMetrics() {
	var metrics MetricData
	metrics.Submitted = false
	metrics.CollectionTime = time.Now()

	if a.metricGetToken() {
		metrics.About = a.metricGetAbout()
		metrics.Info = a.metricGetClusterInfo()
		metrics.Savings, metrics.Nodes = a.metricGetSavings()
		metrics.Platform = a.metricGetPlatform()
		metrics.VM = a.metricGetVirtualMachines()
		metrics.Datastores = a.metricGetDatastores()
		metrics.Summary = a.metricGetSummary()
		metrics.Nodes = a.metricGetFurtherNodeInfo(metrics.Nodes)
	}

	metrics.CollectionDuration = time.Since(metrics.CollectionTime).String()

	a.saveMetricsToLocalDB(metrics)

	a.Logger.Debug("Metrics gathered successfully.")
}

func (a *Application) saveMetricsToLocalDB(m MetricData) {
	err := a.db.data.Insert(bolthold.NextSequence(), m)
	if err == nil {
		a.Logger.WithFields(logrus.Fields{"Saved": "MetricData"}).Debug("Saved information to DB successfully.")
	} else {
		a.Logger.WithFields(logrus.Fields{"Error": err}).Debug("Failed to save information to DB.")
	}
}
