<script lang="ts">
  import { Handle, Position, type NodeProps } from "@xyflow/svelte";

  import { ArrowLeftRight, Box, NetworkIcon, ShipWheel } from "lucide-svelte";

  type $$Props = NodeProps;
  export let data: $$Props["data"];
</script>

<div class="cloud gradient">
  <div>
    <ShipWheel size="16" />
  </div>
</div>
<div class="wrapper gradient">
  <div class="inner">
    <div class="body">
      {#if data.icon}
        <div class="icon">
          {#if data.icon === "ingress"}
            <ArrowLeftRight size="20" />
          {:else if data.icon === "service"}
            <NetworkIcon size="20" />
          {:else if data.icon === "pod"}
            <Box size="20" />
          {/if}
        </div>
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
