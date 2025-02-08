package test

import "time"

//go:generate ../../../bin/cligen --struct-name DisplayServiceNode --pkg-name test --output-file example_gen.go
type DisplayServiceNode struct {
	Deployments      []*DeploymentStats `json:"deployments,omitempty"`
	Name             string             `json:"name" displayName:"Name" columnTag:"name"`
	ServiceMetadata  ServiceMetadata    `json:"service_metadata,omitempty"`
	Resources        Resources          `json:"resources,omitempty"`
	Status           ServiceStatus      `json:"status,omitempty"`
	Incidents        IncidentSummary    `json:"incidents,omitempty"`
	NetworkInfo      NetworkInfo        `json:"network_info,omitempty"`
	ContainerSummary ContainerSummary   `json:"container_summary,omitempty"`
}

type ServiceMetadata struct {
	Name          string    `json:"service_name" displayName:"Service Name" columnTag:"service_name"`
	Region        string    `json:"service_region" displayName:"Region" columnTag:"service_region"`
	ServiceID     string    `json:"service_id" displayName:"Service ID" columnTag:"service_id"`
	Environment   string    `json:"environment" displayName:"Environment" columnTag:"environment"`
	Version       string    `json:"version" displayName:"Version" columnTag:"version"`
	IsActive      bool      `json:"is_active" displayName:"Active" greenTexts:"true" redTexts:"false" columnTag:"is_active"`
	CreatedAt     time.Time `json:"created_at" displayName:"Created At" columnTag:"created_at"`
	LastUpdatedAt time.Time `json:"last_updated_at" displayName:"Last Updated At" columnTag:"last_updated_at"`
}

type Resources struct {
	CPU              string `json:"cpu" displayName:"CPU" columnTag:"cpu"`
	Memory           string `json:"memory" displayName:"Memory" columnTag:"memory"`
	Storage          string `json:"storage" displayName:"Storage" columnTag:"storage"`
	NetworkBandwidth string `json:"network_bandwidth" displayName:"Network Bandwidth" columnTag:"network_bandwidth"`
}

type ServiceStatus struct {
	State       string    `json:"state" displayName:"State" greenTexts:"Running" redTexts:"Stopped,Error" columnTag:"state"`
	Uptime      string    `json:"uptime" displayName:"Uptime" columnTag:"uptime"`
	LastRestart time.Time `json:"last_restart" displayName:"Last Restart" columnTag:"last_restart"`
	HealthCheck string    `json:"health_check" displayName:"Health Check" greenTexts:"Healthy" redTexts:"Unhealthy" columnTag:"health_check"`
}

type IncidentSummary struct {
	TotalIncidents    int        `json:"total_incidents" displayName:"# Total Incidents" columnTag:"total_incidents"`
	OpenIncidents     int        `json:"open_incidents" displayName:"# Open Incidents" columnTag:"open_incidents"`
	ResolvedIncidents int        `json:"resolved_incidents" displayName:"# Resolved Incidents" columnTag:"resolved_incidents"`
	ActiveIncidentID  string     `json:"active_incident_id,omitempty" displayName:"Active Incident ID" columnTag:"active_incident_id"`
	IncidentHistory   []Incident `json:"incident_history,omitempty"`
}

type Incident struct {
	IncidentID  string    `json:"incident_id" displayName:"Incident ID"`
	Status      string    `json:"incident_status" displayName:"Status"`
	ReportedAt  time.Time `json:"reported_at" displayName:"Reported At"`
	ResolvedAt  time.Time `json:"resolved_at,omitempty" displayName:"Resolved At"`
	Description string    `json:"description" displayName:"Description"`
	Priority    string    `json:"priority" displayName:"Priority" redTexts:"High" yellowTexts:"Medium" greenTexts:"Low"`
}

type NetworkInfo struct {
	IPAddress string `json:"ip_address" displayName:"IP Address" columnTag:"ip_address"`
	Subnet    string `json:"subnet" displayName:"Subnet" columnTag:"subnet"`
	Gateway   string `json:"gateway" displayName:"Gateway" columnTag:"gateway"`
	VPN       bool   `json:"vpn" displayName:"VPN Enabled" greenTexts:"true" redTexts:"false" columnTag:"vpn"`
}

type DeploymentStats struct {
	DeploymentID string    `json:"deployment_id" displayName:"Deployment ID"`
	StartTime    time.Time `json:"start_time" displayName:"Start Time"`
	EndTime      time.Time `json:"end_time" displayName:"End Time"`
	Status       string    `json:"status" displayName:"Status" greenTexts:"Success" redTexts:"Failed"`
	Version      string    `json:"version" displayName:"Version"`
}

type ContainerSummary struct {
	RunningContainers int         `json:"running_containers" displayName:"# Running Containers" columnTag:"running_containers"`
	StoppedContainers int         `json:"stopped_containers" displayName:"# Stopped Containers" columnTag:"stopped_containers"`
	ContainerDetails  []Container `json:"container_details,omitempty"`
}

type Container struct {
	Name         string `json:"container_name" displayName:"Container Name"`
	Image        string `json:"container_image" displayName:"Container Image"`
	State        string `json:"container_state" displayName:"State" greenTexts:"Running" redTexts:"Stopped"`
	CPUUsage     string `json:"cpu_usage" displayName:"CPU Usage"`
	MemUsage     string `json:"mem_usage" displayName:"Memory Usage"`
	RestartCount int    `json:"restart_count" displayName:"Restart Count"`
}
