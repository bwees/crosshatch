<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import CreatePrinterDialog from '$lib/components/CreatePrinterDialog.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { printerManager } from '$lib/managers/printers.manager.svelte';
	import { stateColor, stateMessage } from '$lib/utils/printer_status';
	import { cn } from '$lib/utils/utils';
	import { Grid, Plus, Printer } from '@lucide/svelte';
	import type { ComponentProps } from 'svelte';

	let { ref = $bindable(null), ...restProps }: ComponentProps<typeof Sidebar.Root> = $props();

	let createDialogOpen = $state(false);
</script>

<Sidebar.Root {...restProps} bind:ref>
	<Sidebar.Content>
		<Sidebar.Header>
			<div class="flex items-center gap-2 rounded-lg border border-secondary p-1">
				<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
					<Grid class="text-primary" />
				</div>
				<p class="text-lg font-bold">Crosshatch</p>
			</div>
		</Sidebar.Header>

		<Sidebar.MenuItem class="mx-2">
			<Sidebar.MenuButton onclick={() => goto(resolve('/'))}>
				<Printer />All Printers
			</Sidebar.MenuButton>
		</Sidebar.MenuItem>

		<Sidebar.Group class="gap-2">
			<Sidebar.GroupLabel class="text-sm">Printers</Sidebar.GroupLabel>
			<Sidebar.GroupAction
				title="Add Printer"
				onclick={() => (createDialogOpen = true)}
				class="cursor-pointer"
			>
				<Plus /> <span class="sr-only">Add Printer</span>
			</Sidebar.GroupAction>
			<Sidebar.Separator />

			{#each printerManager.printers.entries() as [serial, printer] (serial)}
				{@const printerState = printerManager.printerState.get(serial)}
				<Sidebar.MenuItem>
					<Sidebar.MenuButton
						size="lg"
						onclick={() => goto(resolve('/printer/[id]', { id: printer.serial }))}
					>
						<div class="flex items-center">
							<div
								class={cn(
									'me-4 h-3 w-3 rounded-full',
									stateColor(printerState?.state ?? 'UNKNOWN')
								)}
							></div>
							<div>
								<p class="text-md font-bold">{printer.name}</p>
								<p class="text-xs text-muted-foreground">
									{stateMessage(printerState?.state ?? 'UNKNOWN')}
								</p>
							</div>
						</div>
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			{/each}
		</Sidebar.Group>
	</Sidebar.Content>
	<Sidebar.Rail />
</Sidebar.Root>

<CreatePrinterDialog bind:open={createDialogOpen} />
