// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}

	// Not yet in lib.dom; supported on iOS Safari for MSE playback.
	interface ManagedMediaSourceCtor {
		new (): MediaSource;
		isTypeSupported(type: string): boolean;
	}

	interface Window {
		ManagedMediaSource?: ManagedMediaSourceCtor;
	}
}

export {};
