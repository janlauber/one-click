<script lang="ts">
  import { goto } from "$app/navigation";
  import { client } from "$lib/pocketbase";
  import type {
    ProjectsRecord,
    FrameworksResponse,
    PlansResponse,
    RolloutsRecord
  } from "$lib/pocketbase/generated-types";
  import { frameworks, plans, updateDataStores } from "$lib/stores/data";
  import { frameworkLogoUrl } from "$lib/utils/framework.utils";

  import { Badge, Button, Input, Label } from "flowbite-svelte";
  import { ArrowRight, DollarSign, ExternalLink, XIcon } from "lucide-svelte";
  import toast from "svelte-french-toast";

  export let projectModal: boolean;

  let name: string = "";
  let selectedFramework: FrameworksResponse = $frameworks[0];
  // filter selected plan by selectedFramework
  let selectedPlan: PlansResponse = $plans.filter(
    (plan) => plan.framework === selectedFramework.id
  )[0];
  let tempPlans: PlansResponse[] = $plans.filter((plan) => plan.framework === selectedFramework.id);
  let localTags: Set<string> = new Set();

  $: tempPlans = filterAndSortPlans(selectedFramework.id);
  $: selectedPlan = tempPlans[0];

  function filterAndSortPlans(frameworkId: string): PlansResponse[] {
    // filter plans by frameworkId and sort by price
    return $plans
      .filter((plan) => plan.framework === frameworkId)
      .sort((a, b) => a.price - b.price);
  }

  function formatTag(tag: string): string {
    // remove whitespace and add - between words
    // to lowercase

    return tag.trim().split(" ").join("-").toLowerCase();
  }

  function handleAddTagEnter(event: KeyboardEvent) {
    // if value is empty, do nothing
    if (!(event.target as HTMLInputElement).value) return;

    if (event.key === "Enter") {
      localTags = localTags.add(formatTag((event.target as HTMLInputElement).value));
      (event.target as HTMLInputElement).value = "";
      toast.success("Tag added");
    }
  }

  function handleAddTag(event: any) {
    // if value of #tags is empty, do nothing
    if (!(event.target.previousElementSibling as HTMLInputElement).value) return;

    localTags = localTags.add(
      formatTag((event.target.previousElementSibling as HTMLInputElement).value)
    );
    (event.target.previousElementSibling as HTMLInputElement).value = "";
    toast.success("Tag added");
  }

  function setToString(set: Set<string>): string {
    return Array.from(set).join(",");
  }

  async function handleCreateProject(event: Event) {
    event.preventDefault();

    if (!name) {
      toast.error("Please enter a name");
      return;
    }

    if (!selectedFramework) {
      toast.error("Please select a framework");
      return;
    }

    if (!selectedPlan) {
      toast.error("Please select a plan");
      return;
    }

    const project: ProjectsRecord = {
      name: name,
      framework: selectedFramework.id,
      user: client.authStore.model?.id,
      selectedPlan: selectedPlan.id,
      tags: setToString(localTags)
    };

    await client
      .collection("projects")
      .create(project)
      .then((response) => {
        // create initial rollout
        const rollout: RolloutsRecord = {
          manifest: selectedPlan.manifest,
          startDate: "",
          endDate: "",
          project: response.id,
          user: client.authStore.model?.id
        };

        client
          .collection("rollouts")
          .create(rollout)
          .then((response) => {
            toast.success("Project & initial Rollout created");
            updateDataStores();
          })
          .catch((error) => {
            toast.success("Project created");
            toast.error(error.message);
          })
          .finally(() => {
            localTags = new Set();
            updateDataStores();
            projectModal = false;
          });
      })
      .catch((error) => {
        toast.error(error.message);
      })
      .finally(() => {
        name = "";
        selectedFramework = $frameworks[0];
      });
  }
</script>

<div class="flex flex-col space-y-6">
  <h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">Create your project</h3>
  <Label class="space-y-2">
    <span>Project name *</span>
    <Input
      type="text"
      name="project"
      placeholder="Enter the name of your project"
      required
      bind:value={name}
    />
  </Label>
  <fieldset class="space-y-2">
    <Label class="space-y-2">
      <span>Select a framework *</span>
    </Label>
    <div class="grid grid-cols-2 gap-2">
      {#if $frameworks}
        {#each $frameworks as framework (framework.id)}
          <!-- svelte-ignore a11y-click-events-have-key-events -->
          <!-- svelte-ignore a11y-no-static-element-interactions -->
          <span
            class="cursor-pointer w-full rounded-lg px-6 py-4 sm:flex sm:justify-between border-2
          {selectedFramework?.id === framework?.id
              ? 'border-primary-600 bg-blue-50 dark:bg-transparent'
              : ' border-gray-200'}
          "
            on:click={() => {
              selectedFramework = framework;
            }}
          >
            <input
              type="radio"
              name="server-size"
              value={framework?.id}
              class="sr-only"
              aria-labelledby="server-size-1-label"
              aria-describedby="server-size-1-description-0 server-size-1-description-1"
            />
            <span class="flex items-center">
              <span class="flex flex-col text-sm">
                <span id="server-size-1-label" class="font-medium">
                  {framework?.name}
                </span>

                <span id="server-size-1-description-0" class=" hover:text-gray-600 mt-1">
                  <ExternalLink class="w-4 h-4 mr-1 inline-block" />
                  <a href={framework.url} target="_blank" class="block sm:inline underline">
                    {framework.url}</a
                  >
                </span>
              </span>
            </span>
            <span
              id="server-size-1-description-1"
              class="mt-2 flex text-sm sm:ml-4 sm:mt-0 sm:flex-col sm:text-right"
            >
              <img
                src={frameworkLogoUrl(framework)}
                alt={framework?.name}
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

  <!-- Plans -->

  <fieldset class="space-y-2">
    <Label class="space-y-2">
      <span>Select a plan *</span>
    </Label>
    <div class="grid grid-cols-2 gap-2">
      {#if tempPlans}
        {#each tempPlans as plan (plan.id)}
          <!-- svelte-ignore a11y-click-events-have-key-events -->
          <!-- svelte-ignore a11y-no-static-element-interactions -->
          <div
            class="cursor-pointer w-full rounded-lg px-6 py-4 sm:flex sm:justify-between border-2
      {selectedPlan?.id === plan?.id
              ? 'border-primary-600 bg-blue-50 dark:bg-transparent'
              : ' border-gray-200'}
      "
            on:click={() => {
              selectedPlan = plan;
            }}
          >
            <input
              type="radio"
              name="server-size"
              value={plan?.id}
              class="sr-only"
              aria-labelledby="server-size-1-label"
              aria-describedby="server-size-1-description-0 server-size-1-description-1"
            />
            <span class="flex items-center">
              <span class="flex flex-col text-xs space-y-2">
                <div class="flex flex-col">
                  <span id="server-size-1-label" class="font-medium text-sm">{plan?.name}</span>
                  <!-- Description -->
                  <span id="server-size-1-description-0">
                    {plan?.description}
                  </span>
                </div>
                <!-- <div class="flex flex-col">
                  <span class="font-bold"> Resources: </span>
                  <span>
                    CPU: {plan?.manifest?.spec.resources.limits.cpu} | RAM: {plan.manifest?.spec
                      .resources.limits.memory}
                  </span>
                  {#if plan?.manifest?.spec.volumes}
                    {#if plan?.manifest?.spec.volumes.length > 0}
                      <span class=" font-bold"> Volumes: </span>
                    {/if}
                    {#each plan?.manifest?.spec.volumes as volume (volume.name)}
                      <li class="ml-3">
                        {volume.mountPath}
                        {volume.size}
                      </li>
                    {/each}
                  {/if}
                </div> -->

                <span id="server-size-1-description-0" class="text-primary-600">
                  <span class="font-bold text-sm"
                    ><DollarSign class="w-5 h-5 inline-block" />{Math.round(plan?.price * 100) /
                      100 ===
                    Math.round(plan?.price)
                      ? plan?.price + ".00"
                      : Math.round(plan?.price * 100) / 100}
                  </span>/month
                </span>
              </span>
            </span>
            <span
              class="pointer-events-none absolute -inset-px rounded-lg border-2"
              aria-hidden="true"
            ></span>
          </div>
        {/each}
      {/if}
    </div>
  </fieldset>

  <Label class="space-y-2">
    <span>Tags</span>
    <div class="flex whitespace-nowrap gap-2">
      <Input
        id="tags"
        type="text"
        name="tags"
        size="sm"
        placeholder="Enter the tags of your project"
        on:keydown={handleAddTagEnter}
      />
      <Button color="alternative" size="xs" class="inline" on:click={handleAddTag}>Add tag</Button>
    </div>
  </Label>
  {#key localTags}
    {#if localTags.size > 0}
      <div class="flex flex-wrap">
        {#each [...localTags] as tag (tag)}
          <div
            class="group relative bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400 rounded-lg px-3 py-1 text-sm font-semibold mr-2 mb-2"
          >
            {tag}
            <button
              type="button"
              class="absolute top-0 left-0 right-0 bottom-0 w-full rounded-lg opacity-0 -z-10 group-hover:z-10 group-hover:opacity-100 bg-red-500 transition-opacity text-white"
              on:click={() => {
                localTags.delete(tag);
                localTags = localTags;
              }}
            >
              <XIcon
                class="w-4 h-4 inline-block ml-20 group-hover:ml-0 duration-150 transition-all ease-in-out"
              />
            </button>
          </div>
        {/each}
      </div>
    {/if}
  {/key}
  <Button type="submit" class="w-full1" color="primary" on:click={handleCreateProject}
    >Create project
    <ArrowRight class="w-4 h-4 ml-2 inline-block" />
  </Button>
</div>
