package domain

// ContainerInfo object to send the information of a single container.
type ContainerInfo struct {
	ContainerName string `json:"container_name"`
	PodName       string `json:"pod_name"`
	Namespace     string `json:"namespace"`
	MemoryRequest string `json:"mem_req"`
	MemoryLimit   string `json:"mem_limit"`
	CPURequest    string `json:"cpu_req"`
	CPULimit      string `json:"cpu_limit"`
}
