<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import PrinterSettingsDialog from '$lib/components/PrinterSettingsDialog.svelte';
	import type { Printer } from '$lib/sdk';
	import { Settings } from '@lucide/svelte';

	type Props = {
		printer: Printer;
	};

	let { printer }: Props = $props();

	let open = $state(false);

	function handleClick(e: MouseEvent) {
		e.stopPropagation();
		open = true;
	}
</script>

<Button variant="ghost" size="icon" onclick={handleClick}>
	<Settings />
</Button>

<PrinterSettingsDialog bind:open {printer} />
