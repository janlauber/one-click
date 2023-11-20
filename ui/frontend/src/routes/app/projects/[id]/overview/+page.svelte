<script lang="ts">
  import { goto } from "$app/navigation";
  import { Button, Card, Heading, Indicator } from "flowbite-svelte";
  import { ArrowRight } from "lucide-svelte";
  import selectedProjectId from "$lib/stores/project";
  import CpuChart from "$lib/components/projects/CpuChart.svelte";
  import { rollouts } from "$lib/stores/data";
  import { getRolloutStatus } from "$lib/utils/rollouts";
  import { onMount } from "svelte";
  import type { RolloutStatusResponse } from "$lib/types/status";
  import { navigating } from "$app/stores";

  let current_rollout_status: RolloutStatusResponse | undefined;
  let rollout_status_color:
    | "gray"
    | "red"
    | "yellow"
    | "green"
    | "indigo"
    | "purple"
    | "blue"
    | "dark"
    | "orange"
    | "none"
    | "teal"
    | undefined;

  const determineRolloutColor = (status: string) => {
    switch (status) {
      case "Pending":
        return "orange";
      case "Error":
        return "red";
      case "OK":
        return "green";
      default:
        return "gray";
    }
  };

  const updateCurrentRollout = () => {
    //@ts-ignore
    const currentRollouts = $rollouts.filter((r) => r.expand?.project.id === $selectedProjectId);
    const currentRollout = currentRollouts.find((r) => !r.endDate);
    if (!currentRollout) {
      current_rollout_status = undefined;
      rollout_status_color = "gray";
      return;
    }
    getRolloutStatus($selectedProjectId, currentRollout.id)
      .then((response) => {
        current_rollout_status = response;
        rollout_status_color = determineRolloutColor(
          current_rollout_status?.deployment.status ?? ""
        );
      })
      .catch(() => {
        current_rollout_status = undefined;
        rollout_status_color = "gray";
      });
  };

  onMount(updateCurrentRollout);

  $: if ($navigating) {
    updateCurrentRollout();
  }
</script>

<div class="flex items-center justify-between">
  <Heading tag="h2">Overview</Heading>
  <Button
    color="alternative"
    size="xs"
    class="whitespace-nowrap"
    on:click={() => {
      goto(`/app/projects/${$selectedProjectId}/rollouts`);
    }}
  >
    <div class="relative">
      <Indicator
        size="sm"
        color={rollout_status_color}
        class="mr-1.5 {current_rollout_status ? 'absolute' : ''}"
      />
      {#if current_rollout_status}
        <Indicator size="sm" color={rollout_status_color} class="mr-1.5 animate-ping" />
      {/if}
    </div>
    Current rollout
    <ArrowRight class="w-4 h-4 ml-2" />
  </Button>
</div>

<div class=" gap-4 space-y-4 mt-4">
  <Card size="xl"></Card>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <Card size="xl">
      <!-- <CpuChart /> -->
    </Card>
    <Card size="xl"></Card>
  </div>
</div>
