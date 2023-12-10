<script lang="ts">
  import NewInterface from "$lib/components/networking/NewInterface.svelte";
  import { client } from "$lib/pocketbase";
  import type { RolloutsRecord, RolloutsResponse } from "$lib/pocketbase/generated-types";
  import {
    rollouts,
    currentRollout,
    type Rexpand,
    updateDataStores,
    UpdateFilterEnum
  } from "$lib/stores/data";
  import {
    Accordion,
    AccordionItem,
    Button,
    Heading,
    Input,
    Label,
    Modal,
    P,
    Toggle
  } from "flowbite-svelte";
  import { Network, Plus } from "lucide-svelte";
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

  interface Ingress {
    ingressClass: string;
    annotations: Record<string, string>;
    rules: Rule[];
  }

  interface Rule {
    host: string;
    path: string;
    tls: boolean;
  }

  let interfaces: Interface[] = [];

  $: {
    console.log("currentRollout", $currentRollout?.manifest?.spec?.interfaces);
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
            host: firstRule.host,
            path: firstRule.path,
            tls: firstRule.tls
          };
        } else {
          newInterface = {
            id: interfaceId,
            name: i.name,
            port: i.port,
            host: "",
            path: "",
            tls: false
          };
        }
        interfaces.push(newInterface);
      });
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

    if (!updatedInterface.port) {
      toast.error("Port is required");
      return;
    }

    // Find the index of the current interface based on its unique identifier (id)
    // @ts-ignore
    const currentInterfaceIndex = $currentRollout.manifest.spec.interfaces.findIndex(
      (i: any) => i.id === updatedInterface.id
    );

    if (currentInterfaceIndex !== -1) {
      // Check if there's another interface with the same name, host, or port in the current rollout
      // @ts-ignore
      const existingInterface = $currentRollout.manifest.spec.interfaces.find(
        (i: any) => i.id !== updatedInterface.id && (
          i.name === updatedInterface.name ||
          i.port === updatedInterface.port ||
          (i.ingress && i.ingress.rules.some((rule: any) =>
            rule.host === updatedInterface.host && rule.path === updatedInterface.path))
      ));


      if (existingInterface) {
        toast.error("An interface with the same name, host, or port already exists");
        return;
      }
    }

    // Update interface in $currentRollout
    // @ts-ignore
    const rolloutInterfaceIndex = $currentRollout.manifest.spec.interfaces.findIndex(
      // do not only check for name, but also for port and host (so you can update the name of an interface)
      (i: any) => i.name === updatedInterface.name || i.port === updatedInterface.port
    );

    if (rolloutInterfaceIndex !== -1) {
      // @ts-ignore
      $currentRollout.manifest.spec.interfaces[rolloutInterfaceIndex] = {
        name: updatedInterface.name,
        port: updatedInterface.port,
        ingress: updatedInterface.host
          ? {
              ingressClass: "nginx",
              annotations: {
                "nginx.ingress.kubernetes.io/ssl-redirect": updatedInterface.tls ? "true" : "false"
              },
              rules: [
                {
                  host: updatedInterface.host,
                  path: updatedInterface.path,
                  tls: updatedInterface.tls
                }
              ]
            }
          : undefined
      };

      // if ingress is undefined, remove it from the manifest
      // @ts-ignore
      if (!$currentRollout.manifest.spec.interfaces[rolloutInterfaceIndex].ingress) {
        // @ts-ignore
        delete $currentRollout.manifest.spec.interfaces[rolloutInterfaceIndex].ingress;
      }
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
        user: client.authStore.model?.id
      };

      client
        .collection("rollouts")
        .create(data)
        .then((res) => {
          updateDataStores({
            filter: UpdateFilterEnum.ALL,
            projectId: $currentRollout?.project
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

<Accordion class="gap-2 grid mt-10">
  {#key $rollouts}
    {#each interfaces as inf, i (inf.id)}
      <AccordionItem class="rounded-lg">
        <div slot="header" class="flex">
          <div class="ring-1 p-2 rounded-lg ring-gray-500 mr-2 flex items-center justify-center">
            <Network class="w-4 h-4" />
          </div>
          <span class="pt-1">{inf.name}</span>
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
                </td>
              </tr>
            </tbody>
          </table>
          <!-- Reset & Save Button bottom right -->
          <div class="flex justify-end mt-4 p-4">
            <Button color="red" class="whitespace-nowrap self-start mr-2" on:click={() => {}}>
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

<div>
  <Modal bind:open={modal} size="xs" autoclose={false} class="w-full">
    <NewInterface bind:modal />
  </Modal>
</div>
