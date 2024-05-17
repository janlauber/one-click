package k8s

import (
	"context"
	"flag"
	"log"

	"github.com/janlauber/one-click/pkg/env"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

var (
	Clientset       *kubernetes.Clientset
	DynamicClient   dynamic.Interface
	DiscoveryClient *discovery.DiscoveryClient
	MetricsClient   metricsv.Interface
	Kubeconfig      *rest.Config
	Ctx             context.Context
)

func Init() {
	var err error
	if env.Config.Local {
		kubeconfig := flag.String("kubeconfig", env.Config.LocalKubeConfigFile, "(optional) absolute path to the kubeconfig file")
		flag.Parse()

		Kubeconfig, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Fatalf("Failed to build kubeconfig: %v", err)
		}
	} else {
		Kubeconfig, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("Failed to get in-cluster config: %v", err)
		}
	}

	DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(Kubeconfig)
	if err != nil {
		log.Fatalf("Failed to create discovery client: %v", err)
	}

	MetricsClient, err = metricsv.NewForConfig(Kubeconfig)
	if err != nil {
		log.Fatalf("Failed to create metrics client: %v", err)
	}

	Clientset, err = kubernetes.NewForConfig(Kubeconfig)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes clientset: %v", err)
	}

	Ctx = context.Background()

	DynamicClient, err = dynamic.NewForConfig(Kubeconfig)
	if err != nil {
		log.Fatalf("Failed to create dynamic client: %v", err)
	}
}

func GetClusterVersion() (string, error) {
	if DiscoveryClient != nil {
		clusterVersion, err := DiscoveryClient.ServerVersion()
		if err != nil {
			return "", err
		}
		return clusterVersion.GitVersion, nil
	}

	return "unknown", nil
}

func GetClusterApi() (string, error) {
	var clusterName string

	if Kubeconfig != nil {
		clusterName = Kubeconfig.Host
	} else {
		clusterName = "unknown"
	}

	return clusterName, nil
}
