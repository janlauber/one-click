package k8s

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func DeletePod(namespace string, podName string) error {
	return Clientset.CoreV1().Pods(namespace).Delete(Ctx, podName, metav1.DeleteOptions{})
}
