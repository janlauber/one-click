<script lang="ts">
  import { page } from "$app/stores";
  import { client } from "$lib/pocketbase";
  import type { RolloutsRecord, RolloutsResponse } from "$lib/pocketbase/generated-types";
  import {
    type Rexpand,
    rollouts,
    updateDataStores,
    UpdateFilterEnum,
    autoUpdates
  } from "$lib/stores/data";
  import { Button, Heading, Input, Label, P, Select, Toggle } from "flowbite-svelte";
  import selectedProjectId from "$lib/stores/project";
  import toast from "svelte-french-toast";
  import { goto } from "$app/navigation";
  import { ArrowRight, Clipboard } from "lucide-svelte";
  import { getRandomString } from "$lib/utils/random";

  let current_rollout: RolloutsResponse<Rexpand> | undefined;
  let lastUpdatedRollout: RolloutsResponse<Rexpand> | undefined;
  let hadRolloutsOnLastPage: boolean = false;

  let registry: string = "";
  let username: string = "";
  let password: string = "";
  let repository: string = "";
  let tag: string = "";
  let tagAutoUpdateEnabled: boolean = false;
  let tagAutoUpdateWebhookPath: string = "";
  let tagAutoUpdateInterval: string = "5m";
  let tagAutoUpdatePattern: string = "^\\d+\\.\\d+\\.\\d+$";
  let tagAutoUpdatePolicy: string = "semver";

  let selectIntervals = [
    { value: "1m", name: "1 Minute" },
    { value: "5m", name: "5 Minutes" },
    { value: "10m", name: "10 Minutes" }
  ];

  let initialLoad: boolean = true;

  $: if ($autoUpdates.length > 0 && initialLoad) {
    // set the first autoUpdate for tags
    tagAutoUpdateEnabled = true;
    tagAutoUpdateWebhookPath = "/auto-update/" + $autoUpdates[0].id;
    tagAutoUpdateInterval = $autoUpdates[0].interval;
    tagAutoUpdatePattern = $autoUpdates[0].pattern;
    tagAutoUpdatePolicy = $autoUpdates[0].policy;
    initialLoad = false;
  }

  $: if ($rollouts.length > 0) {
    // get the current rollout on following priority:
    // 1. no endDate set
    // 2. newest endDate

    const temp_rollouts = $rollouts.filter((r) => !r.endDate);
    if (temp_rollouts.length > 0) {
      current_rollout = temp_rollouts[0];
    } else {
      current_rollout = $rollouts.sort((a, b) => {
        if (a.endDate && b.endDate) {
          return b.endDate.localeCompare(a.endDate);
        } else if (a.endDate) {
          return 1;
        } else if (b.endDate) {
          return -1;
        } else {
          return 0;
        }
      })[0];
    }

    if (current_rollout && current_rollout !== lastUpdatedRollout) {
      registry = current_rollout.manifest?.spec.image.registry ?? "";
      username = current_rollout.manifest?.spec.image.username ?? "";
      password = current_rollout.manifest?.spec.image.password ?? "";
      repository = current_rollout.manifest?.spec.image.repository ?? "";
      tag = current_rollout.manifest?.spec.image.tag ?? "";
      lastUpdatedRollout = current_rollout;
    }

    hadRolloutsOnLastPage = true;
  } else {
    // Reset all values when there are no rollouts
    if (hadRolloutsOnLastPage) {
      resetValues();
    }
    current_rollout = undefined;
    hadRolloutsOnLastPage = false;
  }

  // Reactive statement to track page changes
  $: $page,
    () => {
      if (!hadRolloutsOnLastPage) {
        resetValues();
      }
    };

  function resetValues() {
    registry = "";
    username = "";
    password = "";
    repository = "";
    tag = "";
  }

  function resetInput() {
    registry = current_rollout?.manifest?.spec.image.registry ?? "";
    username = current_rollout?.manifest?.spec.image.username ?? "";
    password = current_rollout?.manifest?.spec.image.password ?? "";
    repository = current_rollout?.manifest?.spec.image.repository ?? "";
    tag = current_rollout?.manifest?.spec.image.tag ?? "";
  }

  function handleInputChange(event: any, field: any) {
    switch (field) {
      case "registry":
        registry = event.target.value;
        break;
      case "username":
        username = event.target.value;
        break;
      case "password":
        password = event.target.value;
        break;
      case "repository":
        repository = event.target.value;
        break;
      case "tag":
        tag = event.target.value;
        break;
    }
  }

  async function handleInputSave() {
    if (current_rollout) {
      if (tagAutoUpdateEnabled) {
        const data = {
          interval: tagAutoUpdateInterval,
          pattern: tagAutoUpdatePattern,
          policy: tagAutoUpdatePolicy,
          project: $selectedProjectId,
          user: client.authStore.model?.id
        };

        // only if $autoUpdates is empty, create a new autoUpdate
        if ($autoUpdates.length === 0) {
          toast.promise(
            client
              .collection("autoUpdates")
              .create(data)
              .then(() => {
                updateDataStores({
                  filter: UpdateFilterEnum.ALL,
                  projectId: $selectedProjectId
                });
              }),
            {
              loading: "Creating auto update...",
              success: "Auto update created.",
              error: "Error creating auto update."
            }
          );
        } else {
          if ($autoUpdates[0]) {
            toast.promise(
              client
                .collection("autoUpdates")
                .update($autoUpdates[0].id, data)
                .then(() => {
                  updateDataStores({
                    filter: UpdateFilterEnum.ALL,
                    projectId: $selectedProjectId
                  });
                }),
              {
                loading: "Updating auto update...",
                success: "Auto update updated.",
                error: "Error updating auto update."
              }
            );
          }
        }
      } else {
        if ($autoUpdates[0]) {
          toast.promise(
            client
              .collection("autoUpdates")
              .delete($autoUpdates[0].id)
              .then(() => {
                updateDataStores({
                  filter: UpdateFilterEnum.ALL,
                  projectId: $selectedProjectId
                });
              }),
            {
              loading: "Deleting auto update...",
              success: "Auto update deleted.",
              error: "Error deleting auto update."
            }
          );
        }
      }

      const new_manifest = {
        ...current_rollout.manifest,
        spec: {
          // @ts-ignore
          ...current_rollout.manifest.spec,
          image: {
            registry: registry,
            repository: repository,
            tag: tag,
            username: username,
            password: password
          }
        }
      };

      // check if the manifest has changed
      if (JSON.stringify(current_rollout.manifest) === JSON.stringify(new_manifest)) {
        return;
      }

      const data: RolloutsRecord = {
        manifest: new_manifest,
        startDate: current_rollout.startDate,
        endDate: "",
        project: $selectedProjectId,
        user: client.authStore.model?.id
      };

      toast.promise(
        client
          .collection("rollouts")
          .create(data)
          .then(() => {
            updateDataStores({
              filter: UpdateFilterEnum.ALL,
              projectId: $selectedProjectId
            });
          }),
        {
          loading: "Creating rollout...",
          success: "Rollout created.",
          error: "Error creating rollout."
        }
      );
    }
  }
</script>

<div class="dark:text-white">
  <div class="mt-8 flow-root">
    <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
      <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
        <div class="p-0.5 shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
          <table class="min-w-full divide-y divide-gray-300 dark:divide-gray-600">
            <tbody
              class="divide-y divide-gray-200 dark:divide-gray-600 bg-white dark:bg-transparent"
            >
              <!-- List the current rollout id -->
              <tr class="transition-all hover:bg-gray-50 dark:hover:bg-gray-800">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-xs font-medium sm:pl-6">
                  <Heading tag="h5">Current Rollout</Heading>
                  <P class="text-gray-500 dark:text-gray-400 text-xs">The current rollout id.</P>
                </td>
                <td class="whitespace-nowrap px-3 py-4 text-xs">
                  <Label for="registry" class="block mb-2">ID</Label>
                  <div class="flex gap-2 justify-between w-auto">
                    <Input id={getRandomString(8)} size="sm" value={current_rollout?.id} disabled />
                    <Button
                      color="alternative"
                      size="xs"
                      class="inline"
                      on:click={() => {
                        navigator.clipboard.writeText(current_rollout?.id ?? "");
                        toast.success("Copied to clipboard.");
                      }}
                    >
                      <Clipboard class="w-4 h-4" />
                    </Button>
                    <Button
                      color="alternative"
                      size="xs"
                      class="inline"
                      on:click={() => {
                        $selectedProjectId = $selectedProjectId;
                        goto(`/app/projects/${$selectedProjectId}/rollouts`);
                      }}
                    >
                      <ArrowRight class="w-4 h-4" />
                    </Button>
                  </div>
                </td>
              </tr>

              <tr class="transition-all hover:bg-gray-50 dark:hover:bg-gray-800">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-xs font-medium sm:pl-6">
                  <Heading tag="h5">Image Registry</Heading>
                  <P class="text-gray-500 dark:text-gray-400 text-xs">
                    Define your image registry.
                  </P>
                </td><td class="whitespace-nowrap px-3 py-4 text-xs">
                  <Label for="registry" class="block mb-2"
                    >Registry
                    <span
                      class="
                      {registry === '' ? 'text-red-500' : 'text-green-500'}
                    ">*</span
                    >
                  </Label>
                  <Input
                    id={getRandomString(8)}
                    size="sm"
                    bind:value={registry}
                    on:input={(e) => handleInputChange(e, "registry")}
                    placeholder="docker.io"
                    class="
                    {registry === '' ? 'border-red-500' : 'border-green-500'}
                    "
                  />
                  <Label class="block mb-2 mt-4">Username</Label>
                  <Input
                    id={getRandomString(8)}
                    size="sm"
                    bind:value={username}
                    on:input={(e) => handleInputChange(e, "username")}
                    placeholder="username"
                  />
                  <Label class="block mb-2 mt-4">Password</Label>
                  <Input
                    id={getRandomString(8)}
                    type="password"
                    size="sm"
                    bind:value={password}
                    on:input={(e) => handleInputChange(e, "password")}
                    placeholder="password"
                  />
                </td>
              </tr>
              <tr class="transition-all hover:bg-gray-50 dark:hover:bg-gray-800">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-xs font-medium sm:pl-6">
                  <Heading tag="h5">Image Repository</Heading>
                  <P class="text-gray-500 dark:text-gray-400 text-xs">
                    Define your image repository.
                  </P>
                </td><td class="whitespace-nowrap px-3 py-4 text-xs">
                  <Label for="repository" class="block mb-2"
                    >Repository
                    <span
                      class="
                      {repository === '' ? 'text-red-500' : 'text-green-500'}
                    ">*</span
                    ></Label
                  >
                  <Input
                    id={getRandomString(8)}
                    size="sm"
                    bind:value={repository}
                    on:input={(e) => handleInputChange(e, "repository")}
                    placeholder="nginx"
                    class="
                    {repository === '' ? 'border-red-500' : 'border-green-500'}
                    "
                  />
                </td>
              </tr>
              <tr class="transition-all hover:bg-gray-50 dark:hover:bg-gray-800">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-xs font-medium sm:pl-6">
                  <Heading tag="h5">Image Tag</Heading>
                  <P class="text-gray-500 dark:text-gray-400 text-xs">Define your image tag.</P>
                </td><td class="whitespace-nowrap px-3 py-4 text-xs">
                  <Label for="tagAutoUpdateEnabled" class="block mb-2">Auto Update</Label>
                  <Toggle bind:checked={tagAutoUpdateEnabled} id="tagAutoUpdateEnabled" />

                  {#if tagAutoUpdateEnabled}
                    <Label for="" class="block mb-2 mt-4">Webhook Path</Label>
                    <div class="flex gap-2 justify-between w-auto">
                      <Input
                        id={getRandomString(8)}
                        size="sm"
                        bind:value={tagAutoUpdateWebhookPath}
                        on:input={(e) => handleInputChange(e, "tagAutoUpdateWebhookPath")}
                        placeholder="known after creation"
                        disabled
                      />
                      <Button
                        color="alternative"
                        size="xs"
                        class="inline"
                        on:click={() => {
                          let url_parts = window.location.href.split("/");
                          let url = url_parts[0] + "//" + url_parts[2];
                          // if url contains localhost, then url is http://localhost:8090
                          if (url.includes("localhost")) {
                            url = "http://localhost:8090"
                          }
                          navigator.clipboard.writeText(url + tagAutoUpdateWebhookPath ?? "");
                          toast.success("Copied to clipboard.");
                        }}
                      >
                        <Clipboard class="w-4 h-4" />
                      </Button>
                    </div>

                    <Label for="tagAutoUpdateInterval" class="block mb-2 mt-4">Interval</Label>
                    <Select bind:value={tagAutoUpdateInterval} id="tagAutoUpdateInterval" size="sm">
                      {#each selectIntervals as interval}
                        <option value={interval.value}>{interval.name}</option>
                      {/each}
                    </Select>
                    <Label class="block mb-2 mt-4">Pattern</Label>
                    <Input
                      id={getRandomString(8)}
                      size="sm"
                      bind:value={tagAutoUpdatePattern}
                      on:input={(e) => handleInputChange(e, "tagAutoUpdatePattern")}
                      placeholder="^\d+\.\d+\.\d+$"
                    />
                    <Label for="tagAutoUpdatePolicy" class="block mb-2 mt-4">Policy</Label>
                    <Input
                      id={getRandomString(8)}
                      size="sm"
                      bind:value={tagAutoUpdatePolicy}
                      on:input={(e) => handleInputChange(e, "tagAutoUpdatePolicy")}
                      placeholder="semver"
                    />
                  {/if}

                  <Label for="tag" class="block mb-2 mt-4"
                    >Tag
                    <span
                      class="
                      {tag === '' ? 'text-red-500' : 'text-green-500'}
                    ">*</span
                    ></Label
                  >
                  <Input
                    id={getRandomString(8)}
                    size="sm"
                    bind:value={tag}
                    on:input={(e) => handleInputChange(e, "tag")}
                    placeholder="latest"
                    class="
                    {tag === '' && !tagAutoUpdateEnabled ? 'border-red-500' : 'border-green-500'}
                    "
                    disabled={tagAutoUpdateEnabled}
                  />
                </td>
              </tr>
            </tbody>
          </table>
          <!-- Reset & Save Button bottom right -->
          <div class="flex justify-end mt-4 p-4">
            <Button
              color="alternative"
              class="whitespace-nowrap self-start mr-2"
              on:click={() => resetInput()}
            >
              Reset
            </Button>
            <Button
              color="primary"
              class="whitespace-nowrap self-start"
              on:click={() => handleInputSave()}
            >
              Save
            </Button>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
