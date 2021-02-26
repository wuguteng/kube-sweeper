package models

// KubeAPIStatus interface
type KubeAPIStatus interface {
	GetStatus() string
	GetReason() string
}

// KubeAPIResult for kubernetes api server
type KubeAPIResult struct {
	Kind       string        `json:"kind"`
	Code       int           `json:"code"`
	APIVersion string        `json:"apiVersion"`
	Status     KubePodStatus `json:"status"`
}

// KubePodStatus running status
type KubePodStatus struct {
	Phase     string `json:"phase"`
	HostIP    string `json:"hostIP"`
	PodIP     string `json:"podIP"`
	StartTime string `json:"startTime"`
	QosClass  string `json:"qosClass"`
}

// KubeFailureResult for kubernetes api server
type KubeFailureResult struct {
	Kind       string `json:"kind"`
	Code       int    `json:"code"`
	APIVersion string `json:"apiVersion"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Reason     string `json:"reason"`
}

// GetStatus string
func (r *KubeAPIResult) GetStatus() string {
	return r.Status.Phase
}

// GetReason string
func (r *KubeAPIResult) GetReason() string {
	return ""
}

// GetStatus string
func (r *KubeFailureResult) GetStatus() string {
	return r.Status
}

// GetReason string
func (r *KubeFailureResult) GetReason() string {
	return r.Reason
}
