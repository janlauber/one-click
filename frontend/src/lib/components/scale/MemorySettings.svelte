<script lang="ts">
  import { Label, Radio } from "flowbite-svelte";
  import { MemoryStick } from "lucide-svelte";

  export let memoryRequestsInt = 256;

  interface MemorySettingsProps {
    name: string;
    value: number;
  }

  let memorySettings: MemorySettingsProps[] = [
    {
      name: "Nano",
      value: 256
    },
    {
      name: "Micro",
      value: 512
    },
    {
      name: "Small",
      value: 1024
    },
    {
      name: "Medium",
      value: 2048
    },
    {
      name: "Large",
      value: 4096
    }
  ];
</script>

<Label for="tag" class="block mb-1">
  <MemoryStick class="inline-block mr-1" size="1rem" />
  Memory
</Label>
<form>
  <ul
    class="items-center w-full rounded-lg border border-gray-200 sm:flex dark:bg-gray-800 dark:border-gray-600 divide-x rtl:divide-x-reverse divide-gray-200 dark:divide-gray-600"
  >
    {#each memorySettings as { name, value }, i}
      <li class="w-full">
        <Radio name="hor-list" class="p-3" bind:group={memoryRequestsInt} {value}>
          <div class="block">
            <p class="">{name}</p>
            <p class="text-xs font-light">({value / 1024} GB)</p>
          </div>
        </Radio>
      </li>
    {/each}
  </ul>
</form>
