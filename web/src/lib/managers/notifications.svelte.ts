import { getDeviceId } from '$lib/device';
import { getVapidConfig, subscribeNotifications, unsubscribeNotifications } from '$lib/sdk';

class NotificationsManager {
	permission = $state<NotificationPermission>(
		typeof Notification !== 'undefined' ? Notification.permission : 'default'
	);

	get supported(): boolean {
		return (
			typeof navigator !== 'undefined' &&
			'serviceWorker' in navigator &&
			typeof window !== 'undefined' &&
			'PushManager' in window &&
			'Notification' in window
		);
	}

	async isSubscribed(): Promise<boolean> {
		if (!this.supported) return false;
		const registration = await navigator.serviceWorker.ready;
		const subscription = await registration.pushManager.getSubscription();
		return subscription !== null;
	}

	async enable(): Promise<boolean> {
		if (!this.supported) return false;

		this.permission = await Notification.requestPermission();
		if (this.permission !== 'granted') return false;

		const registration = await navigator.serviceWorker.ready;
		const { publicKey } = await getVapidConfig();

		// A subscription created with a different VAPID key must be dropped first,
		// otherwise subscribe() throws InvalidStateError (e.g. after the server
		// regenerates its keys).
		const existing = await registration.pushManager.getSubscription();
		if (existing) await existing.unsubscribe();

		const subscription = await this.withTimeout(
			registration.pushManager.subscribe({
				userVisibleOnly: true,
				applicationServerKey: this.urlBase64ToUint8Array(publicKey)
			}),
			15000,
			'Push subscription timed out — the browser could not reach its push service.'
		);

		const { endpoint } = subscription;
		const keys = subscription.toJSON().keys;
		if (!keys?.p256dh || !keys?.auth) return false;

		await subscribeNotifications({
			deviceId: getDeviceId(),
			endpoint,
			p256dh: keys.p256dh,
			auth: keys.auth
		});
		return true;
	}

	async disable(): Promise<void> {
		if (!this.supported) return;

		await unsubscribeNotifications({ deviceId: getDeviceId() });

		const registration = await navigator.serviceWorker.ready;
		const subscription = await registration.pushManager.getSubscription();
		if (subscription) await subscription.unsubscribe();
	}

	private withTimeout<T>(promise: Promise<T>, ms: number, message: string): Promise<T> {
		return Promise.race([
			promise,
			new Promise<T>((_, reject) => setTimeout(() => reject(new Error(message)), ms))
		]);
	}

	private urlBase64ToUint8Array(base64String: string): Uint8Array<ArrayBuffer> {
		const padding = '='.repeat((4 - (base64String.length % 4)) % 4);
		const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/');

		const rawData = atob(base64);
		const output = new Uint8Array(new ArrayBuffer(rawData.length));
		for (let i = 0; i < rawData.length; i++) {
			output[i] = rawData.charCodeAt(i);
		}
		return output;
	}
}

export const notifications = new NotificationsManager();
