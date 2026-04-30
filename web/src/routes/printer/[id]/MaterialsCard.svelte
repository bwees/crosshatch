<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import { unloadMaterial, type PrinterDto, type PrinterStatusDto } from '$lib/sdk';
	import { cn } from '$lib/utils';
	import { DropletIcon, SunIcon } from '@lucide/svelte';
	import { Duration } from 'luxon';
	import Separator from '../../../lib/components/ui/separator/separator.svelte';

	type AMSUnit = PrinterStatusDto['ams'][number];
	type Tray = AMSUnit['trays'][number];

	type Props = {
		state: PrinterStatusDto | undefined;
		printer: PrinterDto | undefined;
	};

	let { state, printer }: Props = $props();

	let ams = $derived(state?.ams ?? []);
	let externalSpool = $derived(state?.externalSpool);

	let currentlyLoadedUnitID = $derived(
		ams.find((unit) => unit.trays.some((tray) => tray.loaded))?.id ??
			(externalSpool?.loaded ? 255 : 0)
	);

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

	async function unload() {
		if (!currentlyLoadedUnitID) return;
		await unloadMaterial(printer?.serial ?? '', currentlyLoadedUnitID.toString());
	}
</script>

<Card class="gap-2 p-4 pt-3">
	<div class="flex items-center justify-between">
		<p class="font-bold">Material</p>
		<Button
			variant="outline"
			size="sm"
			disabled={!currentlyLoadedUnitID}
			color="primary"
			onclick={unload}>Unload</Button
		>
	</div>
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
	<div class="flex w-full flex-col gap-2 rounded-lg border border-foreground/10 p-2">
		<div class="flex items-center justify-between">
			<p class="text-sm font-medium">{label}</p>
			<div class="flex h-6 items-center gap-1.5 rounded-full bg-muted px-2 py-1">
				{#if unit.dryingTime > 0}
					{@const time = Duration.fromObject({ seconds: unit.dryingTime }).toFormat('h:mm')}
					<SunIcon class="size-4 text-yellow-500" />
					<Separator orientation="vertical" />
					<p>{time}</p>
				{:else}
					<DropletIcon class="size-4 text-blue-500" />
					<Separator orientation="vertical" />
					<p>{unit.humidity}%</p>
				{/if}
			</div>
		</div>
		<div class="grid grid-cols-4 gap-2">
			{#each unit.trays as tray, i (tray.id)}
				{@render trayCard(tray, `A${i + 1}`)}
			{/each}
		</div>
	</div>
{/snippet}

{#snippet externalUnit(tray: Tray)}
	<div class="w-full rounded-lg border border-foreground/10 p-2">
		<p class="mb-2 text-center text-sm font-medium">Ext</p>
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
				)}
			>
				{label}
			</span>
		{/if}
	</div>
{/snippet}
