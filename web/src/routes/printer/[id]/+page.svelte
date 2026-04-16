<script lang="ts">
	import { page } from '$app/state';
	import PlayPause from '$lib/components/actions/PlayPause.svelte';
	import Stop from '$lib/components/actions/Stop.svelte';
	import Go2RTCPlayer from '$lib/components/Go2RTCPlayer.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import Progress from '$lib/components/ui/progress/progress.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { printerManager } from '$lib/managers/printers.manager.svelte';
	import { stageMessage } from '$lib/utils/printer_status';
	import { DateTime, Duration } from 'luxon';

	const printerSerial = $derived(page.params.id);
	let printer = $derived(printerManager.printerState.get(printerSerial || ''));

	let durationRemaining = $derived(Duration.fromObject({ minutes: printer?.timeRemaining ?? 0 }));
	let timeRemaining = $derived(
		durationRemaining.hours == 0
			? durationRemaining.toFormat("m'm'")
			: durationRemaining.toFormat("h'h 'm'm'")
	);
	let endTime = $derived(
		DateTime.now().plus(durationRemaining).toLocaleString(DateTime.TIME_SIMPLE)
	);
</script>

<div class="grid h-full w-full grid-cols-1 place-items-start gap-4 md:grid-cols-2">
	<div class="flex w-full flex-col gap-4">
		<Card class="w-full p-0">
			<Go2RTCPlayer url={`ws://${page.url.host}/api/go2rtc?src=${printerSerial}`} />
		</Card>
		<Card class="w-full p-4">
			<div class="flex flex-col gap-3">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-xl font-bold text-primary">{printer?.progress ?? '--'}%</p>
						<p class="text-sm">{stageMessage(printer?.stage ?? -1)}</p>
					</div>
					{#if printer?.state === 'RUNNING'}
						<div class="text-right">
							<p class="text-muted-foreground">{timeRemaining} remaining</p>
							<p class="text-muted-foreground">{endTime}</p>
						</div>
					{/if}
				</div>

				<Progress value={printer?.progress ?? 0} max={100} />

				<Separator />

				<div class="flex items-center gap-2">
					<PlayPause printerSerial={printerSerial!} />
					<Stop printerSerial={printerSerial!} />
				</div>
			</div>
		</Card>
	</div>
</div>
