/*
Apache Score
Copyright 2020 The Apache Software Foundation

This product includes software developed at
The Apache Software Foundation (http://www.apache.org/).
*/
package types

// ResourcesSpecs is a map of workload resources specifications.
type ResourcesSpecs map[string]ResourceSpec

// ResourceSpec is a resource specification.
type ResourceSpec struct {
	Metadata   ResourceMeta                    `json:"metadata,omitempty"`
	Type       string                          `json:"type"`
	Properties map[string]ResourcePropertySpec `json:"properties"`
}

// ResourceMeta ia an additional resource metadata.
type ResourceMeta struct {
	Annotations map[string]string `json:"annotations,omitempty"`
}

// ResourcePropertySpec is a resource property specification.
type ResourcePropertySpec struct {
	Type     string      `json:"type"`
	Default  interface{} `json:"default"`
	Required bool        `json:"required"`
	Secret   bool        `json:"secret"`
}
