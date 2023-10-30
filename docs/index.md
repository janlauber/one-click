# One-Click Deployment

- [Abstract](#abstract)
- [Research](./research/index.md)
- [Decisions](./decisions/index.md)

## Abstract

*tbd*

## Introduction

### Current Situation

Over the past few decades, open-source software (OSS) has emerged as a transformative force in the tech ecosystem, empowering end-users with a gamut of applications that are both economically viable and customizable. However, the complexity involved in setting up, deploying, and managing these OSS tools can deter many potential users. While hyperscalers have addressed the issue of customers not wanting to own hardware, the issue of effectively operating OSS remains unaddressed.
Despite the many advantages of OSS, the burden of deployment and associated maintenance often falls upon internal development and operations teams. This situation not only increases the overhead for businesses but can also slow down adaptation and implementation of potentially transformative tools.

### Problem Statement

The primary challenge is the complexity in deploying and managing an open-source tool. Notably, tasks such as initial setup, configuration, updates, and troubleshooting generally require specialised skill-sets that may not always be readily available in-house.
Furthermore, while the 'as-a-service' model is proliferating, it is not universally applicable. For instance, certain applications may not be compatible due to their specific operational requirements, dependencies or architecture. Therefore, creating a 'one-size-fits-all' model for OSS deployment is not straightforward.
The present scenario underscores a clear need for a solution that can comprehensively simplify and democratise OSS deployment.

### Objectives

This investigation aims to define a methodology, backed by a technical solution, which simplifies the deployment and management of open-source software. More specifically, we propose the development of a One-Click Deployment system that abstracts complexities, standardises deployment procedures and thereby, delivers OSS 'as-a-service'.
By addressing these challenges, we aim to overcome the prevalent barriers in OSS adoption and make them more accessible and manageable for businesses and developers alike.

### Methodological Approach

Our approach focuses on leveraging existing technologies, such as Kubernetes for container orchestration, Helm for package management, and Golang for backend programming to develop a prototype of the One-Click Deployment system. The research phase will systematically analyse the architecture of various OSS to understand their unique attributes and requirements and then extrapolate key insights for the design phase. The end goal is to build a robust, scalable, and flexible deployment system that reduces the burden on DevOps teams and democratizes access to open-source solutions.
