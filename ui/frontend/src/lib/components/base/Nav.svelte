<script lang="ts">
  import { goto } from "$app/navigation";
  import { navigating, page } from "$app/stores";
  import { client, logout } from "$lib/pocketbase";
  import { currentRollout, projects } from "$lib/stores/data";
  import selectedProjectId from "$lib/stores/project";
  import { avatarUrl } from "$lib/utils/user.utils";
  import {
    Avatar,
    Dropdown,
    DropdownItem,
    DropdownHeader,
    DropdownDivider,
    Indicator,
    Tooltip
  } from "flowbite-svelte";
  import { recordLogoUrl } from "$lib/utils/blueprint.utils";
  import type { ProjectsResponse } from "$lib/pocketbase/generated-types";
  import { fade } from "svelte/transition";
  import { getRolloutStatus } from "$lib/utils/rollouts";
  import { onMount } from "svelte";
  import type { RolloutStatusResponse } from "$lib/types/status";

  let avatarUrlString: any = avatarUrl();
  let current_rollout_status: RolloutStatusResponse | undefined;
  let rollout_status_color:
    | "gray"
    | "red"
    | "yellow"
    | "green"
    | "indigo"
    | "purple"
    | "blue"
    | "dark"
    | "orange"
    | "none"
    | "teal"
    | undefined;

  const determineRolloutColor = (status: string) => {
    switch (status) {
      case "Pending":
        return "yellow";
      case "Not Ready":
        return "yellow";
      case "Error":
        return "red";
      case "OK":
        return "green";
      default:
        return "gray";
    }
  };

  const updateCurrentRollout = () => {
    getRolloutStatus($selectedProjectId, $currentRollout?.id ?? "")
      .then((response) => {
        current_rollout_status = response;
        rollout_status_color = determineRolloutColor(
          current_rollout_status?.deployment?.status ?? ""
        );
      })
      .catch(() => {
        current_rollout_status = undefined;
        rollout_status_color = "yellow";
      });
  };

  onMount(updateCurrentRollout);

  $: if ($navigating) {
    updateCurrentRollout();
  }

  // update rollout status every 5 seconds
  onMount(() => {
    setInterval(() => {
      updateCurrentRollout();
    }, 5000);
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
    {#if $page.url.pathname.startsWith("/app/projects/")}
      <div in:fade={{ duration: 100 }} out:fade={{ duration: 100 }}>
        {#key $selectedProjectId}
          <div class="flex items-center">
            <div class="relative">
              {#if selectedProject?.avatar}
                <img
                  src={recordLogoUrl(selectedProject)}
                  alt="Tuple"
                  class="h-12 w-12 flex-none rounded-lg object-cover ring-2 ring-{rollout_status_color}-500 p-1"
                />
              {:else}
                <img
                  src={recordLogoUrl(selectedProject?.expand.blueprint)}
                  alt="Tuple"
                  class="h-12 w-12 flex-none rounded-lg object-cover ring-2 ring-{rollout_status_color}-500 p-1"
                />
              {/if}

              <Tooltip>
                Status: {current_rollout_status?.deployment.status ?? "Unknown"}
              </Tooltip>
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
