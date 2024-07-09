<script lang="ts">
  import { drawerHidden, selectedNode } from "$lib/stores/drawer";
  import { Handle, Position, type NodeProps } from "@xyflow/svelte";

  import { ArrowLeftRight, Box, Database, Lock, NetworkIcon, Play, Timer } from "lucide-svelte";

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  type $$Props = NodeProps;

  export let data: any = {
    object: {
      metadata: {
        deletionTimestamp: null // or some timestamp value
      }
    },
    status: "Running", // This would be dynamically set
    containerStatuses: [
      {
        ready: true // This would be dynamically set
      }
    ],
    icon: "pod", // Example icon type
    kind: "pod", // Example kind
    name: "example-name", // Example name
    namespace: "default", // Example namespace
    labels: {} // Example labels
  };

  function getStatusClass(data: any) {
    if (data.object.metadata.deletionTimestamp) {
      return "status-deleting";
    }
    // if kind Job
    if (data.kind === "job") {
      // check first if conditions exist
      if (data.object.status.conditions && data.object.status.conditions.length > 0) {
        if (data.object.status.conditions[0].type === "Complete") {
          return "status-succeeded";
        }
        if (data.object.status.conditions[0].type === "Failed") {
          return "status-problematic";
        }
      }
    }
    switch (data.status) {
      case "Succeeded":
        return "status-succeeded";
      case "Pending":
        return "status-pending";
      case "Running":
        return "status-ok";
      default:
        return "status-problematic";
    }
  }

  function getContainerStatusClass(data: any) {
    if (data.object.metadata.deletionTimestamp) {
      return "status-deleting";
    }
    if (data.containerStatuses) {
      for (let containerStatus of data.containerStatuses) {
        if (
          containerStatus.state.terminated &&
          containerStatus.state.terminated.reason === "Completed"
        ) {
          return "status-succeeded";
        }
        if (containerStatus.state.terminated && containerStatus.state.terminated.exitCode !== 0) {
          return "status-problematic";
        }
        if (!containerStatus.ready) {
          return "status-pending";
        }
        if (containerStatus.ready) {
          return "status-ok";
        }
        return "status-problematic";
      }
    }
    // job
    if (data.kind === "job") {
      if (data.object.status.conditions && data.object.status.conditions.length > 0) {
        if (data.object.status.conditions[0].type === "Complete") {
          return "status-succeeded";
        }
        if (data.object.status.conditions[0].type === "Failed") {
          return "status-problematic";
        }
      }
    }
    // pvc status
    if (data.kind === "pvc") {
      if (data.object.status.phase === "Bound") {
        return "status-ok";
      }
      return "status-problematic";
    }

    return "status-ok";
  }

  $: statusClass = getStatusClass(data);
  $: containerStatusClass = getContainerStatusClass(data);
</script>

<div class="cloud {containerStatusClass}">
  <div>
    {#if data.icon === "ingress"}
      <ArrowLeftRight size="16" />
    {:else if data.icon === "service"}
      <NetworkIcon size="16" />
    {:else if data.icon === "pod"}
      <Box size="16" />
    {:else if data.icon === "secret"}
      <Lock size="16" />
    {:else if data.icon === "pvc"}
      <Database size="16" />
    {:else if data.icon === "cronjob"}
      <Timer size="16" />
    {:else if data.icon === "job"}
      <Play size="16" />
    {/if}
  </div>
</div>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
  class="wrapper {statusClass}"
  on:click={() => {
    $drawerHidden = false;
    if (data) {
      $selectedNode = {
        kind: data.kind,
        name: data.name,
        namespace: data.namespace,
        labels: data.labels,
        icon: data.icon,
        object: data.object
      };
    }
  }}
>
  <div class="inner">
    <div class="body">
      {#if data.icon}
        <div class="icon"></div>
      {/if}
      <div>
        {#if data.kind === "pod"}
          {#if data.object.metadata.deletionTimestamp}
            <div class="status">Terminating</div>
          {:else}
            <div class="status">
              {data.status}
            </div>
          {/if}
        {/if}
        {#if data.kind === "job"}
          <div class="status">
            <!-- check if conditions exists -->
            {#if data.object.status.conditions && data.object.status.conditions.length > 0}
              {data.object.status.conditions[0]?.type}
            {/if}
            <!-- {data.object.status.conditions[0]?.type} -->
          </div>
        {/if}
        <div class="title">{data.name}</div>
        {#if data.kind}
          <div class="subline">{data.kind}</div>
        {/if}
      </div>
    </div>
    <Handle type="target" position={Position.Left} />
    <Handle type="source" position={Position.Right} />
  </div>
</div>
