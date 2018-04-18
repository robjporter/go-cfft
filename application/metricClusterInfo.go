package application

import (
	"strconv"
	"time"

	"../packages/hxconnect"

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
	err := a.HX.ClusterInfo()

	if err != nil {
		a.Logger.Debug("We were unable to collect the cluster information from HX Connect API.")
		a.LastError = err
		return &[]MetricInfo{}
	}

	if a.HX.GetResponseOK() {
		if a.HX.GetResponseCode() == 200 {
			var data []MetricInfo
			a.Logger.Debug("Querying HX Connect for Cluster information.")
			count := a.HX.GetResponseItemInt("#")
			for i := 0; i < count; i++ {
				var tmp MetricInfo
				tmp.ClusterName = a.HX.GetResponseItemString(strconv.Itoa(i) + ".config.name")
				tmp.VcenterDCName = a.HX.GetResponseItemString(strconv.Itoa(i) + ".config.vCenterDatacenter")
				tmp.VcenterClusterName = a.HX.GetResponseItemString(strconv.Itoa(i) + ".config.vCenterClusterName")
				tmp.UCSMOrgName = a.HX.GetResponseItemString(strconv.Itoa(i) + ".config.ucsmOrg")
				tmp.ClusterState = a.HX.GetResponseItemString(strconv.Itoa(i) + ".cluster.state")

				tmp.ClusterNodeSize = a.HX.GetResponseItemInt(strconv.Itoa(i) + ".stNodesSize")
				tmp.ClusterNodesActive = a.HX.GetResponseItemInt(strconv.Itoa(i) + ".cluster.activeNodes")
				tmp.ClusterModelNumbers = a.HX.GetResponseItemString(strconv.Itoa(i) + ".about.modelNumber")
				tmp.ClusterSerialNumbers = a.HX.GetResponseItemString(strconv.Itoa(i) + ".about.serialNumber")
				tmp.ClusterUpgradeState = a.HX.GetResponseItemString(strconv.Itoa(i) + ".upgradeState")

				boot := a.HX.GetResponseItemInt64(strconv.Itoa(i) + ".cluster.boottime")
				tmp.ClusterUptime = hxconnect.DiffString(boot)
				tmp.ClusterBootTime = time.Unix(boot, 0)

				tmp.ClusterRawCapacity = a.HX.GetResponseItemFloat(strconv.Itoa(i) + ".cluster.rawCapacity")
				tmp.ClusterCapacity = a.HX.GetResponseItemFloat(strconv.Itoa(i) + ".cluster.capacity")
				tmp.ClusterUsedCapacity = a.HX.GetResponseItemFloat(strconv.Itoa(i) + ".cluster.usedCapacity")
				tmp.ClusterFreeCapacity = a.HX.GetResponseItemFloat(strconv.Itoa(i) + ".cluster.freeCapacity")

				tmp.ClusterDowntime = a.HX.GetResponseItemFloat(strconv.Itoa(i) + ".cluster.downtime")
				tmp.ClusterAllFlash = a.HX.GetResponseItemBool(strconv.Itoa(i) + ".cluster.allFlash")
				tmp.ClusterRF = a.HX.GetResponseItemInt(strconv.Itoa(i) + ".config.dataReplicationFactor")
				tmp.ClusterPolicy = a.HX.GetResponseItemString(strconv.Itoa(i) + ".config.clusterAccessPolicy")

				data = append(data, tmp)
			}
			a.Logger.WithFields(logrus.Fields{"Cluster Count": count}).Debug("Querying HX Connect for Cluster information complete.")
			return &data
		}
		a.Logger.WithFields(logrus.Fields{"ResponseCode": a.HX.GetResponseCode()}).Warning("An unexpected response code was received for Cluster information.")
	} else {
		a.Logger.WithFields(logrus.Fields{"ResponseOK": false}).Warning("We received a failed attempt at connecting to the Cluster endpoint.")
	}
	return &[]MetricInfo{}
}
