<script lang="ts">
  import SideNav from "$lib/components/projects/SideNav.svelte";
</script>

<div class="absolute top-10 left-10 bottom-10 overflow-y-auto">
  <SideNav />
</div>

<div class="absolute top-1 pt-8 left-64 right-10 bottom-0 overflow-y-auto scrollbar-none pb-8 px-2 overflow-x-hidden">
  <slot />
</div>
