<script lang="ts">
  import { page } from "$app/stores";
  import ProjectSideNav from "$lib/components/projects/ProjectSideNav.svelte";
  import selectedProjectId from "$lib/stores/project";
  import { Modal } from "flowbite-svelte";
  import NewDeployment from "$lib/components/deployments/NewDeployment.svelte";

  export let data: any;
  let modal = false;

  const projectPathRegex = new RegExp(`/app/projects/${$selectedProjectId}(/)?(?!.*deployments).*`);
</script>

{#if projectPathRegex.test($page.url.pathname)}
  <div class="relative h-full max-w-screen-2xl mx-auto top-16">
    <div class="absolute top-0 left-3 right-3 bottom-0">
      <Modal bind:open={modal} size="lg" autoclose={false} class="w-full">
        <NewDeployment bind:modal />
      </Modal>
      <div class="absolute top-0 left-0 bottom-10 overflow-y-hidden">
        <ProjectSideNav bind:modal />
      </div>
      {#key data.url}
        <div
          class="absolute top-0 left-64 right-0 bottom-0 overflow-y-auto scrollbar-none pb-8 overflow-x-hidden"
        >
          <slot />
        </div>
      {/key}
    </div>
  </div>
{:else}
  <slot />
{/if}
