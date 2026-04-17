<script lang="ts">
	import type { SVGAttributes } from 'svelte/elements';

	type Props = SVGAttributes<SVGSVGElement> & {
		size?: number | string;
		color?: string;
		strokeWidth?: number | string;
	};

	let {
		size = 24,
		color = 'currentColor',
		strokeWidth = 2,
		class: className,
		...rest
	}: Props = $props();
</script>

<svg
	xmlns="http://www.w3.org/2000/svg"
	width={size}
	height={size}
	viewBox="0 0 24 24"
	fill="none"
	stroke={color}
	stroke-width={strokeWidth}
	stroke-linecap="round"
	stroke-linejoin="round"
	class={['lucide-icon lucide lucide-printer-nozzle', className]}
	{...rest}
>
	<rect x="5" y="3" width="14" height="4" rx="0.5" />
	<path d="M5 7L12 17L19 7" />
	<path d="M5 13q1 1 0 2t0 2" />
	<path d="M19 13q1 1 0 2t0 2" />
</svg>
