<script lang="ts">
  import darkTheme from "$lib/stores/theme";
  import type { ComponentType, SvelteComponent } from "svelte";
  // eslint-disable-next-line @typescript-eslint/naming-convention
  export let Console: ComponentType<SvelteComponent>;
  export let podName: string;
  export let projectId: string;
</script>

{#key $darkTheme}
  <div class="h-full overflow-y-auto">
    <svelte:component this={Console} {podName} {projectId} />
  </div>
{/key}
