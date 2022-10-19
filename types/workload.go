package types

// WorkloadSpec is a single workload specification.
//
// YAML example:
//
//	apiVersion: score.sh/v1b1
//	metadata:
//	  name: hello-world
//	service:
//	  ports:
//	    www:
//	      hostPort: 80
//	      protocol: TCP
//	      port: 8080
//	containers:
//	  hello:
//	    image: busybox
//	    command: ["/bin/echo"]
//	    args: ["Hello $(FRIEND)"]
//	    variables:
//	      FRIEND: World!
//	    resources:
//	      limits:
//	        memory: "128Mi"
//	        cpu: "500m"
//	      requests:
//	        memory: "64Mi"
//	        cpu: "250m"
//	    livenessProbe:
//	      httpGet:
//	        path: /health
//	        port: 8080
//	        httpHeaders:
//	        - name: Custom-Header
//	          value: Awesome
//	    files:
//	      - target: etc/hello-world/config.yaml
//	        mode: "666"
//	        content: ${resources.env.APP_CONFIG}
//	    volumes:
//	      - source: ${resources.data}
//	        path: sub/path
//	        target: /mnt/data
//	        read_only: true
//	resources:
//	  env:
//	    type: environment
//	    properties:
//	      APP_CONFIG:
//	  dns:
//	    type: dns
//	    properties:
//	      domain:
//	  data:
//	    type: volume
//	  db:
//	    type: postgres
//	    properties:
//	      host:
//	        default: localhost
//	      port:
//	        default: 5432
type WorkloadSpec struct {
	ApiVersion string          `json:"apiVersion"`
	Metadata   WorkloadMeta    `json:"metadata"`
	Service    ServiceSpec     `json:"service"`
	Containers ContainersSpecs `json:"containers"`
	Resources  ResourcesSpecs  `json:"resources"`
}

// WorkloadMeta is a workload metadata.
type WorkloadMeta struct {
	Name string `json:"name"`
}
