package k8s

import (
	"strings"

	"github.com/natrontech/one-click/pkg/util"
	"github.com/pocketbase/pocketbase/models"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceParams struct {
	Name       string
	UserRecord *models.Record
}

func CreateNamespace(params NamespaceParams) error {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: params.Name,
			Labels: map[string]string{
				"project.one-click.io/name":        params.Name,
				"project.one-click.io/userId":      params.UserRecord.GetString("id"),
				"project.one-click.io/username":    params.UserRecord.GetString("username"),
				"project.one-click.io/displayName": util.StringParser(params.UserRecord.GetString("name")),
			},
		},
	}
	_, err := Clientset.CoreV1().Namespaces().Create(Ctx, ns, metav1.CreateOptions{})

	// if err already exists, update
	if err != nil && strings.Contains(err.Error(), "already exists") {
		_, err = Clientset.CoreV1().Namespaces().Update(Ctx, ns, metav1.UpdateOptions{})
	}

	return err
}

func DeleteNamespace(namespace string) error {
	return Clientset.CoreV1().Namespaces().Delete(Ctx, namespace, metav1.DeleteOptions{})
}
