<script lang="ts">
  import { client } from "$lib/pocketbase";
  import type { RolloutsRecord, RolloutsResponse } from "$lib/pocketbase/generated-types";
  import {
    updateDataStores,
    type Rexpand,
    UpdateFilterEnum,
    currentRollout,
    clusterInfo
  } from "$lib/stores/data";
  import { Accordion, AccordionItem, Button, Input, Label, Select, Toggle } from "flowbite-svelte";
  import selectedProjectId from "$lib/stores/project";
  import toast from "svelte-french-toast";

  export let modal: boolean;

  interface Volume {
    id: string;
    name: string;
    mountPath: string;
    size: string;
    storageClass: string;
  }

  let volume: Volume = {
    id: "",
    name: "",
    mountPath: "",
    size: "1Gi",
    storageClass: ""
  };

  function validateSize(size: string) {
    // validate size Gi or Mi
    const regex = new RegExp("^[0-9]+(Gi|Mi)$");
    return regex.test(size);
  }

  async function handleCreateVolume() {
    if (!$currentRollout) {
      toast.error("No rollout selected");
      return;
    }

    if (!volume.name) {
      toast.error("Volume name is required");
      return;
    }

    if (!volume.mountPath) {
      toast.error("Mount path is required");
      return;
    }

    if (!volume.size) {
      toast.error("Size is required");
      return;
    }

    if (!validateSize(volume.size)) {
      toast.error("Size must be in Gi or Mi");
      return;
    }

    if (!volume.storageClass) {
      toast.error("Storage class is required");
      return;
    }

    // Check for existing volume with same name, host, or port
    // @ts-ignore
    if ($currentRollout.manifest.volumes) {
      // @ts-ignore
      const existingVolume = $currentRollout.manifest.volumes.find(
        (v: Volume) => v.name === volume.name
      );
      if (existingVolume) {
        toast.error("Volume with same name already exists");
        return;
      }
    }

    let new_manifest: any = {
      ...$currentRollout.manifest,
      spec: {
        // @ts-ignore
        ...$currentRollout.manifest.spec,
        volumes: [
          // @ts-ignore
          ...$currentRollout.manifest.spec.volumes,
          {
            name: volume.name,
            mountPath: volume.mountPath,
            size: volume.size,
            storageClass: volume.storageClass
          }
        ]
      }
    };

    const data: RolloutsRecord = {
      manifest: new_manifest,
      startDate: $currentRollout.startDate,
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
          modal = false;
        }),
      {
        loading: "Creating rollout...",
        success: "Rollout created.",
        error: "Error creating rollout."
      }
    );
  }
</script>

<div class="flex flex-col space-y-6">
  <h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">Create a new volume</h3>
  <Label class="space-y-2">
    <span>Volume name *</span>
    <Input type="text" name="name" placeholder="data" size="sm" required bind:value={volume.name} />
  </Label>
  <Label class="space-y-2">
    <span>Mount path *</span>
    <Input
      type="text"
      name="mountPath"
      placeholder="/data"
      size="sm"
      required
      bind:value={volume.mountPath}
    />
  </Label>
  <Label class="space-y-2">
    <span>Size *</span>
    <Input type="text" name="size" placeholder="1Gi" size="sm" required bind:value={volume.size} />
  </Label>
  <Label class="space-y-2">
    <span>Storage Class *</span>
    <Select name="storageClass" size="sm" required bind:value={volume.storageClass}>
      {#if !$clusterInfo}
        <option value="">No storage class available</option>
      {:else}
        {#each $clusterInfo.storageClasses as storageClass}
          <option value={storageClass}>{storageClass}</option>
        {/each}
      {/if}
    </Select>
  </Label>
  <Button type="submit" class="w-full1" color="primary" on:click={handleCreateVolume}>
    Create volume
  </Button>
</div>
