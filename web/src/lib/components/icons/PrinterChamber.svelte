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
	class={['lucide-icon lucide lucide-printer-chamber', className]}
	{...rest}
>
	<rect x="3" y="3" width="18" height="18" rx="2" />
	<path d="M3 8h18" />
	<path d="M12 8v5" />
	<path d="M7 17h10" />
</svg>
