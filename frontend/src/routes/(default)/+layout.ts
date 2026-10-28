import { API_ORIGIN } from '$app/env/public';
import type { LayoutLoad } from './$types';
import type { User } from '$lib/types/user';

export const ssr = false;

export const load: LayoutLoad = async ({ fetch }) => {
	try {
		const refreshResponse = await fetch(`${API_ORIGIN}/auth/refresh`, {
			method: 'POST',
			credentials: 'include',
			headers: {
				Accept: 'application/json'
			}
		});

		if (!refreshResponse.ok) {
			return { user: null };
		}

		const refreshResult = await refreshResponse.json();
		const accessToken = refreshResult.data?.accessToken;
		if (!accessToken) {
			return { user: null };
		}

		const userResponse = await fetch(`${API_ORIGIN}/users/current`, {
			method: 'GET',
			headers: {
				Authorization: `Bearer ${accessToken}`,
				Accept: 'application/json'
			}
		});

		if (!userResponse.ok) {
			return { user: null };
		}

		const userResult = await userResponse.json();
		return { user: (userResult.data as User) ?? null };
	} catch {
		return { user: null };
	}
};
