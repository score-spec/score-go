// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package types

import "encoding/json"
import "fmt"
import "reflect"

// The specification of a Container within the Workload.
type Container struct {
	// If specified, overrides the arguments passed to the container entrypoint.
	Args []string `json:"args,omitempty" yaml:"args,omitempty" mapstructure:"args,omitempty"`

	// If specified, overrides the entrypoint defined in the container image.
	Command []string `json:"command,omitempty" yaml:"command,omitempty" mapstructure:"command,omitempty"`

	// Files corresponds to the JSON schema field "files".
	Files ContainerFiles `json:"files,omitempty" yaml:"files,omitempty" mapstructure:"files,omitempty"`

	// The container image name and tag.
	Image string `json:"image" yaml:"image" mapstructure:"image"`

	// The liveness probe for the container.
	LivenessProbe *ContainerProbe `json:"livenessProbe,omitempty" yaml:"livenessProbe,omitempty" mapstructure:"livenessProbe,omitempty"`

	// The readiness probe for the container.
	ReadinessProbe *ContainerProbe `json:"readinessProbe,omitempty" yaml:"readinessProbe,omitempty" mapstructure:"readinessProbe,omitempty"`

	// The compute resources for the container.
	Resources *ContainerResources `json:"resources,omitempty" yaml:"resources,omitempty" mapstructure:"resources,omitempty"`

	// The environment variables for the container.
	Variables ContainerVariables `json:"variables,omitempty" yaml:"variables,omitempty" mapstructure:"variables,omitempty"`

	// Volumes corresponds to the JSON schema field "volumes".
	Volumes ContainerVolumes `json:"volumes,omitempty" yaml:"volumes,omitempty" mapstructure:"volumes,omitempty"`
}

// The details of a file to mount in the container. One of 'source', 'content', or
// 'binaryContent' must be provided.
type ContainerFile struct {
	// Inline standard-base64 encoded content for the file. Does not support
	// placeholder expansion.
	BinaryContent *string `json:"binaryContent,omitempty" yaml:"binaryContent,omitempty" mapstructure:"binaryContent,omitempty"`

	// The inline content for the file. Only supports valid utf-8.
	Content *string `json:"content,omitempty" yaml:"content,omitempty" mapstructure:"content,omitempty"`

	// The optional file access mode in octal encoding. For example 0600.
	Mode *string `json:"mode,omitempty" yaml:"mode,omitempty" mapstructure:"mode,omitempty"`

	// If set to true, the placeholders expansion will not occur in the contents of
	// the file.
	NoExpand *bool `json:"noExpand,omitempty" yaml:"noExpand,omitempty" mapstructure:"noExpand,omitempty"`

	// The relative or absolute path to the content file.
	Source *string `json:"source,omitempty" yaml:"source,omitempty" mapstructure:"source,omitempty"`
}

type ContainerFiles map[string]ContainerFile

// The probe may be defined as either http, command execution, or both. The
// execProbe should be preferred if the Score implementation supports both types.
type ContainerProbe struct {
	// Exec corresponds to the JSON schema field "exec".
	Exec *ExecProbe `json:"exec,omitempty" yaml:"exec,omitempty" mapstructure:"exec,omitempty"`

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

type ContainerVolume struct {
	// An optional sub path in the volume.
	Path *string `json:"path,omitempty" yaml:"path,omitempty" mapstructure:"path,omitempty"`

	// Indicates if the volume should be mounted in a read-only mode.
	ReadOnly *bool `json:"readOnly,omitempty" yaml:"readOnly,omitempty" mapstructure:"readOnly,omitempty"`

	// The external volume reference.
	Source string `json:"source" yaml:"source" mapstructure:"source"`
}

type ContainerVolumes map[string]ContainerVolume

// An executable health probe.
type ExecProbe struct {
	// The command and arguments to execute within the container.
	Command []string `json:"command" yaml:"command" mapstructure:"command"`
}

// An HTTP probe details.
type HttpProbe struct {
	// Host name to connect to. Defaults to the workload IP. The is equivalent to a
	// Host HTTP header.
	Host *string `json:"host,omitempty" yaml:"host,omitempty" mapstructure:"host,omitempty"`

	// Additional HTTP headers to send with the request
	HttpHeaders []HttpProbeHttpHeadersElem `json:"httpHeaders,omitempty" yaml:"httpHeaders,omitempty" mapstructure:"httpHeaders,omitempty"`

	// The path to access on the HTTP server.
	Path string `json:"path" yaml:"path" mapstructure:"path"`

	// The port to access on the workload.
	Port int `json:"port" yaml:"port" mapstructure:"port"`

	// Scheme to use for connecting to the host (HTTP or HTTPS). Defaults to HTTP.
	Scheme *HttpProbeScheme `json:"scheme,omitempty" yaml:"scheme,omitempty" mapstructure:"scheme,omitempty"`
}

type HttpProbeHttpHeadersElem struct {
	// The HTTP header name.
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// The HTTP header value.
	Value string `json:"value" yaml:"value" mapstructure:"value"`
}

type HttpProbeScheme string

const HttpProbeSchemeHTTP HttpProbeScheme = "HTTP"
const HttpProbeSchemeHTTPS HttpProbeScheme = "HTTPS"

// The set of Resources associated with this Workload.
type Resource struct {
	// An optional specialisation of the Resource type.
	Class *string `json:"class,omitempty" yaml:"class,omitempty" mapstructure:"class,omitempty"`

	// An optional Resource identifier. The id may be up to 63 characters, including
	// one or more labels of a-z, 0-9, '-' not starting or ending with '-' separated
	// by '.'. When two resources share the same type, class, and id, they are
	// considered the same resource when used across related Workloads.
	Id *string `json:"id,omitempty" yaml:"id,omitempty" mapstructure:"id,omitempty"`

	// The metadata for the Resource.
	Metadata ResourceMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty" mapstructure:"metadata,omitempty"`

	// Optional parameters used to provision the Resource in the environment.
	Params ResourceParams `json:"params,omitempty" yaml:"params,omitempty" mapstructure:"params,omitempty"`

	// The Resource type. This should be a type supported by the Score implementations
	// being used.
	Type string `json:"type" yaml:"type" mapstructure:"type"`
}

// The metadata for the Resource.
type ResourceMetadata map[string]interface{}

// Optional parameters used to provision the Resource in the environment.
type ResourceParams map[string]interface{}

// The compute and memory resource limits.
type ResourcesLimits struct {
	// The CPU limit as whole or fractional CPUs. 'm' indicates milli-CPUs. For
	// example 2 or 125m.
	Cpu *string `json:"cpu,omitempty" yaml:"cpu,omitempty" mapstructure:"cpu,omitempty"`

	// The memory limit in bytes with optional unit specifier. For example 125M or
	// 1Gi.
	Memory *string `json:"memory,omitempty" yaml:"memory,omitempty" mapstructure:"memory,omitempty"`
}

// The network port description.
type ServicePort struct {
	// The public service port.
	Port int `json:"port" yaml:"port" mapstructure:"port"`

	// The transport level protocol. Defaults to TCP.
	Protocol *ServicePortProtocol `json:"protocol,omitempty" yaml:"protocol,omitempty" mapstructure:"protocol,omitempty"`

	// The internal service port. This will default to 'port' if not provided.
	TargetPort *int `json:"targetPort,omitempty" yaml:"targetPort,omitempty" mapstructure:"targetPort,omitempty"`
}

type ServicePortProtocol string

const ServicePortProtocolTCP ServicePortProtocol = "TCP"
const ServicePortProtocolUDP ServicePortProtocol = "UDP"

// Score workload specification
type Workload struct {
	// The declared Score Specification version.
	ApiVersion string `json:"apiVersion" yaml:"apiVersion" mapstructure:"apiVersion"`

	// The set of named containers in the Workload. The container name must be a valid
	// RFC1123 Label Name of up to 63 characters, including a-z, 0-9, '-' but may not
	// start or end with '-'.
	Containers WorkloadContainers `json:"containers" yaml:"containers" mapstructure:"containers"`

	// The metadata description of the Workload.
	Metadata WorkloadMetadata `json:"metadata" yaml:"metadata" mapstructure:"metadata"`

	// The Resource dependencies needed by the Workload. The resource name must be a
	// valid RFC1123 Label Name of up to 63 characters, including a-z, 0-9, '-' but
	// may not start or end with '-'.
	Resources WorkloadResources `json:"resources,omitempty" yaml:"resources,omitempty" mapstructure:"resources,omitempty"`

	// The service that the workload provides.
	Service *WorkloadService `json:"service,omitempty" yaml:"service,omitempty" mapstructure:"service,omitempty"`
}

// The set of named containers in the Workload. The container name must be a valid
// RFC1123 Label Name of up to 63 characters, including a-z, 0-9, '-' but may not
// start or end with '-'.
type WorkloadContainers map[string]Container

// The metadata description of the Workload.
type WorkloadMetadata map[string]interface{}

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
	if plain.Class != nil && len(*plain.Class) < 2 {
		return fmt.Errorf("field %s length: must be >= %d", "class", 2)
	}
	if plain.Class != nil && len(*plain.Class) > 63 {
		return fmt.Errorf("field %s length: must be <= %d", "class", 63)
	}
	if plain.Id != nil && len(*plain.Id) < 2 {
		return fmt.Errorf("field %s length: must be >= %d", "id", 2)
	}
	if plain.Id != nil && len(*plain.Id) > 63 {
		return fmt.Errorf("field %s length: must be <= %d", "id", 63)
	}
	if len(plain.Type) < 2 {
		return fmt.Errorf("field %s length: must be >= %d", "type", 2)
	}
	if len(plain.Type) > 63 {
		return fmt.Errorf("field %s length: must be <= %d", "type", 63)
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
	if v, ok := raw["port"]; !ok || v == nil {
		return fmt.Errorf("field port in HttpProbe: required")
	}
	type Plain HttpProbe
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.Host != nil && len(*plain.Host) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "host", 1)
	}
	*j = HttpProbe(plain)
	return nil
}

var enumValues_ServicePortProtocol = []interface{}{
	"TCP",
	"UDP",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ServicePortProtocol) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_ServicePortProtocol {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_ServicePortProtocol, v)
	}
	*j = ServicePortProtocol(v)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *HttpProbeHttpHeadersElem) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name in HttpProbeHttpHeadersElem: required")
	}
	if v, ok := raw["value"]; !ok || v == nil {
		return fmt.Errorf("field value in HttpProbeHttpHeadersElem: required")
	}
	type Plain HttpProbeHttpHeadersElem
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if len(plain.Value) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "value", 1)
	}
	*j = HttpProbeHttpHeadersElem(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ContainerVolume) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["source"]; !ok || v == nil {
		return fmt.Errorf("field source in ContainerVolume: required")
	}
	type Plain ContainerVolume
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ContainerVolume(plain)
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
	if len(plain.Image) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "image", 1)
	}
	*j = Container(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ExecProbe) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["command"]; !ok || v == nil {
		return fmt.Errorf("field command in ExecProbe: required")
	}
	type Plain ExecProbe
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ExecProbe(plain)
	return nil
}

// The Resource dependencies needed by the Workload. The resource name must be a
// valid RFC1123 Label Name of up to 63 characters, including a-z, 0-9, '-' but may
// not start or end with '-'.
type WorkloadResources map[string]Resource

// The set of named network ports published by the service. The service name must
// be a valid RFC1123 Label Name of up to 63 characters, including a-z, 0-9, '-'
// but may not start or end with '-'.
type WorkloadServicePorts map[string]ServicePort

// The service that the workload provides.
type WorkloadService struct {
	// The set of named network ports published by the service. The service name must
	// be a valid RFC1123 Label Name of up to 63 characters, including a-z, 0-9, '-'
	// but may not start or end with '-'.
	Ports WorkloadServicePorts `json:"ports,omitempty" yaml:"ports,omitempty" mapstructure:"ports,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ContainerFile) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	type Plain ContainerFile
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if plain.Source != nil && len(*plain.Source) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "source", 1)
	}
	*j = ContainerFile(plain)
	return nil
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
