/// <reference types="@sveltejs/kit" />
/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />

import { build, files, version } from '$service-worker';

const sw = self as unknown as ServiceWorkerGlobalScope;

const CACHE = `cache-${version}`;
const ASSETS = [...build, ...files];

sw.addEventListener('install', (event) => {
	// Cache assets individually so one failed request (common in dev) doesn't
	// reject the whole install and block the worker from activating.
	event.waitUntil(
		(async () => {
			const cache = await caches.open(CACHE);
			await Promise.allSettled(ASSETS.map((asset) => cache.add(asset)));
			await sw.skipWaiting();
		})()
	);
});

sw.addEventListener('activate', (event) => {
	event.waitUntil(
		caches.keys().then(async (keys) => {
			await Promise.all(keys.filter((k) => k !== CACHE).map((k) => caches.delete(k)));
			await sw.clients.claim();
		})
	);
});

type PushPayload = {
	printerSerial: string;
	title: string;
	body: string;
	tag: string;
};

sw.addEventListener('push', (event) => {
	const payload = event.data?.json() as PushPayload | undefined;
	if (!payload) return;

	event.waitUntil(
		sw.registration.showNotification(payload.title, {
			body: payload.body,
			icon: '/icon-192.png',
			badge: '/icon-192.png',
			tag: payload.tag,
			data: { serial: payload.printerSerial }
		})
	);
});

sw.addEventListener('notificationclick', (event) => {
	event.notification.close();

	const serial = event.notification.data?.serial as string | undefined;
	const target = `/printer/${serial}`;

	event.waitUntil(
		(async () => {
			const clients = await sw.clients.matchAll({ type: 'window', includeUncontrolled: true });
			for (const client of clients) {
				if ('focus' in client) return client.focus();
			}
			return sw.clients.openWindow(target);
		})()
	);
});

sw.addEventListener('fetch', (event) => {
	const { request } = event;
	if (request.method !== 'GET') return;

	const url = new URL(request.url);
	if (url.origin !== sw.location.origin) return;
	if (url.pathname.startsWith('/api')) return;

	event.respondWith(
		(async () => {
			const cache = await caches.open(CACHE);

			if (ASSETS.includes(url.pathname)) {
				const cached = await cache.match(request);
				if (cached) return cached;
			}

			try {
				const response = await fetch(request);
				if (response.ok && response.type === 'basic') {
					cache.put(request, response.clone());
				}
				return response;
			} catch {
				const cached = await cache.match(request);
				if (cached) return cached;
				throw new Error('Network error and no cache available');
			}
		})()
	);
});
