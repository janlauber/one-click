<script>
  import Nav from "$lib/components/base/Nav.svelte";
</script>

<div>
  <div class="">
    <Nav />
  </div>
  <div
    class="absolute top-14 left-0 right-0 bottom-0 overflow-y-auto scrollbar-none"
  >
    <slot />
  </div>
</div>
