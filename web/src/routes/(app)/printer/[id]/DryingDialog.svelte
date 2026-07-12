<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import Input from '$lib/components/ui/input/input.svelte';
	import Label from '$lib/components/ui/label/label.svelte';
	import { startDrying, stopDrying, type Printer, type PrinterStatus } from '$lib/sdk';
	import { SunIcon } from '@lucide/svelte';
	import { Duration } from 'luxon';

	type AMSUnit = PrinterStatus['ams'][number];

	type Props = {
		open: boolean;
		printer: Printer | undefined;
		amsId: number;
		label: string;
		unit: AMSUnit | undefined;
	};

	let { open = $bindable(), printer, amsId, label, unit }: Props = $props();

	const COOLING_TEMP = 50;

	let temperature = $state(55);
	let duration = $state(8);
	let busy = $state(false);

	let drying = $derived(unit?.drying ?? false);
	let remaining = $derived(
		unit && unit.dryingTime > 0
			? Duration.fromObject({ minutes: unit.dryingTime }).toFormat("h'h' mm'm'")
			: null
	);

	async function start() {
		if (!printer?.serial) return;
		busy = true;
		try {
			await startDrying(printer.serial, amsId.toString(), {
				temperature,
				duration,
				coolingTemp: COOLING_TEMP
			});
			open = false;
		} finally {
			busy = false;
		}
	}

	async function stop() {
		if (!printer?.serial) return;
		busy = true;
		try {
			await stopDrying(printer.serial, amsId.toString());
			open = false;
		} finally {
			busy = false;
		}
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>Dry Filament - {label}</Dialog.Title>
			<Dialog.Description>
				Run the AMS heater to dry the loaded spools. Empty the AMS of any filament you don't want
				heated first.
			</Dialog.Description>
		</Dialog.Header>

		{#if drying}
			<div class="flex flex-col items-center gap-2 py-4">
				<SunIcon class="size-8 text-yellow-500" />
				<p class="font-medium">Drying in progress</p>
				{#if remaining}
					<p class="text-sm text-muted-foreground">{remaining} remaining</p>
				{/if}
			</div>
			<Dialog.Footer>
				<Button variant="outline" onclick={() => (open = false)}>Close</Button>
				<Button variant="destructive" disabled={busy} onclick={stop}>Stop Drying</Button>
			</Dialog.Footer>
		{:else}
			<div class="grid grid-cols-2 gap-3">
				<div class="grid gap-2">
					<Label for="dry-temp">Temperature (°C)</Label>
					<Input id="dry-temp" type="number" min={45} max={85} bind:value={temperature} />
				</div>
				<div class="grid gap-2">
					<Label for="dry-duration">Duration (hours)</Label>
					<Input id="dry-duration" type="number" min={1} max={24} bind:value={duration} />
				</div>
			</div>

			<Dialog.Footer>
				<Button variant="outline" onclick={() => (open = false)}>Cancel</Button>
				<Button disabled={busy || temperature < 45 || duration < 1} onclick={start}>
					Start Drying
				</Button>
			</Dialog.Footer>
		{/if}
	</Dialog.Content>
</Dialog.Root>
