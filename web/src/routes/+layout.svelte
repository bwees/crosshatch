<script lang="ts">
	import { page } from '$app/state';
	import favicon from '$lib/assets/favicon.svg';
	import AppSidebar from '$lib/components/Sidebar.svelte';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { printerManager } from '$lib/managers/printers.manager.svelte';
	import { onMount } from 'svelte';
	import './layout.css';

	let { children } = $props();

	onMount(() => {
		printerManager.initialize();
	});

	let selectedPrinter = $derived(printerManager.printers.get(page.params.id || '') ?? null);
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<Sidebar.Provider>
	<AppSidebar />
	<Sidebar.Inset>
		<header class="flex h-16 shrink-0 items-center gap-2 border-b px-4">
			<Sidebar.Trigger class="-ms-1" />
			<Separator orientation="vertical" class="me-2 h-4" />
			<div>
				<p>{selectedPrinter?.name}</p>
				<p class="text-xs text-muted-foreground">{selectedPrinter?.serial}</p>
			</div>
		</header>
		<div class="flex-1 p-4">
			{@render children()}
		</div>
	</Sidebar.Inset>
</Sidebar.Provider>
