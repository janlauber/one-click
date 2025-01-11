<script lang="ts">
  import NewInterface from "$lib/components/networking/NewInterface.svelte";
  import { client } from "$lib/pocketbase";
  import type { RolloutsRecord, RolloutsResponse } from "$lib/pocketbase/generated-types";
  import {
    rollouts,
    currentRollout,
    type Rexpand,
    updateDataStores,
    UpdateFilterEnum,
    clusterInfo
  } from "$lib/stores/data";
  import { isValidName } from "$lib/utils/string-validation";
  import {
    Accordion,
    AccordionItem,
    Badge,
    Button,
    Heading,
    Input,
    Label,
    Modal,
    P,
    Select,
    Toggle
  } from "flowbite-svelte";
  import { CircleAlert, Copy, ExternalLink, Lock, Network, Plus } from "lucide-svelte";
  import toast from "svelte-french-toast";

  export let modal: boolean;
  let deleteModal: boolean = false;
  let selectedInterfaceId: string;

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

  let interfaces: Interface[] = [];

  $: {
    parseManifestsToInterfaces($currentRollout);
  }

  function parseManifestsToInterfaces(rollout: RolloutsResponse<Rexpand> | undefined) {
    interfaces = []; // Clear existing interfaces
    if (rollout?.manifest && rollout.manifest.spec && rollout.manifest.spec.interfaces) {
      rollout?.manifest.spec.interfaces.forEach((i: any, index: number) => {
        const interfaceId = `${rollout?.id}_${index}`; // Combining rollout ID with index for uniqueness

        let newInterface;
        if (i.ingress && Array.isArray(i.ingress.rules) && i.ingress.rules.length > 0) {
          const firstRule = i.ingress.rules[0]; // Taking the first rule as a primary one
          newInterface = {
            id: interfaceId,
            name: i.name,
            port: i.port,
            ingressClassName: i.ingress.ingressClass,
            host: firstRule.host,
            path: firstRule.path,
            tls: firstRule.tls,
            tlsSecretName: firstRule.tlsSecretName
          };
        } else {
          newInterface = {
            id: interfaceId,
            name: i.name,
            port: i.port,
            ingressClassName: "",
            host: "",
            path: "",
            tls: false,
            tlsSecretName: ""
          };
        }
        interfaces.push(newInterface);
      });
    }
  }

  function handleDeleteClick(id: string) {
    selectedInterfaceId = id;
    deleteModal = true;
  }

  async function handleDelete() {
    const interfaceIndex = interfaces.findIndex((inf) => inf.id === selectedInterfaceId);
    if (interfaceIndex === -1) {
      return; // Interface not found, do nothing
    }

    interfaces.splice(interfaceIndex, 1); // Remove the interface from the array

    if (!$currentRollout) {
      toast.error("No rollout selected");
      return;
    }

    if (!$currentRollout.manifest) {
      toast.error("No manifest found");
      return;
    }

    // Remove the interface from the manifest
    $currentRollout.manifest.spec.interfaces.splice(interfaceIndex, 1);

    // Save the manifest
    try {
      await updateManifest($currentRollout.manifest);
      toast.success("Interface deleted successfully");
    } catch (error) {
      toast.error("Failed to delete interface");
      console.error(error);
    }
  }

  async function handleInputSave(id: string) {
    const interfaceIndex = interfaces.findIndex((inf) => inf.id === id);
    if (!$currentRollout) {
      toast.error("No rollout selected");
      return;
    }

    const updatedInterface = interfaces[interfaceIndex];

    if (!updatedInterface.name) {
      toast.error("Interface name is required");
      return;
    }

    if (!isValidName(updatedInterface.name)) {
      toast.error(
        "Interface name should only contain lowercase alphanumeric characters or '-' (max 63 characters)"
      );
      return;
    }

    if (!updatedInterface.port) {
      toast.error("Port is required");
      return;
    }

    // Find the index of the current interface based on its unique identifier (id)
    const currentInterfaceIndex = $currentRollout.manifest.spec.interfaces.findIndex(
      (i: any) => i.id === updatedInterface.id
    );

    if (currentInterfaceIndex !== -1) {
      // Check if there's another interface with the same name, host, or port in the current rollout
      const existingInterface = $currentRollout.manifest.spec.interfaces.find(
        (i: any) =>
          i.id !== updatedInterface.id &&
          (i.name === updatedInterface.name ||
            i.port === updatedInterface.port ||
            (i.ingress &&
              i.ingress.rules.some(
                (rule: any) =>
                  rule.host === updatedInterface.host && rule.path === updatedInterface.path
              )))
      );

      if (existingInterface) {
        toast.error("An interface with the same name, host, or port already exists");
        return;
      }
    }

    // Update interface in $currentRollout
    const rolloutInterfaceIndex = $currentRollout.manifest.spec.interfaces.findIndex(
      // do not only check for name, but also for port and host (so you can update the name of an interface)
      (i: any) => i.name === updatedInterface.name || i.port === updatedInterface.port
    );

    if (rolloutInterfaceIndex !== -1) {
      $currentRollout.manifest.spec.interfaces[rolloutInterfaceIndex] = {
        name: updatedInterface.name,
        port: parseInt(String(updatedInterface.port)),
        ingress: updatedInterface.host
          ? {
              ingressClass:
                updatedInterface.ingressClassName ||
                $currentRollout.manifest?.spec.interfaces[rolloutInterfaceIndex].ingress
                  ?.ingressClass,
              annotations:
                $currentRollout.manifest?.spec.interfaces[rolloutInterfaceIndex].ingress
                  ?.annotations,
              rules: [
                {
                  host: updatedInterface.host,
                  path: updatedInterface.path,
                  tls: updatedInterface.tls,
                  tlsSecretName: updatedInterface.tls ? updatedInterface.tlsSecretName : ""
                }
              ]
            }
          : undefined
      };

      // if ingress is undefined, remove it from the manifest
      if (!$currentRollout.manifest.spec.interfaces[rolloutInterfaceIndex].ingress) {
        delete $currentRollout.manifest.spec.interfaces[rolloutInterfaceIndex].ingress;
      }
    }

    // Validate when the host is set then the ingress class should be set
    if (updatedInterface.host && !updatedInterface.ingressClassName) {
      toast.error("Ingress class is required when host is set");
      return;
    }

    // Save changes to the server
    await updateManifest($currentRollout.manifest);

    toast.success("Interface updated successfully.");
  }

  async function updateManifest(manifest: any) {
    try {
      if (!$currentRollout) {
        toast.error("No rollout selected");
        return;
      }
      const data: RolloutsRecord = {
        manifest: manifest,
        startDate: $currentRollout?.startDate,
        endDate: "",
        project: $currentRollout?.project,
        deployment: $currentRollout?.deployment ?? "",
        user: client.authStore.record?.id ?? ""
      };

      client
        .collection("rollouts")
        .create(data)
        .then(() => {
          updateDataStores({
            filter: UpdateFilterEnum.ALL,
            projectId: $currentRollout?.project,
            deploymentId: $currentRollout?.deployment
          });
        });

      // Update the rollout in the store

      // update the $rollouts store
      rollouts.update((rollouts) => {
        const rolloutIndex = rollouts.findIndex((r) => r.id === $currentRollout?.id);
        if (rolloutIndex !== -1) {
          rollouts[rolloutIndex] = {
            ...rollouts[rolloutIndex],
            manifest: manifest
          };
        }
        return rollouts;
      });

      $currentRollout.manifest = manifest;

      // Update the rollout in the store
    } catch (error) {
      console.error("Failed to update manifest:", error);
      toast.error("Failed to update interface.");
    }
  }
</script>

<div class="flex items-start justify-between"></div>

<div class="flex items-start justify-between">
  <div class="flex flex-col">
    <Heading tag="h2">Network</Heading>
    <P class="text-gray-500 dark:text-gray-400 text-sm">Network interface for your rollout.</P>
  </div>
  <Button
    color="primary"
    size="xs"
    class="whitespace-nowrap self-start"
    on:click={() => {
      modal = true;
    }}
  >
    <Plus class="w-4 h-4 mr-2" />
    New interface
  </Button>
</div>

<Accordion class="grid mt-10 p-1">
  {#key $rollouts}
    {#each interfaces as inf (inf.id)}
      <AccordionItem class="">
        <div slot="header" class="flex">
          <div
            class="ring-1 p-2 rounded-lg ring-gray-500 mr-2 flex items-center justify-center relative"
          >
            <Network class="w-4 h-4" />

            {#if inf.tls}
              <Lock
                class="w-4 h-4 inline-block absolute -right-1 -top-2 bg-white dark:bg-slate-800"
              />
            {/if}
          </div>

          <span class="pt-1 inline">
            {inf.name}
            <span class=" font-normal text-sm">
              :{inf.port}
            </span>
          </span>
          <Badge class="ml-2" color="green">
            DNS:
            {inf.name}-{$currentRollout?.deployment}-svc
            <button
              class="ml-2 hover:scale-105 transition-transform duration-200 ease-in-out"
              on:click={() => {
                navigator.clipboard.writeText(`${inf.name}-${$currentRollout?.deployment}-svc`);
                toast.success("Service name copied to clipboard");
              }}
            >
              <Copy class="w-4 h-4 inline-block" />
            </button>
          </Badge>
          {#if inf.host}
            <a
              href={(inf.tls ? "https://" : "http://") + inf.host}
              target="_blank"
              rel="noopener noreferrer"
              class="ml-2 text-blue-500 hover:underline mt-1"
            >
              {inf.host}
              <ExternalLink class="w-4 h-4 inline-block ml-1" />
            </a>
          {/if}
        </div>
        <div class="">
          <table class="min-w-full divide-y divide-gray-300 dark:divide-gray-600">
            <tbody class="divide-y divide-gray-200 dark:divide-gray-600">
              <tr class="transition-all hover:bg-gray-50 dark:hover:bg-gray-800">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-xs font-medium sm:pl-6">
                  <Heading tag="h5">Details</Heading>
                  <P class="text-gray-500 dark:text-gray-400 text-xs">Details of your interface.</P>
                </td><td class="whitespace-nowrap px-3 py-4 text-xs space-y-2">
                  <Label for="tag" class="block ">Interface name *</Label>
                  <Input
                    id="name"
                    size="sm"
                    type="text"
                    bind:value={inf.name}
                    placeholder="Enter the name of your interface"
                    class=""
                  />
                  <Label for="tag" class="block ">Port *</Label>
                  <Input
                    id="port"
                    size="sm"
                    type="number"
                    bind:value={inf.port}
                    placeholder="8080"
                    class=""
                  />
                </td>
              </tr>
              <tr class="transition-all hover:bg-gray-50 dark:hover:bg-gray-800">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-xs font-medium sm:pl-6">
                  <Heading tag="h5">Ingress</Heading>
                  <P class="text-gray-500 dark:text-gray-400 text-xs">Ingress of your interface.</P>
                </td><td class="whitespace-nowrap px-3 py-4 text-xs space-y-2">
                  <Label for="tag" class="block ">Ingress Class</Label>
                  <Select
                    id="ingressClassName"
                    size="sm"
                    bind:value={inf.ingressClassName}
                    placeholder="Select the ingress class"
                    class=""
                  >
                    {#if !$clusterInfo}
                      <option value="">No ingress classes found</option>
                    {:else}
                      <option value="">None</option>
                      {#each $clusterInfo.ingressClasses as ingressClass}
                        <option value={ingressClass}>{ingressClass}</option>
                      {/each}
                    {/if}
                  </Select>
                  <Label for="tag" class="block ">Host</Label>
                  <Input
                    id="host"
                    size="sm"
                    type="text"
                    bind:value={inf.host}
                    placeholder="Enter the host"
                    class=""
                  />
                  <Label for="tag" class="block ">Path</Label>
                  <Input
                    id="path"
                    size="sm"
                    type="text"
                    bind:value={inf.path}
                    placeholder="Enter the path"
                    class=""
                  />
                  <Label for="tag" class="block ">TLS</Label>
                  <Toggle id="tls" size="small" bind:checked={inf.tls} class="" />
                  {#if inf.tls}
                    <Label for="tag" class="block ">TLS Secret Name</Label>
                    <Input
                      id="tlsSecretName"
                      size="sm"
                      type="text"
                      bind:value={inf.tlsSecretName}
                      placeholder="Enter the TLS secret name"
                      class=""
                    />
                  {/if}
                </td>
              </tr>
            </tbody>
          </table>
          <!-- Reset & Save Button bottom right -->
          <div class="flex justify-end mt-4 p-4">
            <Button
              color="red"
              class="whitespace-nowrap self-start mr-2"
              on:click={() => handleDeleteClick(inf.id)}
            >
              Delete
            </Button>
            <Button
              color="primary"
              class="whitespace-nowrap self-start"
              on:click={() => handleInputSave(inf.id)}
            >
              Save
            </Button>
          </div>
        </div>
      </AccordionItem>
    {/each}
  {/key}
</Accordion>

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
    <CircleAlert class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      Are you sure you want to delete this interface?
    </h3>
    <Button color="red" class="me-2" on:click={() => handleDelete()}>Yes, I'm sure</Button>
    <Button color="alternative">No, cancel</Button>
  </div>
</Modal>

<div>
  <Modal bind:open={modal} size="xs" autoclose={false} class="w-full">
    <NewInterface bind:modal />
  </Modal>
</div>
