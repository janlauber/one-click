<script lang="ts">
  import DeploymentCard from "$lib/components/deployments/DeploymentCard.svelte";
  import { deployments } from "$lib/stores/data";
  import { Heading } from "flowbite-svelte";
</script>

<div class="absolute w-full top-0 bottom-0 overflow-y-scroll scrollbar-none">
  <div class="max-w-6xl mx-auto px-5 pb-5 py-5 flex flex-col">
    <div class="flex mb-5">
      <Heading tag="h5" class="flex font-normal items-center w-auto"
        >Your Deployments ({$deployments.length})</Heading
      >
    </div>

    <ul role="list" class="divide-y dark:divide-white/5 divide-gray/5">
      {#each $deployments as deployment (deployment.id)}
        <DeploymentCard {deployment} />
      {/each}
      {#if $deployments.length === 0}
        <li class="p-4 text-center text-gray-500 dark:text-gray-400">No deployments found.</li>
      {/if}
    </ul>
  </div>
</div>
