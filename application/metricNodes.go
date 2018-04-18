package application

import (
	"strconv"

	"github.com/Sirupsen/logrus"
)

type MetricNode struct {
	NodeID                 string
	NodeName               string
	NodePosition           string
	NodeFullName           string
	NodeVersion            int
	NodeBuild              string
	NodeBuildName          string
	NodeDiskCount          int
	NodeOSVersion          string
	NodeOSName             string
	NodeSerial             string
	NodeOSBuild            string
	NodeStatus             string
	NodeType               string
	NodeSize               int
	NodePhyNic             int
	NodeCompressionSavings float64
	NodeDedupSavings       float64
	NodeTotalSavings       float64
}

type MetricNodes struct {
	NodeCount int
	Nodes     *[]MetricNode
}

func (m *MetricNode) setFullName(name string) {
	m.NodeFullName = name
}

func (m *MetricNode) setVersion(version int) {
	m.NodeVersion = version
}

func (a *Application) metricGetFurtherNodeInfo(nodes *MetricNodes) *MetricNodes {
	err := a.HX.ClusterAppliances()

	if err != nil {
		a.Logger.Debug("We were unable to collect the Appliance information from HX Connect API.")
		a.LastError = err
		return &MetricNodes{}
	}

	if a.HX.GetResponseOK() {
		if a.HX.GetResponseCode() == 200 {
			a.Logger.Debug("Querying HX Connect for Appliance information.")

			if a.HX.GetResponseItemInt("#") == nodes.NodeCount {
				var nodes2 MetricNodes
				var nodess []MetricNode
				nodes2.NodeCount = nodes.NodeCount
				for i := 0; i < nodes.NodeCount; i++ {
					id := a.HX.GetResponseItemString(strconv.Itoa(i) + ".nodes.A.entityRef.id")

					for _, e := range *nodes.Nodes {
						if e.NodeID == id {
							n := a.addNode(e, i, "A")
							n2 := a.addNode(e, i, "B")

							if n.NodeID != "" {
								nodess = append(nodess, n)
							}
							if n2.NodeID != "" {
								nodess = append(nodess, n2)
							}

							break
						}
					}
				}
				nodes2.Nodes = &nodess
				a.Logger.WithFields(logrus.Fields{}).Debug("Querying HX Connect for Appliance information complete.")
				return &nodes2
			}
			a.Logger.WithFields(logrus.Fields{"Number received": a.HX.GetResponseItemInt("#"), "Number expected": nodes.NodeCount}).Warning("The number of nodes reported is different to the number expected.")
		}
		a.Logger.WithFields(logrus.Fields{"ResponseCode": a.HX.GetResponseCode()}).Warning("An unexpected response code was received for Appliance information.")
	} else {
		a.Logger.WithFields(logrus.Fields{"ResponseOK": false}).Warning("We received a failed attempt at connecting to the Appliance endpoint.")
	}
	return &MetricNodes{}
}

func (a *Application) addNode(e MetricNode, count int, pos string) MetricNode {
	var n MetricNode
	counter := strconv.Itoa(count)
	n.NodeFullName = a.HX.GetResponseItemString(counter + ".nodes." + pos + ".entityRef.name")

	if n.NodeFullName == "" {
		return MetricNode{}
	}
	// Preset Values
	n.NodeID = e.NodeID
	n.NodeType = e.NodeType
	n.NodeName = e.NodeName
	n.NodeCompressionSavings = e.NodeCompressionSavings
	n.NodeDedupSavings = e.NodeDedupSavings
	n.NodeTotalSavings = e.NodeTotalSavings

	// Dynamic Values
	n.NodePosition = pos
	n.NodeVersion = a.HX.GetResponseItemInt(counter + ".nodes." + pos + ".pNode.version")
	n.NodeBuild = a.HX.GetResponseItemString(counter + ".nodes." + pos + ".pNode.about.build")
	n.NodeBuildName = a.HX.GetResponseItemString(counter + ".nodes." + pos + ".pNode.about.fullName")
	n.NodeStatus = a.HX.GetResponseItemString(counter + ".nodes." + pos + ".pNode.state")
	n.NodeDiskCount = a.HX.GetResponseItemInt(counter + ".disksSize")
	n.NodeOSVersion = a.HX.GetResponseItemString(counter + ".nodes." + pos + ".host.about.productVersion")
	n.NodeOSName = a.HX.GetResponseItemString(counter + ".nodes." + pos + ".host.about.name")
	n.NodeSerial = a.HX.GetResponseItemString(counter + ".nodes." + pos + ".host.about.serialNumber")
	n.NodeOSBuild = a.HX.GetResponseItemString(counter + ".nodes." + pos + ".host.about.build")
	n.NodeSize = a.HX.GetResponseItemInt(counter + ".nodesSize")
	n.NodePhyNic = a.HX.GetResponseItemInt(counter + ".pnicsSize")

	return n
}
