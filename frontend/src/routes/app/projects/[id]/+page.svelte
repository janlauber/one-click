<script lang="ts">
  import { goto } from "$app/navigation";
  import DeploymentCard from "$lib/components/deployments/DeploymentCard.svelte";
  import { deployments, loading } from "$lib/stores/data";
  import { Button, Heading, Spinner } from "flowbite-svelte";
  import { BookDashed } from "lucide-svelte";
</script>

<div class="absolute w-full top-0 bottom-0 overflow-y-scroll scrollbar-none">
  <div class="max-w-screen-2xl mx-auto pt-5 pb-20 flex flex-col">
    <div class="flex mb-5">
      <Heading tag="h5" class="flex font-normal items-center w-auto"
        >Deployments ({$deployments.length})</Heading
      >
      <div class="justify-self-end ml-auto">
        <Button
          color="primary"
          outline
          class="dark:text-white dark:border-white"
          size="sm"
          on:click={() => {
            goto("/app/blueprints/my-blueprints");
          }}
        >
          <BookDashed class="w-4 h-4 mr-2 inline-block" />
          Blueprints
        </Button>
      </div>
    </div>

    <ul role="list" class="divide-y dark:divide-white/5 divide-gray/5">
      {#if $loading}
        <div
          class="absolute top-0 left-0 right-0 bottom-0 flex justify-center items-center bg-gray-50 dark:bg-slate-800 z-20"
        >
          <span class="">
            <Spinner />
          </span>
        </div>
      {:else}
        {#each $deployments as deployment (deployment.id)}
          <DeploymentCard {deployment} />
        {/each}
        {#if $deployments.length === 0}
          <li class="p-4 text-center text-gray-500 dark:text-gray-400">No deployments found.</li>
        {/if}
      {/if}
    </ul>
  </div>
</div>
