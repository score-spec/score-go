package types

// ServiceSpec is a workload service specification.
type ServiceSpec struct {
	Ports ServicePortsSpecs `json:"ports"`
}

// ServicePortsSpecs is a map of named service ports specifications.
type ServicePortsSpecs map[string]ServicePortSpec

// ServicePortSpec is a service port specification.
type ServicePortSpec struct {
	Port       int    `json:"port"`
	TargetPort int    `json:"targetPort"`
	Protocol   string `json:"protocol"`
}
