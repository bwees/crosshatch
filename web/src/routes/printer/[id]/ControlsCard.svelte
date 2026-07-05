<script lang="ts">
	import DetentSlider from '$lib/components/DetentSlider.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import Switch from '$lib/components/ui/switch/switch.svelte';
	import { setLight, setPrintSpeed, type Printer, type PrinterStatus } from '$lib/sdk';
	import { FanIcon, LightbulbIcon, RabbitIcon } from '@lucide/svelte';
	import FanControl from './FanControl.svelte';

	type Props = {
		state: PrinterStatus | undefined;
		printer: Printer | undefined;
	};

	let { state: printerState, printer }: Props = $props();

	const isPrinting = $derived(printerState?.state === 'RUNNING');

	// Slider index 0-3 maps to Bambu speed levels 1-4.
	let speed = $state(1);

	$effect(() => {
		if (printerState?.speedLevel !== undefined) {
			speed = printerState.speedLevel - 1;
		}
	});

	function commitSpeed(index: number) {
		if (!printer?.serial) return;
		setPrintSpeed(printer.serial, { level: index + 1 });
	}

	let chamberLightOn = $state(false);

	// The switch reflects the reported light state, but once the user toggles it
	// we hold that value (ignoring stale reports) until the printer confirms,
	// avoiding flicker without ever pushing a command on load.
	let pendingLight: boolean | null = null;

	$effect(() => {
		const reported = printerState?.chamberLight;
		if (reported === undefined) return;
		if (pendingLight !== null) {
			if (reported === pendingLight) pendingLight = null;
			return;
		}
		chamberLightOn = reported;
	});

	function toggleLight(checked: boolean) {
		pendingLight = checked;
		if (!printer?.serial) return;
		setLight(printer.serial, { state: checked });
	}
</script>

<Card class="w-full gap-3 p-4 lg:w-2/3">
	<div class="flex items-center justify-between">
		<p class="flex items-center gap-2">
			<LightbulbIcon class="size-5" />
			Light
		</p>
		<Switch size="lg" bind:checked={chamberLightOn} onCheckedChange={toggleLight} />
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
			bind:value={speed}
			onValueCommit={commitSpeed}
			disabled={!isPrinting}
		/>
	</div>
</Card>
