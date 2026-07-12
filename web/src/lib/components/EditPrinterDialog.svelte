<script lang="ts">
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { printerManager } from '$lib/managers/printers.manager.svelte';
	import { apiErrorMessage, updatePrinter, type Printer } from '$lib/sdk';
	import Spinner from './ui/spinner/spinner.svelte';

	type Props = {
		open: boolean;
		printer: Printer;
	};

	let { open = $bindable(), printer }: Props = $props();

	let loading = $state(false);
	let error = $state('');

	let name = $state('');
	let hostIp = $state('');
	let accessCode = $state('');

	$effect(() => {
		if (open) {
			name = printer.name;
			hostIp = printer.hostIp;
			accessCode = printer.accessCode;
			error = '';
			loading = false;
		}
	});

	async function handleSubmit() {
		error = '';
		loading = true;

		try {
			await updatePrinter(printer.serial, { name, hostIp, accessCode });
			await printerManager.refreshPrinters();
			open = false;
		} catch (err) {
			error = apiErrorMessage(err, 'Failed to update printer');
		} finally {
			loading = false;
		}
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-lg">
		<Dialog.Header>
			<Dialog.Title>Edit Printer</Dialog.Title>
			<Dialog.Description>Update the details of your Bambu Lab printer.</Dialog.Description>
		</Dialog.Header>
		<div class="grid gap-4">
			<div class="grid gap-3">
				<Label for="edit-printer-serial">Serial Number</Label>
				<Input id="edit-printer-serial" value={printer.serial} disabled />
			</div>
			<div class="grid gap-3">
				<Label for="edit-printer-name">Name</Label>
				<Input id="edit-printer-name" bind:value={name} placeholder="My Printer" required />
			</div>
			<div class="grid gap-3">
				<Label for="edit-printer-host">Host IP</Label>
				<Input id="edit-printer-host" bind:value={hostIp} placeholder="192.168.1.100" required />
			</div>
			<div class="grid gap-3">
				<Label for="edit-printer-access-code">Access Code</Label>
				<Input
					id="edit-printer-access-code"
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
			<Dialog.Close type="button" class={buttonVariants({ variant: 'outline' })}
				>Cancel</Dialog.Close
			>
			<Button type="submit" disabled={loading} onclick={handleSubmit}>
				{#if loading}
					<Spinner class="me-2 size-4" />
				{/if}
				Save Changes
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
