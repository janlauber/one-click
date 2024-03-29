<script lang="ts">
  import { Label, Range } from "flowbite-svelte";
  import { MemoryStick } from "lucide-svelte";

  export let memoryRequestsInt = 128;
  export let memoryLimitsInt = 128;

  $: {
    if (memoryLimitsInt < memoryRequestsInt) {
      memoryLimitsInt = memoryRequestsInt;
    }
  }
</script>

<Label for="tag" class="block mb-1">
  <MemoryStick class="inline-block mr-1" size="1rem" />
  Memory
</Label>
<form>
  <Label for="memory-requests" class="block mb-1">
    Memory Requests ({memoryRequestsInt} MB)
  </Label>
  <Range id="memory-requests" min="128" max="8192" step="128" bind:value={memoryRequestsInt} />
  <p class="text-xs text-gray-500 dark:text-gray-400">
    Memory requests are the guaranteed amount of memory resources that a container will get.
  </p>
  <br />
  <Label for="memory-limit" class="block mb-1">
    Memory Limit ({memoryLimitsInt} MB)
  </Label>
  <Range id="memory-limit" min="128" max="16384" step="128" bind:value={memoryLimitsInt} />
  <p class="text-xs text-gray-500 dark:text-gray-400">
    Memory limit is the maximum amount of memory resources that a container can use.
  </p>
</form>
