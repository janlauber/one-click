<script>
  import { Tabs, TabItem } from "flowbite-svelte";
  import { Boxes, Cog, Group } from "lucide-svelte";
</script>

<Tabs style="underline">
  <TabItem open>
    <div slot="title" class="flex items-center gap-2">
      <Group class="w-5 h-5" />
      Overview
    </div>
    <p class="text-sm text-gray-500 dark:text-gray-400">
      <b>Profile:</b>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore
      et dolore magna aliqua.
    </p>
  </TabItem>
  <TabItem>
    <div slot="title" class="flex items-center gap-2">
      <Boxes class="w-5 h-5" />
      Deployments
    </div>
    <p class="text-sm text-gray-500 dark:text-gray-400">
      <b>Dashboard:</b>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore
      et dolore magna aliqua.
    </p>
  </TabItem>
  <TabItem>
    <div slot="title" class="flex items-center gap-2">
      <Cog class="w-5 h-5" />
      Settings
    </div>
    <p class="text-sm text-gray-500 dark:text-gray-400">
      <b>Settings:</b>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore
      et dolore magna aliqua.
    </p>
  </TabItem>
</Tabs>
