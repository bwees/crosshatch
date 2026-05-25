/// <reference types="@sveltejs/kit" />
/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />

import { build, files, version } from '$service-worker';

const sw = self as unknown as ServiceWorkerGlobalScope;

const CACHE = `cache-${version}`;
const ASSETS = [...build, ...files];

sw.addEventListener('install', (event) => {
	event.waitUntil(caches.open(CACHE).then((cache) => cache.addAll(ASSETS)));
});

sw.addEventListener('activate', (event) => {
	event.waitUntil(
		caches.keys().then(async (keys) => {
			await Promise.all(keys.filter((k) => k !== CACHE).map((k) => caches.delete(k)));
			await sw.clients.claim();
		})
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
