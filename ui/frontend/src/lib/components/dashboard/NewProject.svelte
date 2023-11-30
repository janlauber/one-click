<script lang="ts">
  import { client } from "$lib/pocketbase";
  import type { ProjectsRecord, FrameworksResponse } from "$lib/pocketbase/generated-types";
  import { frameworks, updateDataStores } from "$lib/stores/data";
  import { frameworkLogoUrl } from "$lib/utils/framework.utils";

  import { Button, Input, Label } from "flowbite-svelte";
  import { ArrowRight, ExternalLink, XIcon } from "lucide-svelte";
  import toast from "svelte-french-toast";

  export let projectModal: boolean;

  let name: string = "";
  let selectedFramework: FrameworksResponse = $frameworks[0];
  let localTags: Set<string> = new Set();

  function handleAddTag(event: KeyboardEvent) {
    if (event.key === "Enter") {
      localTags = localTags.add((event.target as HTMLInputElement).value);
      (event.target as HTMLInputElement).value = "";
    }
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

    const project: ProjectsRecord = {
      name: name,
      framework: selectedFramework.id,
      user: client.authStore.model?.id,
      tags: setToString(localTags),
    };

    await client
      .collection("projects")
      .create(project)
      .then((response) => {
        toast.success("Project created");
        localTags = new Set();
        updateDataStores();
        projectModal = false;
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
      <span>Select a framework</span>
    </Label>

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
              <span id="server-size-1-label" class="font-medium">{framework?.name}</span>
              <span id="server-size-1-description-0" class=" hover:text-gray-600">
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
  </fieldset>
  <Label class="space-y-2">
    <span>Tags</span>
    <Input
      type="text"
      name="tags"
      placeholder="Enter the tags of your project"
      on:keydown={handleAddTag}
    />
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
              class="absolute top-0 left-0 right-0 bottom-0 rounded-lg opacity-0 -z-10 group-hover:z-10 group-hover:opacity-100 bg-red-500 transition-opacity text-white"
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
  <Button type="submit" class="w-full1" color="alternative" on:click={handleCreateProject}
    >Create project
    <ArrowRight class="w-4 h-4 ml-2 inline-block" />
  </Button>
</div>
