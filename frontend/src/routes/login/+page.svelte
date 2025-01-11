<script lang="ts">
  import { goto } from "$app/navigation";
  import { client, login } from "$lib/pocketbase";
  import { alertOnFailure } from "$lib/pocketbase/ui";
  import toast from "svelte-french-toast";
  import { Label, Input, Button } from "flowbite-svelte";
  import { onMount } from "svelte";
  import type { AuthProviderInfo } from "pocketbase";
  import colorTheme from "$lib/stores/theme";

  const DEFAULTS = {
    email: "",
    password: ""
  };
  let user = { ...DEFAULTS };

  async function submit() {
    await alertOnFailure(async function () {
      await login(user.email, user.password);
      toast.success("Logged in successfully!");
      goto("/app");
    });
  }

  let authProviders: AuthProviderInfo[] = [];

  async function getAuthMethods() {
    const response = (await client.collection("users").listAuthMethods()) as any;
    authProviders = response.authProviders ?? [];
  }

  onMount(() => {
    getAuthMethods();
  });

  async function loginWithProvider(provider: AuthProviderInfo) {
    try {
      const response = await client.collection("users").authWithOAuth2({ provider: provider.name });

      const meta: any = response.meta;

      if (meta.isNew && client.authStore.record?.id) {
        const formData = new FormData();

        const response = await fetch(meta.avatarUrl);

        if (response.ok) {
          const file = await response.blob();
          formData.append("avatar", file);
        }

        formData.append("name", meta.name);

        await client.collection("users").update(client.authStore.record.id, formData);
      }

      toast.success("Logged in successfully!");
      goto("/app");
    } catch (error) {
      toast.error("Failed to log in");
    }
  }
</script>

<div
  class="flex min-h-full flex-col justify-center py-12 sm:px-6 lg:px-8 bg-gradient-to-r from-gray-100 dark:from-gray-500 to-gray-200 dark:to-gray-600"
>
  <div class="mt-10 sm:mx-auto sm:w-full sm:max-w-[480px]">
    <div class="bg-background dark:bg-gray-800 px-6 py-12 shadow sm:rounded-lg sm:px-12">
      <div class="sm:mx-auto sm:w-full sm:max-w-md mb-5">
        {#if $colorTheme == "light"}
          <img class="h-20 w-auto mx-auto" src="/images/logo_primary_typo.png" alt="logo" />
        {/if}
        {#if $colorTheme == "dark"}
          <img class="h-20 w-auto mx-auto" src="/images/logo_background_typo.png" alt="logo" />
        {/if}
      </div>
      <form class="space-y-6" on:submit|preventDefault={submit} method="POST">
        <div>
          <Label class="space-y-2">
            <span>Email address</span>
            <Input type="email" size="md" bind:value={user.email} />
          </Label>
        </div>

        <div>
          <Label class="space-y-2">
            <span>Password</span>
            <Input type="password" size="md" bind:value={user.password} />
          </Label>
        </div>

        <div>
          <Button class="w-full" color="alternative" type="submit" on:submit={() => submit()}
            >Log in</Button
          >
        </div>
      </form>

      {#if authProviders.length > 0}
        <div>
          <div class="relative mt-10">
            <div class="absolute inset-0 flex items-center" aria-hidden="true">
              <div class="w-full border-t border-gray-200" />
            </div>
            <div class="relative flex justify-center text-sm font-medium leading-6">
              <span class="bg-background dark:bg-gray-800 px-6 dark:text-gray-400">or</span>
            </div>
          </div>

          <!-- if authProviders contains something with .name = "github" -->
          {#if authProviders.find((provider) => provider.name === "github")}
            <div class="mt-2 grid grid-cols-1 gap-4">
              <button
                on:click={() => {
                  const temp = authProviders.find((provider) => provider.name === "github");
                  if (temp) {
                    loginWithProvider(temp);
                  }
                }}
                class="flex w-full items-center justify-center gap-3 rounded-md bg-black px-3 py-1.5 text-background focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#1D9BF0]"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="16"
                  height="16"
                  fill="currentColor"
                  class="bi bi-github"
                  viewBox="0 0 16 16"
                >
                  <path
                    d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27s1.36.09 2 .27c1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.01 8.01 0 0 0 16 8c0-4.42-3.58-8-8-8"
                  />
                </svg>
                <span class="text-sm font-semibold leading-6">GitHub</span>
              </button>
            </div>
          {/if}

          <!-- if authProviders contains something with .name = "google" -->
          {#if authProviders.find((provider) => provider.name === "google")}
            <div class="mt-2 grid grid-cols-1 gap-4">
              <button
                on:click={() => {
                  const temp = authProviders.find((provider) => provider.name === "google");
                  if (temp) {
                    loginWithProvider(temp);
                  }
                }}
                class="flex w-full items-center justify-center gap-3 rounded-md bg-red-500 px-3 py-1.5 text-background focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#1D9BF0]"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="16"
                  height="16"
                  fill="currentColor"
                  class="bi bi-google"
                  viewBox="0 0 16 16"
                >
                  <path
                    d="M15.545 6.558a9.42 9.42 0 0 1 .139 1.626c0 2.434-.87 4.492-2.384 5.885h.002C11.978 15.292 10.158 16 8 16A8 8 0 1 1 8 0a7.689 7.689 0 0 1 5.352 2.082l-2.284 2.284A4.347 4.347 0 0 0 8 3.166c-2.087 0-3.86 1.408-4.492 3.304a4.792 4.792 0 0 0 0 3.063h.003c.635 1.893 2.405 3.301 4.492 3.301 1.078 0 2.004-.276 2.722-.764h-.003a3.702 3.702 0 0 0 1.599-2.431H8v-3.08h7.545z"
                  />
                </svg>
                <span class="text-sm font-semibold leading-6">Google</span>
              </button>
            </div>
          {/if}

          <!-- if authProviders contains something with .name = "microsoft" -->
          {#if authProviders.find((provider) => provider.name === "microsoft")}
            <div class="mt-2 grid grid-cols-1 gap-4">
              <button
                on:click={() => {
                  const temp = authProviders.find((provider) => provider.name === "microsoft");
                  if (temp) {
                    loginWithProvider(temp);
                  }
                }}
                class="flex w-full items-center justify-center gap-3 rounded-md bg-blue-500 px-3 py-1.5 text-background focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#1D9BF0]"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="16"
                  height="16"
                  fill="currentColor"
                  class="bi bi-microsoft"
                  viewBox="0 0 16 16"
                >
                  <path
                    d="M7.462 0H0v7.19h7.462V0zM16 0H8.538v7.19H16V0zM7.462 8.211H0V16h7.462V8.211zm8.538 0H8.538V16H16V8.211z"
                  />
                </svg>
                <span class="text-sm font-semibold leading-6">Microsoft</span>
              </button>
            </div>
          {/if}
        </div>
      {/if}
    </div>
  </div>
</div>
