<script lang="ts">
  import { client } from "$lib/pocketbase";
  import { avatarUrlString } from "$lib/stores/avatar";
  import { avatarUrl } from "$lib/utils/user.utils";
  import { Card, Avatar, Button, Input } from "flowbite-svelte";
  import toast from "svelte-french-toast";

  let fileInput: any; // This will be used to store the file input element
  let displayName: string = client.authStore.model?.name ?? "";

  function triggerFileSelect(node: any) {
    // This function is an action that will be attached to the Avatar
    node.addEventListener("click", () => {
      if (fileInput) fileInput.click(); // Trigger the hidden file input's click event
    });

    return {
      destroy() {
        // Cleanup if the component is destroyed
        node.removeEventListener("click", () => {
          if (fileInput) fileInput.click();
        });
      }
    };
  }

  async function updateProfile(event: any) {
    const files = event.target.files;
    if (!files?.length) {
      toast.error("No file selected");
      return;
    }

    let file = files[0];
    const formData = new FormData();
    formData.append("avatar", file);

    try {
      if (!client.authStore.model) {
        throw new Error("No user found");
      }
      await client.collection("users").update(client.authStore.model?.id, formData);
      await client.collection("users").authRefresh();

      toast.success("Successfully updated profile");
      $avatarUrlString = avatarUrl(); // Update the avatar URL to reflect the new avatar
    } catch (error: any) {
      toast.error(error.message);
    }
  }

  async function updateDisplayName() {
    try {
      if (!client.authStore.model) {
        throw new Error("No user found");
      }
      await client.collection("users").update(client.authStore.model?.id, { name: displayName });
      await client.collection("users").authRefresh();
      toast.success("Successfully updated display name");
    } catch (error: any) {
      toast.error(error.message);
    }
  }

  function getValueFromEvent(e: Event) {
    return (e.target as HTMLInputElement).value;
  }
</script>

<div class="relative h-full max-w-screen-2xl mx-auto p-5">
  <div class="flex flex-col h-full w-full space-y-4">
    <Card size="xl">
      <div class="flex flex-col space-y-4">
        <div class="flex justify-between">
          <div class="flex flex-col space-y-2">
            <h3 class="text-xl font-medium text-gray-900 dark:text-gray-400">Avatar</h3>
            <p>
              This is your avatar. <br /> Click on the avatar to upload a custom one from your files.
            </p>
            <p class="text-sm text-gray-500 dark:text-gray-400">
              Accepted file types: .png, .jpg, .jpeg
            </p>
          </div>
          <span use:triggerFileSelect class="w-20">
            <Avatar src={$avatarUrlString} size="lg" class="cursor-pointer " />
          </span>
        </div>
        <input
          type="file"
          id="user_avatar"
          accept="image/*"
          class="hidden"
          bind:this={fileInput}
          on:change={updateProfile}
        />
      </div>
    </Card>
    <Card size="xl">
      <div class="flex flex-col space-y-4">
        <h3 class="text-xl font-medium text-gray-900 dark:text-gray-400">Profile</h3>
        <p>Update your profile information</p>
        <div class="mb-6">
          <span class="text-sm text-gray-500 dark:text-gray-400">Display Name</span>
          <Input
            id="small-input"
            size="sm"
            placeholder="Your display name"
            class="mt-1"
            value={displayName}
            on:input={(e) => (displayName = getValueFromEvent(e))}
          />
          <Button color="primary" size="sm" class="mt-4" on:click={updateDisplayName}>Save</Button>
        </div>
      </div>
    </Card>
  </div>
</div>
