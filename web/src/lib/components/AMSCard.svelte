<script lang="ts">
	import Card from '$lib/components/ui/card/card.svelte';
	import type { PrinterStatusDto } from '$lib/sdk';
	import { cn } from '$lib/utils';

	type AMSUnit = PrinterStatusDto['ams'][number];
	type Tray = AMSUnit['trays'][number];

	type Props = {
		ams: AMSUnit[];
		externalSpool?: Tray;
	};

	let { ams = $bindable(), externalSpool = $bindable() }: Props = $props();

	function bambuColorToCss(color?: string): string {
		if (!color) return 'transparent';
		const hex = color.length === 8 ? color.slice(0, 6) : color;
		return `#${hex}`;
	}

	function isLightColor(color?: string): boolean {
		if (!color) return true;
		const hex = color.length === 8 ? color.slice(0, 6) : color;
		if (hex.length !== 6) return true;
		const r = parseInt(hex.slice(0, 2), 16);
		const g = parseInt(hex.slice(2, 4), 16);
		const b = parseInt(hex.slice(4, 6), 16);
		return r * 0.299 + g * 0.587 + b * 0.114 > 150;
	}
</script>

<Card class="gap-2 p-4 pt-3">
	<p class="font-bold">Material</p>
	{#if ams.length === 0 && !externalSpool}
		<p class="text-sm text-muted-foreground">No AMS connected</p>
	{:else}
		<div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
			{#each ams as unit, i (unit.id)}
				{@render amsUnit(unit, `AMS ${i + 1}`)}
			{/each}
			{#if externalSpool}
				{@render externalUnit(externalSpool)}
			{/if}
		</div>
	{/if}
</Card>

{#snippet amsUnit(unit: AMSUnit, label: string)}
	<div class="w-full rounded-lg border border-foreground/10 p-2">
		<p class="mb-2 text-center text-xs font-medium text-muted-foreground">{label}</p>
		<div class="grid grid-cols-4 gap-2">
			{#each unit.trays as tray, i (tray.id)}
				{@render trayCard(tray, `A${i + 1}`)}
			{/each}
		</div>
	</div>
{/snippet}

{#snippet externalUnit(tray: Tray)}
	<div class="w-full rounded-lg border border-foreground/10 p-2">
		<p class="mb-2 text-center text-xs font-medium text-muted-foreground">Ext</p>
		<div class="flex justify-center">
			<div class="w-1/4">
				{@render trayCard(tray, '')}
			</div>
		</div>
	</div>
{/snippet}

{#snippet trayCard(tray: Tray, label: string)}
	<div class="flex flex-col items-center gap-1">
		<div
			class="flex aspect-3/4 w-full flex-col items-center justify-center rounded-md border px-2 py-3"
			style:background-color={bambuColorToCss(tray.color)}
			style:color={isLightColor(tray.color) ? '#000' : '#fff'}
		>
			{#if tray.empty}
				<span class="text-2xl font-semibold text-muted-foreground opacity-60">?</span>
			{:else}
				<span class="text-sm leading-tight font-bold">{tray.material ?? '—'}</span>
				{#if tray.kValue !== undefined}
					<span class="mt-0.5 text-[10px] font-medium">K {tray.kValue.toFixed(3)}</span>
				{/if}
			{/if}
		</div>
		{#if label}
			<span
				class={cn(
					'text-[10px] font-medium text-muted-foreground',
					tray.loaded ? 'text-primary' : 'text-muted-foreground'
				)}>{label}</span
			>
		{/if}
	</div>
{/snippet}
