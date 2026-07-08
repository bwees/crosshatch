import { getFilaments, type Filament } from '$lib/sdk';
import { SvelteSet } from 'svelte/reactivity';

class FilamentManager {
	presets: Filament[] = $state([]);
	brands = $derived(new SvelteSet(this.presets.map((p) => p.brand)));

	loading = $state(false);

	async init() {
		if (this.loading) return;
		this.loading = true;

		try {
			this.presets = await getFilaments();
		} finally {
			this.loading = false;
		}
	}
}

export const filamentManager = new FilamentManager();
