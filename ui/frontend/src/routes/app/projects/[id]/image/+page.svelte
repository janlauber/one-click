<script lang="ts">
  import ImagesSettings from "$lib/components/image/ImageSettings.svelte";
  import { Heading, P } from "flowbite-svelte";
</script>

<div class="flex items-start justify-between">
  <div class="flex flex-col">
    <Heading tag="h2">Image</Heading>
    <P class="text-gray-500 dark:text-gray-400 text-sm">
      Define your image and tag for your project.
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
  <ImagesSettings />
</div>
