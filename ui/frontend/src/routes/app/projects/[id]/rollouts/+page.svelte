<script>
  import { goto } from "$app/navigation";
  import { Button, Heading, P } from "flowbite-svelte";
  import { ArrowRight, Plus } from "lucide-svelte";
  import selectedProjectId from "$lib/stores/project";
</script>

<div class="flex items-start justify-between">
  <div class="flex flex-col">
    <Heading tag="h2">Rollouts</Heading>
    <P class="text-gray-500 dark:text-gray-400 text-sm">All rollouts for your project.</P>
  </div>
  <Button
    color="dark"
    size="xs"
    class="whitespace-nowrap self-start"
    on:click={() => {
      goto(`/app/projects/${$selectedProjectId}/rollouts`);
    }}
  >
    <Plus class="w-4 h-4 mr-2" />
    New rollout
  </Button>
</div>


