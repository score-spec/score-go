/*
Apache Score
Copyright 2020 The Apache Software Foundation

This product includes software developed at
The Apache Software Foundation (http://www.apache.org/).
*/
package types

// ContainersSpecs is a map of workload containers specifications.
type ContainersSpecs map[string]ContainerSpec

// ContainerSpec is a workload container specification.
type ContainerSpec struct {
	Image          string                             `json:"image"`
	Command        []string                           `json:"command"`
	Args           []string                           `json:"args"`
	Variables      map[string]string                  `json:"variables"`
	Files          []FileMountSpec                    `json:"files"`
	Volumes        []VolumeMountSpec                  `json:"volumes"`
	Resources      ContainerResourcesRequirementsSpec `json:"resources"`
	LivenessProbe  ContainerProbeSpec                 `json:"livenessProbe"`
	ReadinessProbe ContainerProbeSpec                 `json:"readinessProbe"`
}

// ContainerResourcesRequirementsSpec is a container resources requirements.
type ContainerResourcesRequirementsSpec struct {
	Limits   map[string]interface{} `json:"limits"`
	Requests map[string]interface{} `json:"requests"`
}

// FileMountSpec is a container's file mount specification.
type FileMountSpec struct {
	// The mounted file path and name.
	Target string `json:"target"`
	// The mounted file access mode.
	Mode string `json:"mode"`
	// File content, if scpecified, could be a single string or an array of strings (parsed as []interface{}, DEPRECATED).
	// This property can't be used if 'source' property is used.
	Content interface{} `json:"content"`
	// If specified, file content should be read from the local file, referenced by the source property.
	// The file path is always relative to the source Score file.
	// This property can't be used if 'content' property is used.
	Source string `json:"source"`
	// If set to true, the placeholders expansion will not occur in the contents of the file.
	NoExpand bool `json:"noExpand"`
}

// VolumeMountSpec is a container volume mount point specification.
type VolumeMountSpec struct {
	Source   string `json:"source"`
	Path     string `json:"path"`
	Target   string `json:"target"`
	ReadOnly bool   `json:"read_only"`
}

// ContainerProbeSpec is a container probe specification.
type ContainerProbeSpec struct {
	HTTPGet HTTPGetActionSpec `json:"httpGet"`
}

// HTTPGetActionSpec is an HTTP GET Action specification.
type HTTPGetActionSpec struct {
	Scheme      string           `json:"scheme"`
	Host        string           `json:"host"`
	Port        int              `json:"port"`
	Path        string           `json:"path"`
	HTTPHeaders []HTTPHeaderSpec `json:"httpHeaders"`
}

// HTTPHeaderSpec is an HTTP Header specification.
type HTTPHeaderSpec struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
