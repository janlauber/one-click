# One-Click OSS Deployment Platform

## Introduction

Welcome to the One-Click Deployment Platform. This project aims to provide a user-friendly interface for deploying dockerized containers on a Kubernetes cluster.  
It's built with the following technologies:

- **Frontend**: [Sveltekit](https://kit.svelte.dev/)
- **Backend**: [Go](https://golang.org/) & [Pocketbase](https://pocketbase.io/)
- **Operator**: [Operator SDK](https://sdk.operatorframework.io/) -> [one-click-operator](https://github.com/janlauber/one-click-operator)

![One-Click](./docs/assets/images/gif/one-click.gif)

## Features

- **Simplified Deployment**: Deploy your favorite dockerized container with a single click.
- **Kubernetes Integration**: Utilizes Kubernetes for efficient container orchestration.
- **User-Friendly Interface**: Intuitive UI for users without a technical background.
- **Customizable**: Customize your deployment with a variety of options.
- **Blueprints**: Create and share your own deployment blueprints.

## Getting Started

### Prerequisites

You will need the following to run this project:

- [Kubernetes](https://kubernetes.io/) cluster
- [Docker](https://www.docker.com/) daemon
- [Node.js](https://nodejs.org/en/) v18.16.0 or higher
- [npm](https://www.npmjs.com/) v9.5.1 or higher
- [Go](https://golang.org/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Kustomize](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/)

### Installation

1. Install the Operator
   Follow the installation instructions provided in the [one-click-operator repository](https://github.com/janlauber/one-click-operator).

2. Install the UI & Backend
    Check out the [deployment](./deployment/) folder and change the values for your environment. Then run the following commands:
    ```sh
    cd deployment
    kubectl apply -k .
    ```

3. Access the UI
    ```sh
    # if you are using an ingress
    kubectl get ingress -n one-click
    # if you want to use port-forwarding
    kubectl port-forward -n one-click svc/one-click-ui 8080:80
    ```

4. Access Pocketbase on your URL or localhost:8080 with `/_` as the path. Example: `localhost:8080/_`. You should see the Pocketbase UI and set your admin user. Then create a new user under `users` collection. You can now login with your new user.

## Usage

### Blueprints

Blueprints are an abstraction of a project. They contain some predefined values and can be used to deploy a certain container with a single click. You can create your own blueprints and share them with others.

### Projects

Projects are the actual deployments. They are based on blueprints. You can create a project from a blueprint and customize it to your needs.
Each configuration is stored in a rollout. A rollout is a version of a project configuration. So each time you change a configuration, a new rollout is created. You can then rollback to a previous rollout.

## Roadmap

See the [open issues](https://github.com/janlauber/one-click/issues) for a list of proposed features (and known issues).

## Contributing

Contributions are what make the open-source community such an amazing place to be, learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Submit your PR
