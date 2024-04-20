<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import { client, logout } from "$lib/pocketbase";
  import { blueprints, currentRolloutStatus, deployments, projects } from "$lib/stores/data";
  import selectedProjectId from "$lib/stores/project";
  import { selectedProject } from "$lib/stores/data";
  import {
    Avatar,
    Dropdown,
    DropdownItem,
    DropdownHeader,
    DropdownDivider,
    Indicator
  } from "flowbite-svelte";
  import { recordLogoUrl } from "$lib/utils/blueprint.utils";
  import { fade } from "svelte/transition";
  import { FileQuestion } from "lucide-svelte";
  import { avatarUrlString } from "$lib/stores/avatar";
  import selectedDeploymentId from "$lib/stores/deployment";

  const determineRolloutColor = (status?: string) => {
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

  $selectedProject = $projects.find((p) => p.id === $selectedProjectId);

  // determine if rollout status is pending, not ready, error, or ok $currentRolloutStatus.deployment.status

  let selectedDeployment = $deployments.find((d) => d.id === $selectedDeploymentId);
  let selectedDeploymentBlueprint = $blueprints.find((b) => b.id === selectedDeployment?.blueprint);

  $: $selectedProject = $projects.find((p) => p.id === $selectedProjectId);
  $: selectedDeployment = $deployments.find((d) => d.id === $selectedDeploymentId);
  $: selectedDeploymentBlueprint = $blueprints.find((b) => b.id === selectedDeployment?.blueprint);
</script>

<nav class="bg-primary-600 dark:bg-primary-600 flex py-2">
  <div class="flex justify-between w-full max-w-6xl mx-auto px-4">
    <a href="/app" class="justify-start">
      <img src="/images/logo_background.png" class="mr-3 h-10" alt="Flowbite Logo" />
    </a>
    <!-- only /app/projects/${id}/deployments/ -->
    {#if $page.url.pathname.startsWith(`/app/projects/${$selectedProjectId}`)}
      <div in:fade={{ duration: 100 }} out:fade={{ duration: 100 }}>
        {#key $selectedProjectId}
          <div class="flex items-center">
            <div class="relative border-2 rounded-lg">
              {#if $page.url.pathname.startsWith(`/app/projects/${$selectedProjectId}/deployments`)}
                {#if selectedDeployment?.avatar}
                  <img
                    src={recordLogoUrl(selectedDeployment)}
                    alt="Tuple"
                    class="h-9 w-9 flex-none rounded-lg object-cover p-1"
                  />
                {:else if selectedDeploymentBlueprint?.avatar}
                  <img
                    src={recordLogoUrl(selectedDeploymentBlueprint)}
                    alt="Tuple"
                    class="h-9 w-9 flex-none rounded-lg object-cover p-1"
                  />
                {:else}
                  <FileQuestion class="h-9 w-9 flex-none text-white rounded-lg object-cover p-1" />
                {/if}

                <Indicator
                  color={determineRolloutColor($currentRolloutStatus?.deployment.status)}
                  size="md"
                  class="text-xs font-bold text-white cursor-default absolute -top-1 -right-1"
                />
                <Indicator
                  color={determineRolloutColor($currentRolloutStatus?.deployment.status)}
                  size="md"
                  class="text-xs font-bold text-white cursor-default animate-ping absolute -top-1 -right-1"
                />
              {:else if $page.url.pathname.startsWith(`/app/projects/${$selectedProjectId}`)}
                <img
                  src={recordLogoUrl($selectedProject)}
                  alt="Tuple"
                  class="h-9 w-9 flex-none rounded-lg object-cover p-1"
                />
              {:else}
                <FileQuestion class="h-9 w-9 flex-none text-white rounded-lg object-cover p-1" />
              {/if}
            </div>
            <div class="text-sm font-medium leading-6 text-white ml-4">
              {$selectedProject?.name}
            </div>
          </div>
        {/key}
      </div>
    {/if}

    <!-- <ComboBox /> -->
    <div class="flex items-center md:order-2 cursor-pointer active:scale-105">
      <Avatar id="avatar-menu" src={$avatarUrlString} />
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
