<script lang="ts">
  import { goto } from "$app/navigation";
  import { Button, Card, Heading, Indicator } from "flowbite-svelte";
  import { ArrowRight } from "lucide-svelte";
  import selectedProjectId from "$lib/stores/project";
    import CpuChart from "$lib/components/projects/CpuChart.svelte";
</script>

<div class="flex items-center justify-between">
  <Heading tag="h2">Overview</Heading>
  <Button
    color="alternative"
    size="xs"
    class="whitespace-nowrap"
    on:click={() => {
      goto(`/app/projects/${$selectedProjectId}/deployments`);
    }}
  >
    <div class="relative">
      <Indicator size="sm" color="green" class="mr-1.5 absolute" />
      <Indicator size="sm" color="green" class="mr-1.5 animate-ping" />
    </div>
    Current deployment
    <ArrowRight class="w-4 h-4 ml-2" />
  </Button>
</div>

<div class=" gap-4 space-y-4 mt-4">
  <Card size="xl"></Card>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <Card size="xl">
      <CpuChart />
    </Card>
    <Card size="xl"></Card>
  </div>
</div>
