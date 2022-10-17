package types

// ResourcesSpecs is a map of workload resources specifications.
type ResourcesSpecs map[string]ResourceSpec

// ResourceSpec is a resource specification.
type ResourceSpec struct {
	Type       string                          `json:"type"`
	Properties map[string]ResourcePropertySpec `json:"properties"`
}

// ResourcePropertySpec is a resource property specification.
type ResourcePropertySpec struct {
	Type     string      `json:"type"`
	Default  interface{} `json:"default"`
	Required bool        `json:"required"`
	Secret   bool        `json:"secret"`
}
