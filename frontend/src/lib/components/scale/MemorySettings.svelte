<script lang="ts">
  import { Label, Radio } from "flowbite-svelte";
    import { MemoryStick } from "lucide-svelte";

  export let memoryRequestsInt = 128;
</script>

<Label for="tag" class="block mb-1">
  <MemoryStick class="inline-block mr-1" size="1rem" />
  Memory
</Label>
<form>
  <ul
    class="items-center w-full rounded-lg border border-gray-200 sm:flex dark:bg-gray-800 dark:border-gray-600 divide-x rtl:divide-x-reverse divide-gray-200 dark:divide-gray-600"
  >
    <li class="w-full">
      <Radio name="hor-list" class="p-3" bind:group={memoryRequestsInt} value={256}>
        <div class="block">
          <p class="">Micro</p>
          <p class="text-xs font-light">(256 MB)</p>
        </div>
      </Radio>
    </li>
    <li class="w-full">
      <Radio name="hor-list" class="p-3" bind:group={memoryRequestsInt} value={512}>
        <div class="block">
          <p class="">Small</p>
          <p class="text-xs font-light">(512 MB)</p>
        </div>
      </Radio>
    </li>
    <li class="w-full">
      <Radio name="hor-list" class="p-3" bind:group={memoryRequestsInt} value={1024}>
        <div class="block">
          <p class="">Medium</p>
          <p class="text-xs font-light">(1 GB)</p>
        </div>
      </Radio>
    </li>
    <li class="w-full">
      <Radio name="hor-list" class="p-3" bind:group={memoryRequestsInt} value={2048}>
        <div class="block">
          <p class="">Large</p>
          <p class="text-xs font-light">(2 GB)</p>
        </div>
      </Radio>
    </li>
    <li class="w-full">
      <Radio name="hor-list" class="p-3" bind:group={memoryRequestsInt} value={4096}>
        <div class="block">
          <p class="">X-Large</p>
          <p class="text-xs font-light">(4 GB)</p>
        </div>
      </Radio>
    </li>
  </ul>
</form>
