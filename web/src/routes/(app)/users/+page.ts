import { authenticate } from '$lib/utils/auth';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	await authenticate({ admin: true });
};
