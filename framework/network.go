// Copyright 2026 The Score Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package framework

import score "github.com/score-spec/score-go/types"

// ServicePort is a resolved port from the workload specification. It differs from types.ServicePort by having all
// fields fully resolved (no optional pointers) and a Name field for the port's key in the service ports map.
type ServicePort struct {
	// Name is the name of the port from the workload specification.
	Name string `yaml:"name"`
	// Port is the numeric port intended to be published.
	Port int `yaml:"port"`
	// TargetPort is the port on the workload that hosts the actual traffic.
	TargetPort int `yaml:"target_port"`
	// Protocol is TCP or UDP.
	Protocol score.ServicePortProtocol `yaml:"protocol"`
}

// NetworkService describes how to contact ports exposed by another workload. Provisioners store this in SharedState
// so that downstream workloads can discover service endpoints without target-specific knowledge.
type NetworkService struct {
	// ServiceName is the DNS-resolvable or container-resolvable hostname for the workload.
	ServiceName string `yaml:"service_name"`
	// Ports maps a port key (name or numeric string) to its resolved ServicePort.
	Ports map[string]ServicePort `yaml:"ports"`
}
