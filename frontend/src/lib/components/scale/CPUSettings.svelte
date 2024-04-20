<script lang="ts">
  import { Label, Range } from "flowbite-svelte";
  import { Cpu } from "lucide-svelte";
  export let cpuRequestsFloat = 0.1;
  export let cpuLimitFloat = 0.1;

  $: {
    if (cpuLimitFloat < cpuRequestsFloat) {
      cpuLimitFloat = cpuRequestsFloat;
    }
  }
</script>

<Label for="tag" class="block mb-1">
  <Cpu class="inline-block mr-1" size="1rem" />
  CPU
</Label>
<form>
  <Label for="cpu-requests" class="block mb-1">
    CPU Requests ({cpuRequestsFloat} vCPU)
  </Label>
  <Range id="cpu-requests" min="0.1" max="2" step="0.1" bind:value={cpuRequestsFloat} />
  <p class="text-xs text-gray-500 dark:text-gray-400">
    CPU requests are the guaranteed amount of CPU resources that a container will get.
  </p>
  <br />
  <Label for="cpu-limit" class="block mb-1">
    CPU Limit ({cpuLimitFloat} vCPU)
  </Label>
  <Range id="cpu-limit" min="0.1" max="4" step="0.1" bind:value={cpuLimitFloat} />
  <p class="text-xs text-gray-500 dark:text-gray-400">
    CPU limit is the maximum amount of CPU resources that a container can use.
  </p>
</form>
