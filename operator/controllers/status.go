package controllers

import (
	"context"

	oneclickiov1 "github.com/janlauber/one-click/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *FrameworkReconciler) updateStatus(ctx context.Context, f *oneclickiov1.Framework) error {
	log := log.FromContext(ctx)

	// add some test data to the status
	f.Status.Deployment.Replicas = 1
	f.Status.Deployment.PodNames = []string{"test1", "test2"}
	f.Status.Deployment.Status = "OK"

	f.Status.Ingresses = []oneclickiov1.IngressStatus{
		{
			Hosts:  []string{"test1.example.com", "test2.example.com"},
			Status: "OK",
		},
	}

	f.Status.Services = []oneclickiov1.ServiceStatus{
		{
			Name:   "test1",
			Ports:  []int32{80, 443},
			Status: "OK",
		},
	}

	f.Status.Volumes = []oneclickiov1.VolumeStatus{
		{
			Name:   "test1",
			Status: "OK",
		},
	}

	// Update the Framework status
	if err := r.Status().Update(ctx, f); err != nil {
		log.Error(err, "Failed to update Framework status")
		return err
	}

	return nil
}
