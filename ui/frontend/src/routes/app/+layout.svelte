<script>
  import { page } from "$app/stores";
  import Nav from "$lib/components/base/Nav.svelte";
  import Stats from "$lib/components/dashboard/Stats.svelte";
  import { Badge, Heading } from "flowbite-svelte";
  import { cubicOut } from "svelte/easing";
  import { fly, slide } from "svelte/transition";
</script>

<div>
  <div class="">
    <Nav />
  </div>

  <div class="absolute top-14 left-0 right-0 bottom-0">
    {#if !$page.url.pathname.startsWith("/app/projects/")}
      <div
        class="bg-primary-600 w-full p-4 relative"
        in:slide={{ duration: 200, easing: cubicOut }}
        out:slide={{ duration: 200, easing: cubicOut }}
      >
        <div class="text-center flex justify-center items-center flex-col space-y-4">
          <Heading tag="h3" class="flex text-background font-extralight items-center w-auto">
            Welcome to&nbsp;<b class="font-semibold">One</b>Click <Badge
              color="primary"
              class="text-xl font-semibold ml-2">prototype</Badge
            >
          </Heading>
          <span
            class="bg-black text-background text-sm font-extralight px-2 py-1 rounded-full mt-2"
          >
            The <b>Open Source Platform</b> to manage your <b>Software Rollouts</b>
          </span>
        </div>
        <!-- <div
          class="absolute left-1/2 -translate-x-1/2 -bottom-12"
          in:fly={{ duration: 200, delay: 200, y: -100, easing: cubicOut }}
          out:fly={{ duration: 0, opacity: 0 }}
        >
          <Stats />
        </div> -->
      </div>
    {/if}
    <slot />
  </div>
</div>
