<script lang="ts">
  import { currentRollout, currentRolloutStatus, rollouts } from "$lib/stores/data";
  import { autoScroll } from "$lib/utils/autoScroll";
  import { Accordion, AccordionItem, Button, Heading, P } from "flowbite-svelte";
  import { Box, FileDown, RefreshCcw } from "lucide-svelte";
  import { onMount } from "svelte";
  import Highlight, { LineNumbers } from "svelte-highlight";
  import prolog from "svelte-highlight/languages/prolog";
  import atomOneDark from "svelte-highlight/styles/atom-one-dark";

  function downloadLogs(podName: string) {
    const logsStream = logs[podName].join("\n");
    const blob = new Blob([logsStream], { type: "text/plain;charset=utf-8" });
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
  let logs: { [key: string]: string[] } = {};

  async function fetchLogs(podName: string) {
    const projectId = $currentRollout?.project;
    const baseUrl = window.location.hostname === "localhost" ? "http://localhost:8090" : "";
    const url = `${baseUrl}/rollouts/${projectId}/${podName}/logs`;
    const token = localStorage.getItem("pocketbase_auth");
    if (!token) {
      console.error("No token found");
      return;
    }
    const authHeader = { Authorization: `Bearer ${JSON.parse(token).token}` };

    // fetch logs from pod
    const res = await fetch(url, { headers: authHeader });

    // if response is not ok, throw error
    if (!res.ok) {
      throw new Error("Error fetching logs");
    }

    logs[podName] = [];

    // add the response to the log stream

    const reader = res.body?.getReader();
    if (!reader) {
      throw new Error("Error reading logs");
    }

    const decoder = new TextDecoder("utf-8");

    reader.read().then(function processLog({ done, value }) {
      if (done) {
        return;
      }

      if (value) {
        const decoded = decoder.decode(value);
        const lines = decoded.split("\n");
        lines.forEach((line) => {
          if (line) {
            logs[podName].push(line);
          }
        });
      }

      reader.read().then(processLog);
    });
  }

  onMount(() => {
    // fetch logs for each podName
    $currentRolloutStatus?.deployment?.podNames?.forEach((podName) => {
      fetchLogs(podName);
    });
  });

  // Trigger log stream for each podName
</script>

<svelte:head>
  {@html atomOneDark}
</svelte:head>

<div class="flex items-start justify-between">
  <div class="flex flex-col">
    <Heading tag="h2">Instances</Heading>
    <P class="text-gray-500 dark:text-gray-400 text-sm">Instances of your rollout.</P>
  </div>
</div>

<Accordion class="gap-2 grid mt-10" multiple flush>
  {#key $rollouts}
    {#each $currentRolloutStatus?.deployment?.podNames ?? [] as podName, i (podName)}
      <AccordionItem class="rounded-lg">
        <div slot="header" class="flex">
          <div class="ring-1 p-2 rounded-lg ring-gray-500 mr-2 flex items-center justify-center">
            <Box class="w-4 h-4" />
          </div>
          <span class="pt-1">{podName}</span>
        </div>
        <div class="log-container px-2 rounded-lg bg-gray-800">
          {#if logs[podName]}
            <div class="log-scroll text-sm scrollbar-none max-h-96" use:autoScroll>
              <Highlight language={prolog} code={logs[podName].join("\n")} let:highlighted>
                <LineNumbers {highlighted} wrapLines />
              </Highlight>
            </div>
          {/if}
        </div>
        <div class="flex justify-end">
          <Button
            size="sm"
            color="alternative"
            class="mt-4 mr-2"
            on:click={() => fetchLogs(podName)}
          >
            <RefreshCcw class="mr-2" />
            Refresh</Button
          >
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
