<script lang="ts">
  import { client } from "$lib/pocketbase";
  import type {
    BlueprintsResponse,
    RolloutsRecord,
    DeploymentsRecord
  } from "$lib/pocketbase/generated-types";
  import { blueprints, updateDataStores } from "$lib/stores/data";
  import selectedProjectId from "$lib/stores/project";
  import { recordLogoUrl } from "$lib/utils/blueprint.utils";

  import { Button, Input, Label } from "flowbite-svelte";
  import { ArrowRight, BookLock, BookUser, XIcon } from "lucide-svelte";
  import toast from "svelte-french-toast";

  export let deploymentModal: boolean;

  let name: string = "";

  let filteredBlueprints: BlueprintsResponse[] = [];
  let selectedBlueprint: BlueprintsResponse;

  $: filteredBlueprints = $blueprints.filter(
    (blueprint) => blueprint.owner === client.authStore.model?.id
  );
  $: selectedBlueprint = filteredBlueprints[0];

  async function handleCreateDeployment(event: Event) {
    event.preventDefault();

    if (!name) {
      toast.error("Please enter a name");
      return;
    }

    if (!selectedBlueprint) {
      toast.error("Please select a blueprint");
      return;
    }

    const deployment: DeploymentsRecord = {
      name: name,
      blueprint: selectedBlueprint.id,
      user: client.authStore.model?.id,
      project: $selectedProjectId
    };

    // remove any ingresses from the blueprint manifest
    // first, check if the manifest has interfaces:
    if ((selectedBlueprint.manifest as any).spec.interfaces) {
      // then, remove the ingress object from the interfaces array
      for (let i = 0; i < (selectedBlueprint.manifest as any).spec.interfaces.length; i++) {
        // remove the ingress object from the interfaces array
        if ((selectedBlueprint.manifest as any).spec.interfaces[i].ingress) {
          delete (selectedBlueprint.manifest as any).spec.interfaces[i].ingress;
        }
      }
    }

    await client
      .collection("deployments")
      .create(deployment)
      .then((response) => {
        // create initial rollout
        const rollout: RolloutsRecord = {
          manifest: selectedBlueprint.manifest,
          startDate: "",
          endDate: "",
          project: $selectedProjectId,
          deployment: response.id,
          user: client.authStore.model?.id
        };

        client
          .collection("rollouts")
          .create(rollout)
          .then((response) => {
            toast.success("Deployment & initial Rollout created");
            updateDataStores();
          })
          .catch((error) => {
            toast.success("Deployment created");
            toast.error(error.message);
          })
          .finally(() => {
            updateDataStores();
            deploymentModal = false;
          });
      })
      .catch((error) => {
        toast.error(error.message);
      })
      .finally(() => {
        name = "";
        selectedBlueprint = $blueprints[0];
      });
  }
</script>

<div class="flex flex-col space-y-6">
  <h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">Create your deployment</h3>
  <Label class="space-y-2">
    <span>Deployment name *</span>
    <Input
      type="text"
      name="deployment"
      placeholder="Enter the name of your deployment"
      required
      bind:value={name}
    />
  </Label>
  <fieldset class="space-y-2">
    <Label class="space-y-2">
      <span>Select a blueprint *</span>
    </Label>
    <div class="grid grid-cols-2 gap-2">
      {#if filteredBlueprints}
        {#each filteredBlueprints as blueprint (blueprint.id)}
          <!-- svelte-ignore a11y-click-events-have-key-events -->
          <!-- svelte-ignore a11y-no-static-element-interactions -->
          <span
            class="cursor-pointer w-full rounded-lg px-6 py-4 sm:flex sm:justify-between border-2
          {selectedBlueprint?.id === blueprint?.id
              ? 'border-primary-600 bg-gray-50 dark:bg-transparent'
              : ' border-gray-200'}
          "
            on:click={() => {
              selectedBlueprint = blueprint;
            }}
          >
            <input
              type="radio"
              name="server-size"
              value={blueprint?.id}
              class="sr-only"
              aria-labelledby="server-size-1-label"
              aria-describedby="server-size-1-description-0 server-size-1-description-1"
            />
            <span class="flex items-center">
              <span class="flex flex-col text-sm">
                <span id="server-size-1-label" class="font-medium">
                  {blueprint?.name}
                </span>

                <span id="server-size-1-description-0" class=" hover:text-gray-600 mt-1">
                  {#if blueprint?.owner === client.authStore.model?.id}
                    <BookLock class="w-4 h-4 mr-1 inline-block" />
                  {:else}
                    <BookUser class="w-4 h-4 mr-1 inline-block" />
                  {/if}
                  <p class="block sm:inline">
                    {blueprint?.description}
                  </p>
                </span>
              </span>
            </span>
            <span
              id="server-size-1-description-1"
              class="mt-2 flex text-sm sm:ml-4 sm:mt-0 sm:flex-col sm:text-right"
            >
              <img
                src={recordLogoUrl(blueprint)}
                alt={blueprint?.name}
                class="h-12 w-12 flex-none rounded-lg object-cover ring-1 ring-gray-900/10"
              />
            </span>
            <span
              class="pointer-events-none absolute -inset-px rounded-lg border-2"
              aria-hidden="true"
            ></span>
          </span>
        {/each}
      {/if}
    </div>
  </fieldset>
  <Button type="submit" class="w-full1" color="primary" on:click={handleCreateDeployment}
    >Create deployment
    <ArrowRight class="w-4 h-4 ml-2 inline-block" />
  </Button>
</div>
