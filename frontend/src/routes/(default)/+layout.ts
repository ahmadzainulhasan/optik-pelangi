import { API_ORIGIN } from '$app/env/public';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ fetch }) => {
	try {
		const response = await fetch(`${API_ORIGIN}/users/current`, {
			credentials: 'include'
		});
		if (!response.ok) {
			return { user: null };
		}
		const result = await response.json();
		return { user: result.data ?? null };
	} catch {
		return { user: null };
	}
};
