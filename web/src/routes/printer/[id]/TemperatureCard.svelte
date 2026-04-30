<script lang="ts">
	import PrinterBuildPlate from '$lib/components/icons/PrinterBuildPlate.svelte';
	import PrinterChamber from '$lib/components/icons/PrinterChamber.svelte';
	import PrinterNozzle from '$lib/components/icons/PrinterNozzle.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import type { PrinterStatusDto } from '$lib/sdk';

	type Props = {
		state: PrinterStatusDto | undefined;
	};

	let { state }: Props = $props();
</script>

<Card class="w-1/3 gap-2 p-4">
	<div class="flex items-center gap-2 font-mono text-xl text-muted-foreground">
		<PrinterNozzle class="mr-1 size-6" />
		<p class="font-bold">{state?.nozzle.temperature ?? '--'}</p>
		<p>/</p>
		<p>{state?.nozzle.targetTemperature ?? '--'} °C</p>
	</div>

	<Separator class="my-2" />

	<div class="flex items-center gap-2 font-mono text-xl text-muted-foreground">
		<PrinterBuildPlate class="mr-1 size-6" />
		<p class="font-bold">{state?.buildPlate.temperature ?? '--'}</p>
		<p>/</p>
		<p>{state?.buildPlate.targetTemperature ?? '--'} °C</p>
	</div>

	<Separator class="my-2" />

	<div class="flex items-center gap-2 font-mono text-xl text-muted-foreground">
		<PrinterChamber class="mr-1 size-6" />
		<p class="font-bold">{state?.chamber.temperature ?? '--'}</p>
		{#if state?.chamber.controllable}
			<p>/</p>
			<p>{state?.chamber.targetTemperature ?? '--'} °C</p>
		{:else}
			<p>°C</p>
		{/if}
	</div>
</Card>
