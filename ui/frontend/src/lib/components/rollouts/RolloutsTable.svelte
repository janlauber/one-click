<script lang="ts">
  import { client } from "$lib/pocketbase";
  import type { RolloutsRecord, RolloutsResponse } from "$lib/pocketbase/generated-types";
  import { rollouts, type Rexpand, updateDataStores, UpdateFilterEnum } from "$lib/stores/data";
  import selectedProjectId from "$lib/stores/project";
  import type { RolloutStatusResponse } from "$lib/types/status";
  import { formatDateTime, timeAgo } from "$lib/utils/date.utils";
  import { getRolloutEvents, getRolloutStatus } from "$lib/utils/rollouts";
  import {
    Badge,
    Button,
    CloseButton,
    Drawer,
    Indicator,
    Modal,
    P,
    TableSearch,
    Tooltip
  } from "flowbite-svelte";

  import { Dropdown, DropdownItem } from "flowbite-svelte";
  import { DotsHorizontalOutline, InfoCircleSolid } from "flowbite-svelte-icons";
  import { Copy, HardDrive, PanelRightOpen } from "lucide-svelte";
  import toast from "svelte-french-toast";
  import { sineIn } from "svelte/easing";
  import DiffLines from "../base/DiffLines.svelte";

  let hidden6 = true;
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
  let currentRollout: RolloutsResponse<Rexpand> | undefined;
  let modalTitle1: string = "";
  let modalTitle2: string = "";

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
    currentRollout = $rollouts.find((r) => !r.endDate);

    if (currentRollout == undefined) {
      toast.error("There is no rollout to rollback to.");
      return;
    }

    if (currentRollout.manifest == undefined) {
      toast.error("The current rollout has no manifest.");
      return;
    }

    if (rollout.manifest == undefined) {
      toast.error("This rollout has no manifest.");
      return;
    }

    modalTitle1 = currentRollout.id;
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
      user: client.authStore.model?.id
    };

    toast.promise(
      client
        .collection("rollouts")
        .update(rollout.id, data)
        .then(() => {
          updateDataStores({
            filter: UpdateFilterEnum.ALL,
            projectId: $selectedProjectId
          });
        }),
      {
        loading: "Deploying rollout...",
        success: "Rollout deployed.",
        error: "Error deploying rollout."
      }
    );
  }

  async function handleDelete(rollout: RolloutsResponse<Rexpand>) {
    if (rollout.endDate == "") {
      toast.error("This rollout is currently deployed.");
      return;
    }

    toast.promise(
      client
        .collection("rollouts")
        .delete(rollout.id)
        .then(() => {
          updateDataStores({
            filter: UpdateFilterEnum.ALL,
            projectId: $selectedProjectId
          });
        }),
      {
        loading: "Deleting rollout...",
        success: "Rollout deleted.",
        error: "Error deleting rollout."
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

  $: filteredRollouts = $rollouts.filter((rollout) => {
    const searchTermLower = searchTerm.toLowerCase();
    const manifestString = rollout.manifest ? flattenObject(rollout.manifest).toLowerCase() : "";
    return (
      manifestString.includes(searchTermLower) || rollout.id.toLowerCase().includes(searchTermLower)
    );
  });
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
        <InfoCircleSolid class="w-4 h-4 mr-2.5" /> Rollout details
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
    {:else}{/if}
  {/key}
</Drawer>

<Modal title="Compare Rollouts" bind:open={defaultModal} size="xl" autoclose>
  <DiffLines
    jsonObject1={currentRollout?.manifest ?? {}}
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
          <TableSearch
            placeholder="Search rollouts..."
            hoverable={true}
            divClass="shadow-none"
            bind:inputValue={searchTerm}
          />
          {#if $rollouts.length > 0}
            <table class="min-w-full divide-y divide-gray-300 dark:divide-gray-600">
              <thead class="bg-gray-50 dark:bg-gray-800">
                <tr>
                  <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold sm:pl-6"
                    >ID</th
                  >
                  <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">Image</th>
                  <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">Created</th>
                  <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">Started</th>
                  <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold">Ended</th>
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
                      {#if rollout.manifest}
                        {rollout.manifest.spec.image.registry}/{rollout.manifest.spec.image
                          .repository}:{rollout.manifest.spec.image.tag}
                      {/if}
                    </td>
                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                      <div>{formatDateTime(rollout.created)}</div>
                    </td>
                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                      <div>{timeAgo(rollout.startDate)}</div>
                      <Tooltip>{formatDateTime(rollout.startDate)}</Tooltip>
                    </td>
                    <td class="whitespace-nowrap px-3 py-4 text-sm">
                      {#if rollout.endDate == ""}
                        <Badge border color="green" class="relative pl-6 mr-2">
                          <Indicator size="sm" color="green" class="absolute left-2" />
                          <Indicator size="sm" color="green" class="absolute animate-ping left-2" />
                          Deployed</Badge
                        >
                      {/if}
                      <div>{timeAgo(rollout.endDate)}</div>
                      <Tooltip>{formatDateTime(rollout.endDate)}</Tooltip>
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
                        <DotsHorizontalOutline
                          class="dots-menu-{idx} dark:text-white inline-block cursor-pointer"
                        />
                        <Dropdown triggeredBy=".dots-menu-{idx}">
                          <DropdownItem on:click={() => confirmRollback(rollout)}
                            >Rollback</DropdownItem
                          >
                          <DropdownItem class="text-red-500" on:click={() => handleDelete(rollout)}
                            >Delete</DropdownItem
                          >
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
