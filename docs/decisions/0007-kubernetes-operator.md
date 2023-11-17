# Kubernetes Operator

Date: 17.11.2023

## Status

Accepted

## Context

To manage the applications and frameworks we need to decide which tool we will use. The tool should be able to manage the lifecycle of the applications and frameworks. It should be able to create, update and delete the applications and frameworks. It should also be able to scale the applications and frameworks.

## Decision

We will use a Kubernetes Operator to manage the applications and frameworks. The Kubernetes Operator will be written in Golang. The Kubernetes Operator will be able to manage the lifecycle of the applications and frameworks. It will be able to create, update and delete the applications and frameworks. It will also be able to scale the applications and frameworks.

We will use the [Operator SDK](https://sdk.operatorframework.io/) to create the Kubernetes Operator. The Operator SDK will help us to create the Kubernetes Operator. It will also help us to test the Kubernetes Operator.

## Consequences

Because of using a Kubernetes Operator we will have more time to focus on the research and the prototype.
