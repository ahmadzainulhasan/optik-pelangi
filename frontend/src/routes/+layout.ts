import { API_ORIGIN } from '$app/env/public';
import type { LayoutLoad } from './$types';
import type { WebResponse } from '#lib/types/api';
import type { LoginResponse } from '#lib/types/auth';
import type { User } from '#lib/types/user';

export const ssr = false;

export const load: LayoutLoad = async ({ fetch }) => {
	try {
		const refreshResponse: Response = await fetch(`${API_ORIGIN}/auth/refresh`, {
			method: 'POST',
			credentials: 'include',
			headers: {
				Accept: 'application/json'
			}
		});

		if (!refreshResponse.ok) {
			return { user: null };
		}

		const refreshResult: WebResponse<LoginResponse> = await refreshResponse.json();
		const userResponse: Response = await fetch(`${API_ORIGIN}/users/current`, {
			method: 'GET',
			headers: {
				Authorization: `Bearer ${refreshResult.data.accessToken}`,
				Accept: 'application/json'
			}
		});

		if (!userResponse.ok) {
			return { user: null };
		}

		const userResult: WebResponse<User> = await userResponse.json();
		return { user: userResult.data };
	} catch {
		return { user: null };
	}
};
