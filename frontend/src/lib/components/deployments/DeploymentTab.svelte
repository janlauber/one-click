<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import type { BlueprintsResponse, DeploymentsResponse } from "$lib/pocketbase/generated-types";
  import { blueprints } from "$lib/stores/data";
  import selectedDeploymentId from "$lib/stores/deployment";
  import { recordLogoUrl } from "$lib/utils/blueprint.utils";
  import { FileQuestion } from "lucide-svelte";

  export let deployments: DeploymentsResponse[];

  // /app/projects/:projectId/deployments/:deploymentId/:suffix or no suffix
  const DEPLOYMENT_ID_INDEX = 5;
  let currentDeploymentId = $page.url.pathname.split("/")[DEPLOYMENT_ID_INDEX];

  let currentDeployment = deployments.find((d) => d.id === currentDeploymentId);

  function getDeploymentBlueprint(deploymentId: string): BlueprintsResponse | undefined {
    return $blueprints.find((b) => b.id === deploymentId);
  }
</script>

<div class="flex gap-4 dark-text-white mb-2 overflow-x-auto whitespace-nowrap">
  {#each deployments as deployment}
    <button
      on:click={() => {
        $selectedDeploymentId = deployment.id;
        // check if there is a suffix in the current $page.url.pathname after /app/projects/${deployment.project}/deployments/${deployment.id} e.g. /app/projects/${deployment.project}/deployments/${deployment.id}/map
        let suffix = $page.url.pathname
          .split("/")
          .slice(DEPLOYMENT_ID_INDEX + 1)
          .join("/");
        if (suffix) {
          goto(`/app/projects/${deployment.project}/deployments/${deployment.id}/${suffix}`);
          return;
        }
        goto(`/app/projects/${deployment.project}/deployments/${deployment.id}`);
      }}
      class={`flex items-center gap-2 px-2 py-2 transition-all group cursor-pointer border-b-2 ${
        currentDeployment?.id === deployment.id
          ? "border-primary-500 dark:border-white"
          : "hover:border-primary-500 dark:hover:border-white border-transparent"
      }`}
    >
      <div
        class={`relative border-2 rounded-lg group-hover:-translate-y-0.5 transition-all ${
          currentDeployment?.id === deployment.id
            ? "border-primary-500 dark:border-white"
            : "group-hover:border-primary-500 dark:group-hover:border-white"
        }`}
      >
        {#if deployment.avatar || getDeploymentBlueprint(deployment.blueprint)}
          <img
            src={recordLogoUrl(
              deployment.avatar ? deployment : getDeploymentBlueprint(deployment.blueprint)
            )}
            alt="Tuple"
            class="h-7 w-7 min-w-7 flex-none rounded-lg object-cover p-1"
          />
        {:else}
          <div class="h-7 w-7">
            <FileQuestion />
          </div>
        {/if}
      </div>
      <span class="group-hover:-translate-y-0.5 transition-all">
        {deployment.name}
      </span>
    </button>
  {/each}
</div>
