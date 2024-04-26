<script lang="ts">
  import { client } from "$lib/pocketbase";
  import type { RolloutsRecord } from "$lib/pocketbase/generated-types";
  import {
    updateDataStores,
    UpdateFilterEnum,
    currentRollout,
    clusterInfo
  } from "$lib/stores/data";
  import { Button, Input, Label, Select, Toggle } from "flowbite-svelte";
  import selectedProjectId from "$lib/stores/project";
  import toast from "svelte-french-toast";
  import { isValidName } from "$lib/utils/string-validation";

  export let modal: boolean;

  interface Interface {
    id: string;
    name: string;
    port: number;
    ingressClassName?: string;
    host: string;
    path: string;
    tls: boolean;
    tlsSecretName?: string;
  }

  let inf: Interface = {
    id: "",
    name: "",
    port: 80,
    ingressClassName: "",
    host: "",
    path: "/",
    tls: false,
    tlsSecretName: ""
  };

  async function handleCreateInterface() {
    if (!$currentRollout) {
      toast.error("No rollout selected");
      return;
    }

    if (!inf.name) {
      toast.error("Interface name is required");
      return;
    }

    if (!isValidName(inf.name)) {
      toast.error(
        "Interface name should only contain lowercase alphanumeric characters or '-' (max 63 characters)"
      );
      return;
    }

    if (!inf.port) {
      toast.error("Port is required");
      return;
    }

    // Check for existing interface with same name, host, or port
    const existingInterface = $currentRollout.manifest.spec.interfaces.find(
      (i: any) =>
        i.name === inf.name ||
        i.port === inf.port ||
        (i.ingress && i.ingress.rules.some((rule: any) => rule.host === inf.host))
    );

    if (existingInterface) {
      toast.error("An interface with the same name, host, or port already exists");
      return;
    }

    let new_manifest: any;

    if (inf.host && inf.ingressClassName !== "") {
      new_manifest = {
        ...($currentRollout.manifest as any),
        spec: {
          ...$currentRollout.manifest.spec,
          interfaces: [
            ...$currentRollout.manifest.spec.interfaces,
            {
              name: inf.name,
              port: parseInt(inf.port.toString()),
              ingress: {
                ingressClass: inf.ingressClassName,
                annotations: {
                  "nginx.ingress.kubernetes.io/ssl-redirect": "false"
                },
                rules: [
                  {
                    host: inf.host,
                    path: inf.path,
                    tls: inf.tls,
                    tlsSecretName: inf.tls ? inf.tlsSecretName : ""
                  }
                ]
              }
            }
          ]
        }
      };
    } else {
      new_manifest = {
        ...($currentRollout.manifest as any),
        spec: {
          ...$currentRollout.manifest.spec,
          interfaces: [
            ...$currentRollout.manifest.spec.interfaces,
            {
              name: inf.name,
              port: parseInt(inf.port.toString())
            }
          ]
        }
      };
    }

    const data: RolloutsRecord = {
      manifest: new_manifest,
      startDate: $currentRollout.startDate,
      endDate: "",
      project: $selectedProjectId,
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
            projectId: $selectedProjectId,
            deploymentId: $currentRollout?.deployment
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
  <h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-gray-400">Create a new interface</h3>
  <Label class="space-y-2">
    <span>Interface name *</span>
    <Input type="text" name="name" placeholder="http" size="sm" required bind:value={inf.name} />
  </Label>
  <Label class="space-y-2">
    <span>Port *</span>
    <Input type="number" name="port" placeholder="8080" size="sm" required bind:value={inf.port} />
  </Label>
  <Label class="space-y-2">
    <span>Ingress Class</span>
    <Select id="ingressClassName" size="sm" bind:value={inf.ingressClassName} class="">
      {#if !$clusterInfo}
        <option value="">No ingress classes found</option>
      {:else}
        <option value="">None</option>
        {#each $clusterInfo.ingressClasses as ingressClass}
          <option value={ingressClass}>{ingressClass}</option>
        {/each}
      {/if}
    </Select>
  </Label>
  <Label class="space-y-2">
    <span>Host</span>
    <Input type="text" name="host" placeholder="app.example.com" size="sm" bind:value={inf.host} />
  </Label>
  <Label class="space-y-2">
    <span>Path</span>
    <Input type="text" name="path" placeholder="/" size="sm" bind:value={inf.path} />
  </Label>
  <Label class="space-y-2">
    <span>TLS</span>
    <Toggle name="tls" size="small" bind:checked={inf.tls} />
  </Label>

  {#if inf.tls}
    <Label class="space-y-2">
      <span>TLS Secret Name</span>
      <Input
        id="tlsSecretName"
        size="sm"
        type="text"
        bind:value={inf.tlsSecretName}
        placeholder="Enter the TLS secret name"
        class=""
      /></Label
    >
  {/if}

  <Button type="submit" class="w-full1" color="primary" on:click={handleCreateInterface}>
    Create interface
  </Button>
</div>
