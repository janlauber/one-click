<script lang="ts">
  import { Heading, P } from "flowbite-svelte";
  import RolloutsTable from "$lib/components/rollouts/RolloutsTable.svelte";
</script>

<div class="flex items-start justify-between">
  <div class="flex flex-col">
    <Heading tag="h2">Rollouts</Heading>
    <P class="text-gray-500 dark:text-gray-400 text-sm">
      Each time you change the configuration of a project, a new rollout is created. You can see maximum 10 rollouts.
    </P>
  </div>
</div>

<div>
  <RolloutsTable />
</div>
