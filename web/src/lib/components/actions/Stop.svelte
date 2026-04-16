<script lang="ts">
	import { printerManager } from '$lib/managers/printers.manager.svelte';
	import { stopPrint } from '$lib/sdk';
	import { CircleXIcon } from '@lucide/svelte';
	import Button from '../ui/button/button.svelte';

	type Props = {
		printerSerial: string;
	};

	let { printerSerial = $bindable() }: Props = $props();
	let printerState = $derived(printerManager.printerState.get(printerSerial)!);

	async function stop() {
		if (!printerSerial) return;
		await stopPrint(printerSerial);
	}

	let canStop = $derived(
		['PAUSE', 'RUNNING', 'SLICING', 'PREPARE'].includes(printerState?.state ?? '')
	);
</script>

<Button variant="outline" size="sm" class="flex-1" onclick={stop} disabled={!canStop}>
	<CircleXIcon /> Stop
</Button>
