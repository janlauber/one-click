<script lang="ts">
  import { projects, rollouts } from "$lib/stores/data";
  import { Folder, History } from "lucide-svelte";

  let rolloutsCount: number = $rollouts.length;

  // update with delay of 200ms after $rollouts changes
  $: setTimeout(() => {
    rolloutsCount = $rollouts.length;
  }, 200);

</script>

<div class="stats flex bg-gray-50 border dark:bg-primary-900 dark:text-white dark:shadow-none">
  <div class="stat">
    <div class="stat-figure text-primary">
      <Folder class="inline-block w-8 h-8 stroke-current" />
    </div>
    <div class="stat-title dark:text-white">Projects</div>
    <div class="stat-value">{$projects.length}</div>
    <div class="stat-desc dark:text-white">Total</div>
  </div>

  <div class="stat">
    <div class="stat-figure text-primary">
      <History class="inline-block w-8 h-8 stroke-current" />
    </div>
    <div class="stat-title dark:text-white">Rollouts</div>
    <div class="stat-value">{rolloutsCount}</div>
    <div class="stat-desc dark:text-white">Total</div>
  </div>
</div>
