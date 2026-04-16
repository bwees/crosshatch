<script lang="ts">
	import { resolve } from '$app/paths';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { printerManager } from '$lib/managers/printers.manager.svelte';
	import type { ComponentProps } from 'svelte';

	let { ref = $bindable(null), ...restProps }: ComponentProps<typeof Sidebar.Root> = $props();
</script>

<Sidebar.Root {...restProps} bind:ref>
	<Sidebar.Content>
		<Sidebar.Group>
			<Sidebar.GroupLabel class="px-4 py-2 text-sm font-semibold text-muted-foreground">
				Printers
			</Sidebar.GroupLabel>
			{#each printerManager.printers.entries() as [serial, printer] (serial)}
				<Sidebar.MenuItem>
					<Sidebar.MenuButton>
						<a
							href={resolve('/printer/[id]', { id: printer.serial })}
							class="flex items-center gap-3"
						>
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
