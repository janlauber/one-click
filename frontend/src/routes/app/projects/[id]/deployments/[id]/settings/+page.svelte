<script lang="ts">
  import { goto } from "$app/navigation";
  import { client } from "$lib/pocketbase";
  import type {
    BlueprintsRecord,
    DeploymentsRecord,
    DeploymentsResponse,
    RolloutsRecord
  } from "$lib/pocketbase/generated-types";
  import {
    UpdateFilterEnum,
    currentRollout,
    deployments,
    selectedProject,
    updateDataStores
  } from "$lib/stores/data";
  import { Button, Fileupload, Heading, Input, Label, Modal, P } from "flowbite-svelte";
  import { BookDashed, Code2, Image, Trash } from "lucide-svelte";
  import toast from "svelte-french-toast";
  import MonacoEditor from "svelte-monaco";
  // @ts-expect-error - no types available
  import yaml from "js-yaml";
  import selectedDeploymentId from "$lib/stores/deployment";
  import selectedProjectId from "$lib/stores/project";

  let initialLoad = true;
  let deploymentName: string = "";
  let localSelectedDeployment: DeploymentsResponse | undefined = undefined;
  let inFocus = false;

  let modalAdvancedOpen = false;
  let modalBluprintOpen = false;
  let modalDeleteOpen = false;
  let avatar: File;

  let blueprintName: string = $selectedProject?.name || "";
  let blueprintDescription: string = $selectedProject?.description || "";
  let blueprintAvatar: string = $selectedProject?.avatar || "";
  let blueprintAvatarFile: File;
  let blueprintManifest: any = jsonToYaml($currentRollout?.manifest) || "";
  let advancedManifest: any = jsonToYaml($currentRollout?.manifest) || "";

  function jsonToYaml(json: any): string {
    return yaml.dump(json);
  }

  $: (inFocus || !inFocus) && saveName();

  $: {
    if (initialLoad) {
      localSelectedDeployment = $deployments.find((d) => d.id === $selectedDeploymentId);
      if (localSelectedDeployment) deploymentName = localSelectedDeployment.name;

      initialLoad = false;
    }
  }

  async function saveName() {
    // check if the name changed and input is not in focus
    if (deploymentName === localSelectedDeployment?.name || inFocus) return;

    if (!localSelectedDeployment) return;

    const deployment: DeploymentsRecord = {
      ...localSelectedDeployment,
      name: deploymentName
    };

    client
      .collection("deployments")
      .update($selectedDeploymentId, deployment)
      .then(() => {
        // update the selected project
        if ($currentRollout) {
          updateDataStores({
            filter: UpdateFilterEnum.ALL,
            projectId: $selectedProject?.id,
            deploymentId: $selectedDeploymentId
          });
        }
        toast.success("Name updated");
      })
      .catch((error) => {
        toast.error(error.message);
      });
  }

  async function handleDelete() {
    if (!$selectedDeploymentId) return;

    client
      .collection("deployments")
      .delete($selectedDeploymentId)
      .then(() => {
        toast.success("Deployment deleted");
        modalDeleteOpen = false;
        updateDataStores({
          filter: UpdateFilterEnum.ALL,
          projectId: $selectedProjectId
        });
        goto(`/app/projects/${$selectedProject?.id}`);
      })
      .catch((error) => {
        toast.error(error.message);
      });
  }

  async function handleAvatarUpload(event: any) {
    if (!event.target.files[0]) return;
    if (!$selectedDeploymentId) return;

    avatar = event.target.files[0];

    let formData = new FormData();

    formData.append("avatar", avatar);

    client
      .collection("deployments")
      .update($selectedDeploymentId, formData)
      .then(() => {
        toast.success("Avatar updated");
        updateDataStores({
          filter: UpdateFilterEnum.ALL,
          projectId: $selectedProject?.id,
          deploymentId: $currentRollout?.deployment
        });
      })
      .catch((error) => {
        toast.error(error.message);
      });
  }

  async function handleCreateBlueprint() {
    if (!$selectedProject) return;

    if (!blueprintName) {
      toast.error("Blueprint name is required");
      return;
    }

    if (!blueprintDescription) {
      toast.error("Blueprint description is required");
      return;
    }

    if (!blueprintAvatar) {
      toast.error("Blueprint avatar is required");
      return;
    }

    if (!blueprintManifest) {
      toast.error("Blueprint manifest is required");
      return;
    }

    let formData = new FormData();
    formData.append("avatar", blueprintAvatarFile);

    // parse the manifest yaml to json
    const parsedManifest = yaml.load(blueprintManifest);

    let data: BlueprintsRecord = {
      name: blueprintName,
      description: blueprintDescription,
      manifest: parsedManifest,
      owner: client.authStore?.model?.id
    };

    client
      .collection("blueprints")
      .create(data)
      .then((response) => {
        client
          .collection("blueprints")
          .update(response?.id ?? "", formData)
          .catch((error) => {
            toast.error(error.message);
          });

        toast.success("Blueprint created");

        goto(`/app/blueprints/my-blueprints`);
      })
      .catch((error) => {
        toast.error(error.message);
      })
      .finally(() => {
        // update the selected project
        if ($currentRollout) {
          updateDataStores({
            filter: UpdateFilterEnum.ALL,
            projectId: $currentRollout.project
          });
        }
        modalBluprintOpen = false;
      });
  }

  async function handleSaveManifest() {
    if (!$selectedProject) return;

    if (!advancedManifest) {
      toast.error("Manifest is required");
      return;
    }

    // TODO: Validate the manifest

    if (!$currentRollout) return;
    if (!$selectedDeploymentId) return;

    // parse the manifest yaml to json
    const parsedManifest = yaml.load(advancedManifest);

    const data: RolloutsRecord = {
      manifest: parsedManifest,
      startDate: $currentRollout.startDate,
      endDate: "",
      project: $selectedProject.id,
      deployment: $currentRollout.deployment,
      user: client.authStore.model?.id
    };

    toast.promise(
      client
        .collection("rollouts")
        .create(data)
        .then(() => {
          updateDataStores({
            filter: UpdateFilterEnum.ALL,
            projectId: $selectedProject?.id ?? "",
            deploymentId: $currentRollout?.deployment
          });
          modalAdvancedOpen = false;
        }),
      {
        loading: "Creating rollout...",
        success: "Rollout created.",
        error: "Error creating rollout."
      }
    );
  }
</script>

<div class="flex items-start justify-between">
  <div class="flex flex-col">
    <Heading tag="h2">Settings</Heading>
    <P class="text-gray-500 dark:text-gray-400 text-sm">Settings of your deployment.</P>
  </div>
</div>

<div class="mt-8 space-y-4">
  <Label class="space-y-2">
    <span>Deployment name</span>
    <div class="flex whitespace-nowrap gap-2">
      <Input
        id="name"
        type="text"
        name="name"
        size="sm"
        placeholder="Enter the name of your deployment"
        bind:value={deploymentName}
        on:focus={() => {
          inFocus = true;
        }}
        on:blur={() => {
          inFocus = false;
        }}
      />
    </div>
  </Label>
  <div>
    <Label class="pb-2">Change Avatar</Label>
    <label
      for="avatar"
      class="cursor-pointer flex justify-center items-center w-full px-4 py-2 border-primary-700 dark:border-gray-400 border-2 rounded-lg text-sm font-medium text-primary-700 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800"
    >
      <Image class="w-6 h-6 inline-block mr-1" />
      Choose a file (max 1MB)
      <input
        type="file"
        name="avatar"
        id="avatar"
        class="w-full border-gray-300 border-2"
        on:change={handleAvatarUpload}
      />
    </label>
  </div>
</div>

<!-- Advanced Editing -->

<div class="mt-4 p-0.5 shadow ring-1 ring-black ring-opacity-5 rounded-lg bg-white">
  <table class="min-w-full divide-y divide-gray-300 dark:divide-gray-600">
    <tbody class="divide-y divide-gray-200 dark:divide-gray-600 dark:bg-transparent">
      <tr class="">
        <td class="whitespace-nowrap py-4 pl-4 pr-3 text-xs font-medium sm:pl-6">
          <Heading tag="h5" color="text-black">Advanced Editing</Heading>
          <P color="text-black" class="text-xs">Edit the manifest of your deployment.</P>
        </td>
        <td class="whitespace-nowrap px-3 py-4 text-xs text-right">
          <!-- Modified: Added 'text-right' class -->
          <Button
            color="light"
            size="xs"
            class="whitespace-nowrap"
            on:click={() => {
              if (!$selectedProject) return;

              modalAdvancedOpen = true;
            }}
          >
            <Code2 class="w-4 h-4 inline-block mr-1" />
            Edit manifest
          </Button>
        </td>
      </tr>
    </tbody>
  </table>
</div>

<!-- Create a blueprint -->

<div class="mt-4 p-0.5 shadow ring-1 ring-black ring-opacity-5 rounded-lg bg-primary-500">
  <table class="min-w-full divide-y divide-gray-300 dark:divide-gray-600">
    <tbody class="divide-y divide-gray-200 dark:divide-gray-600 dark:bg-transparent">
      <tr class="">
        <td class="whitespace-nowrap py-4 pl-4 pr-3 text-xs font-medium sm:pl-6">
          <Heading tag="h5" color="text-white">Blueprint</Heading>
          <P color="text-white" class="text-xs">Create a blueprint from your deployment.</P>
        </td>
        <td class="whitespace-nowrap px-3 py-4 text-xs text-right">
          <!-- Modified: Added 'text-right' class -->
          <Button
            color="light"
            size="xs"
            class="whitespace-nowrap"
            on:click={() => {
              if (!$selectedProject) return;

              modalBluprintOpen = true;
            }}
          >
            <BookDashed class="w-4 h-4 inline-block mr-1" />
            New blueprint
          </Button>
        </td>
      </tr>
    </tbody>
  </table>
</div>

<!-- Danger Zone -> Delete Deployment -->

<div
  class="mt-4 p-0.5 shadow ring-1 ring-black ring-opacity-5 rounded-lg bg-red-100 dark:bg-red-800"
>
  <table class="min-w-full divide-y divide-gray-300 dark:divide-gray-600">
    <tbody class="divide-y divide-gray-200 dark:divide-gray-600 dark:bg-transparent">
      <tr class="">
        <td class="whitespace-nowrap py-4 pl-4 pr-3 text-xs font-medium sm:pl-6">
          <Heading tag="h5" color="text-red-600 dark:text-red-100">Danger Zone</Heading>
          <P color="text-red-600 dark:text-red-100" class="text-xs">Delete your deployment.</P>
        </td>
        <td class="whitespace-nowrap px-3 py-4 text-xs text-right">
          <!-- Modified: Added 'text-right' class -->
          <Button
            color="red"
            size="xs"
            class="whitespace-nowrap"
            on:click={() => {
              if (!$selectedProject) return;

              modalDeleteOpen = true;
            }}
          >
            <Trash class="w-4 h-4 inline-block mr-1" />
            Delete deployment
          </Button>
        </td>
      </tr>
    </tbody>
  </table>
</div>

<Modal bind:open={modalAdvancedOpen} size="lg">
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
  <div class=" h-96 overflow-y-auto rounded-lg p-2" style="background-color: #1E1E1E;">
    <MonacoEditor
      bind:value={advancedManifest}
      options={{ language: "yaml", automaticLayout: true, minimap: { enabled: false } }}
      theme="vs-dark"
    />
  </div>
  <div class="flex justify-between">
    <Button color="primary" class="me-2" on:click={() => handleSaveManifest()}>Save</Button>
    <Button color="alternative" on:click={() => (modalAdvancedOpen = false)}>Cancel</Button>
  </div>
</Modal>

<Modal bind:open={modalDeleteOpen} size="xs" autoclose>
  <div class="text-center">
    <Trash class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      Are you sure you want to delete this deployment?
    </h3>
    <Button color="red" class="me-2" on:click={() => handleDelete()}>Yes, I'm sure</Button>
    <Button color="alternative">No, cancel</Button>
  </div>
</Modal>

<Modal bind:open={modalBluprintOpen} size="lg">
  <div class="text-center">
    <BookDashed class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      Create a blueprint from this deployment
    </h3>
  </div>
  <div class="space-y-2">
    <Label class="">Name*</Label>
    <Input bind:value={blueprintName} required />
    <Label class="">Description*</Label>
    <Input bind:value={blueprintDescription} required />
    <Label class="">Avatar*</Label>
    <Fileupload
      bind:value={blueprintAvatar}
      on:change={(event) => {
        // @ts-expect-error - event.target.files is a FileList
        blueprintAvatarFile = event.target.files[0];
      }}
    />

    <Label class="">Manifest</Label>
    <span class="text-xs text-gray-500 dark:text-gray-400">
      The manifest is a YAML file that describes the blueprint. <br />
      On project creation, the ingress resources won't be created.
    </span>
    <div class="h-64 overflow-y-auto rounded-lg p-2" style="background-color: #1E1E1E;">
      <MonacoEditor
        bind:value={blueprintManifest}
        options={{ language: "yaml", automaticLayout: true, minimap: { enabled: false } }}
        theme="vs-dark"
      />
    </div>
  </div>
  <div class="flex justify-between">
    <Button color="primary" class="me-2" on:click={() => handleCreateBlueprint()}
      >Yes, I'm sure</Button
    >
    <Button
      color="alternative"
      on:click={() => {
        modalBluprintOpen = false;
      }}>No, cancel</Button
    >
  </div>
</Modal>

<style>
  input[type="file"] {
    display: none;
  }
</style>
