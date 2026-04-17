<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import Card from '$lib/components/ui/card/card.svelte';
	import { printerManager } from '$lib/managers/printers.manager.svelte';
	import { Printer } from '@lucide/svelte';
</script>

<svelte:head>
	<title>Printers</title>
</svelte:head>

<div class="lg-grid-cols-3 grid grid-cols-1 gap-4 md:grid-cols-2">
	{#each printerManager.printers.entries() as [serial, printer] (serial)}
		<Card
			class="w-full cursor-pointer p-4 transition-colors hover:bg-secondary/70"
			onclick={() => goto(resolve('/printer/[id]', { id: printer.serial }))}
		>
			<div class="flex items-center gap-4">
				<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
					<Printer class="text-primary" />
				</div>
				<div>
					<p class="text-md font-bold">{printer.name}</p>
					<p class="text-xs text-muted-foreground">{printer.serial}</p>
				</div>
			</div>
		</Card>
	{/each}
</div>
