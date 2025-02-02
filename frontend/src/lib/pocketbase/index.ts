import PocketBase from "pocketbase";
import { writable } from "svelte/store";
import toast from "svelte-french-toast";
import { goto } from "$app/navigation";
import type { TypedPocketBase } from "./generated-types";

export const client = new PocketBase() as TypedPocketBase;

export const currentUser = writable(client.authStore.record);

export async function login(
    email: string,
    password: string,
    register = false,
    rest: { [key: string]: unknown } = {}
) {
    if (register) {
        const user = { ...rest, email, password, confirmPassword: password };
        await client.collection("users").create(user);
    }
    await client.collection("users").authWithPassword(email, password);
}

export function logout() {
    client.authStore.clear();
    currentUser.set(null);
    goto("/login/");
    toast.success("Successfully logged out.");
}
