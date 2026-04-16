<script lang="ts">
	import { page } from '$app/state';
	import Go2RTCPlayer from '$lib/components/Go2RTCPlayer.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import Progress from '$lib/components/ui/progress/progress.svelte';
	import { printerManager } from '$lib/managers/printers.manager.svelte';

	const printerSerial = page.params.id;
	let printer = $derived(printerManager.printerState.get(printerSerial || ''));
</script>

<div class="grid h-full w-full grid-cols-1 place-items-start gap-4 md:grid-cols-2">
	<div class="flex w-full flex-col gap-4">
		<Card class="w-full p-0">
			<Go2RTCPlayer url={`ws://${page.url.host}/api/go2rtc?src=${printerSerial}`} />
		</Card>
		<Card class="w-full p-4">
			<div class="flex flex-col gap-3">
				<div class="flex items-center justify-between">
					<p class="text-xl font-bold text-primary">{printer?.progress}%</p>
					<p class="text-muted-foreground">~6h5m remaining</p>
				</div>

				<Progress value={65} max={100} />
			</div>
		</Card>
	</div>
</div>
