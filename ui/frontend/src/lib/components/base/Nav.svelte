<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import { client, logout } from "$lib/pocketbase";
  import { projects, rollouts } from "$lib/stores/data";
  import selectedProjectId from "$lib/stores/project";
  import { avatarUrl } from "$lib/utils/user.utils";
  import {
    Avatar,
    Dropdown,
    DropdownItem,
    DropdownHeader,
    DropdownDivider,
    type SelectOptionType,
    Indicator,
    Tooltip
  } from "flowbite-svelte";
  import { frameworkLogoUrl } from "$lib/utils/framework.utils";
  import type { ProjectsResponse } from "$lib/pocketbase/generated-types";
  import { fade } from "svelte/transition";
  import { onMount } from "svelte";

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

  let selectedProject: ProjectsResponse | undefined = $projects.find(
    (p) => p.id === $selectedProjectId
  );

  $: selectedProject = $projects.find((p) => p.id === $selectedProjectId);
</script>

<nav class="bg-primary-600 dark:bg-primary-600 flex py-3">
  <div class="flex justify-between w-full max-w-6xl mx-auto px-4">
    <a href="/app" class="justify-start">
      <img src="/images/logo_background.png" class="mr-3 h-10" alt="Flowbite Logo" />
    </a>
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
    {#if $page.url.pathname.startsWith("/app/projects/")}
      <div in:fade={{ duration: 100 }} out:fade={{ duration: 100 }}>
        {#key $selectedProjectId}
          <div class="flex items-center">
            <div class="relative">
              <img
                src={frameworkLogoUrl(selectedProject?.expand.framework)}
                alt="Tuple"
                class="h-12 w-12 flex-none rounded-lg object-cover ring-1 ring-white p-1"
              />
              <Indicator
                size="xl"
                placement="top-right"
                class="text-xs font-bold bg-white cursor-default"
                >{$rollouts.filter((r) => r.expand?.project.id === $selectedProjectId).length || 0}
              </Indicator>
              <Tooltip>Rollouts</Tooltip>
            </div>
            <div class="text-sm font-medium leading-6 text-white ml-4">{selectedProject?.name}</div>
          </div>
        {/key}
      </div>
    {/if}

    <!-- <ComboBox /> -->
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
  </div>
</nav>
