package application

import (
	"strconv"

	"github.com/Sirupsen/logrus"
)

type MetricVirtualMachine struct {
	Group             string
	PathName          string
	PowerState        bool
	MemoryUsage       int64
	StorageUncommited int64
	OverallCPUUsage   int64
	IPAddress         string
	GuestOS           string
	GuestID           string
	ConnectionState   bool
	Version           string
	CPUNumber         int
	InstanceID        string
	MemoryMB          int64
	GuestState        bool
	StorageCommitted  int64
	Name              string
	Host              string
	ID                string
	FullName          string
	UUID              string
	ResourcePool      string
}

type MetricVirtualMachines struct {
	VMs               *[]MetricVirtualMachine
	VMCount           int
	PoweredOn         int
	MemoryUsage       int64
	OverrallCPUUsage  int64
	ConnectedCount    int
	CommittedCPU      int
	CommittedMemory   int64
	StorageUncommited int64
	StorageCommitted  int64
}

func (a *Application) metricGetVirtualMachines() *MetricVirtualMachines {
	err := a.HX.ClusterVM()

	if err != nil {
		a.Logger.Debug("We were unable to collect the Virtual Machine information from HX Connect API.")
		a.LastError = err
		return &MetricVirtualMachines{}
	}

	if a.HX.GetResponseOK() {
		if a.HX.GetResponseCode() == 200 {
			a.Logger.Debug("Querying HX Connect for Virtual Machine information.")
			var metric MetricVirtualMachines
			var vms []MetricVirtualMachine
			metric.VMCount = a.HX.GetResponseItemInt("#")
			for i := 0; i < metric.VMCount; i++ {
				var vm MetricVirtualMachine

				vm.ConnectionState = getConnectionStateBool(a.HX.GetResponseItemString(strconv.Itoa(i) + ".runtime?connectionState"))
				vm.CPUNumber = a.HX.GetResponseItemInt(strconv.Itoa(i) + ".config?hardware?numCPU")
				vm.FullName = a.HX.GetResponseItemString(strconv.Itoa(i) + ".config?guestFullName")
				vm.Group = a.HX.GetResponseItemString(strconv.Itoa(i) + ".parent")
				vm.GuestID = a.HX.GetResponseItemString(strconv.Itoa(i) + ".config?guestId")
				vm.GuestOS = a.HX.GetResponseItemString(strconv.Itoa(i) + ".guest?guestFamily")
				vm.GuestState = getGuestStateBool(a.HX.GetResponseItemString(strconv.Itoa(i) + ".guest?guestState"))
				vm.Host = a.HX.GetResponseItemString(strconv.Itoa(i) + ".runtime?host")
				vm.ID = a.HX.GetResponseItemString(strconv.Itoa(i) + ".id")
				vm.InstanceID = a.HX.GetResponseItemString(strconv.Itoa(i) + ".config?instanceUuid")
				vm.IPAddress = a.HX.GetResponseItemString(strconv.Itoa(i) + ".guest?ipAddress")
				vm.MemoryMB = a.HX.GetResponseItemInt64(strconv.Itoa(i) + ".config?hardware?memoryMB")
				vm.MemoryUsage = a.HX.GetResponseItemInt64(strconv.Itoa(i) + ".summary?quickStats?guestMemoryUsage")
				vm.Name = a.HX.GetResponseItemString(strconv.Itoa(i) + ".name")
				vm.OverallCPUUsage = a.HX.GetResponseItemInt64(strconv.Itoa(i) + ".summary?quickStats?overallCpuUsage")
				vm.PathName = a.HX.GetResponseItemString(strconv.Itoa(i) + ".summary?config?vmPathName")
				vm.PowerState = getPowerStateBool(a.HX.GetResponseItemString(strconv.Itoa(i) + ".runtime?powerState"))
				vm.ResourcePool = a.HX.GetResponseItemString(strconv.Itoa(i) + ".resourcePool")
				vm.StorageCommitted = a.HX.GetResponseItemInt64(strconv.Itoa(i) + ".summary?storage?committed")
				vm.StorageUncommited = a.HX.GetResponseItemInt64(strconv.Itoa(i) + ".summary?storage?uncommitted")
				vm.UUID = a.HX.GetResponseItemString(strconv.Itoa(i) + ".config?uuid")
				vm.Version = a.HX.GetResponseItemString(strconv.Itoa(i) + ".config?version")

				vms = append(vms, vm)

				metric.CommittedCPU += vm.CPUNumber
				metric.CommittedMemory += vm.MemoryMB
				if vm.ConnectionState {
					metric.ConnectedCount++
				}
				metric.MemoryUsage += vm.MemoryUsage
				metric.OverrallCPUUsage += vm.OverallCPUUsage
				if vm.PowerState {
					metric.PoweredOn++
				}
				metric.StorageCommitted += vm.StorageCommitted
				metric.StorageUncommited += vm.StorageUncommited
			}
			metric.VMs = &vms
			a.Logger.WithFields(logrus.Fields{"Virtual Machines": metric.VMCount, "Powered Machines": metric.PoweredOn, "Committed vCPU": metric.CommittedCPU, "Committed vMemory": metric.CommittedMemory, "Committed Storage": metric.StorageCommitted}).Debug("Querying HX Connect for Virtual Machine information complete.")
			return &metric
		}
		a.Logger.WithFields(logrus.Fields{"ResponseCode": a.HX.GetResponseCode()}).Warning("An unexpected response code was received for Virtual Machine information.")
	} else {
		a.Logger.WithFields(logrus.Fields{"ResponseOK": false}).Warning("We received a failed attempt at connecting to the Virtual Machine endpoint.")
	}
	return &MetricVirtualMachines{}
}

func getPowerStateBool(state string) bool {
	if state == "poweredOn" {
		return true
	}
	return false
}

func getGuestStateBool(state string) bool {
	if state == "running" {
		return true
	}
	return false
}

func getConnectionStateBool(state string) bool {
	if state == "connected" {
		return true
	}
	return false
}
