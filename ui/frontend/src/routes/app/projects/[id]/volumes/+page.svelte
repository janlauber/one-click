<script>
  import { Button, Heading, P } from "flowbite-svelte";
    import { Plus } from "lucide-svelte";
</script>

<div class="flex items-start justify-between">
</div>

<div class="flex items-start justify-between">
  <div class="flex flex-col">
    <Heading tag="h2">Volumes</Heading>
    <P class="text-gray-500 dark:text-gray-400 text-sm">
      Persistent volumes for your rollout.
    </P>
  </div>
  <Button
    color="primary"
    size="xs"
    class="whitespace-nowrap self-start"
    on:click={() => {
    }}
  >
    <Plus class="w-4 h-4 mr-2" />
    New volume
  </Button>
</div>

<div>
</div>
