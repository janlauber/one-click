package k8s

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func DeleteJob(namespace string, jobName string) error {
	// Delete all pods controlled by the job
	err := Clientset.CoreV1().Pods(namespace).DeleteCollection(Ctx, metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: "job-name=" + jobName,
	})
	if err != nil {
		return err
	}

	propagadtionPolicy := metav1.DeletePropagationForeground

	return Clientset.BatchV1().Jobs(namespace).Delete(Ctx, jobName, metav1.DeleteOptions{
		PropagationPolicy: &propagadtionPolicy,
	})
}
