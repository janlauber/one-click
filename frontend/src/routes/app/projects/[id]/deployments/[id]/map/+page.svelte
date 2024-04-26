<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { writable } from "svelte/store";
  import { SvelteFlow, type Node, type Edge } from "@xyflow/svelte";
  import "@xyflow/svelte/dist/style.css";
  import "./turbo.css";
  import { initialNodes, initialEdges } from "./nodes-and-edges";
  import TurboNode from "./TurboNode.svelte";
  import TurboEdge from "./TurboEdge.svelte";
  import { Button, Heading, P, TabItem, Tabs } from "flowbite-svelte";
  import { selectedProject } from "$lib/stores/data";
  import { Drawer } from "flowbite-svelte";
  import { sineIn } from "svelte/easing";
  import { drawerHidden, selectedNode, type NodeObject } from "$lib/stores/drawer";
  import {
    ArrowLeftRight,
    Bell,
    Box,
    Database,
    FileCode,
    Lock,
    NetworkIcon,
    ScrollText,
    Trash,
    X
  } from "lucide-svelte";
  import LogStream from "$lib/components/map/LogStream.svelte";
  import KubernetesObject from "$lib/components/map/KubernetesObject.svelte";
  import EventStream from "$lib/components/map/EventStream.svelte";
  import toast from "svelte-french-toast";
  import selectedDeploymentId from "$lib/stores/deployment";

  let transitionParamsRight = {
    x: 320,
    duration: 200,
    easing: sineIn
  };
  let initialLoadComplete = false;

  const nodes = writable<Node[]>(initialNodes);
  const edges = writable<Edge[]>(initialEdges);

  const nodeTypes = { turbo: TurboNode };
  const edgeTypes = { turbo: TurboEdge };
  const defaultEdgeOptions = { type: "turbo", markerEnd: "edge-circle" };

  let ws: WebSocket;

  onMount(() => {
    // host
    let host = window.location.host;

    if (host.includes("localhost")) {
      host = "localhost:8090";
    }

    // check for tls
    let protocol = window.location.protocol === "https:" ? "wss" : "ws";

    ws = new WebSocket(`${protocol}://${host}/ws/k8s/deployments`);

    ws.onopen = () => {
      // Send any initial message if required, for example:
      // {
      // deploymentId: "deploymentId"
      // projectId: "projectId"
      // }
      ws.send(
        JSON.stringify({
          deploymentId: $selectedDeploymentId ?? "",
          projectId: $selectedProject?.id
        })
      );
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      updateObjects(data);
    };

    // set initialLoadComplete to true after 0.3s
    setTimeout(() => {
      initialLoadComplete = true;
    }, 300);
  });

  onDestroy(() => {
    // Close the WebSocket connection
    ws.close();

    // Reset the nodes and edges
    nodes.set([]);
    edges.set([]);
  });

  const pods = writable<NodeObject[]>([]);
  const services = writable<NodeObject[]>([]);
  const ingresses = writable<NodeObject[]>([]);
  const secrets = writable<NodeObject[]>([]);
  const pvc = writable<NodeObject[]>([]);

  function updateObjects(data: NodeObject) {
    if (data.status === "ADDED") {
      if (data.kind === "pod") {
        pods.update((p) => [...p, data]);
      } else if (data.kind === "service") {
        services.update((s) => [...s, data]);
      } else if (data.kind === "ingress") {
        ingresses.update((i) => [...i, data]);
      } else if (data.kind === "secret") {
        secrets.update((s) => [...s, data]);
      } else if (data.kind === "pvc") {
        pvc.update((p) => [...p, data]);
      }
    } else if (data.status === "MODIFIED") {
      updateNode(data);
      if (data.kind === "pod") {
        // update with the uid of object metadata
        pods.update((p) => {
          const index = p.findIndex((pod) => pod.object.metadata.uid === data.object.metadata.uid);
          p[index] = data;
          return p;
        });
      } else if (data.kind === "service") {
        services.update((s) => {
          const index = s.findIndex(
            (service) => service.object.metadata.uid === data.object.metadata.uid
          );
          s[index] = data;
          return s;
        });
      } else if (data.kind === "ingress") {
        ingresses.update((i) => {
          const index = i.findIndex(
            (ingress) => ingress.object.metadata.uid === data.object.metadata.uid
          );
          i[index] = data;
          return i;
        });
      } else if (data.kind === "secret") {
        secrets.update((s) => {
          const index = s.findIndex(
            (secret) => secret.object.metadata.uid === data.object.metadata.uid
          );
          s[index] = data;
          return s;
        });
      } else if (data.kind === "pvc") {
        pvc.update((p) => {
          const index = p.findIndex((pvc) => pvc.object.metadata.uid === data.object.metadata.uid);
          p[index] = data;
          return p;
        });
      }
    } else if (data.status === "DELETED") {
      removeNodeAndEdges(data.name);
      if (data.kind === "pod") {
        pods.update((p) => p.filter((pod) => pod.name !== data.name));
      } else if (data.kind === "service") {
        services.update((s) => s.filter((service) => service.name !== data.name));
      } else if (data.kind === "ingress") {
        ingresses.update((i) => i.filter((ingress) => ingress.name !== data.name));
      } else if (data.kind === "secret") {
        secrets.update((s) => s.filter((secret) => secret.name !== data.name));
      } else if (data.kind === "pvc") {
        pvc.update((p) => p.filter((pvc) => pvc.name !== data.name));
      }
    } else if (data.status === "ERROR") {
      console.error(data);
    }
  }

  function createNode(data: NodeObject): Node {
    return {
      id: data.name, // Ensure this is unique
      type: "turbo", // Using the custom node type
      data: {
        // kind: "pod" | "service" | "ingress" | "secret" | "pvc";
        // name: string;
        // namespace: string;
        // labels: Map<string, string>;
        // status: "ADDED" | "MODIFIED" | "DELETED" | "ERROR";
        // object: any;
        kind: data.kind,
        name: data.name,
        namespace: data.namespace,
        labels: data.labels,
        status: data.kind === "pod" ? data.object.status.phase : "Running",
        object: data.object,
        icon: data.kind.toLowerCase(), // Adjust icon based on kind,
        containerStatuses: data.kind === "pod" ? data.object.status.containerStatuses : []
      },
      position: calculatePosition(data)
    };
  }

  function calculatePosition(data: NodeObject): { x: number; y: number } {
    // Calculate the position of the nodes
    // left to right: ingress -> service -> pod -> secret
    const BASE_X_POSITIONS = {
      ingress: 0,
      service: 250,
      pod: 600,
      secret: 1050,
      pvc: 1050
    };

    const Y_OFFSET = 100;

    if (data.kind === "pod") {
      let tempPosition = { x: BASE_X_POSITIONS.pod, y: 0 };
      // check if there is a node with the same position and kind in $nodes
      nodes.subscribe((n) => {
        n.forEach((node) => {
          if (node.position.x === tempPosition.x && node.position.y === tempPosition.y) {
            if (node.data.icon === "pod") {
              tempPosition = { x: BASE_X_POSITIONS.pod, y: tempPosition.y + Y_OFFSET };
            }
          }
        });
      });
      return tempPosition;
    } else if (data.kind === "service") {
      let tempPosition = { x: BASE_X_POSITIONS.service, y: 0 };
      // check if there is a node with the same position and kind in $nodes
      nodes.subscribe((n) => {
        n.forEach((node) => {
          if (node.position.x === tempPosition.x && node.position.y === tempPosition.y) {
            if (node.data.icon === "service") {
              tempPosition = { x: BASE_X_POSITIONS.service, y: tempPosition.y + Y_OFFSET };
            }
          }
        });
      });
      return tempPosition;
    } else if (data.kind === "ingress") {
      let tempPosition = { x: BASE_X_POSITIONS.ingress, y: 0 };
      // check if there is a node with the same position and kind in $nodes
      nodes.subscribe((n) => {
        n.forEach((node) => {
          if (node.position.x === tempPosition.x && node.position.y === tempPosition.y) {
            if (node.data.icon === "ingress") {
              tempPosition = { x: BASE_X_POSITIONS.ingress, y: tempPosition.y + Y_OFFSET };
            }
          }
        });
      });
      return tempPosition;
    } else if (data.kind === "secret") {
      let tempPosition = { x: BASE_X_POSITIONS.secret, y: 0 };
      // check if there is a node with the same position and kind in $nodes
      nodes.subscribe((n) => {
        n.forEach((node) => {
          if (node.position.x === tempPosition.x && node.position.y === tempPosition.y) {
            if (node.data.icon === "secret" || node.data.icon === "pvc") {
              tempPosition = { x: BASE_X_POSITIONS.secret, y: tempPosition.y + Y_OFFSET };
            }
          }
        });
      });
      return tempPosition;
    } else if (data.kind === "pvc") {
      let tempPosition = { x: BASE_X_POSITIONS.secret, y: 0 };
      // check if there is a node with the same position and kind in $nodes
      nodes.subscribe((n) => {
        n.forEach((node) => {
          if (node.position.x === tempPosition.x && node.position.y === tempPosition.y) {
            if (node.data.icon === "pvc" || node.data.icon === "secret") {
              tempPosition = { x: BASE_X_POSITIONS.secret, y: tempPosition.y + Y_OFFSET };
            }
          }
        });
      });
      return tempPosition;
    }

    return { x: 0, y: 0 };
  }

  function addNode(data: NodeObject) {
    const newNode = createNode(data);
    nodes.update((n) => {
      if (!n.find((node) => node.id === newNode.id)) {
        return [...n, newNode];
      }
      return n;
    });
  }

  function createEdge(source: string, target: string): Edge {
    return {
      id: `${source}-${target}`, // Ensure this is unique
      source,
      target,
      type: "turbo", // Using the custom edge type
      markerEnd: "edge-arrow",
      animated: source.includes("ingress") || source.includes("svc") ? true : false
    };
  }

  function addEdge(source: string, target: string) {
    const newEdge = createEdge(source, target);
    edges.update((e) => {
      if (!e.find((edge) => edge.id === newEdge.id)) {
        return [...e, newEdge];
      }
      return e;
    });
  }

  function removeNodeAndEdges(objectName: string) {
    // Remove the node
    nodes.update((n) => n.filter((node) => node.id !== objectName));

    // Optionally, remove edges connected to the node
    edges.update((e) =>
      e.filter((edge) => edge.source !== objectName && edge.target !== objectName)
    );
  }

  function updateNode(data: NodeObject) {
    nodes.update((n) => {
      const nodeIndex = n.findIndex((node) => node.id === data.name);
      if (nodeIndex !== -1) {
        // Assuming the structure of your data object matches the node data structure
        const updatedNode = {
          ...n[nodeIndex],
          data: { ...n[nodeIndex].data, ...createNode(data).data }
        };
        return [...n.slice(0, nodeIndex), updatedNode, ...n.slice(nodeIndex + 1)];
      }
      return n; // Return the original nodes array if the node wasn't found
    });
  }

  async function handleDeletePod(podName: string) {
    const token = localStorage.getItem("pocketbase_auth");
    if (!token) {
      return;
    }
    const authHeader = { Authorization: `Bearer ${JSON.parse(token).token}` };

    // if localhost, use localhost:8090 as base url
    // POST /pb/:projectId/:podName
    if (window.location.hostname === "localhost") {
      try {
        const response = await fetch(
          `http://localhost:8090/pb/${$selectedProject?.id}/${podName}`,
          {
            method: "DELETE",
            headers: {
              ...authHeader,
              "Content-Type": "application/json"
            }
          }
        );
        if (response.status === 200) {
          toast.success("Pod terminating...");
          $drawerHidden = true;
        } else {
          console.error("Error deleting pod", response);
          toast.error("Error deleting pod");
        }
      } catch (error) {
        console.error("Error deleting pod", error);
        toast.error("Error deleting pod");
      }
    } else {
      try {
        const response = await fetch(`/pb/${$selectedProject?.id}/${podName}`, {
          method: "DELETE",
          headers: {
            ...authHeader,
            "Content-Type": "application/json"
          }
        });
        if (response.status === 200) {
          toast.success("Pod terminating...");
          $drawerHidden = true;
        } else {
          console.error("Error deleting pod", response);
          toast.error("Error deleting pod");
        }
      } catch (error) {
        console.error("Error deleting pod", error);
        toast.error("Error deleting pod");
      }
    }
  }

  $: {
    // Update nodes and edges
    $pods = $pods ?? [];
    $services = $services ?? [];
    $ingresses = $ingresses ?? [];
    $secrets = $secrets ?? [];
    $pvc = $pvc ?? [];

    // connect secrets to pods left (pod) to right (secret)
    $pods.forEach((pod) => {
      addNode(pod);
      $secrets.forEach((secret) => {
        if (pod.object.metadata.labels && secret.object.metadata.labels) {
          if (
            pod.object.metadata.labels["app.kubernetes.io/name"] ===
            secret.object.metadata.labels["app.kubernetes.io/name"]
          ) {
            addNode(secret);
            addEdge(pod.name, secret.name);
          }
        }
      });
    });

    // connect pvcs to pods left (pod) to right (pvc)
    $pods.forEach((pod) => {
      addNode(pod);
      $pvc.forEach((pvc) => {
        if (pod.object.metadata.labels && pvc.object.metadata.labels) {
          if (
            pod.object.metadata.labels["app.kubernetes.io/name"] ===
            pvc.object.metadata.labels["app.kubernetes.io/name"]
          ) {
            addNode(pvc);
            addEdge(pod.name, pvc.name);
          }
        }
      });
    });

    // connect pods to services left (service) to right (pod)
    $pods.forEach((pod) => {
      addNode(pod);
      $services.forEach((service) => {
        if (pod.object.metadata.labels && service.object.metadata.labels) {
          if (
            pod.object.metadata.labels["app.kubernetes.io/name"] ===
            service.object.metadata.labels["app.kubernetes.io/name"]
          ) {
            addNode(service);
            addEdge(service.name, pod.name);
          }
        }
      });
    });

    // connect services to ingresses left (ingress) to right (service) but only if they have the same prefix until "-ingress" & "-svc" suffix
    $services.forEach((service) => {
      addNode(service);
      // get the prefix of the service name
      const servicePrefix = service.object.metadata.name.split("-svc")[0];
      $ingresses.forEach((ingress) => {
        // get the prefix of the ingress name
        const ingressPrefix = ingress.object.metadata.name.split("-ingress")[0];
        if (servicePrefix === ingressPrefix) {
          addNode(ingress);
          addEdge(ingress.name, service.name);
        }
      });
    });
  }
</script>

<div class="flex items-start justify-between">
  <div class="flex flex-col">
    <Heading tag="h2">Map</Heading>
    <P class="text-gray-500 dark:text-gray-400 text-sm">Map of your resources.</P>
  </div>
</div>

<!-- load 1s -->
{#if initialLoadComplete}
  <div class="mt-8" style="height: 50vh;">
    <SvelteFlow {nodes} {nodeTypes} {edges} {edgeTypes} {defaultEdgeOptions} fitView>
      <svg>
        <defs>
          <linearGradient id="edge-gradient">
            <stop offset="0%" stop-color="#28CAA3" />
            <stop offset="100%" stop-color="#28CAA3" />
          </linearGradient>
          <marker
            id="edge-arrow"
            markerWidth="10"
            markerHeight="10"
            refX="4"
            refY="3"
            orient="auto"
            markerUnits="strokeWidth"
          >
            <path d="M0,0 L0,6 L6,3 z" fill="url(#edge-gradient)" />
          </marker>
        </defs>
      </svg>
    </SvelteFlow>
  </div>
{:else}
  <div class="flex items-center justify-center h-96">
    <div class="flex items-center justify-center">
      <div
        class="w-8 h-8 border-2 border-t-primary-500 border-b-primary-500 rounded-full animate-spin"
      ></div>
    </div>
  </div>
{/if}

<Drawer
  placement="right"
  transitionType="fly"
  width="w-1/2"
  transitionParams={transitionParamsRight}
  bind:hidden={$drawerHidden}
>
  <div class="flex items-center h-full relative">
    <div class="absolute top-2 left-2">
      <p class="text-sm text-gray-500 dark:text-gray-400">
        {#if $selectedNode?.icon === "ingress"}
          <ArrowLeftRight class="inline" size="16" />
        {:else if $selectedNode?.icon === "service"}
          <NetworkIcon class="inline" size="16" />
        {:else if $selectedNode?.icon === "pod"}
          <Box class="inline" size="16" />
        {:else if $selectedNode?.icon === "secret"}
          <Lock class="inline" size="16" />
        {:else if $selectedNode?.icon === "pvc"}
          <Database class="inline" size="16" />
        {/if}

        {$selectedNode?.kind ?? "Node"}
      </p>
      <h5
        id="drawer-label"
        class="inline-flex items-center mb-4 text-base font-semibold text-gray-500 dark:text-gray-400"
      >
        {$selectedNode?.name ?? "Node"}
      </h5>
    </div>
    <div class="mb-4 dark:text-white absolute top-2 right-2 z-50 space-x-2">
      {#if $selectedNode?.kind == "pod"}
        <Button
          color="red"
          size="xs"
          on:click={() => {
            // api call to delete the object
            handleDeletePod($selectedNode?.name ?? "");
          }}><Trash class="inline" size="16" /></Button
        >
      {/if}
      <Button color="none" size="xs" on:click={() => ($drawerHidden = true)}>
        <X class="inline" size="16" />
      </Button>
    </div>
    <div class="absolute left-0 bottom-0 top-0 pt-20 w-full h-full">
      <Tabs style="underline">
        <TabItem open>
          <div slot="title" class="flex items-center gap-2">
            <FileCode />
            Manifest
          </div>
          <KubernetesObject manifest={$selectedNode?.object} />
        </TabItem>
        {#if $selectedNode?.kind == "pod"}
          <TabItem>
            <div slot="title" class="flex items-center gap-2">
              <ScrollText />
              Logs
            </div>
            <LogStream podName={$selectedNode?.name ?? ""} />
          </TabItem>
          <TabItem>
            <div slot="title" class="flex items-center gap-2">
              <Bell />
              Events
            </div>
            <EventStream name={$selectedNode?.name ?? ""} kind={$selectedNode?.kind ?? ""} />
          </TabItem>
        {/if}
      </Tabs>
    </div>
  </div>
</Drawer>
