<script>
  import NewProject from "$lib/components/dashboard/NewProject.svelte";
  import ProjectCard from "$lib/components/dashboard/ProjectCard.svelte";
  import { projects } from "$lib/stores/data";
  import { Badge, Button, Heading, Modal } from "flowbite-svelte";
  import { Plus } from "lucide-svelte";

  let projectModal = false;
</script>



<div class="absolute w-full top-52 bottom-0 overflow-y-scroll scrollbar-none">
  <Modal bind:open={projectModal} size="md" autoclose={false} class="w-full">
    <NewProject bind:projectModal />
  </Modal>

  <div class="max-w-6xl mx-auto px-5 pb-5 py-10 flex flex-col">
    <div class="flex mb-10">
      <Heading tag="h5" class="flex font-normal items-center w-auto"
        >Your Projects</Heading
      >
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
    <div class="grid grid-cols-1 gap-x-6 gap-y-8 lg:grid-cols-3 xl:gap-x-8">
      {#each $projects as project (project.id)}
        <ProjectCard {project} />
      {/each}
    </div>
  </div>
</div>
