# One-Click Operator

## Description

The purpose of this Kubernetes operator is to significantly streamline the deployment process of single-container applications on Kubernetes. Traditionally, deploying such an application requires a series of complex steps, including setting up a deployment, a service, an ingress, a horizontal pod autoscaler, a service account, and optionally, a certificate. Each of these steps involves intricate configurations and a deep understanding of Kubernetes' workings.

However, this Kubernetes operator simplifies the entire process by serving as an efficient abstraction layer over the Kubernetes API. It eliminates the need for manually creating each of the aforementioned resources. Instead, users only need to define a custom resource. Upon the creation of this custom resource, the Kubernetes operator automatically undertakes the task of generating all the necessary Kubernetes resources.

The primary advantage of using this operator lies in its simplification and automation of the deployment process. It allows users, particularly those who may not be deeply versed in Kubernetes intricacies, to deploy single-container applications with ease and reliability. This streamlined process not only saves time and reduces the potential for human error but also ensures a consistent deployment experience.

Furthermore, the [pocketbase-backend](../ui/pocketbase/) of the project will leverage this operator to deploy applications to Kubernetes. This integration signifies a shift towards a more efficient, less error-prone deployment methodology, emphasizing automation and ease of use. In essence, this Kubernetes operator is a transformative tool, designed to make Kubernetes more accessible and manageable, particularly for deployments involving single-container applications.

## Build

```bash
make build
```

## Run

```bash
make install run # also runs the operator locally
```

## Deploy

```bash
make deploy IMG=<image> # deploy to cluster
```

## uninstall

```bash
make uninstall # uninstall from cluster
```
