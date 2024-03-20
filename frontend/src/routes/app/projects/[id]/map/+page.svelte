<script lang="ts">
  import { Heading, P } from "flowbite-svelte";
  import { writable } from "svelte/store";
  import { SvelteFlow, Controls, type Node, type Edge } from "@xyflow/svelte";

  import "@xyflow/svelte/dist/style.css";
  import "./turbo.css";

  import { initialNodes, initialEdges } from "./nodes-and-edges";
  import TurboNode from "./TurboNode.svelte";
  import TurboEdge from "./TurboEdge.svelte";

  const nodes = writable<Node[]>(initialNodes);
  const edges = writable<Edge[]>(initialEdges);

  const nodeTypes = {
    turbo: TurboNode
  };

  const edgeTypes = {
    turbo: TurboEdge
  };

  const defaultEdgeOptions = {
    type: "turbo",
    markerEnd: "edge-circle"
  };
</script>

<div class="flex items-start justify-between">
  <div class="flex flex-col">
    <Heading tag="h2">Map</Heading>
    <P class="text-gray-500 dark:text-gray-400 text-sm">Map of your resources.</P>
  </div>
</div>

<div style="height:50vh;">
  <SvelteFlow {nodes} {nodeTypes} {edges} {edgeTypes} {defaultEdgeOptions} fitView>
    <Controls showLock={false} />
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
