<script lang="ts">
  import NewProject from "$lib/components/dashboard/NewProject.svelte";
  import ProjectCard from "$lib/components/dashboard/ProjectCard.svelte";
  import type { ProjectsResponse } from "$lib/pocketbase/generated-types";
  import { projects, type Pexpand, rollouts } from "$lib/stores/data";
  import { Badge, Button, Heading, Modal } from "flowbite-svelte";
  import { Plus, Tag } from "lucide-svelte";

  let projectModal = false;

  let tags: Set<string> = new Set();
  let selectedTags: Set<string> = new Set();
  let filteredProjects: ProjectsResponse<Pexpand>[] = [];

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
</script>

<div class="absolute w-full top-28 bottom-0 overflow-y-scroll scrollbar-none">
  <Modal bind:open={projectModal} size="lg" autoclose={false} class="w-full">
    <NewProject bind:projectModal />
  </Modal>

  <div class="max-w-6xl mx-auto px-5 pb-5 py-5 flex flex-col">
    <div class="flex mb-5">
      <Heading tag="h5" class="flex font-normal items-center w-auto">Your Projects ({$projects.length})</Heading>
      <Button
        color="primary"
        class="justify-self-end ml-auto"
        size="sm"
        on:click={() => {
          projectModal = true;
        }}
      >
        <Plus class="w-4 h-4 mr-2 inline-block" />
        Create Project
      </Button>
    </div>
    <!-- Filter by Tags -->
    <div class="flex flex-wrap gap-x-4 gap-y-2 mb-10">
      {#key selectedTags}
        <button
          on:click={() => {
            selectedTags = new Set();

            filteredProjects = $projects;

            console.log(selectedTags);
          }}
        >
          <Badge
            color="alternative"
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
            color={selectedTags.has(tag) ? "primary" : "alternative"}
            class="text-sm cursor-pointer
          {selectedTags.has(tag) ? 'bg-primary-600 text-white' : ''}
        "
          >
            {tag}
          </Badge>
        </button>
      {/each}
    </div>
    <div class="grid grid-cols-1 gap-x-6 gap-y-8 lg:grid-cols-3 xl:gap-x-8">
      {#each filteredProjects as project (project.id)}
        <ProjectCard {project} />
      {/each}
    </div>
  </div>
</div>
