const DEVICE_ID_KEY = 'crosshatch_device_id';

// getDeviceId returns a stable per-device identifier persisted in localStorage.
// It survives login/logout so a user's notification settings stay tied to the
// device rather than the account.
export function getDeviceId(): string {
	let id = localStorage.getItem(DEVICE_ID_KEY);
	if (!id) {
		id = crypto.randomUUID();
		localStorage.setItem(DEVICE_ID_KEY, id);
	}
	return id;
}
