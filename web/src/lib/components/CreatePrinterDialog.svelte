<script lang="ts">
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { printerManager } from '$lib/managers/printers.manager.svelte';
	import { createPrinter } from '$lib/sdk';
	import Spinner from './ui/spinner/spinner.svelte';

	type Props = {
		open: boolean;
	};

	let { open = $bindable() }: Props = $props();

	let loading = $state(false);
	let error = $state('');

	let serial = $state('');
	let name = $state('');
	let hostIp = $state('');
	let accessCode = $state('');

	function reset() {
		serial = '';
		name = '';
		hostIp = '';
		accessCode = '';
		error = '';
		loading = false;
	}

	async function handleSubmit() {
		error = '';
		loading = true;

		console.log('Creating printer with', { serial, name, hostIp, accessCode });

		try {
			await createPrinter({ serial, name, hostIp, accessCode });
			await printerManager.refreshPrinters();
			open = false;
			reset();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to add printer';
		} finally {
			loading = false;
		}
	}
</script>

<Dialog.Root
	bind:open
	onOpenChange={(v) => {
		if (!v) reset();
	}}
>
	<Dialog.Content class="sm:max-w-lg">
		<Dialog.Header>
			<Dialog.Title>Add Printer</Dialog.Title>
			<Dialog.Description>Enter the details of your Bambu Lab printer.</Dialog.Description>
		</Dialog.Header>
		<div class="grid gap-4">
			<div class="grid gap-3">
				<Label for="printer-serial">Serial Number</Label>
				<Input
					id="printer-serial"
					bind:value={serial}
					placeholder="e.g. 01P00A000000000"
					required
				/>
			</div>
			<div class="grid gap-3">
				<Label for="printer-name">Name</Label>
				<Input id="printer-name" bind:value={name} placeholder="e.g. My Printer" required />
			</div>
			<div class="grid gap-3">
				<Label for="printer-host">Host IP</Label>
				<Input id="printer-host" bind:value={hostIp} placeholder="e.g. 192.168.1.100" required />
			</div>
			<div class="grid gap-3">
				<Label for="printer-access-code">Access Code</Label>
				<Input
					id="printer-access-code"
					bind:value={accessCode}
					placeholder="Access code"
					required
				/>
			</div>
		</div>
		{#if error}
			<p class="text-sm text-destructive">{error}</p>
		{/if}
		<Dialog.Footer>
			<Dialog.Close type="button" class={buttonVariants({ variant: 'outline' })}>
				Cancel
			</Dialog.Close>
			<Button type="submit" disabled={loading} onclick={handleSubmit}>
				{#if loading}
					<Spinner class="me-2 size-4" />
				{/if}
				Add Printer
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
