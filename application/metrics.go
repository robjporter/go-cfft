package application

import (
	"time"

	"../packages/uuid"
	"../packages/xTools/carbon"

	"github.com/Sirupsen/logrus"
	"github.com/timshannon/bolthold"
)

type MetricData struct {
	UUID                     string
	CollectionTime           int64
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
	//metrics.CollectionTime = carbon.Now().ToDateTimeString()
	metrics.CollectionTime = carbon.Now().ToTimeStamp()
	u2, err := uuid.NewV4()
	if err != nil {
		a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Warn("There was an error producing a UUID for this update.  Therefore it will not be submitted.")
		metrics.UUID = "FAILED-TO-GENERATE-UUID-" + string(time.Now().Unix())
	} else {
		metrics.UUID = u2.String()
	}

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

	metrics.CollectionDuration = time.Since(time.Unix(metrics.CollectionTime, 0)).String()

	a.saveMetricsToLocalDB(metrics)

	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Metrics gathered successfully.")
}

func (a *Application) saveMetricsToLocalDB(m MetricData) {
	err := a.db.data.Insert(bolthold.NextSequence(), m)
	if err == nil {
		a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks"), "Saved": "MetricData"}).Debug("Saved information to DB successfully.")
	} else {
		a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks"), "Error": err}).Debug("Failed to save information to DB.")
	}
}
