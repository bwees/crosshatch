<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import Spinner from '$lib/components/ui/spinner/spinner.svelte';
	import { auth } from '$lib/managers/auth.svelte';
	import { apiErrorMessage } from '$lib/sdk';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const setupMode = $derived(data.setupRequired);

	let username = $state('');
	let password = $state('');
	let loading = $state(false);
	let error = $state('');

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		error = '';
		loading = true;

		try {
			await (setupMode ? auth.setup(username, password) : auth.login(username, password));
			await goto('/');
		} catch (err) {
			error = apiErrorMessage(err);
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>{setupMode ? 'Set up Crosshatch' : 'Sign in'}</title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center p-4">
	<Card.Root class="w-full sm:max-w-md">
		<Card.Header>
			<Card.Title>{setupMode ? 'Create admin account' : 'Sign in'}</Card.Title>
			<Card.Description>
				{setupMode
					? 'No accounts exist yet. Create the first administrator to get started.'
					: 'Enter your credentials to access Crosshatch.'}
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<form class="grid gap-4" onsubmit={handleSubmit}>
				<div class="grid gap-3">
					<Label for="username">Username</Label>
					<Input id="username" bind:value={username} autocomplete="username" required />
				</div>
				<div class="grid gap-3">
					<Label for="password">Password</Label>
					<Input
						id="password"
						type="password"
						bind:value={password}
						autocomplete={setupMode ? 'new-password' : 'current-password'}
						required
					/>
				</div>
				{#if error}
					<p class="text-sm text-destructive">{error}</p>
				{/if}
				<Button type="submit" class="w-full" disabled={loading}>
					{#if loading}
						<Spinner class="me-2 size-4" />
					{/if}
					{setupMode ? 'Create account' : 'Sign in'}
				</Button>
			</form>
		</Card.Content>
	</Card.Root>
</div>
