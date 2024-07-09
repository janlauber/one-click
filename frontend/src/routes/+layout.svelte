<script lang="ts">
  import "../styles/app.postcss";
  import "../styles/app.css";
  import "../styles/xterm.css";
  import { metadata } from "$lib/stores/metadata";
  import { site } from "$lib/config";
  import { beforeNavigate } from "$app/navigation";
  import { Toaster } from "svelte-french-toast";
  import { DarkMode } from "flowbite-svelte";
  import colorTheme from "$lib/stores/theme";
  import { version } from "$lib/stores/version";
  import { BookOpen } from "lucide-svelte";

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

<div class="absolute top-0 left-0 right-0 bottom-0 overflow-hidden dark:text-gray-400">
  <Toaster position="top-center" />
  <!-- only display nav when not on /login -->
  <main class="absolute top-0 left-0 right-0 bottom-0 overflow-hidden scrollbar-none">
    <slot />
  </main>

  <p
    class="absolute bottom-4 right-16
    text-xs text-gray-400
  "
  >
    <a
      class="hover:text-gray-500 dark:text-gray-500 dark:hover:text-gray-400 transition-all ease-in-out duration-200"
      href="https://github.com/janlauber/one-click/releases/tag/{$version}"
      target="_blank"
    >
      One-Click {$version}
    </a>
    |
    <a
      class="hover:text-gray-500 dark:text-gray-500 dark:hover:text-gray-400 transition-all ease-in-out duration-200"
      href="https://docs.one-click.dev"
      target="_blank"
    >
      <BookOpen class="w-4 h-4 inline-block mr-1" />Docs
    </a>
  </p>
  <button
    class="absolute bottom-2 right-2"
    on:click={() => {
      $colorTheme = $colorTheme === "light" ? "dark" : "light";
      localStorage.setItem("color-theme", $colorTheme);
    }}
  >
    <DarkMode />
  </button>
</div>
