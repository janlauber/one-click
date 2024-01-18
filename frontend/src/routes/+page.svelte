<script lang="ts">
    import {goto} from "$app/navigation";
    import {currentUser} from "$lib/pocketbase/index.js";
    import { get } from 'svelte/store';

    const user = get(currentUser);
    if (user) {
        goto("/app");
    } else {
        goto("/login");
    }

</script>
