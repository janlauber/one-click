<script lang="ts">
  import { navigating, page } from "$app/stores";
  import selectedProjectId from "$lib/stores/project";
  import { Box, Cog, Expand, LineChart, Network, Rocket, Variable } from "lucide-svelte";

  // Function to generate items array
  const generateItems = (projectId: string) => [
    {
      name: "Overview",
      href: `/app/projects/${projectId}/overview`,
      current: false,
      icon: LineChart
    },
    {
      name: "Rollouts",
      href: `/app/projects/${projectId}/rollouts`,
      current: false,
      icon: Rocket
    },
    {
      name: "Scale",
      href: `/app/projects/${projectId}/scale`,
      current: false,
      icon: Expand
    },
    {
      name: "Network",
      href: `/app/projects/${projectId}/network`,
      current: false,
      icon: Network
    },
    {
      name: "Instances",
      href: `/app/projects/${projectId}/instances`,
      current: false,
      icon: Box
    },
    {
      name: "Envs & Secrets",
      href: `/app/projects/${projectId}/envs`,
      current: false,
      icon: Variable
    },
    {
      name: "Settings",
      href: `/app/projects/${projectId}/settings`,
      current: false,
      icon: Cog
    }
  ];

  let items = generateItems($selectedProjectId);

  function setCurrentItem() {
    items = items.map((item) => {
      if ($page.url.pathname.startsWith(item.href)) {
        item.current = true;
      } else {
        item.current = false;
      }
      return item;
    });
  }

  $: items = generateItems($selectedProjectId); // Regenerate items on projectId change
  $: setCurrentItem(); // Call setCurrentItem whenever items are updated

  $: if ($page) {
    setCurrentItem();
  }
</script>

<div class="flex flex-col gap-y-4" role="group" aria-labelledby="projects-headline">
  {#each items as item}
    <a
      href={item.href}
      class=" text-gray-900 hover:text-gray-700 dark:text-gray-100 dark:hover:text-gray-300 pl-4 pr-8 py-2 text-sm font-medium rounded-md transition-all duration-150 ease-in-out hover:bg-gray-100 dark:hover:bg-primary-600 dark:hover:bg-opacity-10
        {item.current ? 'bg-gray-200 dark:bg-primary-600' : ''}
      "
      aria-current={item.current ? "page" : undefined}
    >
      <svelte:component this={item.icon} class="w-5 h-5 mr-2 inline" strokeWidth={2} />
      {item.name}
    </a>
  {/each}
</div>
