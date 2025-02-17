// Code generated by msanath/gondolf/cligen. DO NOT EDIT.

package test

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/msanath/gondolf/pkg/duration"
	"github.com/msanath/gondolf/pkg/printer"
)

const (
	ColumnName              = "name"
	ColumnServiceName       = "service_name"
	ColumnServiceRegion     = "service_region"
	ColumnServiceId         = "service_id"
	ColumnEnvironment       = "environment"
	ColumnVersion           = "version"
	ColumnIsActive          = "is_active"
	ColumnCreatedAt         = "created_at"
	ColumnLastUpdatedAt     = "last_updated_at"
	ColumnCpu               = "cpu"
	ColumnMemory            = "memory"
	ColumnStorage           = "storage"
	ColumnNetworkBandwidth  = "network_bandwidth"
	ColumnState             = "state"
	ColumnUptime            = "uptime"
	ColumnLastRestart       = "last_restart"
	ColumnHealthCheck       = "health_check"
	ColumnTotalIncidents    = "total_incidents"
	ColumnOpenIncidents     = "open_incidents"
	ColumnResolvedIncidents = "resolved_incidents"
	ColumnActiveIncidentId  = "active_incident_id"
	ColumnIpAddress         = "ip_address"
	ColumnSubnet            = "subnet"
	ColumnGateway           = "gateway"
	ColumnVpn               = "vpn"
	ColumnRunningContainers = "running_containers"
	ColumnStoppedContainers = "stopped_containers"
)

func GetDisplayServiceNodeColumnTags() []string {
	return []string{
		ColumnName,
		ColumnServiceName,
		ColumnServiceRegion,
		ColumnServiceId,
		ColumnEnvironment,
		ColumnVersion,
		ColumnIsActive,
		ColumnCreatedAt,
		ColumnLastUpdatedAt,
		ColumnCpu,
		ColumnMemory,
		ColumnStorage,
		ColumnNetworkBandwidth,
		ColumnState,
		ColumnUptime,
		ColumnLastRestart,
		ColumnHealthCheck,
		ColumnTotalIncidents,
		ColumnOpenIncidents,
		ColumnResolvedIncidents,
		ColumnActiveIncidentId,
		ColumnIpAddress,
		ColumnSubnet,
		ColumnGateway,
		ColumnVpn,
		ColumnRunningContainers,
		ColumnStoppedContainers,
	}
}

func ValidateDisplayServiceNodeColumnTags(tags []string) error {
	validTags := GetDisplayServiceNodeColumnTags()
	for _, tag := range tags {
		if !slices.Contains(validTags, tag) {
			return fmt.Errorf("column tag '%s' not found. Valid tags are %v", tag, validTags)
		}
	}
	return nil
}

func (n *DisplayServiceNode) GetDisplayFieldFromColumnTag(columnTag string) (printer.DisplayField, error) {
	switch columnTag {
	case ColumnName:
		return n.GetName(), nil
	case ColumnServiceName:
		return n.ServiceMetadata.GetName(), nil
	case ColumnServiceRegion:
		return n.ServiceMetadata.GetRegion(), nil
	case ColumnServiceId:
		return n.ServiceMetadata.GetServiceID(), nil
	case ColumnEnvironment:
		return n.ServiceMetadata.GetEnvironment(), nil
	case ColumnVersion:
		return n.ServiceMetadata.GetVersion(), nil
	case ColumnIsActive:
		return n.ServiceMetadata.GetIsActive(), nil
	case ColumnCreatedAt:
		return n.ServiceMetadata.GetCreatedAt(), nil
	case ColumnLastUpdatedAt:
		return n.ServiceMetadata.GetLastUpdatedAt(), nil
	case ColumnCpu:
		return n.Resources.GetCPU(), nil
	case ColumnMemory:
		return n.Resources.GetMemory(), nil
	case ColumnStorage:
		return n.Resources.GetStorage(), nil
	case ColumnNetworkBandwidth:
		return n.Resources.GetNetworkBandwidth(), nil
	case ColumnState:
		return n.Status.GetState(), nil
	case ColumnUptime:
		return n.Status.GetUptime(), nil
	case ColumnLastRestart:
		return n.Status.GetLastRestart(), nil
	case ColumnHealthCheck:
		return n.Status.GetHealthCheck(), nil
	case ColumnTotalIncidents:
		return n.Incidents.GetTotalIncidents(), nil
	case ColumnOpenIncidents:
		return n.Incidents.GetOpenIncidents(), nil
	case ColumnResolvedIncidents:
		return n.Incidents.GetResolvedIncidents(), nil
	case ColumnActiveIncidentId:
		return n.Incidents.GetActiveIncidentID(), nil
	case ColumnIpAddress:
		return n.NetworkInfo.GetIPAddress(), nil
	case ColumnSubnet:
		return n.NetworkInfo.GetSubnet(), nil
	case ColumnGateway:
		return n.NetworkInfo.GetGateway(), nil
	case ColumnVpn:
		return n.NetworkInfo.GetVPN(), nil
	case ColumnRunningContainers:
		return n.ContainerSummary.GetRunningContainers(), nil
	case ColumnStoppedContainers:
		return n.ContainerSummary.GetStoppedContainers(), nil
	}
	return printer.DisplayField{}, fmt.Errorf("column tag '%s' not found. Valid tags are %v", columnTag, GetDisplayServiceNodeColumnTags())
}

func (n *IncidentSummary) GetTotalIncidents() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "# Total Incidents",
		ColumnTag:   "total_incidents",
		Value: func() string {
			str := strconv.Itoa(n.TotalIncidents)
			return str
		},
	}
}

func (n *IncidentSummary) GetOpenIncidents() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "# Open Incidents",
		ColumnTag:   "open_incidents",
		Value: func() string {
			str := strconv.Itoa(n.OpenIncidents)
			return str
		},
	}
}

func (n *IncidentSummary) GetResolvedIncidents() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "# Resolved Incidents",
		ColumnTag:   "resolved_incidents",
		Value: func() string {
			str := strconv.Itoa(n.ResolvedIncidents)
			return str
		},
	}
}

func (n *ContainerSummary) GetRunningContainers() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "# Running Containers",
		ColumnTag:   "running_containers",
		Value: func() string {
			str := strconv.Itoa(n.RunningContainers)
			return str
		},
	}
}

func (n *ContainerSummary) GetStoppedContainers() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "# Stopped Containers",
		ColumnTag:   "stopped_containers",
		Value: func() string {
			str := strconv.Itoa(n.StoppedContainers)
			return str
		},
	}
}

func (n *Container) GetRestartCount() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Restart Count",
		ColumnTag:   "",
		Value: func() string {
			str := strconv.Itoa(n.RestartCount)
			return str
		},
	}
}

func (n *ServiceMetadata) GetIsActive() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Active",
		ColumnTag:   "is_active",
		Value: func() string {
			str := strconv.FormatBool(n.IsActive)
			if str == "false" {
				return printer.RedText(str)
			}
			if str == "true" {
				return printer.GreenText(str)
			}
			return str
		},
	}
}

func (n *NetworkInfo) GetVPN() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "VPN Enabled",
		ColumnTag:   "vpn",
		Value: func() string {
			str := strconv.FormatBool(n.VPN)
			if str == "false" {
				return printer.RedText(str)
			}
			if str == "true" {
				return printer.GreenText(str)
			}
			return str
		},
	}
}

func (n *DeploymentStats) GetDeploymentID() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Deployment ID",
		ColumnTag:   "",
		Value: func() string {
			str := n.DeploymentID
			return str
		},
	}
}

func (n *DeploymentStats) GetStatus() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Status",
		ColumnTag:   "",
		Value: func() string {
			str := n.Status
			if str == "Failed" {
				return printer.RedText(str)
			}
			if str == "Success" {
				return printer.GreenText(str)
			}
			return str
		},
	}
}

func (n *DeploymentStats) GetVersion() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Version",
		ColumnTag:   "",
		Value: func() string {
			str := n.Version
			return str
		},
	}
}

func (n *DisplayServiceNode) GetName() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Name",
		ColumnTag:   "name",
		Value: func() string {
			str := n.Name
			return str
		},
	}
}

func (n *ServiceMetadata) GetName() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Service Name",
		ColumnTag:   "service_name",
		Value: func() string {
			str := n.Name
			return str
		},
	}
}

func (n *ServiceMetadata) GetRegion() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Region",
		ColumnTag:   "service_region",
		Value: func() string {
			str := n.Region
			return str
		},
	}
}

func (n *ServiceMetadata) GetServiceID() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Service ID",
		ColumnTag:   "service_id",
		Value: func() string {
			str := n.ServiceID
			return str
		},
	}
}

func (n *ServiceMetadata) GetEnvironment() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Environment",
		ColumnTag:   "environment",
		Value: func() string {
			str := n.Environment
			return str
		},
	}
}

func (n *ServiceMetadata) GetVersion() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Version",
		ColumnTag:   "version",
		Value: func() string {
			str := n.Version
			return str
		},
	}
}

func (n *Resources) GetCPU() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "CPU",
		ColumnTag:   "cpu",
		Value: func() string {
			str := n.CPU
			return str
		},
	}
}

func (n *Resources) GetMemory() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Memory",
		ColumnTag:   "memory",
		Value: func() string {
			str := n.Memory
			return str
		},
	}
}

func (n *Resources) GetStorage() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Storage",
		ColumnTag:   "storage",
		Value: func() string {
			str := n.Storage
			return str
		},
	}
}

func (n *Resources) GetNetworkBandwidth() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Network Bandwidth",
		ColumnTag:   "network_bandwidth",
		Value: func() string {
			str := n.NetworkBandwidth
			return str
		},
	}
}

func (n *ServiceStatus) GetState() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "State",
		ColumnTag:   "state",
		Value: func() string {
			str := n.State
			if str == "Stopped" {
				return printer.RedText(str)
			}
			if str == "Error" {
				return printer.RedText(str)
			}
			if str == "Running" {
				return printer.GreenText(str)
			}
			return str
		},
	}
}

func (n *ServiceStatus) GetUptime() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Uptime",
		ColumnTag:   "uptime",
		Value: func() string {
			str := n.Uptime
			return str
		},
	}
}

func (n *ServiceStatus) GetHealthCheck() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Health Check",
		ColumnTag:   "health_check",
		Value: func() string {
			str := n.HealthCheck
			if str == "Unhealthy" {
				return printer.RedText(str)
			}
			if str == "Healthy" {
				return printer.GreenText(str)
			}
			return str
		},
	}
}

func (n *IncidentSummary) GetActiveIncidentID() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Active Incident ID",
		ColumnTag:   "active_incident_id",
		Value: func() string {
			str := n.ActiveIncidentID
			return str
		},
	}
}

func (n *Incident) GetIncidentID() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Incident ID",
		ColumnTag:   "",
		Value: func() string {
			str := n.IncidentID
			return str
		},
	}
}

func (n *Incident) GetStatus() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Status",
		ColumnTag:   "",
		Value: func() string {
			str := n.Status
			return str
		},
	}
}

func (n *Incident) GetDescription() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Description",
		ColumnTag:   "",
		Value: func() string {
			str := n.Description
			return str
		},
	}
}

func (n *Incident) GetPriority() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Priority",
		ColumnTag:   "",
		Value: func() string {
			str := n.Priority
			if str == "High" {
				return printer.RedText(str)
			}
			if str == "Low" {
				return printer.GreenText(str)
			}
			if str == "Medium" {
				return printer.YellowText(str)
			}
			return str
		},
	}
}

func (n *NetworkInfo) GetIPAddress() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "IP Address",
		ColumnTag:   "ip_address",
		Value: func() string {
			str := n.IPAddress
			return str
		},
	}
}

func (n *NetworkInfo) GetSubnet() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Subnet",
		ColumnTag:   "subnet",
		Value: func() string {
			str := n.Subnet
			return str
		},
	}
}

func (n *NetworkInfo) GetGateway() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Gateway",
		ColumnTag:   "gateway",
		Value: func() string {
			str := n.Gateway
			return str
		},
	}
}

func (n *Container) GetName() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Container Name",
		ColumnTag:   "",
		Value: func() string {
			str := n.Name
			return str
		},
	}
}

func (n *Container) GetImage() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Container Image",
		ColumnTag:   "",
		Value: func() string {
			str := n.Image
			return str
		},
	}
}

func (n *Container) GetState() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "State",
		ColumnTag:   "",
		Value: func() string {
			str := n.State
			if str == "Stopped" {
				return printer.RedText(str)
			}
			if str == "Running" {
				return printer.GreenText(str)
			}
			return str
		},
	}
}

func (n *Container) GetCPUUsage() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "CPU Usage",
		ColumnTag:   "",
		Value: func() string {
			str := n.CPUUsage
			return str
		},
	}
}

func (n *Container) GetMemUsage() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Memory Usage",
		ColumnTag:   "",
		Value: func() string {
			str := n.MemUsage
			return str
		},
	}
}

func (n *DeploymentStats) GetStartTime() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Start Time",
		ColumnTag:   "",
		Value: func() string {
			str := n.StartTime.String()
			curTime := time.Now().UTC()
			agoTime := curTime.Sub(n.StartTime)
			str += " (" + duration.HumanDuration(agoTime) + " ago)"
			return str
		},
	}
}

func (n *DeploymentStats) GetEndTime() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "End Time",
		ColumnTag:   "",
		Value: func() string {
			str := n.EndTime.String()
			curTime := time.Now().UTC()
			agoTime := curTime.Sub(n.EndTime)
			str += " (" + duration.HumanDuration(agoTime) + " ago)"
			return str
		},
	}
}

func (n *ServiceMetadata) GetCreatedAt() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Created At",
		ColumnTag:   "created_at",
		Value: func() string {
			str := n.CreatedAt.String()
			curTime := time.Now().UTC()
			agoTime := curTime.Sub(n.CreatedAt)
			str += " (" + duration.HumanDuration(agoTime) + " ago)"
			return str
		},
	}
}

func (n *ServiceMetadata) GetLastUpdatedAt() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Last Updated At",
		ColumnTag:   "last_updated_at",
		Value: func() string {
			str := n.LastUpdatedAt.String()
			curTime := time.Now().UTC()
			agoTime := curTime.Sub(n.LastUpdatedAt)
			str += " (" + duration.HumanDuration(agoTime) + " ago)"
			return str
		},
	}
}

func (n *ServiceStatus) GetLastRestart() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Last Restart",
		ColumnTag:   "last_restart",
		Value: func() string {
			str := n.LastRestart.String()
			curTime := time.Now().UTC()
			agoTime := curTime.Sub(n.LastRestart)
			str += " (" + duration.HumanDuration(agoTime) + " ago)"
			return str
		},
	}
}

func (n *Incident) GetReportedAt() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Reported At",
		ColumnTag:   "",
		Value: func() string {
			str := n.ReportedAt.String()
			curTime := time.Now().UTC()
			agoTime := curTime.Sub(n.ReportedAt)
			str += " (" + duration.HumanDuration(agoTime) + " ago)"
			return str
		},
	}
}

func (n *Incident) GetResolvedAt() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Resolved At",
		ColumnTag:   "",
		Value: func() string {
			str := n.ResolvedAt.String()
			curTime := time.Now().UTC()
			agoTime := curTime.Sub(n.ResolvedAt)
			str += " (" + duration.HumanDuration(agoTime) + " ago)"
			return str
		},
	}
}
