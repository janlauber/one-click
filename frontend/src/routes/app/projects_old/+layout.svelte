<script lang="ts">
  import { page } from "$app/stores";
  import SideNav from "$lib/components/deployments/SideNav.svelte";
  import { cubicOut } from "svelte/easing";
  import { fly, slide } from "svelte/transition";

  export let data: any;
</script>

{#if $page.url.pathname.startsWith("/app/projects/")}
  <div class="relative h-full max-w-6xl mx-auto">
    <div class="absolute top-10 left-2 bottom-10 overflow-y-hidden">
      <SideNav />
    </div>
    {#key data.url}
      <div
        class="absolute top-4 pt-4 left-64 right-2 bottom-0 overflow-y-auto scrollbar-none pb-8 px-2 overflow-x-hidden"
      >
        <slot />
      </div>
    {/key}
  </div>
{/if}
