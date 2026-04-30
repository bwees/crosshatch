<script lang="ts">
	import ProgressCard from './ProgressCard.svelte';

	import TemperatureCard from './TemperatureCard.svelte';

	import { page } from '$app/state';
	import Go2RTCPlayer from '$lib/components/Go2RTCPlayer.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import { printerManager } from '$lib/managers/printers.manager.svelte';
	import ControlsCard from './ControlsCard.svelte';
	import MaterialsCard from './MaterialsCard.svelte';

	const printerSerial = $derived(page.params.id);
	let printer = $derived(printerManager.printers.get(printerSerial || ''));
	let state = $derived(printerManager.printerState.get(printerSerial || ''));
</script>

<svelte:head>
	<title>{printer?.name}</title>
</svelte:head>

<div class="grid h-full w-full grid-cols-1 place-items-start content-start gap-4 md:grid-cols-2">
	<div class="flex w-full flex-col gap-4">
		<Card class="w-full p-0">
			<Go2RTCPlayer url={`ws://${page.url.host}/api/go2rtc?src=${printerSerial}`} />
		</Card>
		<ProgressCard {state} {printer} />
	</div>
	<div class="flex w-full flex-col gap-4">
		<div class="flex gap-4">
			<TemperatureCard {state} />
			<ControlsCard {state} {printer} />
		</div>
		<MaterialsCard {state} {printer} />
	</div>
</div>
