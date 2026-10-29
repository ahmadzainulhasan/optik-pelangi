import { API_ORIGIN } from '$app/env/public';
import type { PageLoad } from './$types';
import { goto } from '$app/navigation';
import type { WebResponse } from '#lib/types/api';
import type { LoginResponse } from '#lib/types/auth';
import type { User } from '#lib/types/user';

export const ssr = false;
export const csr = true;

export const load: PageLoad = async ({ fetch, url }) => {
	try {
		const response: Response = await fetch(`${API_ORIGIN}/auth/callback${url.search}`, {
			method: 'GET',
			credentials: 'include',
			headers: {
				Accept: 'application/json'
			}
		});
		if (response.ok) {
			const result: WebResponse<LoginResponse> = await response.json();
			const userResponse: Response = await fetch(`${API_ORIGIN}/users/current`, {
				method: 'GET',
				headers: {
					Authorization: `Bearer ${result.data.accessToken}`,
					Accept: 'application/json'
				}
			});
			const userResult: WebResponse<User> = await userResponse.json();

			if (userResult.data.role === 'admin') {
				goto('/admin');
			} else {
				goto('/');
			}
		}

		if (response.status >= 500) {
			const result: WebResponse<unknown> = await response.json();
			goto(`/auth/login?error=${encodeURIComponent(result.errors ?? '')}`);
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
