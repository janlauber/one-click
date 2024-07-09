<script lang="ts">
  import { onDestroy, onMount, type ComponentType, type SvelteComponent } from "svelte";
  import PlaceholderComponent from "./PlaceholderComponent.svelte";
  import Desktop from "./Desktop.svelte";
  export let podName: string;
  export let projectId: string;

  let Console: ComponentType<SvelteComponent> = PlaceholderComponent;
  let clear: number;
  // let loadingCodeEditor = false;

  onMount(async () => {
    // loadingCodeEditor = true;
    Console = (await import("$lib/components/map/Console.svelte")).default;
  });

  onDestroy(() => {
    // loadingCodeEditor = false;
    Console = PlaceholderComponent;
    clearInterval(clear);
  });
</script>

<div class="p-2 rounded-lg bg-gray-800">
  <Desktop {Console} {podName} {projectId} />
</div>
