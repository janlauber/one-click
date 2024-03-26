<script lang="ts">
  import { selectedProject } from "$lib/stores/data";
  import { onDestroy, onMount } from "svelte";
  import MonacoEditor from "svelte-monaco";

  export let podName: string;

  // function downloadLogs(podName: string) {
  //   const logsStream = logs[podName].join("\n");
  //   const blob = new Blob([logsStream], { type: "text/plain;charset=utf-8" });
  //   const url = URL.createObjectURL(blob);
  //   const a = document.createElement("a");
  //   a.href = url;
  //   a.download = `${podName}.log`;
  //   document.body.appendChild(a);
  //   a.click();
  //   document.body.removeChild(a);
  // }

  // // get rollout logs from specific pod as EventSource from /rollouts/{projectId}/{podName}/logs
  // // Object to hold log streams for each pod
  let initialLoadComplete = false;
  let logs: string = "";
  let ws: WebSocket;

  onMount(() => {
    // host
    let host = window.location.host;

    if (host.includes("localhost")) {
      host = "localhost:8090";
    }

    // check for tls
    let protocol = window.location.protocol === "https:" ? "wss" : "ws";

    ws = new WebSocket(`${protocol}://${host}/ws/k8s/logs`);

    ws.onopen = () => {
      type LogMessage = {
        rolloutId: string;
        podName: string;
      };
      let message: LogMessage = {
        rolloutId: $selectedProject?.id ?? "",
        podName: podName
      };

      ws.send(JSON.stringify(message));
    };

    ws.onmessage = (event) => {
      // event.data is a string
      logs += event.data;
    };

    // set initialLoadComplete to true after 0.3s
    setTimeout(() => {
      initialLoadComplete = true;
    }, 300);
  });

  onDestroy(() => {
    // Close the WebSocket connection
    ws.close();

    // Reset the logs
    logs = "";
  });

  function handleEditorReady(event: CustomEvent) {
    const editor = event.detail;
    setTimeout(() => {
      scrollToBottom(editor);
    }, 100); // Adjust delay as necessary
  }

  function scrollToBottom(editor: any) {
    if (editor) {
      const model = editor.getModel();
      if (model) {
        const lastLine = model.getLineCount();
        editor.revealLine(lastLine);
      }
    }
  }
</script>

<div class="log-container px-2 rounded-lg bg-gray-800">
  <div
    class="absolute top-36 bottom-0 right-0 left-0 overflow-y-auto rounded-lg p-2"
    style="background-color: #1E1E1E;"
  >
    <MonacoEditor
      bind:value={logs}
      options={{
        language: "shell",
        automaticLayout: true,
        minimap: { enabled: false },
        readOnly: true,
        scrollBeyondLastLine: false,
        scrollbar: { vertical: "auto", horizontal: "auto" },
        wordWrap: "on"
      }}
      theme="vs-dark"
      on:ready={handleEditorReady}
    />
  </div>
</div>
