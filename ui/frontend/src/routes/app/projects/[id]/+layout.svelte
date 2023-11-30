<script lang="ts">
  import SideNav from "$lib/components/projects/SideNav.svelte";
</script>

<div class="absolute top-10 left-10 bottom-10 overflow-y-auto">
  <SideNav />
</div>

<div class="absolute top-8 left-64 right-10 bottom-10 overflow-y-auto scrollbar-none p-2">
  <slot />
</div>
