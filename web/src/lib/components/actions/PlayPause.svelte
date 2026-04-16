<script lang="ts">
	import { printerManager } from '$lib/managers/printers.manager.svelte';
	import { pausePrint, resumePrint } from '$lib/sdk';
	import { PauseIcon, PlayIcon } from '@lucide/svelte';
	import Button from '../ui/button/button.svelte';

	type Props = {
		printerSerial: string;
	};

	let { printerSerial = $bindable() }: Props = $props();
	let printerState = $derived(printerManager.printerState.get(printerSerial)!);

	async function resume() {
		if (!printerSerial) return;
		await resumePrint(printerSerial);
	}

	async function pause() {
		if (!printerSerial) return;
		await pausePrint(printerSerial);
	}

	let canPause = $derived(printerState?.state === 'RUNNING');
	let canResume = $derived(printerState?.state === 'PAUSE');
</script>

{#if canPause}
	<Button variant="outline" size="sm" class="flex-1" onclick={pause}>
		<PauseIcon fill="#fff" /> Pause
	</Button>
{:else}
	<Button variant="outline" size="sm" class="flex-1" onclick={resume} disabled={!canResume}>
		<PlayIcon fill="#fff" /> Resume
	</Button>
{/if}
