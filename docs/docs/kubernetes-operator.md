# Kubernetes Operator

The Kubernetes operator is a Kubernetes controller that manages the lifecycle of the project's custom resources. It is written in Go and uses the [Operator SDK](https://sdk.operatorframework.io/docs/building-operators/golang/) to generate the boilerplate code.

## Custom Resources

The Kubernetes operator manages the following custom resource:

```yaml
apiVersion: one-click.io/v1alpha1
kind: Rollout
metadata:
  name: <Rollout Name>
  namespace: <Rollout Namespace>
spec:
  image:
    registry: <Image Registry>
    repository: <Image Repository>
    tag: <Image Tag>
  horizontalScale:
    minReplicas: <Minimum Number of Replicas> # 1 for minimum
    maxReplicas: <Maximum Number of Replicas> # n for maximum
    targetCPUUtilizationPercentage: <Target CPU Utilization Percentage> # 0-100
  resources:
    requests:
      cpu: <CPU Request> # 0.1 for 100m
      memory: <Memory Request> # 256Mi for 256Mi
    limits:
      cpu: <CPU Limit> # 0.1 for 100m
      memory: <Memory Limit> # 256Mi for 256Mi
  env:
    - name: <Environment Variable Name> # e.g. DB_HOST
      value: <Environment Variable Value> # e.g. db
  secrets:
    - name: <Secret Name> # e.g. DB_PASSWORD
      value: <Secret Value> # e.g. password
  volumes:
    - name: <Volume Name> # e.g. data
      mountPath: <Volume Mount Path> # e.g. /data
      size: <Volume Size> # e.g. 1Gi
      storageClass: <Volume Storage Class> # e.g. standard
  interfaces:
    - name: <Interface Name> # e.g. http
      port: <Interface Port> # e.g. 80
      ingress: # optional ingress configuration
        ingressClass: <Ingress Class> # e.g. nginx
        annotations: # optional ingress annotations
          nginx.ingress.kubernetes.io/rewrite-target: / # optional
          nginx.ingress.kubernetes.io/ssl-redirect: "false" # optional
        rules:
          - host: <Ingress Host> # e.g. example.com
            path: <Ingress Path> # e.g. /
            tls: <Enable TLS> # e.g. true -> cert-manager will create a certificate (must be installed in the cluster)
  serviceAccountName: <Service Account Name> # e.g. default
```

## Purpose

This Kubernetes operator is an abstraction layer on top of the Kubernetes API. When you want to easily deploy a single container application to Kubernetes you need to create a deployment, a service, an ingress, a horizontal pod autoscaler, a service account, and optionally a certificate. This Kubernetes operator does all of that for you. You just need to create a custom resource and the Kubernetes operator will create all of the Kubernetes resources for you. So the benefit of this Kubernetes operator is to simplify the deployment of single container applications to Kubernetes and the [Backend](backend.md) of the project will use this simplified deployment process to deploy applications to Kubernetes.

## Development

The Kubernetes operator code is located in the `operator` directory. The following sections describe how to build and run the Kubernetes operator.

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Go](https://golang.org/doc/install)
- [kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/)
- [minikube](https://minikube.sigs.k8s.io/docs/start/)
- [operator-sdk](https://sdk.operatorframework.io/docs/installation/install-operator-sdk/)
- [yq](https://mikefarah.gitbook.io/yq/)

### Build

```bash
make build
```

### Run

```bash
make install run # also runs the operator locally
```

### Deploy

```bash
make deploy
```

### Uninstall

```bash
make uninstall
```

### Clean

```bash
make uninstall
make undeploy
```

## References

- [Kubernetes Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [Operator SDK](https://sdk.operatorframework.io/docs/building-operators/golang/)
