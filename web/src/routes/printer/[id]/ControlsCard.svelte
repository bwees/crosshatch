<script lang="ts">
	import DetentSlider from '$lib/components/DetentSlider.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import Switch from '$lib/components/ui/switch/switch.svelte';
	import { setLight, type Printer, type PrinterStatus } from '$lib/sdk';
	import { LightbulbIcon } from '@lucide/svelte';

	type Props = {
		state: PrinterStatus | undefined;
		printer: Printer | undefined;
	};

	let { state: printerState, printer }: Props = $props();

	const isPrinting = $derived(printerState?.state === 'RUNNING');

	let speed = $state(1);

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
		<p>Print Speed</p>
		<DetentSlider
			labels={['Silent', 'Standard', 'Sport', 'Ludicrous']}
			bind:value={speed}
			disabled={!isPrinting}
		/>
	</div>
</Card>
