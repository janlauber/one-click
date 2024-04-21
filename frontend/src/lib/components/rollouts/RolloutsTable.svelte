<script lang="ts">
  import { client } from "$lib/pocketbase";
  import type { RolloutsRecord, RolloutsResponse } from "$lib/pocketbase/generated-types";
  import {
    rollouts,
    type Rexpand,
    updateDataStores,
    UpdateFilterEnum,
    currentRollout,
    currentRolloutStatus,
    type ExpandableResponse
  } from "$lib/stores/data";
  import selectedProjectId from "$lib/stores/project";
  import type { RolloutStatusResponse } from "$lib/types/status";
  import { formatDateTime, timeAgo } from "$lib/utils/date.utils";
  import { getRolloutEvents, getRolloutStatus } from "$lib/api/rollouts";
  import {
    Badge,
    Button,
    CloseButton,
    Drawer,
    Indicator,
    Modal,
    P,
    TableSearch,
    Toggle,
    Tooltip
  } from "flowbite-svelte";

  import { Dropdown, DropdownItem } from "flowbite-svelte";
  import {
    Copy,
    Database,
    Ellipsis,
    Eye,
    EyeOff,
    HardDrive,
    History,
    Info,
    Network,
    PanelRightOpen,
    Pause
  } from "lucide-svelte";
  import toast from "svelte-french-toast";
  import { sineIn } from "svelte/easing";
  import DiffLines from "../base/DiffLines.svelte";
  import { getRandomString } from "$lib/utils/random";
  import { type RolloutEventsResponse } from "$lib/types/events";

  let hidden6 = true;
  let showHiddenRollouts = false;
  let defaultModal = false;
  let searchTerm: string = "";
  let transitionParamsRight = {
    x: 320,
    duration: 200,
    easing: sineIn
  };

  let selectedRollout: RolloutsResponse<Rexpand> | undefined;
  let selectedRolloutStatus: RolloutStatusResponse | undefined;
  let selectedRolloutEvents: RolloutEventsResponse | undefined;
  let loadingStatus: boolean = false;
  let loadingEvents: boolean = false;
  let modalTitle1: string = "";
  let modalTitle2: string = "";

  const determineRolloutColor = (status: string): any => {
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

  async function toggleSidebar(rollout: RolloutsResponse<Rexpand>) {
    selectedRollout = rollout;

    hidden6 = !hidden6;
    loadingStatus = true;
    loadingEvents = true;

    if (rollout.endDate != "") {
      selectedRolloutStatus = undefined;
      selectedRolloutEvents = undefined;
      loadingStatus = false;
      loadingEvents = false;
      return;
    }

    await getRolloutStatus($selectedProjectId, selectedRollout.id)
      .then((res) => {
        selectedRolloutStatus = res;
      })
      .catch((err) => {
        console.log(err);
      })
      .finally(() => {
        loadingStatus = false;
      });

    await getRolloutEvents($selectedProjectId, selectedRollout.id)
      .then((res) => {
        selectedRolloutEvents = res;
      })
      .catch((err) => {
        console.log(err);
      })
      .finally(() => {
        loadingEvents = false;
      });
  }

  function confirmRollback(rollout: RolloutsResponse<Rexpand>) {
    selectedRollout = rollout;
    $currentRollout = $rollouts.find((r) => !r.endDate) as
      | ExpandableResponse<RolloutsResponse, Rexpand>
      | undefined;
    if ($currentRollout == undefined) {
      toast.error("No rollout selected.");
      return;
    }

    if (rollout.manifest == undefined) {
      toast.error("This rollout has no manifest.");
      return;
    }

    modalTitle1 = $currentRollout?.id ?? "None";
    modalTitle2 = rollout.id;

    defaultModal = true;
  }

  async function handleRollback(rollout: RolloutsResponse<Rexpand>) {
    // check if rollout is already deployed (if there is no end date)
    if (rollout.endDate == "") {
      toast.error("This rollout is already deployed.");
      return;
    }

    if (rollout.manifest == undefined) {
      toast.error("This rollout has no manifest.");
      return;
    }

    const data: RolloutsRecord = {
      manifest: rollout.manifest,
      startDate: new Date().toISOString(),
      endDate: "",
      project: $selectedProjectId,
      deployment: rollout.deployment,
      user: client.authStore.model?.id,
      hide: false
    };

    toast.promise(
      client
        .collection("rollouts")
        .update(rollout.id, data)
        .then(() => {
          updateDataStores({
            filter: UpdateFilterEnum.ALL,
            projectId: $selectedProjectId,
            deploymentId: rollout.deployment
          });
        }),
      {
        loading: "Deploying rollout...",
        success: "Rollout deployed.",
        error: "Error deploying rollout."
      }
    );
  }

  async function handleHide(rollout: RolloutsResponse<Rexpand>, hide: boolean = true) {
    if (rollout.endDate == "") {
      toast.error("This rollout is currently deployed.");
      return;
    }

    toast.promise(
      client
        .collection("rollouts")
        .update(rollout.id, { hide: hide })
        .then(() => {
          updateDataStores({
            filter: UpdateFilterEnum.ALL,
            projectId: $selectedProjectId,
            deploymentId: rollout.deployment
          });
        }),
      {
        loading: "Hiding rollout...",
        success: "Rollout hidden.",
        error: "Error hiding rollout."
      }
    );
  }

  // search rollouts by searchTerm in each rollout's manifest and id and name
  function flattenObject(obj: any, prefix = "") {
    return Object.keys(obj).reduce((acc, k) => {
      const pre = prefix.length ? prefix + "." : "";
      if (typeof obj[k] === "object") acc += flattenObject(obj[k], pre + k) + " ";
      else acc += pre + k + ":" + obj[k] + " ";
      return acc;
    }, "");
  }

  let filteredRollouts: RolloutsResponse<Rexpand>[] = [];

  $: filteredRollouts = $rollouts.filter((rollout) => {
    const searchTermLower = searchTerm.toLowerCase();
    const manifestString = rollout.manifest ? flattenObject(rollout.manifest).toLowerCase() : "";
    return (
      (manifestString.includes(searchTermLower) ||
        rollout.id.toLowerCase().includes(searchTermLower)) &&
      (showHiddenRollouts || !rollout.hide)
    );
  }) as RolloutsResponse<Rexpand>[];
</script>

<Drawer
  placement="right"
  transitionType="fly"
  width="max-w-xl w-full"
  transitionParams={transitionParamsRight}
  bind:hidden={hidden6}
  id="sidebar6"
>
  {#key selectedRollout}
    <div class="flex items-center">
      <h5
        id="drawer-label"
        class="inline-flex items-center mb-4 text-base font-semibold text-gray-500 dark:text-white"
      >
        <Info class="w-4 h-4 mr-2.5" /> Rollout details
      </h5>
      <CloseButton on:click={() => (hidden6 = true)} class="mb-4 dark:text-white" />
    </div>
    {#if !loadingStatus && !loadingEvents}
      <p class="mb-6 text-sm text-gray-500 dark:text-white">
        {#if !selectedRollout?.endDate}
          This rollout is currently deployed.
        {:else}
          This rollout has endet on {formatDateTime(selectedRollout?.endDate)}.
        {/if}
      </p>
      <div class="flex items-center mb-6">
        <div
          class="flex items-center justify-center w-10 h-10 mr-4
          {selectedRollout?.endDate ? 'bg-gray-100' : 'bg-green-100 text-green-500'}
         rounded-full"
        >
          <HardDrive class="w-5 h-5 " />
        </div>
        <div class="flex flex-col dark:text-white">
          <span class="text-sm font-light">Image</span>
          <span class="text-sm font-semibold">
            {#if selectedRollout?.manifest}
              {selectedRollout?.manifest.spec.image.registry}/{selectedRollout?.manifest.spec.image
                .repository}:{selectedRollout?.manifest.spec.image.tag}
            {/if}
          </span>
        </div>
      </div>
      {#if selectedRolloutStatus}
        <div class="flex items-center mb-6">
          <div
            class="flex items-center justify-center w-10 h-10 mr-4 {selectedRollout?.endDate
              ? 'bg-gray-100'
              : 'bg-green-100 text-green-500'} rounded-full"
          >
            <Copy class="w-5 h-5" />
          </div>
          <div class="flex flex-col dark:text-white">
            <span class="text-sm font-light">Replicas</span>
            <span class="text-sm font-semibold">
              {#if selectedRolloutStatus}
                {selectedRolloutStatus?.deployment?.replicas}
              {/if}
            </span>
          </div>
        </div>
      {/if}
      {#if selectedRolloutEvents}
        <div class="items-center mb-6 dark:text-white">
          <h5 class="mr-4 text-sm font-semibold text-gray-500 dark:text-white">
            Events <i>(up to 1h)</i>
          </h5>
          <div class="flex flex-col">
            {#if selectedRolloutEvents?.events.length == 0}
              <span class="text-sm font-light">No events</span>
            {:else}
              {#each selectedRolloutEvents?.events as event}
                <div class="flex items-center mb-2">
                  <div class="flex flex-col">
                    <span class="text-sm font-light">{event.typus}</span>
                    <span class="text-sm font-semibold">{event.message}</span>
                  </div>
                </div>
              {/each}
            {/if}
          </div>
        </div>
      {/if}
    {/if}
  {/key}
</Drawer>

<Modal title="Compare Rollouts" bind:open={defaultModal} size="xl" autoclose>
  <DiffLines
    jsonObject1={$currentRollout?.manifest ?? {}}
    jsonObject2={selectedRollout?.manifest ?? {}}
    title1={modalTitle1}
    title2={modalTitle2}
  />

  <svelte:fragment slot="footer">
    <Button
      on:click={() =>
        selectedRollout ? handleRollback(selectedRollout) : toast.error("No rollout selected!")}
      >Confirm</Button
    >
    <Button color="alternative">Cancel</Button>
  </svelte:fragment>
</Modal>

<div class="dark:text-white">
  <div class="mt-8 flow-root">
    <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
      <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
        <div class="p-0.5 shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
          <div
            class="flex items-center justify-between p-4 bg-white dark:bg-gray-800 sm:rounded-t-lg"
          >
            <TableSearch
              placeholder="Search rollouts..."
              hoverable={true}
              divClass="shadow-none"
              id={getRandomString(8)}
              bind:inputValue={searchTerm}
            />
            <div class="flex items-center">
              <Toggle class="ml-4" bind:checked={showHiddenRollouts}>
                Show hidden ({$rollouts.length - filteredRollouts.length})
              </Toggle>
            </div>
          </div>
          {#if $rollouts.length > 0}
            <table class="min-w-full divide-y divide-gray-300 dark:divide-gray-600">
              <thead class="bg-gray-50 dark:bg-gray-800">
                <tr>
                  <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-6"
                    >ID</th
                  >
                  <th scope="col" class="px-3 py-3.5 text-center text-sm font-semibold">
                    <HardDrive class="w-5 h-5 mx-auto" />
                    <Tooltip>Image</Tooltip>
                  </th>
                  <th scope="col" class="px-3 py-3.5 text-center text-sm font-semibold">
                    <Network class="w-5 h-5 mx-auto" />
                    <Tooltip>Interfaces</Tooltip>
                  </th>
                  <th scope="col" class="px-3 py-3.5 text-center text-sm font-semibold">
                    <Database class="w-5 h-5 mx-auto" />
                    <Tooltip>Volumes</Tooltip>
                  </th>
                  <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">Created</th>
                  <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">Status</th>
                  <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-6">
                    <span class="sr-only">Edit</span>
                  </th>
                </tr>
              </thead>
              <tbody
                class="divide-y divide-gray-200 dark:divide-gray-600 bg-white dark:bg-transparent"
              >
                {#each filteredRollouts as rollout, idx (rollout.id)}
                  <tr class="transition-all hover:bg-gray-50 dark:hover:bg-gray-800">
                    <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium sm:pl-6"
                      >{rollout.id}</td
                    ><td class="whitespace-nowrap px-3 py-4 text-sm">
                      {#if rollout.manifest}{rollout.manifest.spec.image.repository}:{rollout
                          .manifest.spec.image.tag}
                      {/if}
                    </td>
                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                      {#if rollout.manifest}
                        {rollout.manifest.spec.interfaces.length ?? "None"}
                      {/if}
                    </td>
                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                      {#if rollout.manifest}
                        {rollout.manifest.spec.volumes.length ?? "None"}
                      {/if}
                    </td>
                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                      <div>{timeAgo(rollout.created)}</div>
                      <Tooltip>{formatDateTime(rollout.created)}</Tooltip>
                    </td>
                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                      {#if rollout.endDate == ""}
                        <Badge
                          border
                          color={determineRolloutColor(
                            $currentRolloutStatus?.deployment?.status ?? ""
                          )}
                          class="relative pl-6"
                        >
                          <Indicator
                            size="sm"
                            color={determineRolloutColor(
                              $currentRolloutStatus?.deployment?.status ?? ""
                            )}
                            class="absolute left-2"
                          />
                          <Indicator
                            size="sm"
                            color={determineRolloutColor(
                              $currentRolloutStatus?.deployment?.status ?? ""
                            )}
                            class="absolute animate-ping left-2"
                          />
                          {$currentRolloutStatus?.deployment?.status ?? "Unknown"}
                        </Badge>

                        <!-- <Badge border color="green" class="relative pl-6 mr-2">
                          <Indicator size="sm" color="green" class="absolute left-2" />
                          <Indicator size="sm" color="green" class="absolute animate-ping left-2" />
                          Deployed</Badge
                        > -->
                      {:else if rollout.hide}
                        <Badge border color={"dark"} class="relative pl-2">
                          <EyeOff class="w-4 h-4 mr-2" />
                          Hidden
                        </Badge>
                        <Tooltip>{formatDateTime(rollout.endDate)}</Tooltip>
                      {:else}
                        <Badge border color={"dark"} class="relative pl-1">
                          <!-- <Indicator size="sm" color={"gray"} class="mr-2" /> -->
                          <Pause class="w-4 h-4 mr-1" />
                          <span class="text-center">Ended</span>
                        </Badge>
                        <Tooltip>{formatDateTime(rollout.endDate)}</Tooltip>
                      {/if}
                    </td>
                    <td class="relative whitespace-nowrap pr-4 text-right text-sm font-medium">
                      {#if rollout.endDate == ""}
                        <button
                          on:click={() => toggleSidebar(rollout)}
                          class="active:scale-95 transition-all duration-150 ease-in-out"
                        >
                          <PanelRightOpen
                            class="inline-block ml-2 cursor-pointer dark:text-white"
                          />
                        </button>
                      {:else}
                        <Ellipsis
                          class="dots-menu-{idx} dark:text-white inline-block cursor-pointer outline-none"
                        />
                        <Dropdown triggeredBy=".dots-menu-{idx}" class="p-0">
                          <DropdownItem
                            class="text-green-500"
                            on:click={() => confirmRollback(rollout)}
                          >
                            <History class="w-4 h-4 mr-2 inline-block" />
                            Rollback</DropdownItem
                          >
                          {#if !rollout.hide}
                            <DropdownItem class="text-red-500" on:click={() => handleHide(rollout)}>
                              <EyeOff class="w-4 h-4 mr-2 inline-block" />
                              Hide</DropdownItem
                            >
                          {:else}
                            <DropdownItem on:click={() => handleHide(rollout, false)}>
                              <Eye class="w-4 h-4 mr-2 inline-block" />
                              Unhide</DropdownItem
                            >
                          {/if}
                          <!-- <DropdownItem class="text-red-500" on:click={() => handleDelete(rollout)}
                            >Delete</DropdownItem
                          > -->
                        </Dropdown>
                      {/if}
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          {:else}
            <div class="flex justify-center items-center h-32">
              <P class="text-gray-500 dark:text-gray-400 text-sm">No rollouts found.</P>
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>
