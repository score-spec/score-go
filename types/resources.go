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
	Type     string                 `json:"type"`
	Metadata ResourceMeta           `json:"metadata,omitempty"`
	Params   map[string]interface{} `json:"params,omitempty"`
}

// ResourceMeta is an additional resource metadata.
type ResourceMeta struct {
	Annotations map[string]string `json:"annotations,omitempty"`
}
