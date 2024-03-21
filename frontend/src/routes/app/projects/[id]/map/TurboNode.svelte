<script lang="ts">
  import { drawerHidden } from "$lib/stores/drawer";
  import { Handle, Position, type NodeProps } from "@xyflow/svelte";

  import { ArrowLeftRight, Box, Database, Lock, NetworkIcon } from "lucide-svelte";

  type $$Props = NodeProps;
  export let data: $$Props["data"];

  let statusClass = "";
  let containerStatusClass = "";
  // if Pending, then statusClass = "status-pending"
  // if Running, then statusClass = "status-ok"
  // else, then statusClass = "status-problematic"

  $: containerStatusClass =
    data.containerStatuses &&
    (data.containerStatuses as Array<any>)[0] &&
    (data.containerStatuses as Array<any>)[0].ready === false
      ? "status-problematic"
      : "status-ok";

  $: statusClass =
    data.status !== "Running"
      ? data.status === "Pending"
        ? "status-pending"
        : "status-problematic"
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
  }}
>
  <div class="inner">
    <div class="body">
      {#if data.icon}
        <div class="icon"></div>
      {/if}
      <div>
        <div class="title">{data.title}</div>
        {#if data.subline}
          <div class="subline">{data.subline}</div>
        {/if}
      </div>
    </div>
    <Handle type="target" position={Position.Left} />
    <Handle type="source" position={Position.Right} />
  </div>
</div>
