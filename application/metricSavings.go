package application

import (
	"strconv"
	"strings"

	"../packages/gjson"

	"github.com/Sirupsen/logrus"
)

type MetricSavings struct {
	ClusterCompressionSavings float64
	ClusterDedupSavings       float64
	ClusterTotalSavings       float64
}

func (a *Application) metricGetSavings() (*MetricSavings, *MetricNodes) {
	err := a.HX.ClusterSavings()

	if err != nil {
		a.Logger.Debug("We were unable to collect the savings information from HX Connect API.")
		a.LastError = err
		return &MetricSavings{}, &MetricNodes{}
	}

	if a.HX.GetResponseOK() {
		if a.HX.GetResponseCode() == 200 {
			var metric MetricSavings
			var nodes MetricNodes
			var nodess []MetricNode
			a.Logger.Debug("Querying HX Connect for Savings information.")
			tmp := strings.TrimLeft(a.HX.GetResponseData().(string), "{")
			tmp = strings.TrimRight(tmp, "}")
			result := gjson.Get(tmp, "..#").String()
			num, err := strconv.Atoi(result)

			if err == nil {
				for i := 0; i < num; i += 2 {
					tmp2 := gjson.Get(tmp, ".."+strconv.Itoa(i))
					tmp3 := gjson.Get(tmp, ".."+strconv.Itoa(i+1))
					tmp4 := strings.Split(tmp2.String(), ",")

					if len(tmp4) == 2 && strings.TrimLeft(tmp4[1], "type:") == "CLUSTER" {
						metric.ClusterCompressionSavings = gjson.Get(tmp3.String(), "compressionSavings").Float()
						metric.ClusterDedupSavings = gjson.Get(tmp3.String(), "dedupSavings").Float()
						metric.ClusterTotalSavings = gjson.Get(tmp3.String(), "totalNodeSavings").Float()
					} else if len(tmp4) == 3 {
						var node MetricNode

						nodes.NodeCount++
						node.NodeID = strings.TrimLeft(tmp4[0], "EntityRef(id:")
						node.NodeType = strings.TrimLeft(tmp4[1], " type:")
						tmp4[2] = strings.TrimLeft(tmp4[2], "name:")
						node.NodeName = strings.TrimRight(tmp4[2], ")")
						node.NodeCompressionSavings = gjson.Get(tmp3.String(), "compressionSavings").Float()
						node.NodeDedupSavings = gjson.Get(tmp3.String(), "dedupSavings").Float()
						node.NodeTotalSavings = gjson.Get(tmp3.String(), "totalNodeSavings").Float()
						nodess = append(nodess, node)
					}
				}
			}
			nodes.Nodes = &nodess
			a.Logger.WithFields(logrus.Fields{"Nodes Count": (num - 2) / 2}).Debug("Querying HX Connect for Savings information complete.")
			return &metric, &nodes
		}
		a.Logger.WithFields(logrus.Fields{"ResponseCode": a.HX.GetResponseCode()}).Warning("An unexpected response code was received for Savings information.")
	} else {
		a.Logger.WithFields(logrus.Fields{"ResponseOK": false}).Warning("We received a failed attempt at connecting to the Savings endpoint.")
	}
	return &MetricSavings{}, &MetricNodes{}
}
