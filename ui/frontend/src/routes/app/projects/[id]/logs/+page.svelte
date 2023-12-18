<script lang="ts">
  import { currentRollout, rollouts } from "$lib/stores/data";
  import { autoScroll } from "$lib/utils/autoScroll";
  import { getRolloutStatus } from "$lib/utils/rollouts";
  import { Accordion, AccordionItem, Button, Heading, P } from "flowbite-svelte";
  import { Box, FileDown, RefreshCcw } from "lucide-svelte";
  import { onDestroy } from "svelte";
  import Highlight, { LineNumbers } from "svelte-highlight";
  import prolog from "svelte-highlight/languages/prolog";
  import atomOneDark from "svelte-highlight/styles/atom-one-dark";

  let podNames: string[] = [];

  $: {
    if ($currentRollout) {
      getRolloutStatus($currentRollout?.project, $currentRollout.id)
        .then((res) => {
          if (res?.deployment.podNames) {
            podNames = res?.deployment.podNames;
          }
        })
        .catch((err) => {
          console.log(err);
        });
    }

    if (podNames.length > 0) {
      podNames.forEach((podName) => {
        if (!logStreams[podName]) {
          startLogStream(podName);
        }
      });
    }
  }

  function downloadLogs(podName: string) {
    const logs = logStreams[podName].join("\n");
    const blob = new Blob([logs], { type: "text/plain;charset=utf-8" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${podName}.log`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
  }

  // get rollout logs from specific pod as EventSource from /rollouts/{projectId}/{podName}/logs
  // Object to hold log streams for each pod
  let logStreams: { [key: string]: string[] } = {};

  function startLogStream(podName: string) {
    const projectId = $currentRollout?.project;
    const baseUrl = window.location.hostname === "localhost" ? "http://localhost:8090" : "";
    const url = `${baseUrl}/rollouts/${projectId}/${podName}/logs`;
    const eventSource = new EventSource(url);

    logStreams[podName] = [];

    eventSource.onmessage = (event) => {
      logStreams[podName].push(event.data);
    };

    eventSource.onerror = (error) => {
      console.error(`Error with EventSource for pod ${podName}:`, error);
      eventSource.close();
    };

    // Clean up the EventSource when the component is destroyed
    onDestroy(() => {
      eventSource.close();
    });
  }

  // Trigger log stream for each podName
</script>

<svelte:head>
  {@html atomOneDark}
</svelte:head>

<div class="flex items-start justify-between">
  <div class="flex flex-col">
    <Heading tag="h2">Logs</Heading>
    <P class="text-gray-500 dark:text-gray-400 text-sm">
      Logs of your application / framework.
    </P>
  </div>
</div>

<Accordion class="gap-2 grid mt-10" multiple flush>
  {#key $rollouts}
    {#each podNames as podName, i (podName)}
      <AccordionItem class="rounded-lg" open>
        <div slot="header" class="flex">
          <div class="ring-1 p-2 rounded-lg ring-gray-500 mr-2 flex items-center justify-center">
            <Box class="w-4 h-4" />
          </div>
          <span class="pt-1">{podName}</span>
        </div>
        <div class="log-container px-2 rounded-lg bg-gray-800">
          {#if logStreams[podName]}
            <div class="log-scroll text-sm scrollbar-none" use:autoScroll>
              <Highlight language={prolog} code={logStreams[podName].join("\n")} let:highlighted>
                <LineNumbers {highlighted} wrapLines />
              </Highlight>
            </div>
          {:else}
            <p class="no-logs">No logs found</p>
          {/if}
        </div>
        <div class="flex justify-end">
          <Button size="sm" color="alternative" class="mt-4 mr-2" on:click={() => startLogStream(podName)}>
            <RefreshCcw class="mr-2" />
            Refresh</Button>
          <Button size="sm" class="mt-4" on:click={() => downloadLogs(podName)}>
            <FileDown class="mr-2" />
            Download Logs</Button
          >
        </div>
      </AccordionItem>
    {/each}
  {/key}
</Accordion>

<style>
  .log-scroll {
    overflow-y: scroll;
  }
</style>
