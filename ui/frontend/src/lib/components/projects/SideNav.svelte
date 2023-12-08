<script lang="ts">
  import { page } from "$app/stores";
  import { projects } from "$lib/stores/data";
  import selectedProjectId from "$lib/stores/project";
  import type { ProjectSettings } from "$lib/utils/interfaces";
  import {
    ArrowLeft,
    Box,
    Boxes,
    Cog,
    Database,
    Expand,
    FileText,
    HardDrive,
    History,
    LineChart,
    Network,
    Package,
    Variable
  } from "lucide-svelte";

  // Get current project settings
  // @ts-ignore
  const projectSettings: ProjectSettings = $projects.find(
    (project) => project.id === $selectedProjectId
  )?.expand?.framework.settings;

  // Return navigation items based on project settings
  const generateItems = (projectId: string) => {
    const items = [
      {
        name: "Overview",
        href: `/app/projects/${projectId}/overview`,
        current: false,
        icon: LineChart
      },
      {
        name: "Logs",
        href: `/app/projects/${projectId}/logs`,
        current: false,
        icon: FileText
      },
      {
        name: "Rollouts",
        href: `/app/projects/${projectId}/rollouts`,
        current: false,
        icon: History
      },
      {
        name: "Builds",
        href: `/app/projects/${projectId}/builds`,
        current: false,
        icon: Package
      },
      {
        name: "Image",
        href: `/app/projects/${projectId}/image`,
        current: false,
        icon: HardDrive
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
        name: "Volumes",
        href: `/app/projects/${projectId}/volumes`,
        current: false,
        icon: Database
      },
      {
        name: "Instances",
        href: `/app/projects/${projectId}/instances`,
        current: false,
        icon: Boxes
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

    // Remove items based on project settings
    if (!projectSettings?.overview) {
      removeItem(
        items,
        items.find((item) => item.name === "Overview")
      );
    }
    if (!projectSettings?.logs) {
      removeItem(
        items,
        items.find((item) => item.name === "Logs")
      );
    }
    if (!projectSettings?.rollouts) {
      removeItem(
        items,
        items.find((item) => item.name === "Rollouts")
      );
    }
    if (!projectSettings?.builds) {
      removeItem(
        items,
        items.find((item) => item.name === "Builds")
      );
    }
    if (!projectSettings?.image) {
      removeItem(
        items,
        items.find((item) => item.name === "Image")
      );
    }
    if (!projectSettings?.scale) {
      removeItem(
        items,
        items.find((item) => item.name === "Scale")
      );
    }
    if (!projectSettings?.network) {
      removeItem(
        items,
        items.find((item) => item.name === "Network")
      );
    }
    if (!projectSettings?.volumes) {
      removeItem(
        items,
        items.find((item) => item.name === "Volumes")
      );
    }
    if (!projectSettings?.instances) {
      removeItem(
        items,
        items.find((item) => item.name === "Instances")
      );
    }
    if (!projectSettings?.envs) {
      removeItem(
        items,
        items.find((item) => item.name === "Envs & Secrets")
      );
    }
    if (!projectSettings?.settings) {
      removeItem(
        items,
        items.find((item) => item.name === "Settings")
      );
    }
    if (!projectSettings?.logs && !projectSettings?.rollouts) {
      removeItem(
        items,
        items.find((item) => item.name === "Logs")
      );
    }

    return items;
  };

  function removeItem(array: any[], item: any) {
    const index = array.indexOf(item);
    if (index > -1) {
      array.splice(index, 1);
    }
  }

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
  <a
    href={"/app"}
    class=" text-white hover:text-primary-700 dark:text-gray-100 dark:hover:text-gray-100 pl-4 pr-10 py-2 text-sm font-medium rounded-md transition-all duration-150 ease-in-out hover:bg-primary-100 dark:hover:bg-primary-600 dark:hover:bg-opacity-10
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
