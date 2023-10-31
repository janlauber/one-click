import { client } from "$lib/pocketbase";
import type { Handle } from "@sveltejs/kit";

export const handle: Handle = async ({ event, resolve }) => {
    if (client.authStore.isValid) {
        try {
            await client.collection("users").authRefresh();
        } catch (_) {
            client.authStore.clear();
        }
    }

    const response = await resolve(event);

    return response;
};
