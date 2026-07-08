<script lang="ts">
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import Label from '$lib/components/ui/label/label.svelte';
	import { FILAMENT_COLORS } from '$lib/filaments';
	import { filamentManager } from '$lib/filaments.svelte';
	import { setFilament, type Filament, type Printer, type PrinterStatus } from '$lib/sdk';
	import { Spinner } from '$lib/components/ui/spinner';
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

	// User edits override the defaults derived from the tray; they win once set and
	// are cleared when the dialog closes so the next open starts fresh.
	let brandOverride = $state<string | undefined>(undefined);
	let idxOverride = $state<string | undefined>(undefined);
	let colorOverride = $state<string | undefined>(undefined);
	let saving = $state(false);

	onMount(() => filamentManager.init());

	function toHexInput(c?: string): string {
		if (!c) return FILAMENT_COLORS[0];
		const hex = c.length >= 6 ? c.slice(0, 6) : c;
		return `#${hex.toUpperCase()}`;
	}

	function toBambuColor(hex: string): string {
		return (hex.replace('#', '') + 'FF').toUpperCase();
	}

	// Prefer an exact match on the reported filament id, falling back to the
	// material family when the loaded spool isn't a known preset.
	const defaultMatch = $derived.by(() => {
		const exact = filamentManager.presets.find((p) => p.trayInfoIdx === tray?.trayInfoIdx);
		return (
			exact ??
			filamentManager.presets.find(
				(p) =>
					p.brand === (tray?.brand?.toLowerCase().includes('bambu') ? 'Bambu' : 'Generic') &&
					p.trayType === tray?.material
			)
		);
	});

	const brand = $derived(brandOverride ?? defaultMatch?.brand ?? 'Generic');
	const presetsForBrand = $derived(filamentManager.presets.filter((p) => p.brand === brand));
	const selectedIdx = $derived(idxOverride ?? defaultMatch?.trayInfoIdx);
	const selectedPreset = $derived(presetsForBrand.find((p) => p.trayInfoIdx === selectedIdx));
	const color = $derived(colorOverride ?? toHexInput(tray?.color));

	$effect(() => {
		if (!open) {
			brandOverride = undefined;
			idxOverride = undefined;
			colorOverride = undefined;
		}
	});

	function setBrand(b: string) {
		brandOverride = b;
	}

	function selectPreset(preset: Filament) {
		idxOverride = preset.trayInfoIdx;
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

		{#if filamentManager.loading}
			<div class="flex items-center justify-center py-12 text-muted-foreground">
				<Spinner class="size-6" />
			</div>
		{:else}
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
								{#each filamentManager.brands as b (b)}
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
								onclick={() => (colorOverride = swatch)}
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
			</div>
		{/if}

		<Dialog.Footer>
			<Button variant="outline" onclick={() => (open = false)}>Cancel</Button>
			<Button disabled={!selectedPreset || saving} onclick={save}>Save</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
