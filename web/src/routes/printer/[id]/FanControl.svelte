<script lang="ts">
	import Slider from '$lib/components/ui/slider/slider.svelte';
	import * as Tabs from '$lib/components/ui/tabs';
	import { setFan, type Printer, type PrinterStatus } from '$lib/sdk';

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

	// Setpoints are seeded once from the first reported speeds, then act purely
	// as user-controlled targets — realtime updates only drive the readout below.
	let setpoints = $state({ part: 0, aux: 0, chamber: 0 });
	let seeded = false;

	$effect(() => {
		if (seeded || !printerState?.fans) return;
		setpoints = {
			part: printerState.fans.part,
			aux: printerState.fans.aux,
			chamber: printerState.fans.chamber
		};
		seeded = true;
	});

	function commit(fan: string, value: number) {
		if (!printer?.serial) return;
		setFan(printer.serial, { fan, speed: value });
	}
</script>

<Tabs.Root bind:value={selected} class="w-full gap-2">
	<Tabs.List class="w-full">
		{#each fans as fan (fan.key)}
			<Tabs.Trigger value={fan.key}>{fan.label}</Tabs.Trigger>
		{/each}
	</Tabs.List>
	{#each fans as fan (fan.key)}
		<Tabs.Content value={fan.key} class="flex items-center gap-3">
			<Slider
				type="single"
				min={0}
				max={100}
				step={10}
				bind:value={setpoints[fan.key]}
				onValueCommit={(v) => commit(fan.key, v as number)}
				class="flex-1"
			/>
			<span class="w-9 text-right text-sm font-medium tabular-nums">{setpoints[fan.key]}%</span>
		</Tabs.Content>
	{/each}
</Tabs.Root>
