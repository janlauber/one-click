<script lang="ts">
  import { goto } from "$app/navigation";
  import { login } from "$lib/pocketbase";
  import { alertOnFailure } from "$lib/pocketbase/ui";
  import toast from "svelte-french-toast";
  import { Label, Input, Checkbox, Button } from "flowbite-svelte";

  const DEFAULTS = {
    email: "",
    password: ""
  };
  let user = { ...DEFAULTS };
  let loading = false;

  async function submit() {
    loading = true;
    await alertOnFailure(async function () {
      await login(user.email, user.password);
      toast.success("Logged in successfully!");
      goto("/app");
    }).finally(() => {
      loading = false;
    });
  }
</script>

<div
  class="flex min-h-full flex-col justify-center py-12 sm:px-6 lg:px-8 bg-gradient-to-r from-primary-600 to-secondary dark:to-secondarydark"
>
  <div class="mt-10 sm:mx-auto sm:w-full sm:max-w-[480px]">
    <div class="bg-background dark:bg-gray-800 px-6 py-12 shadow sm:rounded-lg sm:px-12">
      <div class="sm:mx-auto sm:w-full sm:max-w-md mb-5">
        <img
          class="h-20 w-auto mx-auto"
          src="/images/logo_primary_typo.png"
          alt="logo"
        />
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

        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <Checkbox>Remember me</Checkbox>
          </div>

          <div class="text-sm leading-6">
            <a href="#" class="font-medium dark:text-gray-200">Forgot password?</a>
          </div>
        </div>

        <div>
          <Button class="w-full" color="alternative" type="submit" on:submit={() => submit()}>Log in</Button>
        </div>
      </form>

      <!-- <div>
        <div class="relative mt-10">
          <div class="absolute inset-0 flex items-center" aria-hidden="true">
            <div class="w-full border-t border-gray-200" />
          </div>
          <div class="relative flex justify-center text-sm font-medium leading-6">
            <span class="bg-background px-6 text-gray-900">or</span>
          </div>
        </div>

        <div class="mt-6 grid grid-cols-2 gap-4">
          <a
            href="#"
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
          </a>

          <a
            href="#"
            class="flex w-full items-center justify-center gap-3 rounded-md bg-blue-500 px-3 py-1.5 text-background focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#24292F]"
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
          </a>
        </div>
      </div> -->

    </div>
  </div>
</div>
