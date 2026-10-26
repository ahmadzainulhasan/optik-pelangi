import { API_ORIGIN } from '$app/env/public';
import type { PageLoad } from './$types';
import { goto } from '$app/navigation';

export const ssr = false;
export const csr = true;

export const load: PageLoad = async ({ fetch, url }) => {
	try {
		const response = await fetch(`${API_ORIGIN}/auth/callback${url.search}`, {
			method: 'GET',
			credentials: 'include',
			headers: {
				Accept: 'application/json'
			}
		});
		if (response.ok) {
			goto('/');
		}

		if (response.status >= 500) {
			const result = await response.json();
			goto(`/auth/login?error=${encodeURIComponent(result.errors)}`);
		}

		if (response.status >= 400) {
			goto('/auth/login');
		}
	} catch (err) {
		if (err instanceof Error) {
			goto(
				`/auth/login?error=${encodeURIComponent('Tidak dapat terhubung ke internet. Periksa koneksi Anda.')}`
			);
		}
	}
};
