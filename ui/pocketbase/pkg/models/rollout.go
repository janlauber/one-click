package models

type PodMetrics struct {
	Name   string `json:"name"`
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type PodMetricsResponse struct {
	Metrics []PodMetrics `json:"metrics"`
}

type Rollout struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		Image struct {
			Registry   string `yaml:"registry"`
			Repository string `yaml:"repository"`
			Tag        string `yaml:"tag"`
		} `yaml:"image"`
		HorizontalScale struct {
			MinReplicas                    int `yaml:"minReplicas"`
			MaxReplicas                    int `yaml:"maxReplicas"`
			TargetCPUUtilizationPercentage int `yaml:"targetCPUUtilizationPercentage"`
		} `yaml:"horizontalScale"`
		Resources struct {
			Requests struct {
				CPU    string `yaml:"cpu"`
				Memory string `yaml:"memory"`
			} `yaml:"requests"`
			Limits struct {
				CPU    string `yaml:"cpu"`
				Memory string `yaml:"memory"`
			} `yaml:"limits"`
		} `yaml:"resources"`
		Env []struct {
			Name  string `yaml:"name"`
			Value string `yaml:"value"`
		} `yaml:"env"`
		Secrets []struct {
			Name  string `yaml:"name"`
			Value string `yaml:"value"`
		} `yaml:"secrets"`
		Volumes []struct {
			Name         string `yaml:"name"`
			MountPath    string `yaml:"mountPath"`
			Size         string `yaml:"size"`
			StorageClass string `yaml:"storageClass"`
		} `yaml:"volumes"`
		Interfaces []struct {
			Name    string `yaml:"name"`
			Port    int    `yaml:"port"`
			Ingress struct {
				IngressClass string            `yaml:"ingressClass"`
				Annotations  map[string]string `yaml:"annotations"`
				Rules        []struct {
					Host string `yaml:"host"`
					Path string `yaml:"path"`
					TLS  bool   `yaml:"tls"`
				} `yaml:"rules"`
			} `yaml:"ingress"`
		} `yaml:"interfaces"`
		ServiceAccountName string `yaml:"serviceAccountName"`
	} `yaml:"spec"`
}

type Resources struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type DeploymentResources struct {
	RequestSum Resources `json:"requestSum"`
	LimitSum   Resources `json:"limitSum"`
}

type DeploymentStatus struct {
	Replicas  int32               `json:"replicas"`
	PodNames  []string            `json:"podNames"`
	Resources DeploymentResources `json:"resources"`
	Status    string              `json:"status"`
}

type ServiceStatus struct {
	Name   string  `json:"name"`
	Ports  []int32 `json:"ports"`
	Status string  `json:"status"`
}

type IngressStatus struct {
	Name   string   `json:"name"`
	Hosts  []string `json:"hosts"`
	Status string   `json:"status"`
}

type VolumeStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// RolloutStatus defines the observed state of Rollout
type RolloutStatus struct {
	Deployment DeploymentStatus `json:"deployment"`
	Services   []ServiceStatus  `json:"services,omitempty"`
	Ingresses  []IngressStatus  `json:"ingresses,omitempty"`
	Volumes    []VolumeStatus   `json:"volumes,omitempty"`
}

type Event struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

type EventResponse struct {
	Events []Event `json:"events"`
}
