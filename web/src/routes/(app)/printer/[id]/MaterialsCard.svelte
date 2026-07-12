<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import { unloadMaterial, type Printer, type PrinterStatus } from '$lib/sdk';
	import { cn } from '$lib/utils';
	import { DropletIcon, SunIcon } from '@lucide/svelte';
	import { Duration } from 'luxon';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import DryingDialog from './DryingDialog.svelte';
	import FilamentDialog from './FilamentDialog.svelte';

	type AMSUnit = PrinterStatus['ams'][number];
	type Tray = AMSUnit['trays'][number];

	type Props = {
		state: PrinterStatus | undefined;
		printer: Printer | undefined;
	};

	let { state: printerState, printer }: Props = $props();

	let ams = $derived(printerState?.ams ?? []);
	let externalSpool = $derived(printerState?.externalSpool);

	let dialogOpen = $state(false);
	let target = $state<{ amsId: number; trayId: number; label: string; tray: Tray }>();

	function editTray(amsId: number, trayId: number, label: string, tray: Tray) {
		target = { amsId, trayId, label, tray };
		dialogOpen = true;
	}

	let dryingDialogOpen = $state(false);
	let dryingTarget = $state<{ amsId: number; label: string }>();

	function openDrying(amsId: number, label: string) {
		dryingTarget = { amsId, label };
		dryingDialogOpen = true;
	}

	let dryingUnit = $derived(ams.find((unit) => unit.id === dryingTarget?.amsId));

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

{#if target}
	<FilamentDialog
		bind:open={dialogOpen}
		{printer}
		amsId={target.amsId}
		trayId={target.trayId}
		label={target.label}
		tray={target.tray}
	/>
{/if}

{#if dryingTarget}
	<DryingDialog
		bind:open={dryingDialogOpen}
		{printer}
		amsId={dryingTarget.amsId}
		label={dryingTarget.label}
		unit={dryingUnit}
	/>
{/if}

{#snippet humidityCapsule(unit: AMSUnit)}
	{#if unit.dryingTime > 0}
		{@const time = Duration.fromObject({ minutes: unit.dryingTime }).toFormat('h:mm')}
		<SunIcon class="size-4 text-yellow-500" />
		<Separator orientation="vertical" />
		<p>{time}</p>
	{:else}
		<DropletIcon class="size-4 text-blue-500" />
		<Separator orientation="vertical" />
		<p>{unit.humidity}%</p>
	{/if}
{/snippet}

{#snippet amsUnit(unit: AMSUnit, label: string)}
	<div class="flex w-full flex-col gap-2 rounded-lg border border-foreground/10 p-2">
		<div class="flex items-center justify-between">
			<p class="text-sm font-medium">{label}</p>
			{#if unit.supportsDrying}
				<button
					type="button"
					onclick={() => openDrying(unit.id, label)}
					class="flex h-6 cursor-pointer items-center gap-1.5 rounded-full bg-muted px-2 py-1 transition hover:ring-2 hover:ring-ring focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
				>
					{@render humidityCapsule(unit)}
				</button>
			{:else}
				<div class="flex h-6 items-center gap-1.5 rounded-full bg-muted px-2 py-1">
					{@render humidityCapsule(unit)}
				</div>
			{/if}
		</div>
		<div class="grid grid-cols-4 gap-2">
			{#each unit.trays as tray, i (tray.id)}
				{@render trayCard(tray, {
					slotLabel: `A${i + 1}`,
					amsId: unit.id,
					trayId: tray.id,
					title: `${label} · A${i + 1}`
				})}
			{/each}
		</div>
	</div>
{/snippet}

{#snippet externalUnit(tray: Tray)}
	<div class="w-full rounded-lg border border-foreground/10 p-2">
		<p class="mb-2 text-center text-sm font-medium">Ext</p>
		<div class="flex justify-center">
			<div class="w-1/4">
				{@render trayCard(tray, {
					slotLabel: '',
					amsId: 255,
					trayId: 254,
					title: 'External Spool'
				})}
			</div>
		</div>
	</div>
{/snippet}

{#snippet trayCard(
	tray: Tray,
	opts: { slotLabel: string; amsId: number; trayId: number; title: string }
)}
	<div class="flex flex-col items-center gap-1">
		<button
			type="button"
			onclick={() => editTray(opts.amsId, opts.trayId, opts.title, tray)}
			class="flex aspect-3/4 w-full cursor-pointer flex-col items-center justify-center rounded-md border px-2 py-3 transition hover:ring-2 hover:ring-ring focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
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
		</button>
		{#if opts.slotLabel}
			<span
				class={cn(
					'text-[10px] font-medium text-muted-foreground',
					tray.loaded ? 'text-primary' : 'text-muted-foreground'
				)}
			>
				{opts.slotLabel}
			</span>
		{/if}
	</div>
{/snippet}
