<script lang="ts">
  import { page } from "$app/stores";
  import ProjectSideNav from "$lib/components/projects/ProjectSideNav.svelte";
  import selectedProjectId from "$lib/stores/project";

  export let data: any;

  const projectPathRegex = new RegExp(`/app/projects/${$selectedProjectId}(/)?$`);
</script>

{#if projectPathRegex.test($page.url.pathname)}
  <div class="relative h-full max-w-6xl mx-auto">
    <div class="absolute top-10 left-2 bottom-10 overflow-y-hidden">
      <ProjectSideNav />
    </div>
    {#key data.url}
      <div
        class="absolute top-4 pt-4 left-64 right-2 bottom-0 overflow-y-auto scrollbar-none pb-8 px-2 overflow-x-hidden"
      >
        <slot />
      </div>
    {/key}
  </div>
{:else}
  <slot />
{/if}
