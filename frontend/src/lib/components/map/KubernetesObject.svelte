<script lang="ts">
  import MonacoEditor from "svelte-monaco";
  // @ts-ignore
  import yaml from "js-yaml";

  export let manifest: JSON;

  function jsonToYaml(json: any): string {
    return yaml.dump(json);
  }

  let yamlManifest = jsonToYaml(manifest);
</script>

<div class="log-container px-2 rounded-lg bg-gray-800">
  <div
    class="absolute top-36 bottom-0 right-0 left-0 overflow-y-auto rounded-lg p-2"
    style="background-color: #1E1E1E;"
  >
    <MonacoEditor
      bind:value={yamlManifest}
      options={{
        language: "yaml",
        automaticLayout: true,
        minimap: { enabled: false },
        readOnly: true,
        scrollBeyondLastLine: false,
        scrollbar: { vertical: "auto", horizontal: "auto" }
      }}
      theme="vs-dark"
    />
  </div>
</div>
