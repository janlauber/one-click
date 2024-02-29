<script lang="ts">
  import ProjectCard from "$lib/components/blueprints/ProjectCard.svelte";
import { blueprints } from "$lib/stores/data";
</script>

<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
  {#each $blueprints as blueprint (blueprint.id)}
    <ProjectCard {blueprint} />
  {/each}
</div>
