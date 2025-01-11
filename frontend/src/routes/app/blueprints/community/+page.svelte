<script lang="ts">
  import BlueprintCard from "$lib/components/blueprints/BlueprintCard.svelte";
  import { blueprints } from "$lib/stores/data";
  import { client } from "$lib/pocketbase";
  import type { BlueprintsResponse } from "$lib/pocketbase/generated-types";

  // filter blueprints to only show the ones owned by the user
  let filteredBlueprints: BlueprintsResponse[] = [];

  $: filteredBlueprints = $blueprints.filter((blueprint) =>
    blueprint.users.some((user) => user === client.authStore?.record?.id)
  );
</script>

<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
  {#each filteredBlueprints as blueprint (blueprint.id)}
    <BlueprintCard {blueprint} community={true} />
  {/each}
</div>
