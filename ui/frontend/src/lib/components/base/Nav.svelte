<script lang="ts">
  import { client, logout } from '$lib/pocketbase';
  import { avatarUrl } from '$lib/utils/user.utils';
  import { Navbar, NavBrand, NavLi, NavUl, NavHamburger, Avatar, Dropdown, DropdownItem, DropdownHeader, DropdownDivider } from 'flowbite-svelte';

  let avatarUrlString: any = avatarUrl();


</script>

<Navbar class="bg-primary dark:bg-primary">
  <NavBrand href="/app">
    <img src="/images/logo_background_typo.png" class="mr-3 h-16" alt="Flowbite Logo" />
  </NavBrand>
  <div class="flex items-center md:order-2 cursor-pointer active:scale-105">
    <Avatar id="avatar-menu" src={avatarUrlString} />
    <NavHamburger class1="w-full md:flex md:w-auto md:order-1" />
  </div>
  <Dropdown placement="bottom" triggeredBy="#avatar-menu">
    <DropdownHeader>
      <span class="block text-sm">{client.authStore.model?.name}</span>
      <span class="block truncate text-sm font-medium">{client.authStore.model?.email}</span>
    </DropdownHeader>
    <DropdownItem>Dashboard</DropdownItem>
    <DropdownItem>Settings</DropdownItem>
    <DropdownDivider />
    <DropdownItem
      on:click={
        () => {
          logout();
        }
      }
    >Sign out</DropdownItem>
  </Dropdown>
</Navbar>
