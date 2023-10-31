<script lang="ts">
  import { Boxes, Database, FileCode2, FileLock2, Link2, Network } from "lucide-svelte";
  import type { NodeConfig } from "svelvet";
  export let dragCopy: any;

  const defaultNodeConfig: Partial<NodeConfig> = {
    label: "Node",
  };

  let toolbarItems = [
    {
      icon: Boxes,
      nodeConfig: {
        ...defaultNodeConfig,
        label: "Controller",
      } as NodeConfig // type assertion
    },
    {
      icon: Network,
      nodeConfig: {
        ...defaultNodeConfig,
        label: "Service",
      } as NodeConfig // type assertion
    },
    {
      icon: Link2,
      nodeConfig: {
        ...defaultNodeConfig,
        label: "Ingress",
      } as NodeConfig // type assertion
    },
    {
      icon: Database,
      nodeConfig: {
        ...defaultNodeConfig,
        label: "Volume",
      } as NodeConfig // type assertion
    },
    {
      icon: FileCode2,
      nodeConfig: {
        ...defaultNodeConfig,
        label: "ConfigMap",
      } as NodeConfig // type assertion
    },
    {
      icon: FileLock2,
      nodeConfig: {
        ...defaultNodeConfig,
        label: "Secret",
      } as NodeConfig // type assertion
    }
  ];

  // you can add more items to toolbarItems for more customization
</script>

<div
  class="absolute z-10 left-20 bg-background p-2 shadow-lg rounded-lg bottom-1/2 translate-y-1/2 -translate-x-1/2 border-2 border-primary select-none"
>
  {#each toolbarItems as item}
    <div
      use:dragCopy={item.nodeConfig}
      class="flex flex-col justify-center items-center hover:bg-gray-100 rounded-lg p-2 cursor-auto select-none text-primary draggable z-20"
    >
      <svelte:component this={item.icon} class="w-6 h-6" />
      <div class="text-xs font-semibold">{item.nodeConfig?.label}</div>
    </div>
  {/each}
</div>

<style>
  .draggable {
    cursor: grab;
  }
</style>
