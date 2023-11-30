<script>
  import NewProject from "$lib/components/dashboard/NewProject.svelte";
  import ProjectCard from "$lib/components/dashboard/ProjectCard.svelte";
  import { projects } from "$lib/stores/data";
  import { Badge, Button, Heading, Modal } from "flowbite-svelte";
  import { Plus } from "lucide-svelte";

  let projectModal = false;
</script>

<div class="bg-primary-600 absolute w-full p-10 pb-20 h-52">
  <div class="text-center flex justify-center items-center flex-col space-y-4">
    <Heading tag="h3" class="flex text-background font-extralight items-center w-auto">
      Welcome to&nbsp;<b class="font-semibold">One</b>Click <Badge
        color="primary"
        class="text-xl font-semibold ml-2">prototype</Badge
      >
    </Heading>
    <span class="bg-black text-background text-sm font-extralight px-2 py-1 rounded-full mt-2">
      The <b>Open Source Platform</b> to manage your <b>Software Rollouts</b>
    </span>
  </div>
  <!-- <div class="absolute left-1/2 -translate-x-1/2 -bottom-16">
  <Stats />
  </div> -->
</div>

<div class="absolute w-full top-52 bottom-0 overflow-y-scroll scrollbar-none">
  <Modal bind:open={projectModal} size="xs" autoclose={false} class="w-full">
    <NewProject bind:projectModal />
  </Modal>

  <div class="max-w-6xl mx-auto px-5 pb-5 py-10 flex flex-col">
    <div class="flex mb-10">
      <Heading tag="h5" class="flex font-normal items-center w-auto"
        >Your Projects ({$projects.length})</Heading
      >
      <Button
        color="alternative"
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
