<script lang="ts">
	import PrinterBuildPlate from '$lib/components/icons/PrinterBuildPlate.svelte';
	import PrinterChamber from '$lib/components/icons/PrinterChamber.svelte';
	import PrinterNozzle from '$lib/components/icons/PrinterNozzle.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import type { PrinterStatus } from '$lib/sdk';
	import { cn } from '$lib/utils';

	type Props = {
		state: PrinterStatus | undefined;
	};

	let { state }: Props = $props();

	const cellClass =
		'flex min-w-0 flex-col items-center justify-center gap-1 px-1 py-2 font-mono text-muted-foreground sm:flex-row sm:gap-3 sm:px-4';
	const iconClass = 'size-6 shrink-0 sm:size-8';
	const valueClass = 'text-base sm:text-xl lg:text-2xl';
	const unitClass = 'whitespace-nowrap';
</script>

<Card class="w-full gap-0 p-2">
	<div class="grid grid-cols-3 divide-x divide-border">
		<div class={cellClass}>
			<PrinterNozzle class={cn(iconClass, 'mt-2')} />
			<span class={valueClass}>
				<span class="font-bold text-foreground">{state?.nozzle.temperature ?? '--'}</span>
				<span class={unitClass}> / {state?.nozzle.targetTemperature ?? '--'} °C</span>
			</span>
		</div>

		<div class={cellClass}>
			<PrinterBuildPlate class={cn(iconClass, 'mb-1.5')} />
			<span class={valueClass}>
				<span class="font-bold text-foreground">{state?.buildPlate.temperature ?? '--'}</span>
				<span class={unitClass}> / {state?.buildPlate.targetTemperature ?? '--'} °C</span>
			</span>
		</div>

		<div class={cellClass}>
			<PrinterChamber class={cn(iconClass, 'sm:size-7')} />
			<span class={valueClass}>
				<span class="font-bold text-foreground">{state?.chamber.temperature ?? '--'}</span>
				{#if state?.chamber.controllable}
					<span class={unitClass}> / {state?.chamber.targetTemperature ?? '--'} °C</span>
				{:else}
					<span class={unitClass}> °C</span>
				{/if}
			</span>
		</div>
	</div>
</Card>
