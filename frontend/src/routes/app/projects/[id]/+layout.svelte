<script lang="ts">
  import { updateCurrentRolloutStatus } from "$lib/stores/data";
  import { onDestroy, onMount } from "svelte";
  import selectedProjectId from "$lib/stores/project";
  import { navigating } from "$app/stores";

  const updateCurrentRollout = () => {
    updateCurrentRolloutStatus($selectedProjectId);
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

<slot />
