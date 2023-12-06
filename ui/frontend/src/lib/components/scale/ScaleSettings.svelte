<script lang="ts">
  import { page } from "$app/stores";
  import { client } from "$lib/pocketbase";
  import type { RolloutsRecord, RolloutsResponse } from "$lib/pocketbase/generated-types";
  import { type Rexpand, rollouts, updateDataStores, UpdateFilterEnum } from "$lib/stores/data";
  import { Button, Heading, Input, Label, P } from "flowbite-svelte";
  import selectedProjectId from "$lib/stores/project";
  import toast from "svelte-french-toast";

  let current_rollout: RolloutsResponse<Rexpand> | undefined;
  let lastUpdatedRollout: RolloutsResponse<Rexpand> | undefined;
  let hadRolloutsOnLastPage: boolean = false;

  let minInstances: number = 1;
  let maxInstances: number = 1;
  let targetCPUUtilizationPercentage: number = 80;

  let cpuLimits: string = "200m";
  let cpuLimitsFloat: number = 0.2;
  let memoryLimits: string = "256Mi";
  let memoryLimitsInt: number = 256;
  let cpuRequests: string = "100m";
  let cpuRequestsFloat: number = 0.1;
  let memoryRequests: string = "128Mi";
  let memoryRequestsInt: number = 128;

  function convertFloatToCpuString(value: number): string {
    // 1 = 1000m
    // etc.
    return `${value * 1000}m`;
  }

  function convertCpuStringToFloat(value: string): number {
    // 1000m = 1
    // etc.
    return parseFloat(value.replace("m", "")) / 1000;
  }

  $: if ($rollouts.length > 0) {
    // get the current rollout on following priority:
    // 1. no endDate set
    // 2. newest endDate

    const temp_rollouts = $rollouts.filter((r) => !r.endDate);
    if (temp_rollouts.length > 0) {
      current_rollout = temp_rollouts[0];
    } else {
      current_rollout = $rollouts.sort((a, b) => {
        if (a.endDate && b.endDate) {
          return b.endDate.localeCompare(a.endDate);
        } else if (a.endDate) {
          return 1;
        } else if (b.endDate) {
          return -1;
        } else {
          return 0;
        }
      })[0];
    }

    if (current_rollout && current_rollout !== lastUpdatedRollout) {
      minInstances = current_rollout.manifest?.spec.horizontalScale.minReplicas ?? 1;
      maxInstances = current_rollout.manifest?.spec.horizontalScale.maxReplicas ?? 1;
      targetCPUUtilizationPercentage =
        current_rollout.manifest?.spec.horizontalScale.targetCPUUtilizationPercentage ?? 80;
      lastUpdatedRollout = current_rollout;
      cpuLimits = current_rollout.manifest?.spec.resources?.limits?.cpu ?? "200m";
      memoryLimits = current_rollout.manifest?.spec.resources?.limits?.memory ?? "256Mi";
      cpuRequests = current_rollout.manifest?.spec.resources?.requests?.cpu ?? "100m";
      memoryRequests = current_rollout.manifest?.spec.resources?.requests?.memory ?? "128Mi";
      cpuLimitsFloat = convertCpuStringToFloat(cpuLimits);
      memoryLimitsInt = parseInt(memoryLimits.replace("Mi", ""));
      cpuRequestsFloat = convertCpuStringToFloat(cpuRequests);
      memoryRequestsInt = parseInt(memoryRequests.replace("Mi", ""));
    }

    hadRolloutsOnLastPage = true;
  } else {
    // Reset all values when there are no rollouts
    if (hadRolloutsOnLastPage) {
      resetValues();
    }
    current_rollout = undefined;
    hadRolloutsOnLastPage = false;
  }

  // Reactive statement to track page changes
  $: $page,
    () => {
      if (!hadRolloutsOnLastPage) {
        resetValues();
      }
    };

  function resetValues() {
    minInstances = 1;
    maxInstances = 1;
    targetCPUUtilizationPercentage = 80;
    cpuLimits = "200m";
    memoryLimits = "256Mi";
    cpuRequests = "100m";
    memoryRequests = "128Mi";
    cpuLimitsFloat = 0.2;
    memoryLimitsInt = 256;
    cpuRequestsFloat = 0.1;
    memoryRequestsInt = 128;
  }

  function resetInput() {
    minInstances = current_rollout?.manifest?.spec.horizontalScale.minReplicas ?? 1;
    maxInstances = current_rollout?.manifest?.spec.horizontalScale.maxReplicas ?? 1;
    targetCPUUtilizationPercentage =
      current_rollout?.manifest?.spec.horizontalScale.targetCPUUtilizationPercentage ?? 80;
    cpuLimits = current_rollout?.manifest?.spec.resources?.limits?.cpu ?? "200m";
    memoryLimits = current_rollout?.manifest?.spec.resources?.limits?.memory ?? "256Mi";
    cpuRequests = current_rollout?.manifest?.spec.resources?.requests?.cpu ?? "100m";
    memoryRequests = current_rollout?.manifest?.spec.resources?.requests?.memory ?? "128Mi";
    cpuLimitsFloat = convertCpuStringToFloat(cpuLimits);
    memoryLimitsInt = parseInt(memoryLimits.replace("Mi", ""));
    cpuRequestsFloat = convertCpuStringToFloat(cpuRequests);
    memoryRequestsInt = parseInt(memoryRequests.replace("Mi", ""));
  }

  function handleInputChange(event: any, field: any) {
    switch (field) {
      case "minInstances":
        minInstances = event.target.value;
        break;
      case "maxInstances":
        maxInstances = event.target.value;
        break;
      case "targetCPUUtilizationPercentage":
        targetCPUUtilizationPercentage = event.target.value;
        break;
      case "cpuLimits":
        cpuLimits = event.target.value;
        break;
      case "memoryLimits":
        memoryLimits = event.target.value;
        break;
      case "cpuRequests":
        cpuRequests = event.target.value;
        break;
      case "memoryRequests":
        memoryRequests = event.target.value;
        break;
    }
  }

  async function handleInputSave() {
    if (current_rollout) {
      if (minInstances <= 0) {
        toast.error("Minimum instances must be greater than 0.");
        return;
      }

      if (minInstances > maxInstances) {
        toast.error("Minimum instances must be less or equal than maximum instances.");
        return;
      }

      if (targetCPUUtilizationPercentage <= 0) {
        toast.error("Target CPU utilization percentage must be greater than 0.");
        return;
      }

      if (targetCPUUtilizationPercentage > 100) {
        toast.error("Target CPU utilization percentage must be less or equal than 100.");
        return;
      }

      if (cpuLimitsFloat <= 0) {
        toast.error("CPU limits must be greater than 0.");
        return;
      }

      if (memoryLimitsInt <= 0) {
        toast.error("Memory limits must be greater than 0.");
        return;
      }

      if (cpuRequestsFloat <= 0) {
        toast.error("CPU requests must be greater than 0.");
        return;
      }

      if (memoryRequestsInt <= 0) {
        toast.error("Memory requests must be greater than 0.");
        return;
      }

      if (cpuLimitsFloat < cpuRequestsFloat) {
        toast.error("CPU limits must be greater or equal than CPU requests.");
        return;
      }

      if (memoryLimitsInt < memoryRequestsInt) {
        toast.error("Memory limits must be greater or equal than memory requests.");
        return;
      }

      const new_manifest = {
        ...current_rollout.manifest,
        spec: {
          // @ts-ignore
          ...current_rollout.manifest.spec,
          horizontalScale: {
            minReplicas: parseInt(minInstances.toString()),
            maxReplicas: parseInt(maxInstances.toString()),
            targetCPUUtilizationPercentage: parseInt(targetCPUUtilizationPercentage.toString())
          },
          resources: {
            limits: {
              cpu: convertFloatToCpuString(cpuLimitsFloat),
              memory: `${memoryLimitsInt}Mi`
            },
            requests: {
              cpu: convertFloatToCpuString(cpuRequestsFloat),
              memory: `${memoryRequestsInt}Mi`
            }
          }
        }
      };

      const data: RolloutsRecord = {
        manifest: new_manifest,
        startDate: current_rollout.startDate,
        endDate: "",
        project: $selectedProjectId,
        user: client.authStore.model?.id
      };

      toast.promise(
        client
          .collection("rollouts")
          .create(data)
          .then(() => {
            updateDataStores({
              filter: UpdateFilterEnum.ALL,
              projectId: $selectedProjectId
            });
          }),
        {
          loading: "Creating rollout...",
          success: "Rollout created.",
          error: "Error creating rollout."
        }
      );
    }
  }
</script>

<Heading tag="h3" class="mt-8">Horizontal Scaling</Heading>

<div class="mt-4 dark:text-white">
  <div class="flow-root">
    <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
      <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
        <div class="p-0.5 shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
          <table class="min-w-full divide-y divide-gray-300 dark:divide-gray-600">
            <tbody
              class="divide-y divide-gray-200 dark:divide-gray-600 bg-white dark:bg-transparent"
            >
              <tr class="transition-all hover:bg-gray-50 dark:hover:bg-gray-800">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-xs font-medium sm:pl-6">
                  <Heading tag="h5">Min. Instances</Heading>
                  <P class="text-gray-500 dark:text-gray-400 text-xs">
                    The minimum number of instances to be available at all times.
                  </P>
                </td><td class="whitespace-nowrap px-3 py-4 text-xs">
                  <Label for="tag" class="block "
                    >Instances
                    <span class={minInstances <= 0 ? "text-red-500" : "text-green-500"}> * </span>
                  </Label>
                  <Input
                    id="minInstances"
                    size="sm"
                    type="number"
                    bind:value={minInstances}
                    on:input={(e) => handleInputChange(e, "minInstances")}
                    placeholder="1"
                    class="
                    {minInstances <= 0 ? 'border-red-500' : 'border-green-500'}
                    "
                  />
                </td>
              </tr>
              <tr class="transition-all hover:bg-gray-50 dark:hover:bg-gray-800">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-xs font-medium sm:pl-6">
                  <Heading tag="h5">Max. Instances</Heading>
                  <P class="text-gray-500 dark:text-gray-400 text-xs">
                    The maximum number of instances to be available at all times.
                  </P>
                </td><td class="whitespace-nowrap px-3 py-4 text-xs">
                  <Label for="tag" class="block "
                    >Instances
                    <span class={maxInstances < minInstances ? "text-red-500" : "text-green-500"}>
                      *
                    </span>
                  </Label>
                  <Input
                    id="maxInstances"
                    size="sm"
                    type="number"
                    bind:value={maxInstances}
                    on:input={(e) => handleInputChange(e, "maxInstances")}
                    placeholder="1"
                    class="
                    {maxInstances < minInstances ? 'border-red-500' : 'border-green-500'}
                    "
                  />
                </td>
              </tr>
              <!-- targetCPUUtilizationPercentage scaling -->
              <tr class="transition-all hover:bg-gray-50 dark:hover:bg-gray-800">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-xs font-medium sm:pl-6">
                  <Heading tag="h5">Target CPU Utilization Percentage</Heading>
                  <P class="text-gray-500 dark:text-gray-400 text-xs">
                    The target CPU utilization percentage to trigger a scale up.
                  </P>
                </td><td class="whitespace-nowrap px-3 py-4 text-xs">
                  <Label for="tag" class="block "
                    >Percentage %
                    <span class={targetCPUUtilizationPercentage < 40 ? "text-red-500" : "text-green-500"}>
                      *
                    </span>
                  </Label>
                  <Input
                    id="targetCPUUtilizationPercentage"
                    size="sm"
                    type="number"
                    bind:value={targetCPUUtilizationPercentage}
                    on:input={(e) => handleInputChange(e, "targetCPUUtilizationPercentage")}
                    placeholder="1"
                    class="
                    {targetCPUUtilizationPercentage < 40 ? 'border-red-500' : 'border-green-500'}
                    "
                  />
                </td>
              </tr>
            </tbody>
          </table>
          <!-- Reset & Save Button bottom right -->
          <div class="flex justify-end mt-4 p-4">
            <Button
              color="alternative"
              class="whitespace-nowrap self-start mr-2"
              on:click={() => resetInput()}
            >
              Reset
            </Button>
            <Button
              color="primary"
              class="whitespace-nowrap self-start"
              on:click={() => handleInputSave()}
            >
              Save
            </Button>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<Heading tag="h3" class="mt-8">Vertical Scaling</Heading>

<div class="mt-4 dark:text-white">
  <div class="flow-root">
    <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
      <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
        <div class="p-0.5 shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
          <table class="min-w-full divide-y divide-gray-300 dark:divide-gray-600">
            <tbody
              class="divide-y divide-gray-200 dark:divide-gray-600 bg-white dark:bg-transparent"
            >
              <tr class="transition-all hover:bg-gray-50 dark:hover:bg-gray-800">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-xs font-medium sm:pl-6">
                  <Heading tag="h5">Requests</Heading>
                  <P class="text-gray-500 dark:text-gray-400 text-xs">
                    The reserved resources for each instance.
                  </P>
                </td><td class="whitespace-nowrap px-3 py-4 text-xs">
                  <Label for="tag" class="block ">CPU</Label>
                  <Input
                    id="cpuRequest"
                    size="sm"
                    type="number"
                    bind:value={cpuRequestsFloat}
                    on:input={(e) => handleInputChange(e, "cpuRequest")}
                    placeholder="0.1"
                  />

                  <Label for="tag" class="block mt-2">Memory</Label>
                  <Input
                    id="memoryRequest"
                    size="sm"
                    type="number"
                    bind:value={memoryRequestsInt}
                    on:input={(e) => handleInputChange(e, "memoryRequest")}
                    placeholder="128"
                  />
                </td>
              </tr>
              <tr class="transition-all hover:bg-gray-50 dark:hover:bg-gray-800">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-xs font-medium sm:pl-6">
                  <Heading tag="h5">Limits</Heading>
                  <P class="text-gray-500 dark:text-gray-400 text-xs">
                    The maximum resources for each instance.
                  </P>
                </td><td class="whitespace-nowrap px-3 py-4 text-xs">
                  <Label for="tag" class="block ">CPU</Label>
                  <Input
                    id="cpuLimits"
                    size="sm"
                    type="number"
                    bind:value={cpuLimitsFloat}
                    on:input={(e) => handleInputChange(e, "cpuLimits")}
                    placeholder="0.2"
                  />

                  <Label for="tag" class="block mt-2">Memory</Label>
                  <Input
                    id="memoryLimits"
                    size="sm"
                    type="number"
                    bind:value={memoryLimitsInt}
                    on:input={(e) => handleInputChange(e, "memoryLimits")}
                    placeholder="256"
                  />
                </td>
              </tr>
            </tbody>
          </table>
          <!-- Reset & Save Button bottom right -->
          <div class="flex justify-end mt-4 p-4">
            <Button
              color="alternative"
              class="whitespace-nowrap self-start mr-2"
              on:click={() => resetInput()}
            >
              Reset
            </Button>
            <Button
              color="primary"
              class="whitespace-nowrap self-start"
              on:click={() => handleInputSave()}
            >
              Save
            </Button>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
