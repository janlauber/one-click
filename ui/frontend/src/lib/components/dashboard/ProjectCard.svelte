<script lang="ts">
  import type { ProjectsResponse, TagsResponse } from "$lib/pocketbase/generated-types";
  import type { Pexpand } from "$lib/stores/data";
  import { formatDateTime, timeAgo } from "$lib/utils/date.utils";
  import { technologyLogoUrl } from "$lib/utils/technology.utils";
  import { Badge, Button, Dropdown, DropdownItem, Indicator, Tooltip } from "flowbite-svelte";
    import { EnvelopeSolid } from "flowbite-svelte-icons";
  import { MoreHorizontal, Pencil, Rocket, Tag, Trash, Trash2 } from "lucide-svelte";
  export let project: ProjectsResponse<Pexpand>;

  let tags: TagsResponse[] | undefined = project.expand?.tags;

  if (tags === undefined) {
    tags = [];
  }
</script>

<div class="rounded-xl border border-gray-200">
  <div class="flex items-center gap-x-4 border-b border-gray-900/5 p-6">
    <div class="relative">
      <img
        src={technologyLogoUrl(project.expand?.technology)}
        alt="Tuple"
        class="h-12 w-12 flex-none rounded-lg object-cover ring-1 ring-gray-900/10"
      />
      <Indicator color="dark"  size="xl" placement="top-right" class="text-xs font-bold text-white cursor-default">{project.expand?.deployments.length}</Indicator>
      <Tooltip>Deployments</Tooltip>
    </div>
    <div class="text-sm font-medium leading-6">{project.name}</div>
    <div class="relative ml-auto">
      <div class="flex justify-end">
        <Button color="alternative">
          <MoreHorizontal class="w-4 h-4" />
        </Button>
        <Dropdown class="p-0">
          <DropdownItem>
              <Rocket class="w-4 h-4 mr-2 inline-block" />
              Deploy</DropdownItem
          >
          <DropdownItem>
            <Pencil class="w-4 h-4 mr-2 inline-block" />
            Edit</DropdownItem
          >
          <DropdownItem>
            <Trash2 class="w-4 h-4 mr-2 inline-block" />
            Delete</DropdownItem
          >
        </Dropdown>
      </div>
    </div>
  </div>
  <dl class="-my-3 divide-y divide-gray-100 px-6 py-4 text-sm leading-6">
    <div class="flex justify-between gap-x-4 py-3">
      <dt class="">Last deployment</dt>
      <dd class=" cursor-default">
        <time datetime={formatDateTime(project.expand?.deployments[0].startDate)}>
          {timeAgo(project.expand?.deployments[0].startDate)}
        </time>
        <Tooltip>{formatDateTime(project.expand?.deployments[0].startDate)}</Tooltip>
      </dd>
    </div>
    <!-- <div class="flex justify-between gap-x-4 py-3">
      <dt class="">Status</dt>
      <dd class="flex items-start gap-x-2">
        <span class="font-medium">{project.expand?.status}</span>
      </dd>
    </div> -->
    {#if tags}
      <div class="flex justify-between gap-x-4 py-3">
        <dt class="">
          <Tag class="w-5 h-5 inline-block" strokeWidth={2} /> Tags
        </dt>
        <dd class="flex items-start gap-x-2">
          {#each tags as tag (tag.id)}
            <Badge color={tag.color} large class="cursor-default"
              >{tag.name.charAt(0).toUpperCase() + tag.name.slice(1)}</Badge
            >
          {/each}
        </dd>
      </div>
    {/if}
  </dl>
</div>
