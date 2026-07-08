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
	class={['lucide-icon lucide lucide-printer-build-plate', className]}
	{...rest}
>
	<path d="M4 19h16" />
	<path d="M8 17q2-2 0-4t0-4t0-4" />
	<path d="M12 17q2-2 0-4t0-4t0-4" />
	<path d="M16 17q2-2 0-4t0-4t0-4" />
</svg>
