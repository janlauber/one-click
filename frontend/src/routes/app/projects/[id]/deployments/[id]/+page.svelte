<script lang="ts">
  import { goto } from "$app/navigation";
  import { Button, Card, Heading, Indicator } from "flowbite-svelte";
  import {
    ArrowRight,
    Boxes,
    Database,
    HardDrive,
    History,
    Lock,
    Network,
    Variable
  } from "lucide-svelte";
  import selectedProjectId from "$lib/stores/project";
  import { rollouts, currentRolloutStatus, currentRollout } from "$lib/stores/data";
  import { getRolloutMetrics } from "$lib/api/rollouts";
  import { onDestroy, onMount } from "svelte";
  import { navigating } from "$app/stores";
  import MetricsChart from "$lib/components/deployments/MetricsChart.svelte";
  import RolloutChart from "$lib/components/deployments/RolloutChart.svelte";
  import { determineRolloutColor } from "$lib/utils/color";
  import selectedDeploymentId from "$lib/stores/deployment";

  let cpuRequests = 0;
  let cpuUsage = 0;

  let memoryRequests = 0;
  let memoryUsage = 0;

  const updateStats = () => {
    if (!$currentRolloutStatus?.deployment) {
      return;
    }

    cpuRequests = Number($currentRolloutStatus?.deployment.resources.requestSum.cpu);
    // round to 3 decimal places
    cpuRequests = Math.round((cpuRequests + Number.EPSILON) * 1000) / 1000;
    memoryRequests =
      Number($currentRolloutStatus?.deployment.resources.requestSum.memory) / 1024 / 1024 / 1024;
    // round to 3 decimal places
    memoryRequests = Math.round((memoryRequests + Number.EPSILON) * 1000) / 1000;

    getRolloutMetrics($selectedProjectId, $currentRollout?.deployment ?? "")
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

  let intervalId: any;

  // update rollout status every 5 seconds
  onMount(() => {
    updateStats();
    intervalId = setInterval(() => {
      updateStats();
    }, 5000);
  });

  onDestroy(() => {
    clearInterval(intervalId);
  });

  $: if ($navigating) {
    updateStats();
  }
</script>

<div class="flex items-start justify-between">
  <Heading tag="h2">Overview</Heading>
  <Button
    color="alternative"
    size="xs"
    class="whitespace-nowrap self-start"
    on:click={() => {
      goto(`/app/projects/${$selectedProjectId}/deployments/${$selectedDeploymentId}/rollouts`);
    }}
  >
    <div class="relative">
      {#if $currentRolloutStatus?.deployment?.status === "OK"}
        <Indicator
          size="sm"
          color={determineRolloutColor($currentRolloutStatus?.deployment?.status ?? "")}
          class="mr-1.5 animate-ping absolute "
        />
      {/if}
      <Indicator
        size="sm"
        color={determineRolloutColor($currentRolloutStatus?.deployment?.status ?? "none")}
        class="mr-1.5"
      />
    </div>
    Current rollout (Status: {$currentRolloutStatus?.deployment?.status ?? "Unknown"})
    <ArrowRight class="w-4 h-4 ml-2" />
  </Button>
</div>

<div class=" gap-4 space-y-4 mt-4">
  <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
    <a
      href={`/app/projects/${$selectedProjectId}/deployments/${$selectedDeploymentId}/rollouts`}
      class="flex flex-col justify-between"
    >
      <Card size="xl" class="flex flex-row p-2 bg-primary-500 text-white" padding="none">
        <div
          class="flex items-center justify-center w-10 h-10 bg-white rounded-lg text-primary-500"
        >
          <History
            class="w-5 h-5 text-primary-600
          justify-self-center
        "
          />
        </div>
        <div class="flex flex-col ml-4">
          <span class="text-sm font-light">Rollouts</span>
          <span class="text-sm font-semibold">{$rollouts.length}</span>
        </div>
      </Card>
    </a>
    <a href={`/app/projects/${$selectedProjectId}/deployments/${$selectedDeploymentId}/scale`}>
      <Card size="xl" class="flex flex-row p-2 bg-primary-500 text-white" padding="none">
        <div
          class="flex items-center justify-center w-10 h-10 bg-white rounded-lg text-primary-500"
        >
          <Boxes
            class="w-5 h-5 text-primary-600
          justify-self-center
        "
          />
        </div>
        <div class="flex flex-col ml-4">
          <span class="text-sm font-light">Instances</span>
          <span class="text-sm font-semibold"
            >{$currentRolloutStatus?.deployment?.replicas ?? 0}</span
          >
        </div>
      </Card>
    </a>
    <a href={`/app/projects/${$selectedProjectId}/deployments/${$selectedDeploymentId}/network`}>
      <Card size="xl" class="flex flex-row p-2 bg-primary-500 text-white" padding="none">
        <div
          class="flex items-center justify-center w-10 h-10 bg-white rounded-lg text-primary-500"
        >
          <Network
            class="w-5 h-5 text-primary-600
          justify-self-center
        "
          />
        </div>
        <div class="flex flex-col ml-4">
          <span class="text-sm font-light">Interfaces</span>
          <span class="text-sm font-semibold">{$currentRolloutStatus?.services?.length ?? 0}</span>
        </div>
      </Card>
    </a>
    <a
      href={`/app/projects/${$selectedProjectId}/deployments/${$selectedDeploymentId}/volumes`}
      class="flex flex-col justify-between"
    >
      <Card size="xl" class="flex flex-row p-2 bg-primary-500 text-white" padding="none">
        <div class="flex items-center justify-center w-10 h-10 bg-white rounded-lg text-black">
          <Database
            class="w-5 h-5 text-black
          justify-self-center
        "
          />
        </div>
        <div class="flex flex-col ml-4">
          <span class="text-sm font-light">Volumes</span>
          <span class="text-sm font-semibold"
            >{$currentRollout?.manifest?.spec?.volumes?.length ?? 0}</span
          >
        </div>
      </Card>
    </a>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-4 gap-4" style="grid-template-rows: auto 1fr">
    <a
      href={`/app/projects/${$selectedProjectId}/deployments/${$selectedDeploymentId}/image`}
      class="col-span-2"
    >
      <Card size="xl" class="flex flex-row p-2 text-primary-500" padding="none">
        <div
          class="flex items-center justify-center w-10 h-10 bg-primary-500 rounded-lg text-white"
        >
          <HardDrive
            class="w-5 h-5 text-white
          justify-self-center
        "
          />
        </div>
        <div class="flex flex-col ml-4">
          <span class="text-sm font-light">Image</span>
          <span class="text-sm font-semibold">
            {$currentRollout?.manifest?.spec.image.registry ??
              ""}/{$currentRollout?.manifest?.spec.image.repository.replace(/^library\//, "") ??
              ""}:{$currentRollout?.manifest?.spec.image.tag ?? ""}</span
          >
        </div>
      </Card>
    </a>
    <a
      href={`/app/projects/${$selectedProjectId}/deployments/${$selectedDeploymentId}/envs`}
      class="flex flex-col justify-between"
    >
      <Card size="xl" class="flex flex-row p-2 text-primary-500" padding="none">
        <div
          class="flex items-center justify-center w-10 h-10 bg-primary-500 rounded-lg text-white"
        >
          <Variable
            class="w-5 h-5 text-white
          justify-self-center
        "
          />
        </div>
        <div class="flex flex-col ml-4">
          <span class="text-sm font-light">Envs</span>
          <span class="text-sm font-semibold"
            >{$currentRollout?.manifest?.spec?.env?.length ?? 0}</span
          >
        </div>
      </Card>
    </a>
    <a
      href={`/app/projects/${$selectedProjectId}/deployments/${$selectedDeploymentId}/envs`}
      class="flex flex-col justify-between"
    >
      <Card size="xl" class="flex flex-row p-2 text-primary-500" padding="none">
        <div
          class="flex items-center justify-center w-10 h-10 bg-primary-500 rounded-lg text-white"
        >
          <Lock
            class="w-5 h-5 text-white
          justify-self-center
        "
          />
        </div>
        <div class="flex flex-col ml-4">
          <span class="text-sm font-light">Secrets</span>
          <span class="text-sm font-semibold"
            >{$currentRollout?.manifest?.spec?.secrets?.length ?? 0}</span
          >
        </div>
      </Card>
    </a>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    {#if $currentRollout}
      <MetricsChart usage={cpuUsage} requests={cpuRequests} title="Total CPU (Cores)" />

      <MetricsChart usage={memoryUsage} requests={memoryRequests} title="Total Memory (GB)" />
    {/if}
  </div>
  {#if !$rollouts}
    <div class="flex items-center justify-center">
      <div
        class="animate-spin rounded-full h-32 w-32 border-t-2 border-b-2 border-primary-500"
      ></div>
    </div>
  {:else}
    <RolloutChart />
  {/if}
</div>
