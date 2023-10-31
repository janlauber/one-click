<script lang="ts">
	import { Node, Anchor, Slider } from 'svelvet';
	import { writable } from 'svelte/store';

	const parameter = writable(10);
	export let drop: any;
	export let dynamic: any;
</script>

<Node let:grabHandle let:selected {drop} {dynamic}>
	<div class="node" use:grabHandle class:selected>

	</div>
</Node>
