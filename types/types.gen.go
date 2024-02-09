// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package types

import "encoding/json"
import "fmt"
import "reflect"

// The container name.
type Container struct {
	// If specified, overrides container entry point arguments.
	Args []string `json:"args,omitempty" yaml:"args,omitempty" mapstructure:"args,omitempty"`

	// If specified, overrides container entry point.
	Command []string `json:"command,omitempty" yaml:"command,omitempty" mapstructure:"command,omitempty"`

	// The extra files to mount.
	Files []ContainerFilesElem `json:"files,omitempty" yaml:"files,omitempty" mapstructure:"files,omitempty"`

	// The image name and tag.
	Image string `json:"image" yaml:"image" mapstructure:"image"`

	// The liveness probe for the container.
	LivenessProbe *ContainerProbe `json:"livenessProbe,omitempty" yaml:"livenessProbe,omitempty" mapstructure:"livenessProbe,omitempty"`

	// The readiness probe for the container.
	ReadinessProbe *ContainerProbe `json:"readinessProbe,omitempty" yaml:"readinessProbe,omitempty" mapstructure:"readinessProbe,omitempty"`

	// The compute resources for the container.
	Resources *ContainerResources `json:"resources,omitempty" yaml:"resources,omitempty" mapstructure:"resources,omitempty"`

	// The environment variables for the container.
	Variables ContainerVariables `json:"variables,omitempty" yaml:"variables,omitempty" mapstructure:"variables,omitempty"`

	// The volumes to mount.
	Volumes []ContainerVolumesElem `json:"volumes,omitempty" yaml:"volumes,omitempty" mapstructure:"volumes,omitempty"`
}

type ContainerFilesElem struct {
	// The inline content for the file.
	Content interface{} `json:"content,omitempty" yaml:"content,omitempty" mapstructure:"content,omitempty"`

	// The file access mode.
	Mode *string `json:"mode,omitempty" yaml:"mode,omitempty" mapstructure:"mode,omitempty"`

	// If set to true, the placeholders expansion will not occur in the contents of
	// the file.
	NoExpand *bool `json:"noExpand,omitempty" yaml:"noExpand,omitempty" mapstructure:"noExpand,omitempty"`

	// The relative or absolute path to the content file.
	Source *string `json:"source,omitempty" yaml:"source,omitempty" mapstructure:"source,omitempty"`

	// The file path and name.
	Target string `json:"target" yaml:"target" mapstructure:"target"`
}

type ContainerProbe struct {
	// HttpGet corresponds to the JSON schema field "httpGet".
	HttpGet *HttpProbe `json:"httpGet,omitempty" yaml:"httpGet,omitempty" mapstructure:"httpGet,omitempty"`
}

// The compute resources for the container.
type ContainerResources struct {
	// The maximum allowed resources for the container.
	Limits *ResourcesLimits `json:"limits,omitempty" yaml:"limits,omitempty" mapstructure:"limits,omitempty"`

	// The minimal resources required for the container.
	Requests *ResourcesLimits `json:"requests,omitempty" yaml:"requests,omitempty" mapstructure:"requests,omitempty"`
}

// The environment variables for the container.
type ContainerVariables map[string]string

type ContainerVolumesElem struct {
	// An optional sub path in the volume.
	Path *string `json:"path,omitempty" yaml:"path,omitempty" mapstructure:"path,omitempty"`

	// Indicates if the volume should be mounted in a read-only mode.
	ReadOnly *bool `json:"read_only,omitempty" yaml:"read_only,omitempty" mapstructure:"read_only,omitempty"`

	// The external volume reference.
	Source string `json:"source" yaml:"source" mapstructure:"source"`

	// The target mount on the container.
	Target string `json:"target" yaml:"target" mapstructure:"target"`
}

// An HTTP probe details.
type HttpProbe struct {
	// Host name to connect to. Defaults to the container IP.
	Host *string `json:"host,omitempty" yaml:"host,omitempty" mapstructure:"host,omitempty"`

	// Additional HTTP headers to send with the request
	HttpHeaders []HttpProbeHttpHeadersElem `json:"httpHeaders,omitempty" yaml:"httpHeaders,omitempty" mapstructure:"httpHeaders,omitempty"`

	// The path of the HTTP probe endpoint.
	Path string `json:"path" yaml:"path" mapstructure:"path"`

	// The path of the HTTP probe endpoint.
	Port *int `json:"port,omitempty" yaml:"port,omitempty" mapstructure:"port,omitempty"`

	// Scheme to use for connecting to the host (HTTP or HTTPS). Defaults to HTTP.
	Scheme *HttpProbeScheme `json:"scheme,omitempty" yaml:"scheme,omitempty" mapstructure:"scheme,omitempty"`
}

type HttpProbeHttpHeadersElem struct {
	// The HTTP header name.
	Name *string `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`

	// The HTTP header value.
	Value *string `json:"value,omitempty" yaml:"value,omitempty" mapstructure:"value,omitempty"`
}

type HttpProbeScheme string

const HttpProbeSchemeHTTP HttpProbeScheme = "HTTP"
const HttpProbeSchemeHTTPS HttpProbeScheme = "HTTPS"

// UnmarshalJSON implements json.Unmarshaler.
func (j *Resource) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["type"]; !ok || v == nil {
		return fmt.Errorf("field type in Resource: required")
	}
	type Plain Resource
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Resource(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *HttpProbe) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["path"]; !ok || v == nil {
		return fmt.Errorf("field path in HttpProbe: required")
	}
	type Plain HttpProbe
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.HttpHeaders != nil && len(plain.HttpHeaders) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "httpHeaders", 1)
	}
	*j = HttpProbe(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *HttpProbeScheme) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_HttpProbeScheme {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_HttpProbeScheme, v)
	}
	*j = HttpProbeScheme(v)
	return nil
}

var enumValues_HttpProbeScheme = []interface{}{
	"HTTP",
	"HTTPS",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ContainerVolumesElem) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["source"]; !ok || v == nil {
		return fmt.Errorf("field source in ContainerVolumesElem: required")
	}
	if v, ok := raw["target"]; !ok || v == nil {
		return fmt.Errorf("field target in ContainerVolumesElem: required")
	}
	type Plain ContainerVolumesElem
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ContainerVolumesElem(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ContainerFilesElem) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["target"]; !ok || v == nil {
		return fmt.Errorf("field target in ContainerFilesElem: required")
	}
	type Plain ContainerFilesElem
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.Source != nil && len(*plain.Source) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "source", 1)
	}
	*j = ContainerFilesElem(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Container) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["image"]; !ok || v == nil {
		return fmt.Errorf("field image in Container: required")
	}
	type Plain Container
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.Args != nil && len(plain.Args) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "args", 1)
	}
	if plain.Command != nil && len(plain.Command) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "command", 1)
	}
	if plain.Files != nil && len(plain.Files) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "files", 1)
	}
	if plain.Volumes != nil && len(plain.Volumes) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "volumes", 1)
	}
	*j = Container(plain)
	return nil
}

// The metadata for the resource.
type ResourceMetadata map[string]interface{}

// The parameters used to validate or provision the resource in the environment.
type ResourceParams map[string]interface{}

// The resource name.
type Resource struct {
	// A specialisation of the resource type.
	Class *string `json:"class,omitempty" yaml:"class,omitempty" mapstructure:"class,omitempty"`

	// The metadata for the resource.
	Metadata ResourceMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty" mapstructure:"metadata,omitempty"`

	// The parameters used to validate or provision the resource in the environment.
	Params ResourceParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`

	// The resource in the target environment.
	Type string `json:"type" yaml:"type" mapstructure:"type"`
}

// The compute resources limits.
type ResourcesLimits struct {
	// The CPU limit.
	Cpu *string `json:"cpu,omitempty" yaml:"cpu,omitempty" mapstructure:"cpu,omitempty"`

	// The memory limit.
	Memory *string `json:"memory,omitempty" yaml:"memory,omitempty" mapstructure:"memory,omitempty"`
}

// The network port description.
type ServicePort struct {
	// The public service port.
	Port int `json:"port" yaml:"port" mapstructure:"port"`

	// The transport level protocol. Defaults to TCP.
	Protocol *string `json:"protocol,omitempty" yaml:"protocol,omitempty" mapstructure:"protocol,omitempty"`

	// The internal service port. This will default to 'port' if not provided.
	TargetPort *int `json:"targetPort,omitempty" yaml:"targetPort,omitempty" mapstructure:"targetPort,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ServicePort) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["port"]; !ok || v == nil {
		return fmt.Errorf("field port in ServicePort: required")
	}
	type Plain ServicePort
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ServicePort(plain)
	return nil
}

// The declared Score Specification version.
type WorkloadContainers map[string]Container

// The metadata description of the Workload.
type WorkloadMetadata map[string]interface{}

// The dependencies needed by the Workload.
type WorkloadResources map[string]Resource

// List of network ports published by the service.
type WorkloadServicePorts map[string]ServicePort

// The service that the workload provides.
type WorkloadService struct {
	// List of network ports published by the service.
	Ports WorkloadServicePorts `json:"ports,omitempty" yaml:"ports,omitempty" mapstructure:"ports,omitempty"`
}

// Score workload specification
type Workload struct {
	// The declared Score Specification version.
	ApiVersion string `json:"apiVersion" yaml:"apiVersion" mapstructure:"apiVersion"`

	// The declared Score Specification version.
	Containers WorkloadContainers `json:"containers" yaml:"containers" mapstructure:"containers"`

	// The metadata description of the Workload.
	Metadata WorkloadMetadata `json:"metadata" yaml:"metadata" mapstructure:"metadata"`

	// The dependencies needed by the Workload.
	Resources WorkloadResources `json:"resources,omitempty" yaml:"resources,omitempty" mapstructure:"resources,omitempty"`

	// The service that the workload provides.
	Service *WorkloadService `json:"service,omitempty" yaml:"service,omitempty" mapstructure:"service,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Workload) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["apiVersion"]; !ok || v == nil {
		return fmt.Errorf("field apiVersion in Workload: required")
	}
	if v, ok := raw["containers"]; !ok || v == nil {
		return fmt.Errorf("field containers in Workload: required")
	}
	if v, ok := raw["metadata"]; !ok || v == nil {
		return fmt.Errorf("field metadata in Workload: required")
	}
	type Plain Workload
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Workload(plain)
	return nil
}
