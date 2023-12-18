<script lang="ts">
  import { goto } from "$app/navigation";
  import { client } from "$lib/pocketbase";
  import type { ProjectsRecord } from "$lib/pocketbase/generated-types";
  import {
    UpdateFilterEnum,
    currentRollout,
    selectedProject,
    updateDataStores
  } from "$lib/stores/data";
  import { Button, Heading, Input, Label, Modal, P } from "flowbite-svelte";
  import { ExclamationCircleOutline } from "flowbite-svelte-icons";
  import { XIcon } from "lucide-svelte";
  import toast from "svelte-french-toast";

  let localTags: Set<string> = new Set();
  let initialLoad = true;
  let projectName: string = "";
  let inFocus = false;
  let modalOpen = false;

  $: {
    // update tags
    let tempTags = $selectedProject?.tags;
    // split tags by comma
    if (initialLoad) {
      // split by comma or if there is no comma, then its a single tag
      if (tempTags?.includes(",")) {
        localTags = new Set(tempTags.split(",").map((tag) => formatTag(tag)));
      } else {
        if (tempTags) localTags = new Set([formatTag(tempTags)]);
      }

      // set project name
      if ($selectedProject) projectName = $selectedProject.name;

      initialLoad = false;
    }
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
    }
  }

  function handleAddTag(event: any) {
    // if value of #tags is empty, do nothing
    if (!(event.target.previousElementSibling as HTMLInputElement).value) return;

    localTags = localTags.add(
      formatTag((event.target.previousElementSibling as HTMLInputElement).value)
    );
    (event.target.previousElementSibling as HTMLInputElement).value = "";
  }

  function setToString(set: Set<string>): string {
    return Array.from(set).join(",");
  }

  $: localTags && saveTags();

  $: (inFocus || !inFocus) && saveName();

  async function saveTags() {
    // check if the tags changed
    if (setToString(localTags) === $selectedProject?.tags) return;

    if (!$selectedProject) return;

    const project: ProjectsRecord = {
      ...$selectedProject,
      tags: setToString(localTags)
    };

    console.log(project);

    client
      .collection("projects")
      .update($selectedProject.id, project)
      .then(() => {
        // update the selected project
        if ($currentRollout) {
          updateDataStores({
            filter: UpdateFilterEnum.ALL,
            projectId: $currentRollout.project
          });
        }
        toast.success("Tags updated");
      })
      .catch((error) => {
        toast.error(error.message);
      });
  }

  async function saveName() {
    // check if the name changed and input is not in focus
    if (projectName === $selectedProject?.name || inFocus) return;

    if (!$selectedProject) return;

    const project: ProjectsRecord = {
      ...$selectedProject,
      name: projectName
    };

    client
      .collection("projects")
      .update($selectedProject.id, project)
      .then(() => {
        // update the selected project
        if ($currentRollout) {
          updateDataStores({
            filter: UpdateFilterEnum.ALL,
            projectId: $currentRollout.project
          });
        }
        toast.success("Name updated");
      })
      .catch((error) => {
        toast.error(error.message);
      });
  }

  async function handleDelete() {
    if (!$selectedProject) return;

    client
      .collection("projects")
      .delete($selectedProject.id)
      .then(() => {
        toast.success("Project deleted");
        modalOpen = false;
        updateDataStores({
          filter: UpdateFilterEnum.ALL
        });
        goto("/app");
      })
      .catch((error) => {
        toast.error(error.message);
      });
  }
</script>

<div class="flex items-start justify-between">
  <div class="flex flex-col">
    <Heading tag="h2">Settings</Heading>
    <P class="text-gray-500 dark:text-gray-400 text-sm">Settings of your project.</P>
  </div>
</div>

<div class="mt-8 space-y-4">
  <Label class="space-y-2">
    <span>Project name</span>
    <div class="flex whitespace-nowrap gap-2">
      <Input
        id="name"
        type="text"
        name="name"
        size="sm"
        placeholder="Enter the name of your project"
        bind:value={projectName}
        on:focus={() => {
          inFocus = true;
        }}
        on:blur={() => {
          inFocus = false;
        }}
      />
    </div>
  </Label>

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
      <div class="flex flex-wrap mt-2">
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
</div>

<!-- Danger Zone -> Delete Project -->

<div
  class="mt-4 p-0.5 shadow ring-1 ring-black ring-opacity-5 rounded-lg bg-red-100 dark:bg-red-800"
>
  <table class="min-w-full divide-y divide-gray-300 dark:divide-gray-600">
    <tbody class="divide-y divide-gray-200 dark:divide-gray-600 dark:bg-transparent">
      <tr class="">
        <td class="whitespace-nowrap py-4 pl-4 pr-3 text-xs font-medium sm:pl-6">
          <Heading tag="h5" color="text-red-600 dark:text-red-100">Danger Zone</Heading>
          <P color="text-red-600 dark:text-red-100" class="text-xs">Delete your project.</P>
        </td>
        <td class="whitespace-nowrap px-3 py-4 text-xs text-right">
          <!-- Modified: Added 'text-right' class -->
          <Button
            color="red"
            size="xs"
            class="whitespace-nowrap"
            on:click={() => {
              if (!$selectedProject) return;

              modalOpen = true;
            }}
          >
            Delete project
          </Button>
        </td>
      </tr>
    </tbody>
  </table>
</div>

<Modal bind:open={modalOpen} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleOutline class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      Are you sure you want to delete this project?
    </h3>
    <Button color="red" class="me-2" on:click={() => handleDelete()}>Yes, I'm sure</Button>
    <Button color="alternative">No, cancel</Button>
  </div>
</Modal>
