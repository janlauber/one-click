<script lang="ts">
  import { page } from "$app/stores";
  import BlueprintSideNav from "$lib/components/blueprints/BlueprintSideNav.svelte";
  import { metadata } from "$lib/stores/metadata";
  import { cubicOut } from "svelte/easing";
  import { fly, slide } from "svelte/transition";

  $metadata.title = "Blueprints";

  export let data: any;
</script>

{#if $page.url.pathname.startsWith("/app/blueprints/") && !$page.url.pathname.startsWith("/app/blueprints/shared/")}
  <div
    class="relative h-full max-w-screen-2xl mx-auto"
    in:slide={{ duration: 200, easing: cubicOut }}
    out:slide={{ duration: 200, easing: cubicOut }}
  >
    <div class="absolute top-16 left-3 right-3 bottom-0">
      <div class="absolute top-0 left-0 bottom-10 overflow-y-hidden">
        <BlueprintSideNav />
      </div>
      {#key data.url}
        <div
          class="absolute top-0 py-3 left-64 right-0 bottom-0 overflow-y-auto scrollbar-none pb-8 px-2 overflow-x-hidden"
          in:fly={{ duration: 100, easing: cubicOut }}
          out:fly={{ duration: 100, easing: cubicOut }}
        >
          <slot />
        </div>
      {/key}
    </div>
  </div>
{:else}
  <slot />
{/if}
