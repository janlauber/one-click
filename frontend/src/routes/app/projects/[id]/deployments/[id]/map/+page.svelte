<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { writable, get } from "svelte/store";
  import { SvelteFlow, type Node, type Edge, Background, BackgroundVariant } from "@xyflow/svelte";
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
    TerminalSquare,
    Trash,
    X
  } from "lucide-svelte";
  import LogStream from "$lib/components/map/LogStream.svelte";
  import KubernetesObject from "$lib/components/map/KubernetesObject.svelte";
  import EventStream from "$lib/components/map/EventStream.svelte";
  import toast from "svelte-french-toast";
  import selectedDeploymentId from "$lib/stores/deployment";
  import ShellObject from "$lib/components/map/ShellObject.svelte";
  import selectedProjectId from "$lib/stores/project";

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
    let host = window.location.host;
    if (host.includes("localhost")) {
      host = "localhost:8090";
    }

    let protocol = window.location.protocol === "https:" ? "wss" : "ws";

    ws = new WebSocket(`${protocol}://${host}/ws/k8s/deployments`);

    ws.onopen = () => {
      ws.send(
        JSON.stringify({
          deploymentId: get(selectedDeploymentId) ?? "",
          projectId: get(selectedProject)?.id
        })
      );
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      updateObjects(data);
    };

    setTimeout(() => {
      initialLoadComplete = true;
    }, 500);
  });

  onDestroy(() => {
    ws.close();
    nodes.set([]);
    edges.set([]);
  });

  const pods = writable<NodeObject[]>([]);
  const services = writable<NodeObject[]>([]);
  const ingresses = writable<NodeObject[]>([]);
  const secrets = writable<NodeObject[]>([]);
  const pvc = writable<NodeObject[]>([]);
  const cronJobs = writable<NodeObject[]>([]);
  const jobs = writable<NodeObject[]>([]);
  const jobPods = writable<NodeObject[]>([]);

  function updateObjects(data: NodeObject) {
    const { labels } = data.object.metadata;
    if (!labels || labels["one-click.dev/projectId"] !== get(selectedProject)?.id) return;

    const deploymentIdLabel = labels["one-click.dev/deploymentId"];

    if (deploymentIdLabel && deploymentIdLabel !== get(selectedDeploymentId)) return;

    if (data.status === "ADDED") {
      addResource(data);
    } else if (data.status === "MODIFIED") {
      updateNode(data);
      updateResourceList(data);
    } else if (data.status === "DELETED") {
      removeNodeAndEdges(data.name);
      removeFromResourceList(data);
    } else if (data.status === "ERROR") {
      console.error(data);
    }
  }

  function addResource(data: NodeObject) {
    switch (data.kind) {
      case "pod":
        // if pod is controlled by a replica set add to pods array
        if (
          data.object.metadata.ownerReferences &&
          data.object.metadata.ownerReferences[0].kind === "ReplicaSet"
        ) {
          pods.update((p) => [...p, data]);
        } else {
          // if pod is controlled by a job add to jobPods array
          if (
            data.object.metadata.ownerReferences &&
            data.object.metadata.ownerReferences[0].kind === "Job"
          ) {
            jobPods.update((j) => [...j, data]);
          }
        }
        break;
      case "service":
        services.update((s) => [...s, data]);
        break;
      case "ingress":
        ingresses.update((i) => [...i, data]);
        break;
      case "secret":
        secrets.update((s) => [...s, data]);
        break;
      case "pvc":
        pvc.update((p) => [...p, data]);
        break;
      case "job":
        jobs.update((j) => [...j, data]);
        break;
      case "cronjob":
        cronJobs.update((c) => [...c, data]);
        break;
    }
  }

  function updateResourceList(data: NodeObject) {
    const updateList = (list: any) =>
      list.update((items: any) => {
        const index = items.findIndex(
          (item: any) => item.object.metadata.uid === data.object.metadata.uid
        );
        if (index !== -1) {
          items[index] = data;
        }
        return items;
      });

    switch (data.kind) {
      case "pod":
        updateList(pods);
        break;
      case "service":
        updateList(services);
        break;
      case "ingress":
        updateList(ingresses);
        break;
      case "secret":
        updateList(secrets);
        break;
      case "pvc":
        updateList(pvc);
        break;
      case "job":
        updateList(jobs);
        break;
      case "cronjob":
        updateList(cronJobs);
        break;
    }
  }

  function removeFromResourceList(data: NodeObject) {
    const removeList = (list: any) =>
      list.update((items: any) => items.filter((item: any) => item.name !== data.name));

    switch (data.kind) {
      case "pod":
        removeList(pods);
        break;
      case "service":
        removeList(services);
        break;
      case "ingress":
        removeList(ingresses);
        break;
      case "secret":
        removeList(secrets);
        break;
      case "pvc":
        removeList(pvc);
        break;
      case "job":
        removeList(jobs);
        break;
      case "cronjob":
        removeList(cronJobs);
        break;
    }
  }

  function createNode(data: NodeObject): Node {
    return {
      id: data.name,
      type: "turbo",
      data: {
        kind: data.kind,
        name: data.name,
        namespace: data.namespace,
        labels: data.labels,
        status: data.kind === "pod" ? data.object.status.phase : "Running",
        object: data.object,
        icon: data.kind.toLowerCase(),
        containerStatuses: data.kind === "pod" ? data.object.status.containerStatuses : []
      },
      position: calculatePosition(data)
    };
  }

  function calculatePosition(data: NodeObject): { x: number; y: number } {
    const BASE_X_POSITIONS = {
      ingress: 0,
      service: 450,
      pod: 900,
      secret: 1400,
      pvc: 1400,
      cronjob: 0,
      job: 450,
      jobPod: 900
    };

    const Y_OFFSET = 100;

    let basePosition = BASE_X_POSITIONS[data.kind] || 0;
    if (data.kind === "pod" && data.object.metadata.ownerReferences[0].kind === "Job") {
      basePosition = BASE_X_POSITIONS.jobPod;
    }

    let tempPosition = { x: basePosition, y: 0 };
    if (
      data.kind === "job" ||
      data.kind === "cronjob" ||
      (data.kind === "pod" && data.object.metadata.ownerReferences[0].kind === "Job")
    ) {
      tempPosition.y = 500;
    }

    nodes.subscribe((n) => {
      n.forEach((node) => {
        if (
          node.position.x === tempPosition.x &&
          node.position.y === tempPosition.y &&
          node.data.kind === data.kind
        ) {
          tempPosition.y += Y_OFFSET;
        }
        // kind secret and pvc can have the same x position and need a y offset
        if (
          (node.position.x === BASE_X_POSITIONS.secret ||
            node.position.x === BASE_X_POSITIONS.pvc) &&
          node.position.y === tempPosition.y &&
          (data.kind === "secret" || data.kind === "pvc")
        ) {
          tempPosition.y += Y_OFFSET;
        }
      });
    })();

    return tempPosition;
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
      id: `${source}-${target}`,
      source,
      target,
      type: "turbo",
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
    nodes.update((n) => n.filter((node) => node.id !== objectName));
    edges.update((e) =>
      e.filter((edge) => edge.source !== objectName && edge.target !== objectName)
    );
    // remove also nodes if it's from jobPods
    jobPods.update((j) => j.filter((jobPod) => jobPod.name !== objectName));
  }

  function updateNode(data: NodeObject) {
    nodes.update((n) => {
      const nodeIndex = n.findIndex((node) => node.id === data.name);
      if (nodeIndex !== -1) {
        const updatedNode = {
          ...n[nodeIndex],
          data: { ...n[nodeIndex].data, ...createNode(data).data }
        };
        return [...n.slice(0, nodeIndex), updatedNode, ...n.slice(nodeIndex + 1)];
      }
      return n;
    });
  }

  async function handleDelete(name: string, kind: string) {
    const token = localStorage.getItem("pocketbase_auth");
    if (!token) {
      return;
    }
    const authHeader = { Authorization: `Bearer ${JSON.parse(token).token}` };
    const url =
      window.location.hostname === "localhost"
        ? `http://localhost:8090/pb/${get(selectedProject)?.id}/${kind}/${name}`
        : `/pb/${get(selectedProject)?.id}/${kind}/${name}`;
    try {
      const response = await fetch(url, {
        method: "DELETE",
        headers: {
          ...authHeader,
          "Content-Type": "application/json"
        }
      });
      if (response.status === 200) {
        toast.success(`${kind} ${name} deleted`);
        drawerHidden.set(true);
      } else {
        toast.error(`Failed to delete ${kind} ${name}`);
      }
    } catch (error) {
      toast.error(`Failed to delete ${kind} ${name}`);
    }
  }

  $: {
    // Update nodes and edges
    $pods = $pods ?? [];
    $services = $services ?? [];
    $ingresses = $ingresses ?? [];
    $secrets = $secrets ?? [];
    $pvc = $pvc ?? [];
    $cronJobs = $cronJobs ?? [];
    $jobs = $jobs ?? [];
    $jobPods = $jobPods ?? [];

    // get all cronjobs (left) connect to jobs (right) if the job has the name of the cronjob as ownerReference
    $cronJobs.forEach((cronJob) => {
      addNode(cronJob);
      $jobs.forEach((job) => {
        if (
          job.object.metadata.ownerReferences &&
          job.object.metadata.ownerReferences[0].kind === "CronJob" &&
          job.object.metadata.ownerReferences[0].name === cronJob.name
        ) {
          addNode(job);
          addEdge(cronJob.name, job.name);
        }
      });
    });

    // add pods which are generated by jobs and connect to the job
    // example metadata
    // ownerReferences:
    // - apiVersion: batch/v1
    //   kind: Job
    //   name: some-bash-job-28613928
    //   uid: 69727eb6-d3a7-4f6e-84aa-30947103431a
    //   controller: true
    //   blockOwnerDeletion: true
    $jobs.forEach((job) => {
      // pods which are controlled by a job add to the jobPods array
      $jobPods.forEach((jobPod) => {
        if (
          jobPod.object.metadata.ownerReferences &&
          jobPod.object.metadata.ownerReferences[0].kind === "Job" &&
          jobPod.object.metadata.ownerReferences[0].name === job.name
        ) {
          addNode(jobPod);
          addEdge(job.name, jobPod.name);
        }
      });
    });

    // add pods which are controlled by a replica set
    // ownerReferences:
    // - apiVersion: apps/v1
    //   kind: ReplicaSet
    //   name: ipt658nmh93f0sj-7c8bbb88c6
    //   uid: 37aa0d69-f5c1-4f20-943e-1353b2c317a3
    //   controller: true
    //   blockOwnerDeletion: true
    $pods.forEach((pod) => {
      if (
        pod.object.metadata.ownerReferences &&
        pod.object.metadata.ownerReferences[0].kind === "ReplicaSet"
      ) {
        addNode(pod);
        addEdge(pod.object.metadata.ownerReferences[0].name, pod.name);
      }
    });

    // add services on the left and connect to pods which are controlled by the replica set
    $services.forEach((service) => {
      addNode(service);
      $pods.forEach((pod) => {
        if (
          pod.object.metadata.ownerReferences &&
          pod.object.metadata.ownerReferences[0].kind === "ReplicaSet"
        ) {
          addEdge(service.name, pod.name);
        }
      });
    });

    // add ingresses on the left and connect to services on the right
    // connect services to ingresses left (ingress) to right (service) but only if they have the same prefix until "-ingress" & "-svc" suffix
    $services.forEach((service) => {
      $ingresses.forEach((ingress) => {
        if (
          service.name.includes("-svc") &&
          ingress.name.includes("-ingress") &&
          service.name.split("-svc")[0] === ingress.name.split("-ingress")[0]
        ) {
          addNode(ingress);
          addEdge(ingress.name, service.name);
        }
      });
    });

    // add secrets next to the pods which are controlled by the replica set
    $secrets.forEach((secret) => {
      addNode(secret);
      $pods.forEach((pod) => {
        if (
          pod.object.metadata.ownerReferences &&
          pod.object.metadata.ownerReferences[0].kind === "ReplicaSet"
        ) {
          addEdge(pod.name, secret.name);
        }
      });
    });

    // add pvc next to the pods which are controlled by the replica set
    $pvc.forEach((pvc) => {
      addNode(pvc);
      $pods.forEach((pod) => {
        if (
          pod.object.metadata.ownerReferences &&
          pod.object.metadata.ownerReferences[0].kind === "ReplicaSet"
        ) {
          addEdge(pod.name, pvc.name);
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

{#if initialLoadComplete}
  <div class="absolute left-0 right-0 top-36 bottom-10">
    <SvelteFlow {nodes} {nodeTypes} {edges} {edgeTypes} {defaultEdgeOptions} fitView>
      <Background variant={BackgroundVariant.Dots} />
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
  <div class="absolute left-0 right-0 top-36 bottom-10 flex items-center justify-center">
    <div
      class="w-8 h-8 border-2 border-t-primary-500 border-b-primary-500 rounded-full animate-spin"
    ></div>
  </div>
{/if}

<Drawer
  placement="right"
  transitionType="fly"
  width="w-full sm:w-2/3 lg:w-1/2"
  transitionParams={transitionParamsRight}
  bind:hidden={$drawerHidden}
  activateClickOutside={false}
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
        {:else if $selectedNode?.icon === "job"}
          <FileCode class="inline" size="16" />
        {:else if $selectedNode?.icon === "cronjob"}
          <ScrollText class="inline" size="16" />
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
    <div class="mb-4 dark:text-gray-400 absolute top-2 right-2 space-x-2 z-50">
      {#if $selectedNode?.kind == "pod" || $selectedNode?.kind == "job"}
        <Button
          color="red"
          size="xs"
          on:click={() => handleDelete($selectedNode?.name ?? "", $selectedNode.kind)}
        >
          <Trash class="inline" size="16" />
        </Button>
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
              <TerminalSquare />
              Shell
            </div>
            <ShellObject podName={$selectedNode?.name ?? ""} projectId={$selectedProjectId} />
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
