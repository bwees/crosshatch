// Bambu reports tray colors as 8-char ARGB hex (e.g. "RRGGBBFF"); the CSS/UI
// side only ever wants the 6-char RGB portion.
function stripAlpha(color: string): string {
	return color.length === 8 ? color.slice(0, 6) : color;
}

export function bambuColorToCss(color?: string): string {
	if (!color) return 'transparent';
	return `#${stripAlpha(color)}`;
}

export function isLightColor(color?: string): boolean {
	if (!color) return true;
	const hex = stripAlpha(color);
	if (hex.length !== 6) return true;
	const r = parseInt(hex.slice(0, 2), 16);
	const g = parseInt(hex.slice(2, 4), 16);
	const b = parseInt(hex.slice(4, 6), 16);
	return r * 0.299 + g * 0.587 + b * 0.114 > 150;
}

export function toHexInput(color: string | undefined, fallback: string): string {
	if (!color) return fallback;
	return `#${stripAlpha(color).toUpperCase()}`;
}

export function toBambuColor(hex: string): string {
	return (hex.replace('#', '') + 'FF').toUpperCase();
}
