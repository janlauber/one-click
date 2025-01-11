<script lang="ts">
  import { Card, Avatar, Button } from "flowbite-svelte";
  import { UpdateFilterEnum, blueprints, updateDataStores } from "$lib/stores/data";
  import type { BlueprintsResponse } from "$lib/pocketbase/generated-types";
  import { recordLogoUrl } from "$lib/utils/blueprint.utils";
  import { ArrowRight, BookPlus, Lock } from "lucide-svelte";
  import toast from "svelte-french-toast";
  import { goto } from "$app/navigation";
  import { client } from "$lib/pocketbase";

  export let data: any;
  let blueprint: BlueprintsResponse | undefined = undefined;

  $: blueprint = $blueprints.find((blueprint) => blueprint.id === data.props?.blueprintId);

  async function handleAddToBlueprints(blueprintId: string | undefined) {
    const response = await postToBlueprints(blueprintId);

    if (response) {
      if (response.status === 200) {
        toast.success("Blueprint added to your community blueprints");

        // refresh blueprints
        updateDataStores({
          filter: UpdateFilterEnum.ALL
        }).then(() => {
          goto("/app/blueprints/community");
        });
      } else {
        console.error("Error adding blueprint to your blueprints", response);
        toast.error("Error adding blueprint to your community blueprints");
      }
    }
  }

  async function postToBlueprints(blueprintId: string | undefined) {
    const token = localStorage.getItem("pocketbase_auth");
    if (!token) {
      return;
    }
    const authHeader = { Authorization: `Bearer ${JSON.parse(token).token}` };

    // if localhost, use localhost:8090 as base url
    if (window.location.hostname === "localhost") {
      try {
        const response = await fetch(`http://localhost:8090/pb/blueprints/shared/${blueprintId}`, {
          method: "POST",
          headers: {
            ...authHeader,
            "Content-Type": "application/json"
          }
        });
        return response;
      } catch (error) {
        console.error("Error adding blueprint to your blueprints", error);
      }
    } else {
      try {
        const response = await fetch(`/pb/blueprints/shared/${blueprintId}`, {
          method: "POST",
          headers: {
            ...authHeader,
            "Content-Type": "application/json"
          }
        });
        return response;
      } catch (error) {
        console.error("Error adding blueprint to your blueprints", error);
      }
    }
  }
</script>

<div class="h-full w-full flex items-center justify-center bg-gray-100 dark:bg-gray-800">
  <Card padding="lg" size="xl">
    <div class="flex flex-col items-center">
      {#if recordLogoUrl(blueprint) === ""}
        <Lock class="w-16 h-16 text-gray-500 dark:text-gray-400" />
      {:else}
        <Avatar size="lg" src={recordLogoUrl(blueprint)} />
      {/if}
      <h5 class="mb-1 text-xl font-medium text-gray-900 dark:text-gray-400">
        {blueprint?.name || "Private Blueprint"}
      </h5>
      <span class="text-sm text-gray-500 dark:text-gray-400">
        {blueprint?.description || "Ask the owner of this blueprint to make it public"}
      </span>
      {#if blueprint !== undefined}
        <div class="flex mt-4 space-x-3 rtl:space-x-reverse lg:mt-6">
          {#if blueprint?.owner == client.authStore?.record?.id}
            <Button on:click={() => goto(`/app/blueprints/my-blueprints`)}>
              Edit
              <ArrowRight class="w-5 h-5 ml-2 inline" />
            </Button>
          {:else if blueprint?.users.includes(client.authStore?.record?.id ?? "")}
            <Button
              on:click={() => {
                goto(`/app/blueprints/community`);
              }}
            >
              View
              <ArrowRight class="w-5 h-5 ml-2 inline" />
            </Button>
          {:else}
            <Button on:click={() => handleAddToBlueprints(blueprint?.id)}>
              <BookPlus class="w-5 h-5 mr-2 -ml-1 inline" />
              Add to blueprints
            </Button>
          {/if}
        </div>
      {/if}
    </div>
  </Card>
</div>
