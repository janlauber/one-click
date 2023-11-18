<script lang="ts">
  import { page } from "$app/stores";

  import selectedProjectId from "$lib/stores/project";
  import { Cog, LineChart, Rocket } from "lucide-svelte";

  let items = [
    {
      name: "Overview",
      href: `/app/projects/${$selectedProjectId}/overview`,
      current: false,
      icon: LineChart
    },
    {
      name: "Rollouts",
      href: `/app/projects/${$selectedProjectId}/rollouts`,
      current: false,
      icon: Rocket
    },
    {
      name: "Settings",
      href: `/app/projects/${$selectedProjectId}/settings`,
      current: false,
      icon: Cog
    }
  ];

  $: items = items.map((item) => {
    if ($page.url.pathname.startsWith(item.href)) {
      item.current = true;
    } else {
      item.current = false;
    }
    return item;
  });
</script>

<div class="flex flex-col gap-y-4 " role="group" aria-labelledby="projects-headline">
  {#each items as item}
    <a
      href={item.href}
      class=" text-gray-900 hover:text-gray-700 dark:text-gray-100 dark:hover:text-gray-300 pl-4 pr-8 py-2 text-sm font-medium rounded-md transition-all duration-150 ease-in-out hover:bg-gray-100 dark:hover:bg-primary-600 dark:hover:bg-opacity-10
        {item.current ? 'bg-gray-200 dark:bg-primary-600' : ''}
      "
      aria-current={item.current ? "page" : undefined}
    >
      <svelte:component this={item.icon} class="w-5 h-5 mr-2 inline" strokeWidth={1.5} />
      {item.name}
    </a>
  {/each}
</div>
