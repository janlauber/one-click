<script lang="ts">
  import { Label, Radio } from "flowbite-svelte";
  import { Cpu } from "lucide-svelte";
  export let cpuRequestsFloat = 0.25;

  interface CPUSettingsProps {
    name: string;
    value: number;
  }

  let cpuSettings: CPUSettingsProps[] = [
    {
      name: "Nano",
      value: 0.25
    },
    {
      name: "Micro",
      value: 0.5
    },
    {
      name: "Small",
      value: 0.75
    },
    {
      name: "Medium",
      value: 1
    },
    {
      name: "Large",
      value: 2
    }
  ];


</script>

<Label for="tag" class="block mb-1">
  <Cpu class="inline-block mr-1" size="1rem" />
  CPU
</Label>
<form>
  <ul
    class="items-center w-full rounded-lg border border-gray-200 sm:flex dark:bg-gray-800 dark:border-gray-600 divide-x rtl:divide-x-reverse divide-gray-200 dark:divide-gray-600"
  >
    {#each cpuSettings as { name, value }, i}
      <li class="w-full">
        <Radio name="hor-list" class="p-3" bind:group={cpuRequestsFloat} value={value}>
          <div class="block">
            <p class="">{name}</p>
            <p class="text-xs font-light">({value} vCPU)</p>
          </div>
        </Radio>
      </li>
    {/each}
  </ul>
</form>
