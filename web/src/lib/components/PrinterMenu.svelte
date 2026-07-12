<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { printerManager } from '$lib/managers/printers.manager.svelte';
	import { deletePrinter, type Printer } from '$lib/sdk';
	import { cn } from '$lib/utils.js';
	import { EllipsisVertical, Pencil, Trash2 } from '@lucide/svelte';
	import EditPrinterDialog from './EditPrinterDialog.svelte';

	type Props = {
		printer: Printer;
	};

	let { printer }: Props = $props();

	let editOpen = $state(false);

	async function handleDelete() {
		if (!confirm(`Remove printer "${printer.name}"?`)) return;
		await deletePrinter(printer.serial);
		await printerManager.refreshPrinters();
		if (page.params.id === printer.serial) {
			goto(resolve('/(app)'));
		}
	}
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger
		onclick={(e) => e.stopPropagation()}
		class={cn(buttonVariants({ variant: 'ghost', size: 'icon' }))}
	>
		<EllipsisVertical />
	</DropdownMenu.Trigger>
	<DropdownMenu.Content align="end">
		<DropdownMenu.Item onSelect={() => (editOpen = true)}>
			<Pencil />
			Edit
		</DropdownMenu.Item>
		<DropdownMenu.Item variant="destructive" onSelect={handleDelete}>
			<Trash2 />
			Remove
		</DropdownMenu.Item>
	</DropdownMenu.Content>
</DropdownMenu.Root>

<EditPrinterDialog bind:open={editOpen} {printer} />
