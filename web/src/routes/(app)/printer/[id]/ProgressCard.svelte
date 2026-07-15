<script lang="ts">
	import PlayPause from '$lib/components/actions/PlayPause.svelte';
	import Stop from '$lib/components/actions/Stop.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import Progress from '$lib/components/ui/progress/progress.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import type { Printer, PrinterStatus } from '$lib/sdk';
	import { stageMessage } from '$lib/utils/printer_status';
	import { DateTime, Duration } from 'luxon';
	import { BanIcon, TriangleAlertIcon } from '@lucide/svelte';

	type Props = {
		printer: Printer | undefined;
		state: PrinterStatus | undefined;
	};

	let { state, printer }: Props = $props();

	let durationRemaining = $derived(Duration.fromObject({ minutes: state?.timeRemaining ?? 0 }));
	let timeRemaining = $derived(
		durationRemaining.minutes < 60
			? durationRemaining.toFormat("m'm'")
			: durationRemaining.toFormat("h'h 'm'm'")
	);
	let endTime = $derived(
		DateTime.now().plus(durationRemaining).toLocaleString(DateTime.TIME_SIMPLE)
	);
</script>

<Card class="w-full p-4">
	<div class="flex flex-col gap-3">
		{#if state?.error}
			{@const cancelled = state.error.cancelled}
			<div
				title={state.error.message || undefined}
				class="flex items-start gap-2 rounded-md border p-2 text-sm {cancelled
					? 'border-amber-500/40 bg-amber-500/10 text-amber-700 dark:text-amber-400'
					: 'border-destructive/40 bg-destructive/10 text-destructive'}"
			>
				{#if cancelled}
					<BanIcon class="mt-0.5 size-4 shrink-0" />
				{:else}
					<TriangleAlertIcon class="mt-0.5 size-4 shrink-0" />
				{/if}
				<div class="min-w-0">
					<p class="font-medium">{state.error.summary}</p>
					<p class="text-xs opacity-70">{state.error.code}</p>
				</div>
			</div>
		{/if}

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
			<PlayPause printerSerial={printer?.serial ?? ''} />
			<Stop printerSerial={printer?.serial ?? ''} />
		</div>
	</div>
</Card>
