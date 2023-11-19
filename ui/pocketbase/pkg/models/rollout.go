package models

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
