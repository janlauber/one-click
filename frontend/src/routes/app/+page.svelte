<script lang="ts">
  import { goto } from "$app/navigation";
  import NewProject from "$lib/components/projects/NewProject.svelte";
  import ProjectCard from "$lib/components/projects/ProjectCard.svelte";
  import { client } from "$lib/pocketbase";
  import type { ProjectsResponse } from "$lib/pocketbase/generated-types";
  import { projects, blueprints, loading } from "$lib/stores/data";
  import { getTagColor } from "$lib/utils/tags";
  import { Badge, Button, Heading, Modal, Spinner } from "flowbite-svelte";
  import { BookDashed, Tag, FolderPlus } from "lucide-svelte";

  let projectModal = false;

  let tags: Set<string> = new Set();
  let selectedTags: Set<string> = new Set();
  let filteredProjects: ProjectsResponse[] = [];

  $: {
    if ($projects) {
      $projects.forEach((p) => {
        if (p.tags) {
          p.tags.split(",").forEach((t) => {
            tags.add(t);
          });
        }
      });
    }
  }

  $: if (selectedTags.size === 0) {
    filteredProjects = $projects;
  }

  function handleFilterProjects(tag: string) {
    if (selectedTags.has(tag)) {
      selectedTags.delete(tag);
    } else {
      selectedTags.add(tag);
    }

    if (selectedTags.size === 0) {
      filteredProjects = $projects;
    } else {
      filteredProjects = $projects.filter((p) => {
        if (p.tags) {
          const projectTags = p.tags.split(",");
          return projectTags.some((t) => selectedTags.has(t));
        }
        return false;
      });
    }

    selectedTags = new Set(selectedTags);
  }

  function getOwnedBlueprints() {
    return $blueprints.filter(
      (blueprint) => blueprint.owner === (client.authStore?.model?.id ?? null)
    );
  }

  function getCommunityBlueprints() {
    return $blueprints.filter(
      (blueprint) =>
        blueprint.owner !== client.authStore?.model?.id &&
        blueprint.users.some((user) => user === client.authStore?.model?.id)
    );
  }

  // getOwnedBlueprints() + getCommunityBlueprints() = filtered_blueprints
  let filtered_blueprints = getOwnedBlueprints().concat(getCommunityBlueprints());
</script>

<div class="absolute w-full top-44 bottom-0 overflow-y-scroll scrollbar-none p-3">
  <Modal bind:open={projectModal} size="lg" autoclose={false} class="w-full">
    <NewProject bind:projectModal />
  </Modal>

  <div class="max-w-screen-2xl mx-auto pb-5 flex flex-col">
    <div class="flex mb-5">
      <Heading tag="h5" class="flex font-normal items-center w-auto"
        >Your Projects ({$projects.length})</Heading
      >
      <div class="justify-self-end ml-auto space-x-3">
        <Button
          color="primary"
          outline
          size="sm"
          class="dark:text-white dark:border-white"
          on:click={() => {
            goto("/app/blueprints/my-blueprints");
          }}
        >
          <BookDashed class="w-4 h-4 mr-2 inline-block" />
          Blueprints ({filtered_blueprints.length})
        </Button>
        <Button
          color="primary"
          size="sm"
          on:click={() => {
            projectModal = true;
          }}
        >
          <FolderPlus class="w-4 h-4 mr-2 inline-block" />
          New Project
        </Button>
      </div>
    </div>
    <!-- Filter by Tags -->
    <div class="flex flex-wrap gap-x-4 gap-y-2 mb-10">
      {#key selectedTags}
        <button
          on:click={() => {
            selectedTags = new Set();

            filteredProjects = $projects;
          }}
        >
          <Badge
            color="primary"
            class="text-sm cursor-pointer
        {selectedTags.size === 0 ? 'bg-primary-600 text-white' : ''}
      "
          >
            All
          </Badge>
        </button>
      {/key}
      {#each Array.from(tags) as tag (tag)}
        <button on:click={() => handleFilterProjects(tag)}>
          <Badge
            color={getTagColor(tag)}
            class="text-sm cursor-pointer
          {selectedTags.has(tag) ? 'bg-primary-600 text-white' : ''}
        "
          >
            <Tag class="w-4 h-4 inline-block" strokeWidth={2} />&nbsp;{tag}
          </Badge>
        </button>
      {/each}
    </div>
    <div class="grid grid-cols-1 gap-x-6 gap-y-8 lg:grid-cols-3 xl:gap-x-8">
      {#if $loading}
        <div
          class="absolute top-0 left-0 right-0 bottom-0 flex justify-center items-center bg-gray-50 dark:bg-slate-800 z-20"
        >
          <span class="">
            <Spinner />
          </span>
        </div>
      {:else}
        {#each filteredProjects as project (project.id)}
          <ProjectCard {project} />
        {/each}
        {#if filteredProjects.length === 0}
          <div class="flex justify-center items-center w-full col-span-3">No projects found.</div>
        {/if}
      {/if}
    </div>
  </div>
</div>
