<script lang="ts">
  // @ts-ignore
  import { createPatch } from "diff";
  import { html } from "diff2html";
  // @ts-ignore
  import yaml from "js-yaml";
  import "diff2html/bundles/css/diff2html.min.css"; // Import the CSS for diff2html

  export let jsonObject1 = {}; // The first JSON object
  export let jsonObject2 = {}; // The second JSON object
  export let title1 = "Old Rollout"; // Title for the first JSON object
  export let title2 = "New Rollout"; // Title for the second JSON object

  $: diffHTML = generateDiffHTML(jsonObject1, jsonObject2);

  function generateDiffHTML(obj1: object, obj2: object) {
    const yamlStr1 = yaml.dump(obj1);
    const yamlStr2 = yaml.dump(obj2);
    const diffString = createPatch("rollout.yaml", yamlStr1, yamlStr2);

    return html(diffString, {
      drawFileList: false,
      matching: "lines",
      outputFormat: "side-by-side"
    });
  }
</script>

<div class="flex flex-col md:flex-row -mb-6">
  <div class="md:w-1/2">
    <div class="text-center text-red-600 text-sm p-2 bg-red-200">
      <span class=" font-bold">Current: </span>
      {title1}
    </div>
  </div>
  <div class="md:w-1/2">
    <div class="text-center text-green-600 text-sm p-2 bg-green-200">
      <span class=" font-bold">New: </span>
      {title2}
    </div>
  </div>
</div>

<div class="overflow-auto">{@html diffHTML}</div>
