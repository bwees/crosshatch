<script lang="ts">
	import { page } from '$app/state';
	import PlayPause from '$lib/components/actions/PlayPause.svelte';
	import Stop from '$lib/components/actions/Stop.svelte';
	import AMSCard from '$lib/components/AMSCard.svelte';
	import Go2RTCPlayer from '$lib/components/Go2RTCPlayer.svelte';
	import PrinterBuildPlate from '$lib/components/icons/PrinterBuildPlate.svelte';
	import PrinterChamber from '$lib/components/icons/PrinterChamber.svelte';
	import PrinterNozzle from '$lib/components/icons/PrinterNozzle.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import Progress from '$lib/components/ui/progress/progress.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { printerManager } from '$lib/managers/printers.manager.svelte';
	import { stageMessage } from '$lib/utils/printer_status';
	import { DateTime, Duration } from 'luxon';

	const printerSerial = $derived(page.params.id);
	let printer = $derived(printerManager.printers.get(printerSerial || ''));
	let state = $derived(printerManager.printerState.get(printerSerial || ''));

	let durationRemaining = $derived(Duration.fromObject({ minutes: state?.timeRemaining ?? 0 }));
	let timeRemaining = $derived(
		durationRemaining.hours == 0
			? durationRemaining.toFormat("m'm'")
			: durationRemaining.toFormat("h'h 'm'm'")
	);
	let endTime = $derived(
		DateTime.now().plus(durationRemaining).toLocaleString(DateTime.TIME_SIMPLE)
	);
</script>

<svelte:head>
	<title>{printer?.name}</title>
</svelte:head>

<div class="grid h-full w-full grid-cols-1 place-items-start content-start gap-4 md:grid-cols-2">
	<div class="flex w-full flex-col gap-4">
		<Card class="w-full p-0">
			<Go2RTCPlayer url={`ws://${page.url.host}/api/go2rtc?src=${printerSerial}`} />
		</Card>
		<Card class="w-full p-4">
			<div class="flex flex-col gap-3">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-xl font-bold text-primary">{state?.progress ?? '--'}%</p>
						<p class="text-sm">{stageMessage(state?.stage ?? -1)}</p>
					</div>
					{#if state?.state === 'RUNNING'}
						<div class="text-right">
							<p class="text-muted-foreground">{timeRemaining} remaining</p>
							<p class="text-muted-foreground">{endTime}</p>
						</div>
					{/if}
				</div>

				<Progress value={state?.progress ?? 0} max={100} />

				<Separator />

				<div class="flex items-center gap-2">
					<PlayPause printerSerial={printerSerial!} />
					<Stop printerSerial={printerSerial!} />
				</div>
			</div>
		</Card>
	</div>
	<div class="flex w-full flex-col gap-4">
		<div class="grid grid-cols-3 gap-2">
			<Card class="gap-2 p-4">
				<div class="flex items-center gap-2">
					<PrinterNozzle />
					<span class="text-sm">Nozzle</span>
				</div>
				<p class="font-mono text-xl font-bold text-muted-foreground">
					{state?.nozzle.temperature} / {state?.nozzle.targetTemperature}
				</p>
			</Card>

			<Card class="gap-2 p-4">
				<div class="flex items-center gap-2">
					<PrinterBuildPlate />
					<span class="text-sm">Build Plate</span>
				</div>
				<p class="font-mono text-xl font-bold text-muted-foreground">
					{state?.buildPlate.temperature} / {state?.buildPlate.targetTemperature}
				</p>
			</Card>
			<Card class="gap-2 p-4">
				<div class="flex items-center gap-2">
					<PrinterChamber />
					<span class="text-sm">Chamber</span>
				</div>
				<p class="font-mono text-xl font-bold text-muted-foreground">
					{state?.chamber.temperature ?? '--'} / {state?.chamber.targetTemperature ?? '--'}
				</p>
			</Card>
		</div>
		<AMSCard ams={state?.ams ?? []} externalSpool={state?.externalSpool} />
	</div>
</div>
