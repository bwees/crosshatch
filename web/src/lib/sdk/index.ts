import * as client from './client';

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL ?? 'http://brandon-macbook-pro:5173';

client.defaults.baseUrl = apiBaseUrl;

export * from './client';
export { apiBaseUrl };
