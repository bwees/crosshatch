import { getFilaments, type Filament } from '$lib/sdk';

let presets = $state<Filament[]>([]);
let started = false;

// ensureFilaments fetches the catalog once and caches it for every dialog.
export async function ensureFilaments() {
	if (started) return;
	started = true;
	try {
		presets = await getFilaments();
	} catch (e) {
		started = false;
		throw e;
	}
}

export const filaments = {
	get presets() {
		return presets;
	},
	get brands() {
		return [...new Set(presets.map((p) => p.brand))];
	}
};
