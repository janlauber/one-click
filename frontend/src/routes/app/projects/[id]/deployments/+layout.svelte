<script lang="ts">
  import { page } from "$app/stores";
  import SideNav from "$lib/components/deployments/SideNav.svelte";

  export let data: any;
</script>

<!-- only if startswith /app/projects/[id]/deployments/[id] -->
{#if $page.url.pathname.startsWith("/app/projects/") && $page.url.pathname.includes("/deployments/")}
  <div class="relative h-full max-w-screen-2xl mx-auto">
    <div class="absolute top-16 left-3 right-3 bottom-0">
      <div class="absolute top-0 bottom-10 overflow-y-hidden">
        <SideNav />
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
{/if}
