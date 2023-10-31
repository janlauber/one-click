<script>
  import Nav from "$lib/components/base/Nav.svelte";
</script>

<div>
  <div class="h-20">
    <Nav />
  </div>
  <div
    class="absolute top-20 left-0 right-0 bottom-0"
  >
    <slot />
  </div>
</div>
