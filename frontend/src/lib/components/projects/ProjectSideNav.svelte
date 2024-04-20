<script lang="ts">
  import { page } from "$app/stores";
  import { client } from "$lib/pocketbase";
  import { blueprints } from "$lib/stores/data";
  import { ArrowLeft, BookDashed, Cog, Plus } from "lucide-svelte";
  import selectedProjectId from "$lib/stores/project";

  export let modal: boolean;

  function getOwnedBlueprints() {
    return $blueprints.filter(
      (blueprint) => blueprint.owner === (client.authStore?.model?.id ?? null)
    );
  }

  // Return navigation items based on project settings
  let generateItems = () => {
    let items = [
      {
        name: `Blueprints (${getOwnedBlueprints().length})`,
        href: "/app/blueprints/my-blueprints",
        current: false,
        icon: BookDashed
      },
      {
        name: `Project Settings`,
        href: `/app/projects/${$selectedProjectId}/settings`,
        current: false,
        icon: Cog
      }
    ];

    return items;
  };

  let items = generateItems();

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

  $: items = generateItems(); // Regenerate items on projectId change
  $: setCurrentItem(); // Call setCurrentItem whenever items are updated

  $: if ($page) {
    setCurrentItem();
  }
</script>

<div class="flex flex-col gap-y-4" role="group" aria-labelledby="projects-headline">
  <a
    href={"/app"}
    class=" text-white hover:text-primary-700 dark:text-gray-100 dark:hover:text-gray-100 pl-4 pr-10 py-2 text-sm font-medium rounded-md transition-all duration-150 ease-in-out hover:bg-gray-200 dark:hover:bg-primary-600 dark:hover:bg-opacity-10
     bg-primary-600
    "
  >
    <svelte:component this={ArrowLeft} class="w-5 h-5 mr-2 inline" strokeWidth={2} />
    Back
  </a>
  <!-- Create new blueprint -->
  <button
    on:click={() => (modal = true)}
    class=" text-white hover:text-primary-700 dark:text-gray-100 dark:hover:text-gray-100 pl-4 pr-10 py-2 text-sm font-medium rounded-md transition-all duration-150 ease-in-out hover:bg-gray-200 dark:hover:bg-primary-600 dark:hover:bg-opacity-10
  bg-primary-600"
  >
    <svelte:component this={Plus} class="w-5 h-5 mr-2 inline" strokeWidth={2} />
    New Deployment
  </button>
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
