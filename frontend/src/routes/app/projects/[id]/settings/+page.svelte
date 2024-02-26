<script lang="ts">
  import { goto } from "$app/navigation";
  import { client } from "$lib/pocketbase";
  import type { BlueprintsRecord, ProjectsRecord } from "$lib/pocketbase/generated-types";
  import {
    UpdateFilterEnum,
    currentRollout,
    selectedProject,
    updateDataStores
  } from "$lib/stores/data";
  import { Button, Fileupload, Heading, Input, Label, Modal, P } from "flowbite-svelte";
  import { BookDashed, Image, Trash, XIcon } from "lucide-svelte";
  import toast from "svelte-french-toast";
  import MonacoEditor from "svelte-monaco";
  // @ts-ignore
  import yaml from "js-yaml";

  let localTags: Set<string> = new Set();
  let initialLoad = true;
  let projectName: string = "";
  let inFocus = false;
  let modalBluprintOpen = false;
  let modalDeleteOpen = false;
  let avatar: File;

  let blueprintName: string = $selectedProject?.name || "";
  let blueprintDescription: string = $selectedProject?.description || "";
  let blueprintAvatar: string = $selectedProject?.avatar || "";
  let blueprintAvatarFile: File;
  let blueprintManifest: any = jsonToYaml($currentRollout?.manifest) || "";

  function jsonToYaml(json: any): string {
    return yaml.dump(json);
  }

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

    client
      .collection("projects")
      .update($selectedProject.id, project)
      .then(() => {
        toast.success("Tags updated");
      })
      .catch((error) => {
        toast.error(error.message);
      })
      .finally(() => {
        // update the selected project
        if ($currentRollout) {
          updateDataStores({
            filter: UpdateFilterEnum.ALL,
            projectId: $currentRollout.project
          });
        }
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
        modalDeleteOpen = false;
        updateDataStores({
          filter: UpdateFilterEnum.ALL
        });
        goto("/app");
      })
      .catch((error) => {
        toast.error(error.message);
      });
  }

  async function handleAvatarUpload(event: any) {
    if (!event.target.files[0]) return;
    if (!$selectedProject) return;

    avatar = event.target.files[0];

    let formData = new FormData();

    formData.append("avatar", avatar);

    client
      .collection("projects")
      .update($selectedProject.id, formData)
      .then(() => {
        toast.success("Avatar updated");
        updateDataStores({
          filter: UpdateFilterEnum.ALL
        });
      })
      .catch((error) => {
        toast.error(error.message);
      });
  }

  async function handleCreateBlueprint() {
    if (!$selectedProject) return;

    if (!blueprintName) {
      toast.error("Blueprint name is required");
      return;
    }

    if (!blueprintDescription) {
      toast.error("Blueprint description is required");
      return;
    }

    if (!blueprintAvatar) {
      toast.error("Blueprint avatar is required");
      return;
    }

    if (!blueprintManifest) {
      toast.error("Blueprint manifest is required");
      return;
    }

    let formData = new FormData();
    formData.append("avatar", blueprintAvatarFile);

    // parse the manifest yaml to json
    const parsedManifest = yaml.load(blueprintManifest);

    let data: BlueprintsRecord = {
      name: blueprintName,
      description: blueprintDescription,
      manifest: parsedManifest,
      owner: client.authStore?.model?.id
    };

    client
      .collection("blueprints")
      .create(data)
      .then((response) => {
        client
          .collection("blueprints")
          .update(response?.id ?? "", formData)
          .catch((error) => {
            toast.error(error.message);
          });

        toast.success("Blueprint created");

        goto(`/app/blueprints/my-blueprints`);
      })
      .catch((error) => {
        toast.error(error.message);
      })
      .finally(() => {
        // update the selected project
        if ($currentRollout) {
          updateDataStores({
            filter: UpdateFilterEnum.ALL,
            projectId: $currentRollout.project
          });
        }
        modalBluprintOpen = false;
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

  <div>
    <Label class="pb-2">Change Avatar</Label>
    <label
      for="avatar"
      class="cursor-pointer flex justify-center items-center w-full px-4 py-2 border-primary-700 dark:border-gray-400 border-2 rounded-lg text-sm font-medium text-primary-700 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800"
    >
      <Image class="w-6 h-6 inline-block mr-1" />
      Choose a file (max 1MB)
      <input
        type="file"
        name="avatar"
        id="avatar"
        class="w-full border-gray-300 border-2"
        on:change={handleAvatarUpload}
      />
    </label>
  </div>
</div>

<!-- Create a blueprint -->

<div class="mt-4 p-0.5 shadow ring-1 ring-black ring-opacity-5 rounded-lg bg-primary-500">
  <table class="min-w-full divide-y divide-gray-300 dark:divide-gray-600">
    <tbody class="divide-y divide-gray-200 dark:divide-gray-600 dark:bg-transparent">
      <tr class="">
        <td class="whitespace-nowrap py-4 pl-4 pr-3 text-xs font-medium sm:pl-6">
          <Heading tag="h5" color="text-white">Blueprint</Heading>
          <P color="text-white" class="text-xs">Create a blueprint from your project.</P>
        </td>
        <td class="whitespace-nowrap px-3 py-4 text-xs text-right">
          <!-- Modified: Added 'text-right' class -->
          <Button
            color="light"
            size="xs"
            class="whitespace-nowrap"
            on:click={() => {
              if (!$selectedProject) return;

              modalBluprintOpen = true;
            }}
          >
            <BookDashed class="w-4 h-4 inline-block mr-1" />
            New blueprint
          </Button>
        </td>
      </tr>
    </tbody>
  </table>
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

              modalDeleteOpen = true;
            }}
          >
            <Trash class="w-4 h-4 inline-block mr-1" />
            Delete project
          </Button>
        </td>
      </tr>
    </tbody>
  </table>
</div>

<Modal bind:open={modalDeleteOpen} size="xs" autoclose>
  <div class="text-center">
    <Trash class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      Are you sure you want to delete this project?
    </h3>
    <Button color="red" class="me-2" on:click={() => handleDelete()}>Yes, I'm sure</Button>
    <Button color="alternative">No, cancel</Button>
  </div>
</Modal>

<Modal bind:open={modalBluprintOpen} size="lg">
  <div class="text-center">
    <BookDashed class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      Create a blueprint from this project
    </h3>
  </div>
  <div class="space-y-2">
    <Label class="">Name</Label>
    <Input bind:value={blueprintName} />
    <Label class="">Description</Label>
    <Input bind:value={blueprintDescription} />
    <Label class="">Avatar</Label>
    <Fileupload
      bind:value={blueprintAvatar}
      on:change={(event) => {
        // @ts-ignore
        blueprintAvatarFile = event.target.files[0];
      }}
    />

    <Label class="">Manifest</Label>
    <span class="text-xs text-gray-500 dark:text-gray-400">
      The manifest is a YAML file that describes the blueprint. <br />
      On project creation, the ingress resources won't be created.
    </span>
    <div class="h-64 overflow-y-auto rounded-lg p-2" style="background-color: #1E1E1E;">
      <MonacoEditor
        bind:value={blueprintManifest}
        options={{ language: "yaml", automaticLayout: false, minimap: { enabled: false } }}
        theme="vs-dark"
      />
    </div>
  </div>
  <div class="flex justify-between">
    <Button color="primary" class="me-2" on:click={() => handleCreateBlueprint()}
      >Yes, I'm sure</Button
    >
    <Button
      color="alternative"
      on:click={() => {
        modalBluprintOpen = false;
      }}>No, cancel</Button
    >
  </div>
</Modal>

<style>
  input[type="file"] {
    display: none;
  }
</style>
