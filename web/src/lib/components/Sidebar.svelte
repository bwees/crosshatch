<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import CreatePrinterDialog from '$lib/components/CreatePrinterDialog.svelte';
	import * as Collapsible from '$lib/components/ui/collapsible/index.js';
	import { useSidebar } from '$lib/components/ui/sidebar/context.svelte.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { printerManager } from '$lib/managers/printers.manager.svelte';
	import { cn } from '$lib/utils';
	import { stateColor, stateMessage } from '$lib/utils/printer_status';
	import { ChevronRight, Grid, Printer } from '@lucide/svelte';
	import type { ComponentProps } from 'svelte';

	let { ref = $bindable(null), ...restProps }: ComponentProps<typeof Sidebar.Root> = $props();

	let createDialogOpen = $state(false);
	let printersOpen = $state(true);
	const sidebar = useSidebar();

	function navigate(path: string) {
		if (sidebar.isMobile) sidebar.setOpenMobile(false);
		goto(path);
	}
</script>

<Sidebar.Root {...restProps} bind:ref>
	<Sidebar.Content>
		<Sidebar.Header>
			<a class="flex items-center gap-2 rounded-lg border border-secondary p-1" href="/">
				<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
					<Grid class="text-primary" />
				</div>
				<p class="text-lg font-bold">Crosshatch</p>
			</a>
		</Sidebar.Header>

		<Sidebar.Group class="gap-2">
			<Sidebar.GroupAction
				title="Add Printer"
				onclick={() => (createDialogOpen = true)}
				class="cursor-pointer"
			/>

			<Collapsible.Root bind:open={printersOpen} class="group/collapsible">
				<Sidebar.MenuItem class="list-none">
					<Collapsible.Trigger>
						{#snippet child({ props })}
							<Sidebar.MenuButton {...props} class="cursor-pointer">
								<Printer />
								<span>Printers</span>
								<ChevronRight
									class={cn(
										'ml-auto transition-transform duration-200',
										printersOpen && 'rotate-90'
									)}
								/>
							</Sidebar.MenuButton>
						{/snippet}
					</Collapsible.Trigger>
				</Sidebar.MenuItem>

				<Collapsible.Content>
					<Sidebar.MenuSub class="mt-2 gap-2">
						{#each printerManager.printers.entries() as [serial, printer] (serial)}
							{@const printerState = printerManager.printerState.get(serial)}
							<Sidebar.MenuSubItem>
								<Sidebar.MenuSubButton
									onclick={() => navigate(resolve('/printer/[id]', { id: printer.serial }))}
									class="h-auto cursor-pointer py-1"
								>
									<div class="flex items-center">
										<div
											class={cn(
												'me-3 h-3 w-3 rounded-full',
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
								</Sidebar.MenuSubButton>
							</Sidebar.MenuSubItem>
						{/each}
					</Sidebar.MenuSub>
				</Collapsible.Content>
			</Collapsible.Root>
		</Sidebar.Group>
	</Sidebar.Content>
	<Sidebar.Rail />
</Sidebar.Root>

<CreatePrinterDialog bind:open={createDialogOpen} />
