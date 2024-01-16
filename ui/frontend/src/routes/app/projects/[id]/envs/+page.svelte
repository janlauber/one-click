<script lang="ts">
  import { client } from "$lib/pocketbase";
  import type { RolloutsRecord, RolloutsResponse } from "$lib/pocketbase/generated-types";
  import {
    rollouts,
    type Rexpand,
    currentRollout,
    updateDataStores,
    UpdateFilterEnum
  } from "$lib/stores/data";
  import { Accordion, AccordionItem, Button, Heading, P } from "flowbite-svelte";
  import { Lock, Variable } from "lucide-svelte";
  import toast from "svelte-french-toast";
  import MonacoEditor from "svelte-monaco";

  interface Env {
    id: string;
    name: string;
    value: string;
  }

  type RolloutProperty = "env" | "secrets";

  function parseManifests(
    rollout: RolloutsResponse<Rexpand> | undefined,
    property: RolloutProperty
  ) {
    let items: Env[] = [];
    const specs = rollout?.manifest?.spec;

    if (specs && specs[property]) {
      specs[property].forEach((item: any, index: number) => {
        const itemId = `${rollout?.id}_${index}`;
        let newItem: Env = {
          id: itemId,
          name: item.name,
          value: item.value
        };
        items.push(newItem);
      });
    }

    return items;
  }

  let envs: Env[] = [];
  let secrets: Env[] = [];
  let envValue = "";
  let secretValue = "";

  let initialLoad = true;

  $: {
    if (initialLoad) {
      if ($currentRollout) {
        envs = parseManifests($currentRollout, "env");

        // parse envs to envValue for the code editor NAME=VALUE\nNAME=VALUE
        envValue = envs
          .map((env) => {
            return `${env.name}=${env.value}`;
          })
          .join("\n");

        secrets = parseManifests($currentRollout, "secrets");

        // parse secrets to secretValue for the code editor NAME=VALUE\nNAME=VALUE
        secretValue = secrets
          .map((secret) => {
            return `${secret.name}=${secret.value}`;
          })
          .join("\n");
      }
      initialLoad = false;
    }
  }

  async function handleInputSave() {
    // parse the values from the code editor
    const envs = envValue.split("\n");
    const secrets = secretValue.split("\n");

    // parse the envs and secrets to the correct format
    const parsedEnvs = envs.map((env) => {
      const [name, value] = env.split("=");
      return {
        name,
        value
      };
    });

    const parsedSecrets = secrets.map((secret) => {
      const [name, value] = secret.split("=");
      return {
        name,
        value
      };
    });

    // check if there are any duplicates
    const duplicateEnvs = parsedEnvs.filter((env) => {
      return parsedEnvs.filter((e) => e.name === env.name).length > 1;
    });

    const duplicateSecrets = parsedSecrets.filter((secret) => {
      return parsedSecrets.filter((s) => s.name === secret.name).length > 1;
    });

    // also check if there are any duplicates between envs and secrets
    const duplicateEnvSecrets = parsedEnvs.filter((env) => {
      return parsedSecrets.filter((s) => s.name === env.name).length > 0;
    });

    if (
      duplicateEnvs.length > 0 &&
      !(parsedEnvs.length === 1 && parsedEnvs[0].name === "" && parsedEnvs[0].value === undefined)
    ) {
      toast.error("There are duplicate envs");
      return;
    }

    if (
      duplicateSecrets.length > 0 &&
      !(
        parsedSecrets.length === 1 &&
        parsedSecrets[0].name === "" &&
        parsedSecrets[0].value === undefined
      )
    ) {
      toast.error("There are duplicate secrets");
      return;
    }

    if (
      duplicateEnvSecrets.length > 0 &&
      !(
        parsedEnvs.length === 1 &&
        parsedEnvs[0].name === "" &&
        parsedEnvs[0].value === undefined
      ) &&
      !(
        parsedSecrets.length === 1 &&
        parsedSecrets[0].name === "" &&
        parsedSecrets[0].value === undefined
      )
    ) {
      toast.error("There are duplicate envs and secrets");
      return;
    }

    // check if there are any empty envs or secrets
    const emptyEnvs = parsedEnvs.filter((env) => {
      return env.name === "" || env.value === "";
    });

    const emptySecrets = parsedSecrets.filter((secret) => {
      return secret.name === "" || secret.value === "";
    });

    if (
      emptyEnvs.length > 0 &&
      !(parsedEnvs.length === 1 && parsedEnvs[0].name === "" && parsedEnvs[0].value === undefined)
    ) {
      toast.error("There are empty envs");
      return;
    }

    if (
      emptySecrets.length > 0 &&
      !(
        parsedSecrets.length === 1 &&
        parsedSecrets[0].name === "" &&
        parsedSecrets[0].value === undefined
      )
    ) {
      toast.error("There are empty secrets");
      return;
    }

    // check if there are undefined envs or secrets
    const undefinedEnvs = parsedEnvs.filter((env) => {
      return env.name === undefined || env.value === undefined;
    });

    const undefinedSecrets = parsedSecrets.filter((secret) => {
      return secret.name === undefined || secret.value === undefined;
    });

    if (
      undefinedEnvs.length > 0 &&
      !(parsedEnvs.length === 1 && parsedEnvs[0].name === "" && parsedEnvs[0].value === undefined)
    ) {
      toast.error("There are undefined envs");
      return;
    }

    if (
      undefinedSecrets.length > 0 &&
      !(
        parsedSecrets.length === 1 &&
        parsedSecrets[0].name === "" &&
        parsedSecrets[0].value === undefined
      )
    ) {
      toast.error("There are undefined secrets");
      return;
    }

    // update the envs and secrets in the rollout
    if (!$currentRollout) {
      toast.error("No rollout selected");
      return;
    }

    if (parsedEnvs.length === 1 && parsedEnvs[0].name === "" && parsedEnvs[0].value === undefined) {
      if ($currentRollout && $currentRollout.manifest) {
        if (!$currentRollout.manifest.spec) {
          $currentRollout.manifest.spec = {};
        }
        $currentRollout.manifest.spec.env = [];
      }
    } else {
      if ($currentRollout && $currentRollout.manifest) {
        $currentRollout.manifest.spec.env = parsedEnvs;
      }
    }

    if (
      parsedSecrets.length === 1 &&
      parsedSecrets[0].name === "" &&
      parsedSecrets[0].value === undefined
    ) {
      if ($currentRollout && $currentRollout.manifest) {
        if (!$currentRollout.manifest.spec) {
          $currentRollout.manifest.spec = {};
        }
        $currentRollout.manifest.spec.secrets = [];
      }
    } else {
      if ($currentRollout && $currentRollout.manifest) {
        $currentRollout.manifest.spec.secrets = parsedSecrets;
      }
    }

    try {
      await updateManifest($currentRollout.manifest);
      toast.success("Envs and secrets updated successfully");
    } catch (error) {
      console.error("Failed to update envs and secrets:", error);
      toast.error("Failed to update envs and secrets");
    }
  }

  async function handleEnvInputSave(id?: string) {
    // if id is not provided, save all envs
    if (!id) {
      envs.forEach((env) => {
        handleEnvInputSave(env.id);
      });
      return;
    }

    const envIndex = envs.findIndex((env) => env.id === id);
    if (!$currentRollout) {
      toast.error("No rollout selected");
      return;
    }

    const updatedEnv = envs[envIndex];

    if (!updatedEnv) {
      toast.error("No env found");
      return;
    }

    if (!updatedEnv.name || !updatedEnv.value) {
      toast.error("Please fill in all fields");
      return;
    }

    //@ts-ignore
    const currentEnvIndex = $currentRollout.manifest.spec.env.findIndex(
      (env: any) => env.name === updatedEnv.name
    );

    if (currentEnvIndex === -1) {
      // Check if there is an existing env with the same name
      // @ts-ignore
      const existingEnv = $currentRollout.manifest.spec.env.find(
        (env: any) => env.name === updatedEnv.name
      );

      // exclude if the existing env is the same as the updated env
      // and there is only one env in the list
      if (existingEnv && envs.length > 1) {
        toast.error("An env with this name already exists");
        return;
      }
    }

    // Update the env in $currentRollout
    // @ts-ignore
    const rolloutEnvIndex = $currentRollout.manifest.spec.env.findIndex(
      (env: any) => env.name === updatedEnv.name
    );

    if (rolloutEnvIndex !== -1) {
      // @ts-ignore
      $currentRollout.manifest.spec.env[rolloutEnvIndex] = {
        name: updatedEnv.name,
        value: updatedEnv.value
      };
    } else {
      // @ts-ignore
      $currentRollout.manifest.spec.env.push({
        name: updatedEnv.name,
        value: updatedEnv.value
      });
    }

    // Update the manifest
    if (!$currentRollout.manifest) {
      toast.error("No manifest found");
      return;
    }

    // @ts-ignore
    await updateManifest($currentRollout.manifest);

    toast.success("Env updated successfully");

    initialLoad = true;
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

<div class="flex items-start justify-between">
  <div class="flex flex-col">
    <Heading tag="h2">Envs & Secrets</Heading>
    <P class="text-gray-500 dark:text-gray-400 text-sm">
      Environment variables and secrets for your rollout.
    </P>
  </div>
</div>

<Accordion class="gap-2 grid mt-10" flush multiple>
  <AccordionItem class="rounded-lg" open>
    <div slot="header" class="flex">
      <div class="ring-1 p-2 rounded-lg ring-gray-500 mr-2 flex items-center justify-center">
        <Variable class="w-4 h-4" />
      </div>
      <span class="pt-1">Environment Variables</span>
    </div>
    <!-- Code editor -->

    <div class="h-64 overflow-y-auto rounded-lg p-2" style="background-color: #1E1E1E;">
      {#key $rollouts}
        <MonacoEditor
          bind:value={envValue}
          options={{ language: "shell", automaticLayout: false, minimap: { enabled: false } }}
          theme="vs-dark"
        />
      {/key}
    </div>
  </AccordionItem>
  <AccordionItem class="rounded-lg" open>
    <div slot="header" class="flex">
      <div class="ring-1 p-2 rounded-lg ring-gray-500 mr-2 flex items-center justify-center">
        <Lock class="w-4 h-4" />
      </div>
      <span class="pt-1">Secret Environment Variables</span>
    </div>
    <!-- Code editor -->
    <div class="h-64 overflow-y-auto rounded-lg p-2" style="background-color: #1E1E1E;">
      {#key $rollouts}
        <MonacoEditor
          bind:value={secretValue}
          options={{ language: "shell", automaticLayout: true }}
        />
      {/key}
    </div>
  </AccordionItem>
</Accordion>

<div class="flex justify-end mt-4 gap-4">
  <!-- Reset -->
  <Button color="alternative" on:click={() => (initialLoad = true)}>Reset</Button>
  <Button on:click={() => handleInputSave()}>Save</Button>
</div>
