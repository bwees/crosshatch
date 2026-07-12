<script lang="ts">
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import Card from '$lib/components/ui/card/card.svelte';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import * as Switch from '$lib/components/ui/switch/index.js';
	import Spinner from '$lib/components/ui/spinner/spinner.svelte';
	import { auth } from '$lib/managers/auth.svelte';
	import { apiErrorMessage, createUser, deleteUser, getUsers, type UserDto } from '$lib/sdk';
	import { Plus, Shield, Trash2 } from '@lucide/svelte';
	import { onMount } from 'svelte';

	let users = $state<UserDto[]>([]);
	let dialogOpen = $state(false);

	let username = $state('');
	let password = $state('');
	let isAdmin = $state(false);
	let saving = $state(false);
	let error = $state('');

	onMount(refresh);

	async function refresh() {
		users = await getUsers();
	}

	function resetForm() {
		username = '';
		password = '';
		isAdmin = false;
		error = '';
		saving = false;
	}

	async function handleCreate() {
		error = '';
		saving = true;
		try {
			await createUser({ username, password, isAdmin });
			await refresh();
			dialogOpen = false;
			resetForm();
		} catch (err) {
			error = apiErrorMessage(err, 'Failed to create user');
		} finally {
			saving = false;
		}
	}

	async function handleDelete(user: UserDto) {
		if (!confirm(`Delete user "${user.username}"?`)) return;
		await deleteUser(String(user.id));
		await refresh();
	}
</script>

<svelte:head>
	<title>Users</title>
</svelte:head>

<div class="mb-4 flex items-center justify-between">
	<h1 class="text-2xl font-bold">Users</h1>
	<Button size="icon" variant="outline" onclick={() => (dialogOpen = true)}>
		<Plus />
	</Button>
</div>

<div class="grid gap-3">
	{#each users as user (user.id)}
		<Card class="flex flex-row items-center justify-between p-4">
			<div class="flex items-center gap-2">
				<p class="font-medium">{user.username}</p>
				{#if user.isAdmin}
					<span
						class="inline-flex items-center gap-1 rounded-full bg-primary/10 px-2 py-0.5 text-xs text-primary"
					>
						<Shield class="size-3" /> Admin
					</span>
				{/if}
			</div>
			{#if user.id !== auth.user?.id}
				<Button size="icon" variant="ghost" onclick={() => handleDelete(user)}>
					<Trash2 class="text-destructive" />
				</Button>
			{/if}
		</Card>
	{/each}
</div>

<Dialog.Root
	bind:open={dialogOpen}
	onOpenChange={(v) => {
		if (!v) resetForm();
	}}
>
	<Dialog.Content class="sm:max-w-lg">
		<Dialog.Header>
			<Dialog.Title>Add User</Dialog.Title>
			<Dialog.Description>Create a new account.</Dialog.Description>
		</Dialog.Header>
		<div class="grid gap-4">
			<div class="grid gap-3">
				<Label for="new-username">Username</Label>
				<Input id="new-username" bind:value={username} autocomplete="off" required />
			</div>
			<div class="grid gap-3">
				<Label for="new-password">Password</Label>
				<Input
					id="new-password"
					type="password"
					bind:value={password}
					autocomplete="new-password"
					placeholder="At least 8 characters"
					required
				/>
			</div>
			<div class="flex items-center gap-3">
				<Switch.Root id="new-admin" bind:checked={isAdmin} />
				<Label for="new-admin">Administrator</Label>
			</div>
		</div>
		{#if error}
			<p class="text-sm text-destructive">{error}</p>
		{/if}
		<Dialog.Footer>
			<Dialog.Close type="button" class={buttonVariants({ variant: 'outline' })}
				>Cancel</Dialog.Close
			>
			<Button type="submit" disabled={saving} onclick={handleCreate}>
				{#if saving}
					<Spinner class="me-2 size-4" />
				{/if}
				Add User
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
