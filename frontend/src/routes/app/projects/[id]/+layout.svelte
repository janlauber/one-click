<script lang="ts">
  import { metadata } from "$lib/stores/metadata";
  import { selectedProject } from "$lib/stores/data";

  $metadata.title = $selectedProject?.name || "Deployments";
</script>

<slot />
