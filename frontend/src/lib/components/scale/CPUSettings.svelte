<script lang="ts">
  import { Label, Radio } from "flowbite-svelte";
    import { Cpu } from "lucide-svelte";
  export let cpuRequestsFloat = 0.1;
</script>

<Label for="tag" class="block mb-1">
  <Cpu class="inline-block mr-1" size="1rem" />
  CPU
</Label>
<form>
  <ul
    class="items-center w-full rounded-lg border border-gray-200 sm:flex dark:bg-gray-800 dark:border-gray-600 divide-x rtl:divide-x-reverse divide-gray-200 dark:divide-gray-600"
  >
    <li class="w-full">
      <Radio name="hor-list" class="p-3" bind:group={cpuRequestsFloat} value={0.25}>
        <div class="block">
          <p class="">Micro</p>
          <p class="text-xs font-light">(0.25 vCPU)</p>
        </div>
      </Radio>
    </li>
    <li class="w-full">
      <Radio name="hor-list" class="p-3" bind:group={cpuRequestsFloat} value={0.5}>
        <div class="block">
          <p class="">Small</p>
          <p class="text-xs font-light">(0.5 vCPU)</p>
        </div>
      </Radio>
    </li>
    <li class="w-full">
      <Radio name="hor-list" class="p-3" bind:group={cpuRequestsFloat} value={1}>
        <div class="block">
          <p class="">Medium</p>
          <p class="text-xs font-light">(1 vCPU)</p>
        </div>
      </Radio>
    </li>
    <li class="w-full">
      <Radio name="hor-list" class="p-3" bind:group={cpuRequestsFloat} value={2}>
        <div class="block">
          <p class="">Large</p>
          <p class="text-xs font-light">(2 vCPU)</p>
        </div>
      </Radio>
    </li>
    <li class="w-full">
      <Radio name="hor-list" class="p-3" bind:group={cpuRequestsFloat} value={4}>
        <div class="block">
          <p class="">X-Large</p>
          <p class="text-xs font-light">(4 vCPU)</p>
        </div>
      </Radio>
    </li>
  </ul>
</form>
