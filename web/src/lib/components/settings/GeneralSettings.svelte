<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import Spinner from '$lib/components/ui/spinner/spinner.svelte';
	import { printerManager } from '$lib/managers/printers.manager.svelte';
	import { apiErrorMessage, deletePrinter, updatePrinter, type Printer } from '$lib/sdk';
	import { untrack } from 'svelte';

	type Props = {
		printer: Printer;
		onClose: () => void;
	};

	let { printer, onClose }: Props = $props();

	// Snapshot the printer into editable form state; the component remounts each
	// time the dialog opens, so a one-time initial read is intentional.
	let name = $state(untrack(() => printer.name));
	let hostIp = $state(untrack(() => printer.hostIp));
	let accessCode = $state(untrack(() => printer.accessCode));
	let saving = $state(false);
	let error = $state('');

	async function handleSave() {
		error = '';
		saving = true;
		try {
			await updatePrinter(printer.serial, { name, hostIp, accessCode });
			await printerManager.refreshPrinters();
			onClose();
		} catch (err) {
			error = apiErrorMessage(err, 'Failed to update printer');
		} finally {
			saving = false;
		}
	}

	async function handleDelete() {
		if (!confirm(`Remove printer "${printer.name}"?`)) return;
		await deletePrinter(printer.serial);
		await printerManager.refreshPrinters();
		onClose();
		if (page.params.id === printer.serial) {
			goto(resolve('/(app)'));
		}
	}
</script>

<div class="grid gap-4">
	<div class="grid gap-3">
		<Label for="settings-printer-serial">Serial Number</Label>
		<Input id="settings-printer-serial" value={printer.serial} disabled />
	</div>
	<div class="grid gap-3">
		<Label for="settings-printer-name">Name</Label>
		<Input id="settings-printer-name" bind:value={name} placeholder="My Printer" required />
	</div>
	<div class="grid gap-3">
		<Label for="settings-printer-host">Host IP</Label>
		<Input id="settings-printer-host" bind:value={hostIp} placeholder="192.168.1.100" required />
	</div>
	<div class="grid gap-3">
		<Label for="settings-printer-access-code">Access Code</Label>
		<Input
			id="settings-printer-access-code"
			bind:value={accessCode}
			placeholder="Access code"
			required
		/>
	</div>

	{#if error}
		<p class="text-sm text-destructive">{error}</p>
	{/if}

	<div class="flex justify-end">
		<Button type="button" disabled={saving} onclick={handleSave}>
			{#if saving}
				<Spinner class="me-2 size-4" />
			{/if}
			Save Changes
		</Button>
	</div>

	<div class="mt-2 border-t pt-4">
		<div class="flex items-center justify-between gap-4">
			<div>
				<p class="text-sm font-medium">Remove printer</p>
				<p class="text-xs text-muted-foreground">This removes the printer from Crosshatch.</p>
			</div>
			<Button variant="destructive" type="button" onclick={handleDelete}>Remove</Button>
		</div>
	</div>
</div>
