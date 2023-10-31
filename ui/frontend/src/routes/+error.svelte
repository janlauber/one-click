<script>
    import { page } from '$app/stores';
</script>

<div class="flex h-screen">
    <div class="m-auto">
        <h1>{$page.status}: {$page.error?.message}</h1>
    </div>
</div>
