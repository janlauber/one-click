<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import { client, logout } from "$lib/pocketbase";
  import { projects } from "$lib/stores/data";
  import selectedProjectId from "$lib/stores/project";
  import { avatarUrl } from "$lib/utils/user.utils";
  import {
    Navbar,
    NavBrand,
    Avatar,
    Dropdown,
    DropdownItem,
    DropdownHeader,
    DropdownDivider,
    Select,
    type SelectOptionType
  } from "flowbite-svelte";
    import ComboBox from "./ComboBox.svelte";

  let avatarUrlString: any = avatarUrl();
  let projectsChoices: SelectOptionType<any>[] | undefined;

  projects.subscribe((value) => {
    projectsChoices = value.map((project) => {
      return {
        name: project.name,
        value: project.id
      };
    });
  });

  $: {
    if ($page.url.pathname.startsWith("/app/projects/")) {
      const pathParts = $page.url.pathname.split("/");
      const currentProjectId = pathParts[3];
      const slug = pathParts[4];

      if ($selectedProjectId && slug && currentProjectId !== $selectedProjectId) {
        const targetUrl = `/app/projects/${$selectedProjectId}/${slug}`;
        if ($page.url.pathname !== targetUrl) {
          goto(targetUrl);
        }
      }
    }
  }
</script>

<Navbar class="bg-primary-600 dark:bg-primary-600">
  <NavBrand href="/app" class="justify-start">
    <img src="/images/logo_background.png" class="mr-3 h-10" alt="Flowbite Logo" />
  </NavBrand>
  <!-- display only when under /app/projects/{id} -->
  <!-- {#if $page.url.pathname.startsWith("/app/projects/")}
    <div class="flex items-center justify-center md:order-1 cursor-pointer active:scale-105">
      <Select
        placeholder="Choose Project"
        size="sm"
        items={projectsChoices}
        bind:value={$selectedProjectId}
      ></Select>
    </div>
  {/if} -->
  <ComboBox />
  <div class="flex items-center md:order-2 cursor-pointer active:scale-105">
    <Avatar id="avatar-menu" src={avatarUrlString} />
  </div>
  <Dropdown placement="bottom" triggeredBy="#avatar-menu">
    <DropdownHeader>
      <span class="block text-sm">{client.authStore.model?.name}</span>
      <span class="block truncate text-sm font-medium">{client.authStore.model?.email}</span>
    </DropdownHeader>
    <DropdownItem
      on:click={() => {
        goto("/app/profile");
      }}>Settings</DropdownItem
    >
    <DropdownDivider />
    <DropdownItem
      on:click={() => {
        logout();
      }}>Sign out</DropdownItem
    >
  </Dropdown>
</Navbar>
