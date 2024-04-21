<script lang="ts">
  import { deployments, selectedProject, updateCurrentRolloutStatus } from "$lib/stores/data";
  import { onDestroy, onMount } from "svelte";
  import { navigating } from "$app/stores";
  import DeploymentTab from "$lib/components/deployments/DeploymentTab.svelte";
  import { metadata } from "$lib/stores/metadata";
  import selectedDeploymentId from "$lib/stores/deployment";

  let currentDeployment = $deployments.find((d) => d.id === $selectedDeploymentId);

  $metadata.title = currentDeployment?.name + " | " + $selectedProject?.name || "Rollouts";

  const updateCurrentRollout = () => {
    updateCurrentRolloutStatus();
  };

  $: if ($navigating) {
    updateCurrentRollout();
  }

  let intervalId: any;

  // update rollout status every 5 seconds
  onMount(() => {
    intervalId = setInterval(updateCurrentRollout, 5000);
  });

  onDestroy(() => {
    clearInterval(intervalId);
  });
</script>

<DeploymentTab deployments={$deployments} />

<slot />
