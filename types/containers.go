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
	Target  string   `json:"target"`
	Mode    string   `json:"mode"`
	Content []string `json:"content"`
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
	HTTPGet                       HTTPGetActionSpec `json:"httpGet"`
	InitialDelaySeconds           int32             `json:"initialDelaySeconds"`
	TimeoutSeconds                int32             `json:"timeoutSeconds"`
	PeriodSeconds                 int32             `json:"periodSeconds"`
	SuccessThreshold              int32             `json:"successThreshold"`
	FailureThreshold              int32             `json:"failureThreshold"`
	TerminationGracePeriodSeconds int64             `json:"terminationGracePeriodSeconds"`
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
