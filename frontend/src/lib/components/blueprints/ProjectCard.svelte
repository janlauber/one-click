<script lang="ts">
  import { client } from "$lib/pocketbase";
  import type { BlueprintsRecord, BlueprintsResponse } from "$lib/pocketbase/generated-types";
  import { UpdateFilterEnum, updateDataStores } from "$lib/stores/data";
  import { recordLogoUrl } from "$lib/utils/blueprint.utils";
  import { avatarUrl } from "$lib/utils/user.utils";
  import {
    Button,
    Dropdown,
    DropdownItem,
    Fileupload,
    Input,
    Label,
    Modal,
    Toggle,
    Tooltip
  } from "flowbite-svelte";
  import { DotsVerticalOutline } from "flowbite-svelte-icons";
  import { Clipboard, Code2, FileQuestion, Lock, Share2, Trash, Trash2 } from "lucide-svelte";
  import toast from "svelte-french-toast";
  import MonacoEditor from "svelte-monaco";
  import { getRandomString } from "$lib/utils/random";
  // @ts-ignore
  import yaml from "js-yaml";

  export let blueprint: BlueprintsResponse;
  export let community: boolean = false;

  function jsonToYaml(json: any): string {
    return yaml.dump(json);
  }
  let name: string = blueprint?.name || "";
  let description: string = blueprint?.description || "";
  let isPrivate: boolean = blueprint?.private || false;
  let avatar: string = blueprint?.avatar || "";
  let avatarFile: File;
  let manifest: any = jsonToYaml(blueprint?.manifest) || "";

  let editBlueprintModal = false;
  let confirmShareModal = false;
  let confirmDeleteModal = false;
  let confirmUnsubscribeModal = false;

  async function handleDelete() {
    client
      .collection("blueprints")
      .delete(blueprint.id)
      .then(() => {
        toast.success("Project deleted");
        confirmDeleteModal = false;
        updateDataStores({
          filter: UpdateFilterEnum.ALL
        });
      })
      .catch((error) => {
        toast.error(error.message);
      });
  }

  async function handleUnsubscribe() {
    client
      .collection("blueprints")
      .update(blueprint.id, {
        users: blueprint.users.filter((user) => user !== client.authStore?.model?.id)
      })
      .then(() => {
        toast.success("Unsubscribed from blueprint");
        confirmUnsubscribeModal = false;
        updateDataStores({
          filter: UpdateFilterEnum.ALL
        });
      })
      .catch((error) => {
        toast.error(error.message);
      });
  }

  async function handleSaveManifest() {
    if (!name) {
      toast.error("Name is required");
      return;
    }

    if (!description) {
      toast.error("Description is required");
      return;
    }

    if (!manifest) {
      toast.error("Manifest is required");
      return;
    }

    let formData = new FormData();
    formData.append("avatar", avatarFile);

    // parse the manifest yaml to json
    const parsedManifest = yaml.load(manifest);

    const data: BlueprintsRecord = {
      ...blueprint,
      name,
      description,
      private: isPrivate,
      manifest: parsedManifest
    };

    toast.promise(
      client
        .collection("blueprints")
        .update(blueprint.id, data)
        .then(() => {
          if (avatarFile) {
            client
              .collection("blueprints")
              .update(blueprint.id, formData)
              .then(() => {
                updateDataStores({
                  filter: UpdateFilterEnum.ALL
                });
              })
              .catch((error) => {
                toast.error(error.message);
              });
          }

          updateDataStores({
            filter: UpdateFilterEnum.ALL
          });

          editBlueprintModal = false;
        }),
      {
        loading: "Saving...",
        success: "Blueprint saved",
        error: "Error saving blueprint"
      }
    );
  }
</script>

<div class="rounded-xl border border-gray-200 ov">
  <div class="border-b border-gray-900/5">
    <div class="flex items-center gap-x-4 mx-6 mt-6">
      <div class="relative">
        {#if blueprint}
          <img
            src={recordLogoUrl(blueprint)}
            alt="Tuple"
            class="h-12 w-12 flex-none rounded-lg object-cover ring-1 ring-gray-900/10 p-1"
          />
        {:else}
          <FileQuestion class="h-12 w-12 flex-none rounded-lg object-cover p-1" />
        {/if}
        {#if blueprint.private}
          <Lock class="absolute -bottom-2 -right-2 w-4 h-4 text-gray-500 dark:text-gray-400" />
        {/if}
      </div>
      <div class="text-sm font-medium leading-6">{blueprint.name}</div>
      <DotsVerticalOutline
        class="dots-menu{blueprint.id} dark:text-white ml-auto outline-none cursor-pointer"
      />
      <Dropdown triggeredBy=".dots-menu{blueprint.id}" class="p-0">
        <DropdownItem
          class="w-full text-left"
          on:click={() => {
            editBlueprintModal = true;
          }}
        >
          <Code2 class="w-4 h-4 mr-2 inline-block" />
          {#if community}
            View
          {:else}
            Edit
          {/if}
        </DropdownItem>
        <DropdownItem
          class="w-full text-left"
          on:click={() => {
            confirmShareModal = true;
          }}
        >
          <Share2 class="w-4 h-4 mr-2 inline-block" />
          Share
        </DropdownItem>
        {#if !community}
          <DropdownItem
            slot="footer"
            class="w-full text-left text-red-500 hover:text-red-600"
            on:click={() => {
              confirmDeleteModal = true;
            }}
          >
            <Trash2 class="w-4 h-4 mr-2 inline-block" />
            Delete
          </DropdownItem>
        {:else}
          <DropdownItem
            slot="footer"
            class="w-full text-left text-red-500 hover:text-red-600"
            on:click={() => {
              confirmUnsubscribeModal = true;
            }}
          >
            <Trash class="w-4 h-4 mr-2 inline-block" />
            Unsubscribe
          </DropdownItem>
        {/if}
      </Dropdown>
    </div>
    <div class="text-sm text-gray-500 leading-6 m-6 my-2">{blueprint.description}</div>
  </div>

  {#if !community}
    <dl class="-my-3 divide-y divide-gray-100 px-6 py-4 text-sm leading-6">
      <div class="flex justify-between gap-x-4 py-3">
        <dt class="">Shared</dt>
        <dd class="flex items-start gap-x-2">
          <span class="font-medium">{blueprint.users.length}x</span>
        </dd>
        <Tooltip>
          {#if blueprint.users.length > 0}
            <div class="flex flex-wrap gap-2">
              {#each blueprint.expand?.users as user (user)}
                <div class="flex items-center gap-x-2">
                  {#if user.avatar}
                    <img
                      src={avatarUrl(user)}
                      alt={user.name}
                      class="h-8 w-8 rounded-full object-cover ring-1 ring-gray-900/10"
                    />
                  {:else}
                    <div class="h-8 w-8 rounded-full bg-gray-200 ring-1 ring-gray-900/10"></div>
                  {/if}
                  <div class="text-sm font-medium">{user.name}</div>
                </div>
              {/each}
            </div>
          {:else}
            <div>You haven't shared this blueprint with anyone yet.</div>
          {/if}
        </Tooltip>
      </div>
    </dl>
  {/if}
</div>

<Modal bind:open={editBlueprintModal} size="lg">
  <div class="text-center">
    <Code2 class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      <!-- Refer to https://docs.one-click.dev for advanced documentation about the manifest values -->
      Refer to
      <a href="https://docs.one-click.dev" target="_blank" class="text-primary-500"
        >docs.one-click.dev</a
      >
    </h3>
  </div>
  <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 rounded-lg">
    <div class="col-span-1">
      <Label class="">Name*</Label>
      {#if !community}
        <Input bind:value={name} required />
      {:else}
        <Input bind:value={name} disabled />
      {/if}
    </div>
    <div class="col-span-1">
      <Label class="">Description*</Label>
      {#if !community}
        <Input bind:value={description} required />
      {:else}
        <Input bind:value={description} disabled />
      {/if}
    </div>
    <div />
    {#if !community}
      <div class="flex justify-end">
        <Toggle bind:checked={isPrivate} class="">
          <span class="text-gray-500 dark:text-gray-400">Private</span>
        </Toggle>
      </div>
      <div class="col-span-2">
        <Label class="">Avatar*</Label>
        <Fileupload
          bind:value={avatar}
          on:change={(event) => {
            // @ts-ignore
            avatarFile = event.target.files[0];
          }}
        />
      </div>
    {/if}
  </div>
  <div class=" h-96 overflow-y-auto rounded-lg p-2" style="background-color: #1E1E1E;">
    <MonacoEditor
      bind:value={manifest}
      options={{ language: "yaml", automaticLayout: false, minimap: { enabled: false } }}
      theme="vs-dark"
    />
  </div>
  {#if !community}
    <div class="flex justify-between">
      <Button color="primary" class="me-2" on:click={() => handleSaveManifest()}>Save</Button>
      <Button color="alternative" on:click={() => (editBlueprintModal = false)}>Cancel</Button>
    </div>
  {/if}
</Modal>

{#if !community}
  <Modal bind:open={confirmDeleteModal} size="xs" autoclose>
    <div class="text-center">
      <Trash class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200" />
      <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
        Are you sure you want to delete the <b>{blueprint.name}</b> blueprint?
      </h3>
      <Button color="red" class="me-2" on:click={() => handleDelete()}>Yes, I'm sure</Button>
      <Button color="alternative">No, cancel</Button>
    </div>
  </Modal>
{:else}
  <Modal bind:open={confirmUnsubscribeModal} size="xs" autoclose>
    <div class="text-center">
      <Trash class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200" />
      <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
        Are you sure you want to unsubscribe the <b>{blueprint.name}</b> blueprint?
      </h3>
      <Button color="red" class="me-2" on:click={() => handleUnsubscribe()}>Yes, I'm sure</Button>
      <Button color="alternative">No, cancel</Button>
    </div>
  </Modal>
{/if}

<Modal bind:open={confirmShareModal} size="xs" autoclose>
  <div class="text-center">
    <Share2 class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      Share the following link with your team to give them access to the <b>{blueprint.name}</b> blueprint.
    </h3>
    <div class="flex gap-2 justify-between w-auto">
      <Input
        id={getRandomString(8)}
        size="sm"
        value={window.location.href.split("/").slice(0, 3).join("/") +
          "/app/blueprints/shared/" +
          blueprint.id}
        disabled
      />
      <Button
        color="alternative"
        size="xs"
        class="inline"
        on:click={() => {
          navigator.clipboard.writeText(
            window.location.href.split("/").slice(0, 3).join("/") +
              "/app/blueprints/shared/" +
              blueprint.id
          );
          toast.success("Copied to clipboard.");
        }}
      >
        <Clipboard class="w-4 h-4" />
      </Button>
    </div>
  </div>
</Modal>
