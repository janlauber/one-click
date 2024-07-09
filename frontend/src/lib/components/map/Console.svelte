<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import darkTheme from "$lib/stores/theme";
  import { terminal_size } from "$lib/stores/terminal";

  import { Terminal } from "xterm";
  import { FitAddon } from "xterm-addon-fit";

  export let podName: string;
  export let projectId: string;

  const terminal = new Terminal({
    convertEol: true,
    disableStdin: false,
    cursorBlink: true,
    fontFamily: "monospace",
    fontSize: 14,
    theme: $darkTheme
      ? {
          foreground: "#d2d2d2",
          background: "#2B3441",
          cursor: "#adadad",
          black: "#000000",
          red: "#d81e00",
          green: "#5ea702",
          yellow: "#cfae00",
          blue: "#427ab3",
          magenta: "#89658e",
          cyan: "#00a7aa",
          white: "#dbded8",
          brightBlack: "#686a66",
          brightRed: "#f54235",
          brightGreen: "#99e343",
          brightYellow: "#fdeb61",
          brightBlue: "#84b0d8",
          brightMagenta: "#bc94b7",
          brightCyan: "#37e6e8",
          brightWhite: "#f1f1f0"
        }
      : {
          foreground: "#d2d2d2",
          background: "#2B3441",
          cursor: "#adadad",
          black: "#000000",
          red: "#d81e00",
          green: "#5ea702",
          yellow: "#cfae00",
          blue: "#427ab3",
          magenta: "#89658e",
          cyan: "#00a7aa",
          white: "#dbded8",
          brightBlack: "#686a66",
          brightRed: "#f54235",
          brightGreen: "#99e343",
          brightYellow: "#fdeb61",
          brightBlue: "#84b0d8",
          brightMagenta: "#bc94b7",
          brightCyan: "#37e6e8",
          brightWhite: "#f1f1f0"
        },
    scrollOnUserInput: true
  });

  let socket: WebSocket;
  const fitAddon = new FitAddon();
  terminal.loadAddon(fitAddon);

  const connectWebSocket = () => {
    let host = window.location.host;
    if (host.includes("localhost")) {
      host = "localhost:8090";
    }

    const protocol = window.location.protocol === "https:" ? "wss" : "ws";
    socket = new WebSocket(`${protocol}://${host}/ws/k8s/terminal`);

    socket.binaryType = "arraybuffer";

    socket.onopen = () => {
      const message = JSON.stringify({ projectId, podName });
      socket.send(message);
      terminal.onData((data) => {
        socket.send(new TextEncoder().encode(data));
      });
    };

    socket.onmessage = (event: MessageEvent) => {
      if (event.data instanceof ArrayBuffer) {
        terminal.write(new TextDecoder().decode(event.data));
      }
    };

    socket.onerror = (error: Event) => {
      const errorMessage = (error as ErrorEvent).message || "An error occurred";
      terminal.write(`\r\nWebSocket error: ${errorMessage}\r\n`);
    };

    socket.onclose = () => {
      terminal.write("\r\nWebSocket closed. Attempting to reconnect...\r\n");
      setTimeout(connectWebSocket, 5000);
    };
  };

  terminal.onResize((size) => {
    const terminal_size = {
      cols: size.cols,
      rows: size.rows,
      y: size.rows,
      x: size.cols
    };
    if (socket.readyState === WebSocket.OPEN) {
      socket.send(new TextEncoder().encode("\x01" + JSON.stringify(terminal_size)));
    }
  });

  let div: HTMLDivElement;

  onMount(() => {
    connectWebSocket();
    terminal.open(div);
    setTimeout(() => {
      fitAddon.fit();
    }, 300);
  });

  onDestroy(() => {
    socket?.close();
    terminal.dispose();
  });

  export const update_height = () => {
    fitAddon.fit();
  };

  $: {
    $terminal_size;
    setTimeout(() => {
      update_height();
    }, 300);
  }
</script>

<div bind:this={div} style="height: 100%; width: 100%;" />

<style>
  div {
    height: 100%;
    width: 100%;
  }
  div :global(.xterm) {
    height: 100%;
    padding: 5px;
  }
  div :global(.xterm-viewport) {
    overflow-y: hidden !important;
  }
</style>
