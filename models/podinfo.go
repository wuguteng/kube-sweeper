package models

// PodInfo for kubernetes pod
type PodInfo struct {
	Namespace string            `json:"namespace" yaml:"namespace"`
	Name      string            `json:"name" yaml:"name"`
	Version   string            `json:"version" yaml:"version"`
	PosSuffix string            `json:"pod_suffix" yaml:"pod_suffix"`
	Labels    map[string]string `json:"labels" yaml:"labels"`
	PodName   string            `json:"pod_name" yaml:"pod_name"`
}
