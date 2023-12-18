<script lang="ts">
  import { page } from "$app/stores";
  import { client } from "$lib/pocketbase";
  import type { RolloutsRecord, RolloutsResponse } from "$lib/pocketbase/generated-types";
  import { type Rexpand, rollouts, updateDataStores, UpdateFilterEnum } from "$lib/stores/data";
  import { Button, Heading, Input, Label, P, TableSearch, Toggle } from "flowbite-svelte";
  import selectedProjectId from "$lib/stores/project";
  import toast from "svelte-french-toast";
  import { goto } from "$app/navigation";
  import { ArrowRight, Clipboard } from "lucide-svelte";

  let current_rollout: RolloutsResponse<Rexpand> | undefined;
  let lastUpdatedRollout: RolloutsResponse<Rexpand> | undefined;
  let hadRolloutsOnLastPage: boolean = false;

  let registry: string = "";
  let username: string = "";
  let password: string = "";
  let repository: string = "";
  let tag: string = "";
  let verify: boolean = false;

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
      verify = current_rollout.manifest?.spec.image.verify ?? false;
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
    verify = false;
  }

  function resetInput() {
    registry = current_rollout?.manifest?.spec.image.registry ?? "";
    username = current_rollout?.manifest?.spec.image.username ?? "";
    password = current_rollout?.manifest?.spec.image.password ?? "";
    repository = current_rollout?.manifest?.spec.image.repository ?? "";
    tag = current_rollout?.manifest?.spec.image.tag ?? "";
    verify = current_rollout?.manifest?.spec.image.verify ?? false;
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
      case "verify":
        verify = event.target.checked;
        break;
    }
  }

  async function handleInputSave() {
    if (current_rollout) {
      const new_manifest = {
        ...current_rollout.manifest,
        spec: {
          // @ts-ignore
          ...current_rollout.manifest.spec,
          image: {
            registry: registry,
            username: username,
            password: password,
            repository: repository,
            tag: tag,
            verify: verify
          }
        }
      };

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
                    <Input id="id" size="sm" value={current_rollout?.id} disabled />
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
                    id="registry"
                    size="sm"
                    bind:value={registry}
                    on:input={(e) => handleInputChange(e, "registry")}
                    placeholder="docker.io"
                    class="
                    {registry === '' ? 'border-red-500' : 'border-green-500'}
                    "
                  />
                  <Label for="username" class="block mb-2 mt-4">Username</Label>
                  <Input
                    id="username"
                    size="sm"
                    bind:value={username}
                    on:input={(e) => handleInputChange(e, "username")}
                    placeholder="username"
                  />
                  <Label for="password" class="block mb-2 mt-4">Password</Label>
                  <Input
                    id="password"
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
                    id="repository"
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
                  <Label for="tag" class="block mb-2"
                    >Tag
                    <span
                      class="
                      {tag === '' ? 'text-red-500' : 'text-green-500'}
                    ">*</span
                    ></Label
                  >
                  <Input
                    id="tag"
                    size="sm"
                    bind:value={tag}
                    on:input={(e) => handleInputChange(e, "tag")}
                    placeholder="latest"
                    class="
                    {tag === '' ? 'border-red-500' : 'border-green-500'}
                    "
                  />
                </td>
              </tr>
              <tr class="transition-all hover:bg-gray-50 dark:hover:bg-gray-800">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-xs font-medium sm:pl-6">
                  <Heading tag="h5">Verify Image</Heading>
                  <P class="text-gray-500 dark:text-gray-400 text-xs">
                    Verify your image before deploying.
                  </P>
                </td><td class="whitespace-nowrap px-3 py-4 text-xs">
                  <Toggle bind:checked={verify} id="verify" />
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
