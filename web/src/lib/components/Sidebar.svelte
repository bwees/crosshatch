<script lang="ts">
	import { resolve } from '$app/paths';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { printerManager } from '$lib/managers/printers.manager.svelte';
	import { Plus } from '@lucide/svelte';
	import type { ComponentProps } from 'svelte';

	let { ref = $bindable(null), ...restProps }: ComponentProps<typeof Sidebar.Root> = $props();
</script>

<Sidebar.Root {...restProps} bind:ref>
	<Sidebar.Content>
		<Sidebar.Group class="gap-2">
			<Sidebar.GroupLabel class="text-sm">Printers</Sidebar.GroupLabel>
			<Sidebar.GroupAction
				title="Add Printer"
				onclick={() => alert('Add printer functionality coming soon!')}
				class="cursor-pointer"
			>
				<Plus /> <span class="sr-only">Add Printer</span>
			</Sidebar.GroupAction>
			<Sidebar.Separator />

			{#each printerManager.printers.entries() as [serial, printer] (serial)}
				<Sidebar.MenuItem>
					<Sidebar.MenuButton size="lg">
						<a href={resolve('/printer/[id]', { id: printer.serial })} class="flex items-center">
							<div class="me-4 h-3 w-3 rounded-full bg-green-500"></div>
							<div>
								<p class="text-md font-bold">{printer.name}</p>
								<p class="text-xs text-muted-foreground">{printer.serial}</p>
							</div>
						</a>
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			{/each}
		</Sidebar.Group>
	</Sidebar.Content>
	<Sidebar.Rail />
</Sidebar.Root>
