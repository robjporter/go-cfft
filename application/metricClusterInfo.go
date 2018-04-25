package application

import (
	"strconv"
	"time"

	"../packages/xTools/hxconnect"

	"github.com/Sirupsen/logrus"
)

type MetricInfo struct {
	ClusterName        string
	VcenterDCName      string
	VcenterClusterName string
	UCSMOrgName        string
	ClusterState       string

	ClusterNodeSize      int
	ClusterNodesActive   int
	ClusterModelNumbers  string
	ClusterSerialNumbers string
	ClusterUpgradeState  string

	ClusterUptime       string
	ClusterBootTime     time.Time
	ClusterRawCapacity  float64
	ClusterCapacity     float64
	ClusterUsedCapacity float64
	ClusterFreeCapacity float64

	ClusterDowntime float64
	ClusterAllFlash bool
	ClusterRF       int
	ClusterPolicy   string
}

func (a *Application) metricGetClusterInfo() *[]MetricInfo {
	res, err := a.HX.ClusterInfo()

	if err != nil {
		a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("We were unable to collect the cluster information from HX Connect API.")
		a.LastError = err
		return &[]MetricInfo{}
	}

	if a.HX.GetResponseOK(res) {
		if a.HX.GetResponseCode(res) == 200 {
			var data []MetricInfo
			a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Querying HX Connect for Cluster information.")
			count := a.HX.GetResponseItemInt(res, "#")
			for i := 0; i < count; i++ {
				var tmp MetricInfo
				tmp.ClusterName = a.HX.GetResponseItemString(res, strconv.Itoa(i)+".config.name")
				tmp.VcenterDCName = a.HX.GetResponseItemString(res, strconv.Itoa(i)+".config.vCenterDatacenter")
				tmp.VcenterClusterName = a.HX.GetResponseItemString(res, strconv.Itoa(i)+".config.vCenterClusterName")
				tmp.UCSMOrgName = a.HX.GetResponseItemString(res, strconv.Itoa(i)+".config.ucsmOrg")
				tmp.ClusterState = a.HX.GetResponseItemString(res, strconv.Itoa(i)+".cluster.state")

				tmp.ClusterNodeSize = a.HX.GetResponseItemInt(res, strconv.Itoa(i)+".stNodesSize")
				tmp.ClusterNodesActive = a.HX.GetResponseItemInt(res, strconv.Itoa(i)+".cluster.activeNodes")
				tmp.ClusterModelNumbers = a.HX.GetResponseItemString(res, strconv.Itoa(i)+".about.modelNumber")
				tmp.ClusterSerialNumbers = a.HX.GetResponseItemString(res, strconv.Itoa(i)+".about.serialNumber")
				tmp.ClusterUpgradeState = a.HX.GetResponseItemString(res, strconv.Itoa(i)+".upgradeState")

				boot := a.HX.GetResponseItemInt64(res, strconv.Itoa(i)+".cluster.boottime")
				tmp.ClusterUptime = hxconnect.DiffString(boot)
				tmp.ClusterBootTime = time.Unix(boot, 0)

				tmp.ClusterRawCapacity = a.HX.GetResponseItemFloat(res, strconv.Itoa(i)+".cluster.rawCapacity")
				tmp.ClusterCapacity = a.HX.GetResponseItemFloat(res, strconv.Itoa(i)+".cluster.capacity")
				tmp.ClusterUsedCapacity = a.HX.GetResponseItemFloat(res, strconv.Itoa(i)+".cluster.usedCapacity")
				tmp.ClusterFreeCapacity = a.HX.GetResponseItemFloat(res, strconv.Itoa(i)+".cluster.freeCapacity")

				tmp.ClusterDowntime = a.HX.GetResponseItemFloat(res, strconv.Itoa(i)+".cluster.downtime")
				tmp.ClusterAllFlash = a.HX.GetResponseItemBool(res, strconv.Itoa(i)+".cluster.allFlash")
				tmp.ClusterRF = a.HX.GetResponseItemInt(res, strconv.Itoa(i)+".config.dataReplicationFactor")
				tmp.ClusterPolicy = a.HX.GetResponseItemString(res, strconv.Itoa(i)+".config.clusterAccessPolicy")

				data = append(data, tmp)
			}
			a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks"), "Cluster Count": count}).Debug("Querying HX Connect for Cluster information complete.")
			return &data
		}
		a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks"), "ResponseCode": a.HX.GetResponseCode(res)}).Warning("An unexpected response code was received for Cluster information.")
	} else {
		a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks"), "ResponseOK": false}).Warning("We received a failed attempt at connecting to the Cluster endpoint.")
	}
	return &[]MetricInfo{}
}
