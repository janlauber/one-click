<script lang="ts">
  import { rollouts } from "$lib/stores/data";
  import { formatDateTime, timeAgo } from "$lib/utils/date.utils";
  import { Badge, Indicator, Tooltip } from "flowbite-svelte";

  import { Button, Dropdown, DropdownItem, ToolbarButton, DropdownDivider } from "flowbite-svelte";
  import { DotsHorizontalOutline } from "flowbite-svelte-icons";
</script>

<div class="px-1">
  <div class="mt-8 flow-root">
    <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
      <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
        <div class=" shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
          <table class="min-w-full divide-y divide-gray-300">
            <thead class="bg-gray-50">
              <tr>
                <th
                  scope="col"
                  class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-6"
                  >ID</th
                >
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900"
                  >Image</th
                >
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900"
                  >Created</th
                >
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900"
                  >Ended</th
                >
                <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-6">
                  <span class="sr-only">Edit</span>
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 bg-white">
              {#each $rollouts as rollout (rollout.id)}
                <tr>
                  <td
                    class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6"
                    >{rollout.id}</td
                  ><td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                    {#if rollout.manifest}
                      {rollout.manifest.spec.image.registry}/{rollout.manifest.spec.image
                        .repository}:{rollout.manifest.spec.image.tag}
                    {/if}
                  </td>
                  <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                    <div>{timeAgo(rollout.startDate)}</div>
                    <Tooltip>{formatDateTime(rollout.startDate)}</Tooltip>
                  </td>

                  <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                    <div>{timeAgo(rollout.endDate)}</div>
                    <Tooltip>{formatDateTime(rollout.endDate)}</Tooltip>
                  </td>
                  <td class="relative whitespace-nowrap pr-4 text-right text-sm font-medium">
                    {#if rollout.endDate == ""}
                      <Badge border color="green" class="relative pl-6 mr-2">
                        <Indicator size="sm" color="green" class="absolute left-2" />
                        <Indicator size="sm" color="green" class="absolute animate-ping left-2" />
                        Deployed</Badge
                      >
                    {/if}
                    <DotsHorizontalOutline
                      class="dots-menu dark:text-white inline-block cursor-pointer"
                    />
                    <Dropdown triggeredBy=".dots-menu">
                      <DropdownItem>Deploy</DropdownItem>
                      <DropdownItem class="text-red-500">Delete</DropdownItem>
                    </Dropdown>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</div>
