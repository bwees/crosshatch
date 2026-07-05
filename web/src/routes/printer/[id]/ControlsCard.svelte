<script lang="ts">
	import DetentSlider from '$lib/components/DetentSlider.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import Switch from '$lib/components/ui/switch/switch.svelte';
	import { setLight, setPrintSpeed, type Printer, type PrinterStatus } from '$lib/sdk';
	import { SyncedControl } from '$lib/synced-control.svelte';
	import { FanIcon, LightbulbIcon, RabbitIcon } from '@lucide/svelte';
	import FanControl from './FanControl.svelte';

	type Props = {
		state: PrinterStatus | undefined;
		printer: Printer | undefined;
	};

	let { state: printerState, printer }: Props = $props();

	const isPrinting = $derived(printerState?.state === 'RUNNING');

	const light = new SyncedControl<boolean>({
		initial: false,
		reported: () => printerState?.chamberLight,
		apply: (on) => {
			if (printer?.serial) setLight(printer.serial, { state: on });
		}
	});

	// Slider index 0-3 maps to Bambu speed levels 1-4.
	const speed = new SyncedControl<number>({
		initial: 1,
		reported: () =>
			printerState?.speedLevel === undefined ? undefined : printerState.speedLevel - 1,
		apply: (index) => {
			if (printer?.serial) setPrintSpeed(printer.serial, { level: index + 1 });
		}
	});
</script>

<Card class="w-full gap-3 p-4">
	<div class="flex items-center justify-between">
		<p class="flex items-center gap-2">
			<LightbulbIcon class="size-5" />
			Light
		</p>
		<Switch size="lg" bind:checked={light.current} onCheckedChange={light.set} />
	</div>

	<Separator />

	<div class="flex flex-col gap-2">
		<p class="flex items-center gap-2">
			<FanIcon class="size-5" />
			Fans
		</p>
		<FanControl state={printerState} {printer} />
	</div>

	<Separator />

	<div class="flex flex-col gap-2">
		<p class="flex items-center gap-2">
			<RabbitIcon class="size-5" />
			Print Speed
		</p>
		<DetentSlider
			labels={['Silent', 'Standard', 'Sport', 'Ludicrous']}
			bind:value={speed.current}
			onValueCommit={speed.set}
			disabled={!isPrinting}
		/>
	</div>
</Card>
