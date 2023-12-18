<script lang="ts">
  import { rollouts } from "$lib/stores/data";
  import { Accordion, AccordionItem, Button, Heading, Input, Label, P } from "flowbite-svelte";
  import { Lock, Plus, Variable, X } from "lucide-svelte";

  interface Env {
    id: string;
    name: string;
    value: string;
  }

  let envs: Env[] = [
    {
      id: "1",
      name: "TEST",
      value: "test"
    },
    {
      id: "2",
      name: "TEST2",
      value: "test2"
    }
  ];
  let secrets: Env[] = [];
</script>

<div class="flex items-start justify-between">
  <div class="flex flex-col">
    <Heading tag="h2">Envs & Secrets</Heading>
    <P class="text-gray-500 dark:text-gray-400 text-sm">
      Environment variables and secrets for your rollout.
    </P>
  </div>
</div>

<Accordion class="gap-2 grid mt-10">
  {#key $rollouts}
    <AccordionItem class="rounded-lg">
      <div slot="header" class="flex">
        <div class="ring-1 p-2 rounded-lg ring-gray-500 mr-2 flex items-center justify-center">
          <Variable class="w-4 h-4" />
        </div>
        <span class="pt-1">Environment Variables</span>
      </div>
      {#each envs as env, i (env.id)}
        <div class="flex items-end space-x-4">
          <div class="flex-grow grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <Label class="mt-4">Name</Label>
              <Input
                size="sm"
                type="text"
                placeholder="Name"
                bind:value={env.name}
                class="w-full"
              />
            </div>
            <div>
              <Label class="mt-4">Value</Label>
              <Input
                size="sm"
                type="text"
                placeholder="Value"
                bind:value={env.value}
                class="w-full"
              />
            </div>
          </div>
          <div>
            <Button color="red">
              <X class="w-4 h-4" />
            </Button>
          </div>
        </div>
      {/each}
      <div class="flex justify-between mt-4 gap-4">
        <Button color="alternative">
          <Plus class="mr-2 w-4 h-4" />
          Add</Button>
        <Button>Save</Button>
      </div>
      <!-- Save button -->
      <div class="flex justify-end mt-4">
      </div>
    </AccordionItem>
    <AccordionItem class="rounded-lg">
      <div slot="header" class="flex">
        <div class="ring-1 p-2 rounded-lg ring-gray-500 mr-2 flex items-center justify-center">
          <Lock class="w-4 h-4" />
        </div>
        <span class="pt-1">Secret Environment Variables</span>
      </div>
      <div class=""></div>
    </AccordionItem>
  {/key}
</Accordion>
