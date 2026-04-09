import * as client from './client';

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:3000';

client.defaults.baseUrl = apiBaseUrl;

export * from './client';
export { apiBaseUrl };
