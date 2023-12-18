<script lang="ts">
  import { goto } from "$app/navigation";
  import type { ProjectsResponse, RolloutsResponse } from "$lib/pocketbase/generated-types";
  import type { Pexpand, Rexpand } from "$lib/stores/data";
  import { rollouts } from "$lib/stores/data";
  import selectedProjectId from "$lib/stores/project";
  import { formatDateTime, timeAgo } from "$lib/utils/date.utils";
  import { frameworkLogoUrl } from "$lib/utils/framework.utils";
  import { Badge, Button, Indicator, Tooltip } from "flowbite-svelte";
  import { ArrowRight, ChevronRight, Tag } from "lucide-svelte";
  export let project: ProjectsResponse<Pexpand>;

  let tags: Set<string> = new Set();
  if (project.tags) {
    tags = new Set(project.tags.split(","));
  }

  // filter $rollouts by $rollouts.expand.project
  let these_rollouts: RolloutsResponse<Rexpand>[] = [];
  // @ts-ignore
  $: these_rollouts = $rollouts.filter((r) => r.expand?.project.id === project.id);
</script>

<div class="rounded-xl border border-gray-200 ov">
  <div class="flex items-center gap-x-4 border-b border-gray-900/5 p-6">
    <div class="relative">
      <img
        src={frameworkLogoUrl(project.expand?.framework)}
        alt="Tuple"
        class="h-12 w-12 flex-none rounded-lg object-cover ring-1 ring-gray-900/10 p-1"
      />
      <Indicator
        color="dark"
        size="xl"
        placement="top-right"
        class="text-xs font-bold text-white cursor-default"
        >{these_rollouts.length || 0}
        </Indicator
      >
      <Tooltip>Rollouts</Tooltip>
    </div>
    <div class="text-sm font-medium leading-6">{project.name}</div>
    <div class="relative ml-auto">
      <div class="flex justify-end">
        <Button
          color="alternative"
          on:click={() => {
            $selectedProjectId = project.id;
            goto(`/app/projects/${project.id}/overview`);
          }}
        >
          <ArrowRight class="w-5 h-5" />
        </Button>
      </div>
    </div>
  </div>
  <dl class="-my-3 divide-y divide-gray-100 px-6 py-4 text-sm leading-6">
    {#if these_rollouts.length > 0}
      <div class="flex justify-between gap-x-4 py-3">
        <dt class="">Last rollout</dt>
        <dd class=" cursor-default">
          <time datetime={formatDateTime(these_rollouts[0].startDate)}>
            {timeAgo(these_rollouts[0].startDate)}
          </time>
          <Tooltip>{formatDateTime(these_rollouts[0].startDate)}</Tooltip>
        </dd>
      </div>
    {/if}
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
        <dd class="items-start gap-y-2 space-x-2">
          {#each [...tags] as tag (tag)}
            <Badge color="dark" large class="cursor-default"
              >{tag.charAt(0) + tag.slice(1)}</Badge
            >
          {/each}

          <!-- {#each tags as tag (tag)}
            <Badge color={tag.color} large class="cursor-default"
              >{tag.name.charAt(0).toUpperCase() + tag.name.slice(1)}</Badge
            >
          {/each} -->
        </dd>
      </div>
    {/if}
  </dl>
</div>
