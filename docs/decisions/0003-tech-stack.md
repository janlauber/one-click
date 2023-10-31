# Tech Stack

Date: 31.10.2023

## Status

Accepted

## Context

We need to decide on a tech stack for this project to create a web application prototype.

## Decision

We will use the following tech stack:

- Figma (Mockups): <https://www.figma.com/>
  - The reason for using Figma is that you can create and iterate on mockups very quickly. Also you can share the mockups with other people and get feedback.
- Sveltekit (Frontend): <https://kit.svelte.dev/>
  - We will use Sveltekit for the frontend, because it we got good experiences and boilerplate laying around. Also it is very easy to use and has a lot of features.
  - We will use TailwindCSS and Flowbite-Svelte for the styling of the application.
- Pocketbase / Golang (Backend): <https://pocketbase.io>
  - For the backend we will use Pocketbase and extend it with hooks written in Golang. Pocketbase can function as a standalone backend with a database and authentication.
- Github (Version Control): <https://github.com>
  - We will use Github for version control, because it is the most popular version control system and we got good experiences with it. It also has a lot of features like CI/CD, Issues, Pull Requests, etc. which we will use to automate our workflow.
- Docker (Containerization): <https://www.docker.com/>
  - We will use Docker for containerization to make it easy to deploy the application to different environments.
- Kubernetes (Orchestration): <https://kubernetes.io/>
  - We will use Kubernetes as the platform to run our application on. It is very popular and has a lot of features like auto-scaling, load balancing, etc. Also it is directly connected to our product, because we will extend Pocketbase to deploy Workloads to Kubernetes. (One-Click Deployment)

## Consequences

Because of choosing this tech stack we will have a prototype of the web application in a short amount of time.
