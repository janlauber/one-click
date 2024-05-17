<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import MonacoEditor from "svelte-monaco";
  import { selectedProject } from "$lib/stores/data";

  export let podName: string;
  let logs = "";
  let ws: WebSocket;
  let reconnectInterval: any;

  const reconnectDelay = 5000;

  function setupWebSocket() {
    let host = window.location.host.includes("localhost") ? "localhost:8090" : window.location.host;
    const protocol = window.location.protocol === "https:" ? "wss" : "ws";

    ws = new WebSocket(`${protocol}://${host}/ws/k8s/logs`);

    ws.onopen = () => {
      const message = {
        rolloutId: $selectedProject?.id ?? "",
        podName
      };
      ws.send(JSON.stringify(message));
    };

    ws.onmessage = (event) => {
      logs += event.data; // Consider debouncing this for performance
    };

    ws.onclose = () => {
      console.log("WebSocket closed. Attempting to reconnect...");
      clearTimeout(reconnectInterval);
      reconnectInterval = setTimeout(setupWebSocket, reconnectDelay);
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };
  }

  onMount(() => {
    setupWebSocket();
  });

  onDestroy(() => {
    ws.close();
    clearTimeout(reconnectInterval);
  });

  function handleEditorReady(event: CustomEvent) {
    const editor = event.detail;
    setTimeout(() => scrollToBottom(editor), 100);
  }

  function scrollToBottom(editor: any) {
    const model = editor.getModel();
    const lastLine = model.getLineCount();
    editor.revealLine(lastLine);
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
