<script lang="ts">
  import { drawerHidden, selectedNode } from "$lib/stores/drawer";
  import { Handle, Position, type NodeProps } from "@xyflow/svelte";

  import { ArrowLeftRight, Box, Database, Lock, NetworkIcon } from "lucide-svelte";

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  type $$Props = NodeProps;
  export let data: any;

  let statusClass = "";
  let containerStatusClass = "";
  // if Pending, then statusClass = "status-pending"
  // if Running, then statusClass = "status-ok"
  // else, then statusClass = "status-problematic"

  // if the metadata contains a key called "deletionTimestamp", then statusClass = "status-problematic"
  // else, statusClass = "status-ok"
  $: statusClass = data.object.metadata.deletionTimestamp
    ? "status-deleting"
    : data.status !== "Running"
      ? data.status === "Pending"
        ? "status-pending"
        : "status-problematic"
      : "status-ok";

  $: containerStatusClass = data.object.metadata.deletionTimestamp
    ? "status-deleting"
    : data.containerStatuses &&
        data.containerStatuses[0] &&
        data.containerStatuses[0].ready === false
      ? "status-problematic"
      : "status-ok";
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
    {/if}
  </div>
</div>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore missing-declaration -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
  class="wrapper {statusClass}"
  on:click={() => {
    $drawerHidden = false;
    // check if data is not null
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
          {:else if data.status === "Running"}
            <div class="status">Running</div>
          {:else if data.status === "Pending"}
            <div class="status">Pending</div>
          {:else}
            <div class="status">Problematic</div>
          {/if}
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
