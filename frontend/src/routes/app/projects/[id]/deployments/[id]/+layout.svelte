<script lang="ts">
  import { deployments, updateCurrentRolloutStatus } from "$lib/stores/data";
  import { onDestroy, onMount } from "svelte";
  import { navigating } from "$app/stores";
  import { TabItem, Tabs } from "flowbite-svelte";
  import { recordLogoUrl } from "$lib/utils/blueprint.utils";

  const updateCurrentRollout = () => {
    updateCurrentRolloutStatus();
  };

  $: if ($navigating) {
    updateCurrentRollout();
  }

  let intervalId: any;

  // update rollout status every 5 seconds
  onMount(() => {
    updateCurrentRollout();
    intervalId = setInterval(() => {
      updateCurrentRollout();
    }, 5000);
  });

  onDestroy(() => {
    clearInterval(intervalId);
  });
</script>

<Tabs tabStyle="underline">
  {#each $deployments as deployment}
    <TabItem>
      <div slot="title" class="flex items-center gap-2">
        <img
          src={recordLogoUrl(deployment)}
          alt="Tuple"
          class="h-9 w-9 flex-none rounded-lg object-cover p-1"
        />
        {deployment.name}
      </div>
    </TabItem>
  {/each}
</Tabs>

<slot />
