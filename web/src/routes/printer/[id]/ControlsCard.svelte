<script lang="ts">
	import DetentSlider from '$lib/components/DetentSlider.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import Switch from '$lib/components/ui/switch/switch.svelte';
	import { setLight, type PrinterDto, type PrinterStatusDto } from '$lib/sdk';
	import { LightbulbIcon } from '@lucide/svelte';

	type Props = {
		state: PrinterStatusDto | undefined;
		printer: PrinterDto | undefined;
	};

	let { state: printerState, printer }: Props = $props();

	let speed = $state(1);
	const isPrinting = $derived(printerState?.state === 'RUNNING');

	// svelte-ignore state_referenced_locally
	// light state sometimes takes some time to update after toggling
	let chamberLightOn = $state(printerState?.chamberLight ?? false);

	$effect(() => {
		setLight(printer?.serial ?? '', { state: chamberLightOn });
	});
</script>

<Card class="w-2/3 gap-3 p-4">
	<div class="grid grid-cols-2 items-center justify-between">
		<div class="flex flex-col items-start gap-2">
			<p class="flex">
				<LightbulbIcon />
				Light
			</p>
			<Switch size="lg" bind:checked={chamberLightOn} />
		</div>
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
