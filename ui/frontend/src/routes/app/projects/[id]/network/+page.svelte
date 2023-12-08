<script lang="ts">
  import NewInterface from "$lib/components/networking/NewInterface.svelte";
  import type { RolloutsResponse } from "$lib/pocketbase/generated-types";
  import { rollouts, type Rexpand } from "$lib/stores/data";
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

  export let modal: boolean;
  let current_rollout: RolloutsResponse<Rexpand> | undefined;

  $: if ($rollouts.length > 0) {
    // get the current rollout on following priority:
    // 1. no endDate set
    // 2. newest endDate

    const current_rollouts = $rollouts.filter((r) => !r.endDate);
    if (current_rollouts.length > 0) {
      current_rollout = current_rollouts[0];
    } else {
      current_rollout = $rollouts.reduce((prev, current) => {
        return prev.endDate > current.endDate ? prev : current;
      });
    }
  }

  interface Interface {
    id: string;
    name: string;
    port: number;
    host: string;
    path: string;
    tls: boolean;
  }

  interface PocketbaseInterface {
    name: string;
    port: string;
    ingresses: Ingress[];
  }

  interface Ingress {
    ingressClass: string;
    annotations: Annotation[];
    rules: Rule[];
  }

  interface Annotation {
    key: string;
    value: string;
  }

  interface Rule {
    host: string;
    path: string;
    tls: boolean;
  }

  let interfaces: Interface[] = [];

  $: parseManifestsToInterfaces(current_rollout);

  function parseManifestsToInterfaces(rollout: RolloutsResponse<Rexpand> | undefined) {

    interfaces = []; // Clear existing interfaces
    if (
      rollout?.manifest &&
      rollout.manifest.spec &&
      rollout.manifest.spec.interfaces
    ) {
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

  function handleAddInterface() {
    const newInterface = {
      id: Math.random().toString(36).substring(2, 9),
      name: "",
      port: 0,
      host: "",
      path: "",
      tls: false
    };
    interfaces = [...interfaces, newInterface];
  }

  async function handleInputSave(id: string) {
    const interfaceIndex = interfaces.findIndex((inf) => inf.id === id);
    const ingress: Ingress = {
      ingressClass: "nginx",
      annotations: [
        {
          key: "nginx.ingress.kubernetes.io/rewrite-target",
          value: "/"
        }
      ],
      rules: [
        {
          host: interfaces[interfaceIndex].host,
          path: interfaces[interfaceIndex].path,
          tls: interfaces[interfaceIndex].tls
        }
      ]
    };

    const interfaceData: PocketbaseInterface = {
      name: interfaces[interfaceIndex].name,
      port: interfaces[interfaceIndex].port.toString(),
      ingresses: [ingress]
    };
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
</Accordion>

<div>
  <Modal bind:open={modal} size="xs" autoclose={false} class="w-full">
    <NewInterface {current_rollout} bind:modal />
  </Modal>
</div>
