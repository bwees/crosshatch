<script lang="ts">
	import MenuIcon from '@lucide/svelte/icons/menu';
	import { Slider as SliderPrimitive } from 'bits-ui';

	type Props = {
		labels: string[];
		value?: number;
		disabled?: boolean;
		class?: string;
		onValueCommit?: (value: number) => void;
	};

	let {
		labels,
		value = $bindable(0),
		disabled = false,
		class: className,
		onValueCommit
	}: Props = $props();
</script>

<SliderPrimitive.Root
	type="single"
	bind:value
	min={0}
	max={labels.length - 1}
	step={1}
	trackPadding={8}
	{onValueCommit}
	{disabled}
	class={[
		'relative mb-8 flex w-full touch-none items-center pt-1 pb-1 select-none data-disabled:opacity-50',
		className
	]}
>
	{#snippet children({ thumbItems, tickItems })}
		<span class="relative h-1.5 w-full grow overflow-hidden rounded-full bg-muted">
			<SliderPrimitive.Range class="absolute h-full bg-muted" />
		</span>

		{#each tickItems as tick (tick.index)}
			<SliderPrimitive.Tick
				index={tick.index}
				class="block size-2 -translate-x-1/2 rounded-full bg-muted-foreground/50 data-[selected]:opacity-0"
			/>
			<SliderPrimitive.TickLabel
				index={tick.index}
				position="bottom"
				class="-translate-x-1/2 pt-3 text-sm font-medium text-muted-foreground data-[selected]:text-primary"
			>
				{labels[tick.index]}
			</SliderPrimitive.TickLabel>
		{/each}

		{#each thumbItems as thumb (thumb.index)}
			<SliderPrimitive.Thumb
				index={thumb.index}
				class="flex size-8 shrink-0 items-center justify-center rounded-full border bg-muted text-muted-foreground shadow-sm ring-ring/50 transition-[color,box-shadow] select-none hover:ring-4 focus-visible:ring-4 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
			>
				<MenuIcon class="size-4" />
			</SliderPrimitive.Thumb>
		{/each}
	{/snippet}
</SliderPrimitive.Root>
