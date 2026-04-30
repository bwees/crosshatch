<script lang="ts">
	import MenuIcon from '@lucide/svelte/icons/menu';
	import { Slider as SliderPrimitive } from 'bits-ui';

	type Props = {
		labels: string[];
		value?: number;
		disabled?: boolean;
		class?: string;
	};

	let {
		labels,
		value = $bindable(0),
		disabled = false,
		class: className
	}: Props = $props();
</script>

<SliderPrimitive.Root
	type="single"
	bind:value
	min={0}
	max={labels.length - 1}
	step={1}
	trackPadding={8}
	{disabled}
	class={[
		'relative flex w-full touch-none items-center pt-1 pb-2 select-none data-disabled:opacity-50',
		className
	]}
>
	{#snippet children({ thumbItems, tickItems })}
		<span class="bg-muted relative h-1.5 w-full grow overflow-hidden rounded-full">
			<SliderPrimitive.Range class="bg-muted absolute h-full" />
		</span>

		{#each tickItems as tick (tick.index)}
			<SliderPrimitive.Tick
				index={tick.index}
				class="bg-muted-foreground/50 block size-2 -translate-x-1/2 rounded-full data-[selected]:opacity-0"
			/>
			<SliderPrimitive.TickLabel
				index={tick.index}
				position="bottom"
				class="text-muted-foreground data-[selected]:text-primary -translate-x-1/2 pt-3 text-sm font-medium"
			>
				{labels[tick.index]}
			</SliderPrimitive.TickLabel>
		{/each}

		{#each thumbItems as thumb (thumb.index)}
			<SliderPrimitive.Thumb
				index={thumb.index}
				class="bg-muted text-muted-foreground ring-ring/50 flex size-8 shrink-0 items-center justify-center rounded-full border shadow-sm transition-[color,box-shadow] select-none hover:ring-4 focus-visible:ring-4 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
			>
				<MenuIcon class="size-4" />
			</SliderPrimitive.Thumb>
		{/each}
	{/snippet}
</SliderPrimitive.Root>
