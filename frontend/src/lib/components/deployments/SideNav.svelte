<script lang="ts">
  import { page } from "$app/stores";
  import selectedProjectId from "$lib/stores/project";
  import selectedDeploymentId from "$lib/stores/deployment";
  import {
    ArrowLeft,
    Boxes,
    Cog,
    Database,
    Expand,
    HardDrive,
    History,
    LineChart,
    Map,
    Network,
    Variable
  } from "lucide-svelte";

  // Return navigation items based on project settings
  const generateItems = (projectId: string, deploymentId: string) => {
    const items = [
      {
        name: "Overview",
        href: `/app/projects/${projectId}/deployments/${deploymentId}/overview`,
        current: false,
        icon: LineChart
      },
      {
        name: "Map",
        href: `/app/projects/${projectId}/deployments/${deploymentId}/map`,
        current: false,
        icon: Map
      },
      {
        name: "Rollouts",
        href: `/app/projects/${projectId}/deployments/${deploymentId}/rollouts`,
        current: false,
        icon: History
      },
      {
        name: "Image",
        href: `/app/projects/${projectId}/deployments/${deploymentId}/image`,
        current: false,
        icon: HardDrive
      },
      {
        name: "Scale",
        href: `/app/projects/${projectId}/deployments/${deploymentId}/scale`,
        current: false,
        icon: Expand
      },
      {
        name: "Network",
        href: `/app/projects/${projectId}/deployments/${deploymentId}/network`,
        current: false,
        icon: Network
      },
      {
        name: "Volumes",
        href: `/app/projects/${projectId}/deployments/${deploymentId}/volumes`,
        current: false,
        icon: Database
      },
      {
        name: "Instances",
        href: `/app/projects/${projectId}/deployments/${deploymentId}/instances`,
        current: false,
        icon: Boxes
      },
      {
        name: "Envs & Secrets",
        href: `/app/projects/${projectId}/deployments/${deploymentId}/envs`,
        current: false,
        icon: Variable
      },
      {
        name: "Settings",
        href: `/app/projects/${projectId}/deployments/${deploymentId}/settings`,
        current: false,
        icon: Cog
      }
    ];

    return items;
  };

  let items = generateItems($selectedProjectId, $selectedDeploymentId);

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

  $: items = generateItems($selectedProjectId, $selectedDeploymentId); // Regenerate items on projectId change
  $: setCurrentItem(); // Call setCurrentItem whenever items are updated

  $: if ($page) {
    setCurrentItem();
  }
</script>

<div class="flex flex-col gap-y-4" role="group" aria-labelledby="projects-headline">
  <a
    href={"/app/projects/" + $selectedProjectId}
    class=" text-white hover:text-primary-700 dark:text-gray-100 dark:hover:text-gray-100 pl-4 pr-10 py-2 text-sm font-medium rounded-md transition-all duration-150 ease-in-out hover:bg-gray-200 dark:hover:bg-primary-600 dark:hover:bg-opacity-10
     bg-primary-600
    "
  >
    <svelte:component this={ArrowLeft} class="w-5 h-5 mr-2 inline" strokeWidth={2} />
    Back
  </a>
  <hr class="border-gray-200 dark:border-gray-700" />
  {#each items as item}
    <a
      href={item.href}
      class=" text-gray-900 hover:text-gray-700 dark:text-gray-100 dark:hover:text-gray-100 pl-4 pr-10 py-2 text-sm font-medium rounded-md transition-all duration-150 ease-in-out hover:bg-gray-100 dark:hover:bg-primary-600 dark:hover:bg-opacity-10
        {item.current ? 'bg-gray-200 dark:bg-primary-600' : ''}
      "
      aria-current={item.current ? "page" : undefined}
    >
      <svelte:component this={item.icon} class="w-5 h-5 mr-2 inline" strokeWidth={2} />
      {item.name}
    </a>
  {/each}
</div>
