<script lang="ts">
  import { client } from "$lib/pocketbase";
  import type { RolloutsRecord } from "$lib/pocketbase/generated-types";
  import { updateDataStores, UpdateFilterEnum, currentRollout } from "$lib/stores/data";
  import { Button, Input, Label, Toggle } from "flowbite-svelte";
  import selectedProjectId from "$lib/stores/project";
  import toast from "svelte-french-toast";

  export let modal: boolean;

  interface Interface {
    id: string;
    name: string;
    port: number;
    host: string;
    path: string;
    tls: boolean;
  }

  let inf: Interface = {
    id: "",
    name: "",
    port: 80,
    host: "",
    path: "",
    tls: false
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

    if (!inf.port) {
      toast.error("Port is required");
      return;
    }

    // Check for existing interface with same name, host, or port
    // @ts-ignore
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

    if (inf.host) {
      new_manifest = {
        ...$currentRollout.manifest,
        spec: {
          // @ts-ignore
          ...$currentRollout.manifest.spec,
          interfaces: [
            // @ts-ignore
            ...$currentRollout.manifest.spec.interfaces,
            {
              name: inf.name,
              port: parseInt(inf.port.toString()),
              ingress: {
                ingressClass: "nginx",
                annotations: {
                  "nginx.ingress.kubernetes.io/ssl-redirect": "false"
                },
                rules: [
                  {
                    host: inf.host,
                    path: inf.path,
                    tls: inf.tls
                  }
                ]
              }
            }
          ]
        }
      };
    } else {
      new_manifest = {
        ...$currentRollout.manifest,
        spec: {
          // @ts-ignore
          ...$currentRollout.manifest.spec,
          interfaces: [
            // @ts-ignore
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
  <h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">Create a new interface</h3>
  <Label class="space-y-2">
    <span>Interface name *</span>
    <Input type="text" name="name" placeholder="http" size="sm" required bind:value={inf.name} />
  </Label>
  <Label class="space-y-2">
    <span>Port *</span>
    <Input type="number" name="port" placeholder="8080" size="sm" required bind:value={inf.port} />
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

  <Button type="submit" class="w-full1" color="primary" on:click={handleCreateInterface}>
    Create interface
  </Button>
</div>
