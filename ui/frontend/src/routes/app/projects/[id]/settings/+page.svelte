<script lang="ts">
    import { client } from "$lib/pocketbase";
    import type { ProjectsRecord } from "$lib/pocketbase/generated-types";
  import { selectedProject } from "$lib/stores/data";
  import { Button, Heading, Input, Label, P } from "flowbite-svelte";
  import { XIcon } from "lucide-svelte";
  import toast from "svelte-french-toast";

  let localTags: Set<string> = new Set();
  let initialLoad = true;

  $: {
    // update tags
    let tempTags = $selectedProject?.tags;
    // split tags by comma
    if (initialLoad && tempTags) {
      // split by comma or if there is no comma, then its a single tag
      localTags = new Set(tempTags.split(",").map((tag) => formatTag(tag)));
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

  $: localTags && safeTags();

  async function safeTags() {

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
      });
  }


</script>

<div class="flex items-start justify-between">
  <div class="flex flex-col">
    <Heading tag="h2">Settings</Heading>
    <P class="text-gray-500 dark:text-gray-400 text-sm">Settings of your project.</P>
  </div>
</div>

<div class="mt-8">
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
