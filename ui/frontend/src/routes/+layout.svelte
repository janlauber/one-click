<script lang="ts">
  import "../app.postcss";
  import "../app.css";
  import { metadata } from "$lib/stores/metadata";
  import { site } from "$lib/config";
  import { beforeNavigate } from "$app/navigation";
  import { Toaster } from "svelte-french-toast";
  import { DarkMode } from "flowbite-svelte";
  import colorTheme from "$lib/stores/theme";

  $: title = $metadata.title ? $metadata.title + " | " + site.name : site.name;
  $: description = $metadata.description ?? site.description;

  // reset metadata on navigation so that the new page inherits nothing from the old page
  beforeNavigate(() => {
    $metadata = {};
  });
</script>

<svelte:head>
  <title>{title}</title>
  <meta name="description" content={description} />
</svelte:head>

<div class="absolute top-0 left-0 right-0 bottom-0 overflow-hidden">
  <Toaster position="top-center" />
  <!-- only display nav when not on /login -->
  <main class="absolute top-0 left-0 right-0 bottom-0 overflow-hidden scrollbar-none">
    <slot />
  </main>
  <button class="absolute bottom-2 right-2"
    on:click={() => {
      $colorTheme = $colorTheme === "light" ? "dark" : "light";
      localStorage.setItem("color-theme", $colorTheme);
    }}
  >
    <DarkMode />
  </button>
</div>
