type SyncedControlOptions<T> = {
	/** Value shown before the device has reported anything. */
	initial: T;
	/** Reads the latest value reported by the device, or undefined if unknown. */
	reported: () => T | undefined;
	/** Pushes a user-chosen value to the device. */
	apply: (value: T) => void;
	/** Equality used to detect when the device confirms a pending change. */
	equals?: (a: T, b: T) => boolean;
};

export class SyncedControl<T> {
	current = $state<T>() as T;
	#pending: { value: T } | null = null;
	#apply: (value: T) => void;

	constructor(options: SyncedControlOptions<T>) {
		const equals = options.equals ?? ((a, b) => a === b);
		this.#apply = options.apply;
		this.current = options.initial;

		$effect(() => {
			const reported = options.reported();
			if (reported === undefined) return;
			if (this.#pending !== null) {
				if (equals(reported, this.#pending.value)) this.#pending = null;
				return;
			}
			this.current = reported;
		});
	}

	set = (value: T) => {
		this.#pending = { value };
		this.current = value;
		this.#apply(value);
	};
}
