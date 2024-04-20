<script lang="ts">
  import { client } from "$lib/pocketbase";
  import type { ProjectsRecord } from "$lib/pocketbase/generated-types";
  import { UpdateFilterEnum, updateDataStores } from "$lib/stores/data";

  import { Button, Fileupload, Input, Label, Textarea } from "flowbite-svelte";
  import { ArrowRight, XIcon } from "lucide-svelte";
  import toast from "svelte-french-toast";

  export let projectModal: boolean;

  let name: string = "";
  let description: string = "";
  let avatar: string = "";
  let avatarFile: File;

  let localTags: Set<string> = new Set();

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

    let formData = new FormData();
    formData.append("avatar", avatarFile);

    const project: ProjectsRecord = {
      name: name,
      description: description,
      user: client.authStore.model?.id,
      tags: setToString(localTags)
    };

    await client
      .collection("projects")
      .create(project)
      .then((response) => {
        if (avatarFile) {
          client
            .collection("projects")
            .update(response.id, formData)
            .then(() => {
              updateDataStores({
                filter: UpdateFilterEnum.ALL
              });
            })
            .catch((error) => {
              toast.error(error.message);
            });
        }
        toast.success("Project created");
        projectModal = false;
        localTags = new Set();
        name = "";
        description = "";
        updateDataStores();
      })
      .catch((error) => {
        toast.error(error.message);
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

  <Label class="space-y-2">
    <span>Description</span>
    <Textarea
      type="text"
      name="description"
      placeholder="Enter the description of your project"
      required
      bind:value={description}
    />
  </Label>

  <Label class="space-y-2">
    <span>Avatar*</span>
    <Fileupload
      bind:value={avatar}
      on:change={(event) => {
        // @ts-expect-error - event.target.files is a FileList
        avatarFile = event.target.files[0];
      }}
    />
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
