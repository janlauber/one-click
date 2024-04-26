<script lang="ts">
  import type { DeploymentsResponse, RolloutsResponse } from "$lib/pocketbase/generated-types";
  import { blueprints, type Rexpand } from "$lib/stores/data";
  import { rollouts } from "$lib/stores/data";
  import selectedProjectId from "$lib/stores/project";
  import { timeAgo } from "$lib/utils/date.utils";
  import { Badge, Indicator, Tooltip } from "flowbite-svelte";
  import {
    Box,
    ChevronRight,
    Cog,
    Dot,
    ExternalLink,
    FileQuestion,
    HardDrive,
    Hash
  } from "lucide-svelte";
  import type { RolloutStatusResponse } from "$lib/types/status";
  import { getRolloutStatus } from "$lib/api/rollouts";
  import { onDestroy, onMount } from "svelte";
  import { determineRolloutColor } from "$lib/utils/color";
  import { recordLogoUrl } from "$lib/utils/blueprint.utils";
  export let deployment: DeploymentsResponse;

  let deploymentBlueprint = $blueprints.find((b) => b.id === deployment.blueprint);
  $: deploymentBlueprint = $blueprints.find((b) => b.id === deployment.blueprint);

  let current_rollout_status: RolloutStatusResponse | undefined;

  const updateCurrentRollout = () => {
    getRolloutStatus($selectedProjectId, deployment.id ?? "")
      .then((response) => {
        current_rollout_status = response;
      })
      .catch(() => {
        current_rollout_status = undefined;
      });
  };

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

  // filter $rollouts by $rollouts.expand.project
  let these_rollouts: RolloutsResponse<Rexpand>[] = [];
  $: these_rollouts = $rollouts.filter((r) => r.expand?.deployment?.id === deployment.id);

  let current_rollout: RolloutsResponse<Rexpand> | undefined;

  $: current_rollout = these_rollouts.find((r) => !r.endDate);

  type Ingress = {
    host: string;
    tls: boolean;
  };

  let ingresses: Ingress[] = [];

  $: if (current_rollout && current_rollout.manifest && current_rollout.manifest.spec.interfaces) {
    // @ts-expect-error - TS doesn't like the filter function
    current_rollout.manifest.spec.interfaces.forEach((inf) => {
      if (inf.ingress) {
        // @ts-expect-error - TS doesn't like the forEach function
        inf.ingress.rules.forEach((rule) => {
          ingresses.push({ host: rule.host, tls: inf.tls });
        });
      }
    });
  }
</script>

<a
  class="relative flex items-center space-x-4 py-4 px-4 lg:px-8 dark:text-gray-400 group hover:bg-gray-100 dark:hover:bg-gray-700 transition-all"
  href="deployments/{deployment.id}"
>
  {#if deployment?.avatar}
    <img
      src={recordLogoUrl(deployment)}
      alt="Tuple"
      class="h-9 w-9 flex-none rounded-lg object-cover p-1 border-2"
    />
  {:else if deploymentBlueprint?.avatar}
    <img
      src={recordLogoUrl(deploymentBlueprint)}
      alt="Tuple"
      class="h-9 w-9 flex-none rounded-lg object-cover p-1 border-2"
    />
    <Tooltip class="absolute">
      {deploymentBlueprint.name}
    </Tooltip>
  {:else}
    <FileQuestion class="h-9 w-9 flex-none  rounded-lg object-cover p-1 border-2" />
  {/if}
  <div class="min-w-0 flex-auto">
    <div class="flex items-center gap-x-3">
      <Indicator
        color={determineRolloutColor(current_rollout_status?.deployment?.status ?? "")}
        size="md"
        class="absolute"
      />
      {#if determineRolloutColor(current_rollout_status?.deployment?.status ?? "") === "green"}
        <Indicator
          color={determineRolloutColor(current_rollout_status?.deployment?.status ?? "")}
          size="md"
          class="absolute animate-ping"
        />
      {/if}
      <h2 class="min-w-0 text-sm font-semibold leading-6 ml-5">
        <a
          href="deployments/{deployment.id}"
          class="flex gap-x-2
          dark:text-gray-400 group-hover:text-gray-600 transition-all
        "
        >
          <span class="whitespace-nowrap">{deployment.name}</span>
          <span class="font-normal">/</span>
          <Badge color="dark" class="text-xs ">
            <Hash class="w-3 h-3 inline-block mr-1" />{deployment.id}
          </Badge>
        </a>
      </h2>
    </div>
    <div class="mt-3 flex items-center gap-x-2.5 text-xs leading-5">
      <p class="truncate">
        {#if ingresses.length > 0}
          {#each ingresses as ingress (ingress)}
            <a
              href={(ingress.tls ? "https://" : "http://") + ingress.host}
              target="_blank"
              rel="noopener noreferrer"
              class="text-blue-500 hover:underline"
            >
              {ingress.host}
              <ExternalLink class="w-4 h-4 inline-block ml-1" />
            </a>
            <br />
          {/each}
        {:else}
          <a href="deployments/{deployment.id}/network" class="text-blue-500 hover:underline">
            Configure Hosts
            <Cog class="w-4 h-4 inline-block" />
          </a>
        {/if}
      </p>
      <Dot class="w-2 h-2" />
      <Box class="w-4 h-4 inline-block -mr-2" />{current_rollout_status?.deployment?.replicas ?? 0}
      <Dot class="w-2 h-2" />
      <p class="whitespace-nowrap dark:text-gray-400">
        Last Rollout:
        {current_rollout ? timeAgo(current_rollout.startDate) : "No Rollouts"}
      </p>
    </div>
  </div>
  <div>
    {#if current_rollout}
      <div
        class="rounded-full flex-none py-1 px-2 text-xs font-medium ring-1 ring-inset ring-black dark:ring-white"
      >
        <HardDrive class="w-4 h-4 inline-block mr-1" />
        {current_rollout?.manifest?.spec?.image.repository}:{current_rollout?.manifest?.spec?.image
          .tag}
      </div>
    {:else}
      <div
        class="rounded-full flex-none py-1 px-2 text-xs font-medium ring-1 ring-inset ring-black dark:ring-white"
      >
        No Rollouts
      </div>
    {/if}
  </div>
  <ChevronRight class="w-4 h-4 text-gray-400 group-hover:translate-x-1 transition-all" />
</a>
