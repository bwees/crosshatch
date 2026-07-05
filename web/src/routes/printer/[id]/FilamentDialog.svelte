<script lang="ts">
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import Label from '$lib/components/ui/label/label.svelte';
	import { FILAMENT_COLORS } from '$lib/filaments';
	import { ensureFilaments, filaments } from '$lib/filaments.svelte';
	import { setFilament, type Filament, type Printer, type PrinterStatus } from '$lib/sdk';
	import { cn } from '$lib/utils';
	import { ChevronDownIcon } from '@lucide/svelte';
	import { onMount } from 'svelte';

	type Tray = PrinterStatus['ams'][number]['trays'][number];

	type Props = {
		open: boolean;
		printer: Printer | undefined;
		amsId: number;
		trayId: number;
		label: string;
		tray: Tray | undefined;
	};

	let { open = $bindable(), printer, amsId, trayId, label, tray }: Props = $props();

	let brand = $state<string>('Generic');
	let selectedIdx = $state<string | undefined>(undefined);
	let color = $state('#888888');
	let saving = $state(false);

	onMount(ensureFilaments);

	const presetsForBrand = $derived(filaments.presets.filter((p) => p.brand === brand));
	const selectedPreset = $derived(filaments.presets.find((p) => p.trayInfoIdx === selectedIdx));

	function toHexInput(c?: string): string {
		if (!c) return FILAMENT_COLORS[0];
		const hex = c.length >= 6 ? c.slice(0, 6) : c;
		return `#${hex.toUpperCase()}`;
	}

	function toBambuColor(hex: string): string {
		return (hex.replace('#', '') + 'FF').toUpperCase();
	}

	// Snapshot the loaded filament into the form when the dialog opens, so live
	// status updates don't clobber in-progress edits. onOpenChange doesn't fire
	// for programmatic opens, so key off the open transition instead.
	let wasOpen = false;
	$effect(() => {
		if (open && !wasOpen) {
			// Prefer an exact match on the reported filament id, falling back to the
			// material family when the loaded spool isn't a known preset.
			const exact = filaments.presets.find((p) => p.trayInfoIdx === tray?.trayInfoIdx);
			const match =
				exact ??
				filaments.presets.find(
					(p) =>
						p.brand === (tray?.brand?.toLowerCase().includes('bambu') ? 'Bambu' : 'Generic') &&
						p.trayType === tray?.material
				);
			brand = match?.brand ?? 'Generic';
			selectedIdx = match?.trayInfoIdx;
			color = toHexInput(tray?.color);
		}
		wasOpen = open;
	});

	function setBrand(b: string) {
		brand = b;
		if (selectedPreset?.brand !== b) selectedIdx = undefined;
	}

	function selectPreset(preset: Filament) {
		selectedIdx = preset.trayInfoIdx;
	}

	async function save() {
		if (!printer?.serial || !selectedPreset) return;
		saving = true;
		try {
			await setFilament(printer.serial, {
				amsId,
				trayId,
				trayInfoIdx: selectedPreset.trayInfoIdx,
				trayColor: toBambuColor(color),
				trayType: selectedPreset.trayType,
				nozzleTempMin: selectedPreset.nozzleTempMin,
				nozzleTempMax: selectedPreset.nozzleTempMax
			});
			open = false;
		} finally {
			saving = false;
		}
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>Set Filament — {label}</Dialog.Title>
			<Dialog.Description>Configure the material for this slot.</Dialog.Description>
		</Dialog.Header>

		<div class="grid gap-4">
			<div class="grid grid-cols-2 gap-3">
				<div class="grid gap-2">
					<Label>Brand</Label>
					<DropdownMenu.Root>
						<DropdownMenu.Trigger
							class={cn(buttonVariants({ variant: 'outline' }), 'w-full justify-between')}
						>
							{brand}
							<ChevronDownIcon class="size-4 opacity-50" />
						</DropdownMenu.Trigger>
						<DropdownMenu.Content>
							{#each filaments.brands as b (b)}
								<DropdownMenu.Item onSelect={() => setBrand(b)}>{b}</DropdownMenu.Item>
							{/each}
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				</div>

				<div class="grid gap-2">
					<Label>Type</Label>
					<DropdownMenu.Root>
						<DropdownMenu.Trigger
							class={cn(buttonVariants({ variant: 'outline' }), 'w-full justify-between')}
						>
							<span class="truncate">{selectedPreset?.name ?? 'Select…'}</span>
							<ChevronDownIcon class="size-4 shrink-0 opacity-50" />
						</DropdownMenu.Trigger>
						<DropdownMenu.Content class="max-h-64 overflow-y-auto">
							{#each presetsForBrand as preset (preset.trayInfoIdx)}
								<DropdownMenu.Item onSelect={() => selectPreset(preset)}>
									{preset.name}
								</DropdownMenu.Item>
							{/each}
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				</div>
			</div>

			<div class="grid gap-2">
				<Label>Color</Label>
				<div class="flex flex-wrap gap-2">
					{#each FILAMENT_COLORS as swatch (swatch)}
						<button
							type="button"
							aria-label={swatch}
							onclick={() => (color = swatch)}
							style:background-color={swatch}
							class={cn(
								'size-7 rounded-full border border-foreground/15 transition',
								color.toUpperCase() === swatch
									? 'ring-2 ring-ring ring-offset-2 ring-offset-background'
									: 'hover:scale-110'
							)}
						></button>
					{/each}
				</div>
			</div>

			{#if selectedPreset}
				<p class="text-sm text-muted-foreground">
					Nozzle temperature: {selectedPreset.nozzleTempMin}–{selectedPreset.nozzleTempMax} °C
				</p>
			{/if}
		</div>

		<Dialog.Footer>
			<Button variant="outline" onclick={() => (open = false)}>Cancel</Button>
			<Button disabled={!selectedPreset || saving} onclick={save}>Save</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
