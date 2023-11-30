<script lang="ts">
  import ScaleSettings from "$lib/components/scale/ScaleSettings.svelte";
  import { Heading, P } from "flowbite-svelte";
</script>

<div class="flex items-start justify-between">
  <div class="flex flex-col">
    <Heading tag="h2">Scale</Heading>
    <P class="text-gray-500 dark:text-gray-400 text-sm">
      Scale your rollout to a specific number of instances.
    </P>
  </div>
  <!-- <Button
    color="primary"
    size="xs"
    class="whitespace-nowrap self-start"
    on:click={() => {
    }}
  >
    <Plus class="w-4 h-4 mr-2" />
    New image
  </Button> -->
</div>

<div>
  <ScaleSettings />
</div>
