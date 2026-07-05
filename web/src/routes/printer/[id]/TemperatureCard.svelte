<script lang="ts">
	import PrinterBuildPlate from '$lib/components/icons/PrinterBuildPlate.svelte';
	import PrinterChamber from '$lib/components/icons/PrinterChamber.svelte';
	import PrinterNozzle from '$lib/components/icons/PrinterNozzle.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import type { PrinterStatus } from '$lib/sdk';

	type Props = {
		state: PrinterStatus | undefined;
	};

	let { state }: Props = $props();

	const cellClass =
		'flex min-w-0 flex-col items-center justify-center gap-1 px-1 py-2 font-mono text-muted-foreground lg:flex-row lg:justify-start lg:gap-2 lg:px-3';
	const valueClass = 'text-sm whitespace-nowrap sm:text-base lg:text-lg';
</script>

<Card class="w-full gap-0 p-4 lg:w-1/3">
	<div class="grid grid-cols-3 divide-x divide-border lg:grid-cols-1 lg:divide-x-0 lg:divide-y">
		<div class={cellClass}>
			<PrinterNozzle class="size-5 shrink-0 sm:size-6" />
			<span class={valueClass}>
				<span class="font-bold text-foreground">{state?.nozzle.temperature ?? '--'}</span>
				<span> / {state?.nozzle.targetTemperature ?? '--'} °C</span>
			</span>
		</div>

		<div class={cellClass}>
			<PrinterBuildPlate class="size-5 shrink-0 sm:size-6" />
			<span class={valueClass}>
				<span class="font-bold text-foreground">{state?.buildPlate.temperature ?? '--'}</span>
				<span> / {state?.buildPlate.targetTemperature ?? '--'} °C</span>
			</span>
		</div>

		<div class={cellClass}>
			<PrinterChamber class="size-5 shrink-0 sm:size-6" />
			<span class={valueClass}>
				<span class="font-bold text-foreground">{state?.chamber.temperature ?? '--'}</span>
				{#if state?.chamber.controllable}
					<span> / {state?.chamber.targetTemperature ?? '--'} °C</span>
				{:else}
					<span> °C</span>
				{/if}
			</span>
		</div>
	</div>
</Card>
