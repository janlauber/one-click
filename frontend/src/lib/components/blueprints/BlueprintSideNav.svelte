<script lang="ts">
  import { page } from "$app/stores";
  import { client } from "$lib/pocketbase";
  import { UpdateFilterEnum, blueprints, updateDataStores } from "$lib/stores/data";
  import { Button, Fileupload, Input, Label, Modal, Toggle } from "flowbite-svelte";
  import { ArrowLeft, BookLock, BookPlus, BookUser, Code2 } from "lucide-svelte";
  import MonacoEditor from "svelte-monaco";
  // @ts-expect-error - MonacoEditor types are not available
  import yaml from "js-yaml";
  import toast from "svelte-french-toast";
  import type { BlueprintsRecord } from "$lib/pocketbase/generated-types";
  import { goto } from "$app/navigation";

  let blueprintModalOpen = false;

  let name: string = "";
  let description: string = "";
  let isPrivate: boolean = false;
  let avatar: string = "";
  let avatarFile: File;
  let manifest: any = "";

  // Get current project settings

  function getOwnedBlueprints() {
    return $blueprints.filter(
      (blueprint) => blueprint.owner === (client.authStore?.model?.id ?? null)
    );
  }

  function getCommunityBlueprints() {
    return $blueprints.filter(
      (blueprint) =>
        blueprint.owner !== client.authStore?.model?.id &&
        blueprint.users.some((user) => user === client.authStore?.model?.id)
    );
  }

  // Return navigation items based on project settings
  let generateItems = () => {
    let items = [
      {
        name: `My Blueprints (${getOwnedBlueprints().length})`,
        href: `/app/blueprints/my-blueprints`,
        current: false,
        icon: BookLock
      },
      {
        name: `Community (${getCommunityBlueprints().length})`,
        href: `/app/blueprints/community`,
        current: false,
        icon: BookUser
      }
      // {
      //   name: `Git`,
      //   href: `/app/blueprints/git`,
      //   current: false,
      //   icon: FolderGit
      // }
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

  async function handleSaveManifest() {
    if (!name) {
      toast.error("Name is required");
      return;
    }

    if (!description) {
      toast.error("Description is required");
      return;
    }

    if (!manifest) {
      toast.error("Manifest is required");
      return;
    }

    let formData = new FormData();
    formData.append("avatar", avatarFile);

    // parse the manifest yaml to json
    const parsedManifest = yaml.load(manifest);

    const data: BlueprintsRecord = {
      name,
      description,
      private: isPrivate,
      manifest: parsedManifest,
      owner: client.authStore?.model?.id
    };

    toast.promise(
      client
        .collection("blueprints")
        .create(data)
        .then((response) => {
          if (avatarFile) {
            client
              .collection("blueprints")
              .update(response.id, formData)
              .then(() => {
                updateDataStores({
                  filter: UpdateFilterEnum.ALL
                });
              })
              .catch((error) => {
                toast.error(error.message);
              });
          }

          updateDataStores({
            filter: UpdateFilterEnum.ALL
          });

          blueprintModalOpen = false;
        }),
      {
        loading: "Saving...",
        success: "Blueprint saved",
        error: "Error saving blueprint"
      }
    );
  }
</script>

<div class="flex flex-col gap-y-4 mt-3" role="group" aria-labelledby="projects-headline">
  <button
    on:click={() => {
      goto("/app");
    }}
    class=" text-left text-white hover:text-primary-700 dark:text-gray-100 dark:hover:text-gray-100 pl-4 pr-10 py-2 text-sm font-medium rounded-md transition-all duration-150 ease-in-out hover:bg-gray-200 dark:hover:bg-primary-600 dark:hover:bg-opacity-10
     bg-primary-600
    "
  >
    <svelte:component this={ArrowLeft} class="w-5 h-5 mr-2 inline" strokeWidth={2} />
    Back
  </button>
  <!-- Create new blueprint -->
  <button
    on:click={() => (blueprintModalOpen = true)}
    class=" text-white hover:text-primary-700 dark:text-gray-100 dark:hover:text-gray-100 pr-10 py-2 text-sm font-medium rounded-md transition-all duration-150 ease-in-out hover:bg-gray-200 dark:hover:bg-primary-600 dark:hover:bg-opacity-10
  bg-primary-600"
  >
    <svelte:component this={BookPlus} class="w-5 h-5 mr-2 inline" strokeWidth={2} />
    New Blueprint
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

<Modal bind:open={blueprintModalOpen} size="lg">
  <div class="text-center">
    <Code2 class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      <!-- Refer to https://docs.one-click.dev for advanced documentation about the manifest values -->
      Refer to
      <a href="https://docs.one-click.dev" target="_blank" class="text-primary-500"
        >docs.one-click.dev</a
      >
    </h3>
  </div>
  <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 rounded-lg">
    <div class="col-span-1">
      <Label class="">Name*</Label>
      <Input bind:value={name} required />
    </div>
    <div class="col-span-1">
      <Label class="">Description*</Label>
      <Input bind:value={description} required />
    </div>
    <div />
    <div class="flex justify-end">
      <Toggle bind:checked={isPrivate} class="">
        <span class="text-gray-500 dark:text-gray-400">Private</span>
      </Toggle>
    </div>
    <div class="col-span-2">
      <Label class="">Avatar*</Label>
      <Fileupload
        bind:value={avatar}
        on:change={(event) => {
          // @ts-expect-error - event.target.files is not available
          avatarFile = event.target.files[0];
        }}
      />
    </div>
  </div>
  <div class=" h-96 overflow-y-auto rounded-lg p-2" style="background-color: #1E1E1E;">
    <MonacoEditor
      bind:value={manifest}
      options={{ language: "yaml", automaticLayout: true, minimap: { enabled: false } }}
      theme="vs-dark"
    />
  </div>
  <div class="flex justify-between">
    <Button color="primary" class="me-2" on:click={() => handleSaveManifest()}>Save</Button>
    <Button color="alternative" on:click={() => (blueprintModalOpen = false)}>Cancel</Button>
  </div>
</Modal>
