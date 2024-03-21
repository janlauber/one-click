<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { writable, type Writable } from "svelte/store";
  import { SvelteFlow, Controls, type Node, type Edge } from "@xyflow/svelte";
  import "@xyflow/svelte/dist/style.css";
  import "./turbo.css";
  import { initialNodes, initialEdges } from "./nodes-and-edges";
  import TurboNode from "./TurboNode.svelte";
  import TurboEdge from "./TurboEdge.svelte";
  import { Heading, P } from "flowbite-svelte";
  import { selectedProject } from "$lib/stores/data";

  const nodes = writable<Node[]>(initialNodes);
  const edges = writable<Edge[]>(initialEdges);

  const nodeTypes = { turbo: TurboNode };
  const edgeTypes = { turbo: TurboEdge };
  const defaultEdgeOptions = { type: "turbo", markerEnd: "edge-circle" };

  let ws: WebSocket;

  onMount(() => {
    ws = new WebSocket("ws://localhost:8090/ws/k8s/rollouts");

    ws.onopen = () => {
      // Send any initial message if required, for example:
      ws.send(JSON.stringify({ rolloutId: $selectedProject?.id ?? "" }));
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      updateObjects(data);
    };
  });

  onDestroy(() => {
    ws.close();
  });

  type Object = {
    kind: "pod" | "service" | "ingress" | "secret" | "pvc";
    name: string;
    namespace: string;
    labels: Map<string, string>;
    status: "ADDED" | "MODIFIED" | "DELETED" | "ERROR";
    object: any;
  };

  const pods = writable<Object[]>([]);
  const services = writable<Object[]>([]);
  const ingresses = writable<Object[]>([]);
  const secrets = writable<Object[]>([]);
  const pvc = writable<Object[]>([]);

  function updateObjects(data: Object) {
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

  function createNode(data: Object): Node {
    return {
      id: data.name, // Ensure this is unique
      type: "turbo", // Using the custom node type
      data: {
        title: data.name,
        subline: data.kind,
        icon: data.kind.toLowerCase() // Adjust icon based on kind
      },
      position: calculatePosition(data)
    };
  }

  function calculatePosition(data: Object): { x: number; y: number } {
    // Calculate the position of the nodes
    // left to right: ingress -> service -> pod -> secret
    const BASE_X_POSITIONS = {
      ingress: 0,
      service: 250,
      pod: 500,
      secret: 950,
      pvc: 950
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

  function addNode(data: Object) {
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
      animated: true
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

  function updateNode(data: Object) {
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

    // connect services to ingresses left (ingress) to right (service)
    $services.forEach((service) => {
      addNode(service);
      $ingresses.forEach((ingress) => {
        if (service.object.metadata.labels && ingress.object.metadata.labels) {
          if (
            service.object.metadata.labels["app.kubernetes.io/name"] ===
            ingress.object.metadata.labels["app.kubernetes.io/name"]
          ) {
            addNode(ingress);
            addEdge(ingress.name, service.name);
          }
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
{#if $pods.length > 0 && $services.length > 0 && $ingresses.length > 0 && $secrets.length > 0 }
  <div style="height:50vh;">
    <SvelteFlow {nodes} {nodeTypes} {edges} {edgeTypes} {defaultEdgeOptions} fitView>
      <!-- <Controls showLock={false} /> -->
      <svg>
        <defs>
          <linearGradient id="edge-gradient">
            <stop offset="0%" stop-color="#ae53ba" />
            <stop offset="100%" stop-color="#2a8af6" />
          </linearGradient>
        </defs>
      </svg>
    </SvelteFlow>
  </div>
{:else}
  <div class="flex items-center justify-center h-96">
    <P class="text-gray-500 dark:text-gray-400 text-sm">No resources found.</P>
  </div>
{/if}
