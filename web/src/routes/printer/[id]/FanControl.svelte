<script lang="ts">
	import Slider from '$lib/components/ui/slider/slider.svelte';
	import * as Tabs from '$lib/components/ui/tabs';
	import { setFan, type Printer, type PrinterStatus } from '$lib/sdk';
	import { SyncedControl } from '$lib/synced-control.svelte';

	type Props = {
		state: PrinterStatus | undefined;
		printer: Printer | undefined;
	};

	let { state: printerState, printer }: Props = $props();

	const fans = [
		{ key: 'part', label: 'Part' },
		{ key: 'aux', label: 'Aux' },
		{ key: 'chamber', label: 'Chamber' }
	] as const;

	let selected = $state('part');

	// One synced control per fan: seeds from the reported speed, applies on
	// commit, and holds the setpoint until the printer confirms it. The report is
	// rounded to the slider step so confirmation matches exactly — the printer
	// only resolves fan speed to 0-15 gears.
	const controls = fans.map(
		(fan) =>
			new SyncedControl<number>({
				initial: 0,
				reported: () => {
					const value = printerState?.fans?.[fan.key];
					return value === undefined ? undefined : Math.round(value / 10) * 10;
				},
				apply: (speed) => {
					if (printer?.serial) setFan(printer.serial, { fan: fan.key, speed });
				}
			})
	);
</script>

<Tabs.Root bind:value={selected} class="w-full gap-2">
	<Tabs.List class="w-full">
		{#each fans as fan (fan.key)}
			<Tabs.Trigger value={fan.key}>{fan.label}</Tabs.Trigger>
		{/each}
	</Tabs.List>
	{#each fans as fan, i (fan.key)}
		<Tabs.Content value={fan.key} class="flex items-center gap-3">
			<Slider
				type="single"
				min={0}
				max={100}
				step={10}
				bind:value={controls[i].current}
				onValueCommit={controls[i].set}
				class="flex-1"
			/>
			<span class="w-9 text-right text-sm font-medium tabular-nums">{controls[i].current}%</span>
		</Tabs.Content>
	{/each}
</Tabs.Root>
