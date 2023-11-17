# State of Research

## Overview of Current Research

At present, the realm of automating open-source software deployment has seen various initiatives. The focus has primarily been on developing individual tools that simplify the deployment of specific software packages. This research ranges from simplified deployment tools for specific applications to broad 'Platform as a Service' (PaaS) providers such as Heroku, Google App Engine, and AWS Elastic Beanstalk that offer abstracted infrastructure and generalized deployment procedures for a range of applications.

## Underpinning Theories and Concepts

Our work draws inspiration from these existing solutions but distinguishes its approach by aiming for a universally compatible one-click solution specifically for open-source software. Central to our investigation is the concept of 'Infrastructure as Code' (IaC), emphasizing automation, and maintaining infrastructural consistency. Further, we borrow from the principles of Containerization and Orchestration, which form the underpinnings of modern deployment procedures largely due to their benefits of scalability, isolation, and portability.

## Terminology and Standards

Key terms that need to be defined for understanding this research include:

- **One-Click Deployment**: Refers to the automated procedure where software can be deployed easily with a single operation.
- **Containerization**: A lightweight form of virtualization that packages an application and its dependencies in a 'container', enabling it to run consistently in any environment.
- **Infrastructure as Code (IaC)**: The process of managing, provisioning, and configuring computing infrastructure through machine-readable script files rather than manual processes.

As for standards, our research abides by the principles defined by the 12-factor app methodology, a widely accepted standard for building software-as-a-service applications, with an emphasis on automation and scalability. Other standards involved in this investigation will be Kubernetes for orchestration, Helm for package management, and Docker for containerization.

## Application of Current Research

The intent of our work is not just to echo the previous research but to advance it by addressing the noted gaps and specific nuances of deploying open-source software. Our research will synthesize the lessons from the state of research, mould these insights into a cohesive system, and ultimately present a novel solution to simplify OSS deployment.
