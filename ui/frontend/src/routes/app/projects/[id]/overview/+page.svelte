<script lang="ts">
  import { goto } from "$app/navigation";
  import { Button, Card, Heading, Indicator } from "flowbite-svelte";
  import { ArrowRight, Copy, Network, Rocket } from "lucide-svelte";
  import selectedProjectId from "$lib/stores/project";
  import { rollouts, type Rexpand } from "$lib/stores/data";
  import { getRolloutMetrics, getRolloutStatus } from "$lib/utils/rollouts";
  import { onMount } from "svelte";
  import type { RolloutStatusResponse } from "$lib/types/status";
  import { navigating, page } from "$app/stores";
  import type { RolloutsResponse } from "$lib/pocketbase/generated-types";
  import MetricsChart from "$lib/components/projects/MetricsChart.svelte";

  let current_rollout_status: RolloutStatusResponse | undefined;
  let currentRollouts: RolloutsResponse<Rexpand>[] = [];
  let currentRollout: RolloutsResponse<Rexpand> | undefined;

  let cpuLimits = 0;
  let cpuRequests = 0;
  let cpuUsage = 0;

  let memoryLimits = 0;
  let memoryRequests = 0;
  let memoryUsage = 0;

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
        return "yellow";
      case "Not Ready":
        return "yellow";
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
    currentRollouts = $rollouts.filter((r) => r.expand?.project.id === $selectedProjectId);
    currentRollout = currentRollouts.find((r) => !r.endDate);
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
        cpuLimits = Number(current_rollout_status?.deployment.resources.limitSum.cpu);
        // round to 3 decimal places
        cpuLimits = Math.round((cpuLimits + Number.EPSILON) * 1000) / 1000;
        cpuRequests = Number(current_rollout_status?.deployment.resources.requestSum.cpu);
        // round to 3 decimal places
        cpuRequests = Math.round((cpuRequests + Number.EPSILON) * 1000) / 1000;
        memoryLimits =
          Number(current_rollout_status?.deployment.resources.limitSum.memory) / 1024 / 1024 / 1024;
        // round to 3 decimal places
        memoryLimits = Math.round((memoryLimits + Number.EPSILON) * 1000) / 1000;
        memoryRequests =
          Number(current_rollout_status?.deployment.resources.requestSum.memory) /
          1024 /
          1024 /
          1024;
        // round to 3 decimal places
        memoryRequests = Math.round((memoryRequests + Number.EPSILON) * 1000) / 1000;
      })
      .catch(() => {
        current_rollout_status = undefined;
        rollout_status_color = "gray";
      });

    getRolloutMetrics($selectedProjectId, currentRollout.id)
      .then((response) => {
        // sum up cpu and memory Usage
        cpuUsage = 0;
        memoryUsage = 0;
        if (!response) return;
        response.metrics.forEach((metric) => {
          //convert string to number
          cpuUsage += Number(metric.cpu);
          memoryUsage += Number(metric.memory);
        });
        // round to 3 decimal places
        cpuUsage = Math.round((cpuUsage + Number.EPSILON) * 1000) / 1000;
        // bytes to GiBi
        memoryUsage = memoryUsage / 1024 / 1024 / 1024;
        // round to 3 decimal places
        memoryUsage = Math.round((memoryUsage + Number.EPSILON) * 1000) / 1000;
      })
      .catch(() => {});
  };

  onMount(updateCurrentRollout);

  $: if ($navigating) {
    updateCurrentRollout();
  }
</script>

<div class="flex items-start justify-between">
  <Heading tag="h2">Overview</Heading>
  <Button
    color="alternative"
    size="xs"
    class="whitespace-nowrap self-start"
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
  <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
    <Card size="xl" class="flex flex-row p-2" padding="none">
      <div class="flex items-center justify-center w-10 h-10 bg-blue-100 rounded-lg text-blue-500">
        <Rocket
          class="w-5 h-5 text-primary-600
          justify-self-center
        "
        />
      </div>
      <div class="flex flex-col ml-4">
        <span class="text-sm font-light">Rollouts</span>
        <span class="text-sm font-semibold">{currentRollouts.length}</span>
      </div>
    </Card>
    <Card size="xl" class="flex flex-row p-2" padding="none">
      <div class="flex items-center justify-center w-10 h-10 bg-blue-100 rounded-lg text-blue-500">
        <Copy
          class="w-5 h-5 text-primary-600
          justify-self-center
        "
        />
      </div>
      <div class="flex flex-col ml-4">
        <span class="text-sm font-light">Replicas</span>
        <span class="text-sm font-semibold">{current_rollout_status?.deployment?.replicas ?? 0}</span
        >
      </div>
    </Card>
    <Card size="xl" class="flex flex-row p-2" padding="none">
      <div class="flex items-center justify-center w-10 h-10 bg-blue-100 rounded-lg text-blue-500">
        <Network
          class="w-5 h-5 text-primary-600
          justify-self-center
        "
        />
      </div>
      <div class="flex flex-col ml-4">
        <span class="text-sm font-light">Services</span>
        <span class="text-sm font-semibold">{current_rollout_status?.services?.length ?? 0}</span>
      </div>
    </Card>
  </div>

  <Heading tag="h3">Live Metrics</Heading>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    {#if currentRollout}
      <MetricsChart
        usage={cpuUsage}
        requests={cpuRequests}
        limits={cpuLimits}
        title="CPU (Cores)"
      />

      <MetricsChart
        usage={memoryUsage}
        requests={memoryRequests}
        limits={memoryLimits}
        title="Memory (GB)"
      />
    {/if}
  </div>
</div>
