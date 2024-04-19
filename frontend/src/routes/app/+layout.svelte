<script lang="ts">
  import { page } from "$app/stores";
  import Nav from "$lib/components/base/Nav.svelte";
  import { Heading } from "flowbite-svelte";
  import { cubicOut } from "svelte/easing";
  import { slide } from "svelte/transition";
</script>

<div>
  <div class="">
    <Nav />
  </div>

  <div class="absolute top-12 left-0 right-0 bottom-0">
    {#if !$page.url.pathname.startsWith("/app/projects/") && !$page.url.pathname.startsWith("/app/deployments/") && !$page.url.pathname.startsWith("/app/blueprints/") && !$page.url.pathname.startsWith("/app/admin/")}
      <div
        class="bg-primary-600 w-full p-4 relative"
        in:slide={{ duration: 200, easing: cubicOut }}
        out:slide={{ duration: 200, easing: cubicOut }}
      >
        <div class="text-center flex justify-center items-center flex-col space-y-4">
          <Heading tag="h3" class="flex text-background font-extralight items-center w-auto">
            Welcome to&nbsp;<b class="font-semibold">One</b>Click
          </Heading>
          <span
            class="bg-white text-primary-500 text-sm font-extralight px-2 py-1 rounded-full mt-2"
          >
            The <b>Open Source Platform</b> to manage your <b>Container Deployments</b>
          </span>
        </div>
      </div>
    {/if}
    <slot />
  </div>
</div>
