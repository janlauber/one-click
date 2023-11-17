# Applications Exclusion

Date: 17.11.2023

## Status

Accepted

## Context

We need to decide which frameworks and applications will be supported by the prototype.
Frameworks are lightweight tools like Reflex, Streamlit, Gradio, Anvil, Plotly, etc. that can be used to build applications. The user needs to build a docker image with the framework.
The applications are ready-to-use tools like Node-RED, JupyterHub, etc. that can be used out-of-the-box. The user needs to provide a helm chart for the application.

## Decision

We will exclude the support for applications from the prototype, because it is too much work to implement it in the given time frame.

## Consequences

Because of excluding the support for applications we will have more time to focus on the research and the prototype.
