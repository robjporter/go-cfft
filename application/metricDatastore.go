package application

import (
	"strconv"

	"github.com/Sirupsen/logrus"
)

type MetricDatastore struct {
	DatastoreID                         string
	DatastoreType                       string
	DatastoreName                       string
	DatastoreVirtualName                string
	DatastoreCapacity                   int64
	DatastoreNumMirrors                 int
	DatastoreCreationTime               int64
	DatastoreFreeCapacity               int64
	DatastoreUnsharedBytes              int64
	DatastoreUncompressedBytes          int64
	DatastoreVirtualMountSummary        bool
	DatastoreVirtualAccessibleSummary   bool
	DatastoreReplicationDatastorePaired bool
}

type MetricDatastores struct {
	DS                *[]MetricDatastore
	DatastoreCount    int
	TotalCapacity     int64
	TotalFreeCapacuty int64
	TotalMounted      int
	TotalAccessible   int
	TotalReplicated   int
}

func (a *Application) metricGetDatastores() *MetricDatastores {
	res,err := a.HX.ClusterDatastores()

	if err != nil {
		a.Logger.Debug("We were unable to collect the Datastore information from HX Connect API.")
		a.LastError = err
		return &MetricDatastores{}
	}

	if a.HX.GetResponseOK(res) {
		if a.HX.GetResponseCode(res) == 200 {
			a.Logger.Debug("Querying HX Connect for Datastore information.")
			var metric MetricDatastores
			var dss []MetricDatastore
			metric.DatastoreCount = a.HX.GetResponseItemInt(res,"#")
			for i := 0; i < metric.DatastoreCount; i++ {
				var ds MetricDatastore

				ds.DatastoreID = a.HX.GetResponseItemString(res,strconv.Itoa(i) + ".stPlatDatastore.entityRef.id")
				ds.DatastoreType = a.HX.GetResponseItemString(res,strconv.Itoa(i) + ".stPlatDatastore.entityRef.type")
				ds.DatastoreName = a.HX.GetResponseItemString(res,strconv.Itoa(i) + ".stPlatDatastore.entityRef.name")
				ds.DatastoreVirtualName = a.HX.GetResponseItemString(res,strconv.Itoa(i) + ".stPlatDatastore.virtDatastore.id")
				ds.DatastoreCapacity = a.HX.GetResponseItemInt64(res,strconv.Itoa(i) + ".stPlatDatastore.config.capacity")
				ds.DatastoreNumMirrors = a.HX.GetResponseItemInt(res,strconv.Itoa(i) + ".stPlatDatastore.config.numMirrors")
				ds.DatastoreCreationTime = a.HX.GetResponseItemInt64(res,strconv.Itoa(i) + ".stPlatDatastore.creationTime")
				ds.DatastoreFreeCapacity = a.HX.GetResponseItemInt64(res,strconv.Itoa(i) + ".stPlatDatastore.freeCapacity")
				ds.DatastoreUnsharedBytes = a.HX.GetResponseItemInt64(res,strconv.Itoa(i) + ".stPlatDatastore.unsharedUsedBytes")
				ds.DatastoreUncompressedBytes = a.HX.GetResponseItemInt64(res,strconv.Itoa(i) + ".stPlatDatastore.unCompressedUsedBytes")
				ds.DatastoreVirtualMountSummary = getDatastoreMountSummary(a.HX.GetResponseItemString(res,strconv.Itoa(i) + ".virtDatastore.mountSummary"))
				ds.DatastoreVirtualAccessibleSummary = getDatastoreAccessible(a.HX.GetResponseItemString(res,strconv.Itoa(i) + ".virtDatastore.accessibilitySummary"))
				ds.DatastoreReplicationDatastorePaired = getDatastoreReplication(a.HX.GetResponseItemString(res,strconv.Itoa(i) + ".replDatastore.paired"))

				dss = append(dss, ds)

				metric.TotalCapacity += ds.DatastoreCapacity
				metric.TotalFreeCapacuty += ds.DatastoreFreeCapacity
				if ds.DatastoreVirtualMountSummary {
					metric.TotalMounted++
				}
				if ds.DatastoreVirtualAccessibleSummary {
					metric.TotalAccessible++
				}
				if ds.DatastoreReplicationDatastorePaired {
					metric.TotalReplicated++
				}
			}
			metric.DS = &dss
			a.Logger.WithFields(logrus.Fields{}).Debug("Querying HX Connect for Datastore information complete.")
			return &metric
		}
		a.Logger.WithFields(logrus.Fields{"ResponseCode": a.HX.GetResponseCode(res)}).Warning("An unexpected response code was received for Datastore information.")
	} else {
		a.Logger.WithFields(logrus.Fields{"ResponseOK": false}).Warning("We received a failed attempt at connecting to the Datastore endpoint.")
	}
	return &MetricDatastores{}
}

func getDatastoreReplication(state string) bool {
	if state == "true" {
		return true
	}
	return false
}

func getDatastoreMountSummary(state string) bool {
	if state == "MOUNTED" {
		return true
	}
	return false
}

func getDatastoreAccessible(state string) bool {
	if state == "ACCESSIBLE" {
		return true
	}
	return false
}
