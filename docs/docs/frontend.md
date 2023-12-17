# Frontend

## Mockups

### Introduction

In this section we will show you the mockups for the frontend of the application.
It is important to note that these mockups are not final and will change over time.

For the mockups we used Figma, which is a tool for creating mockups and prototypes.

!!! note "Figma"
    You can find the mockups here:  
    [Figma Design File](https://www.figma.com/file/Bi4OKCqGSgPXN1tvFzVV6Q/Untitled?type=design&mode=design&t=Sv1KYLH6X63M3oaV-1)

## Sveltekit

We used Sveltekit to create the frontend of the application. Sveltekit is a framework for creating web applications. It is based on Svelte, which is a compiler for creating web applications. Sveltekit is still in beta and will be released in the next few months.

!!! note "Sveltekit"
    You can find the Sveltekit documentation here:  
    [Sveltekit Documentation](https://kit.svelte.dev/docs)

### Pages

The frontend of the application consists of the following pages:

- Login
- Overview
- Projects
  - Overview
  - Environment Variables
  - Builds
  - Instances
  - Logs
  - Networking
  - Rollouts
  - Scale
  - Settings
  - Volumes
- Profile

### Authentication

The authentication is done via static Pocketbase account. This means that you can only login with the credentials that are stored in the [backend](backend.md). The authentication is done via the [Pocketbase JS-SDK](https://github.com/pocketbase/js-sdk).
