import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals }) => {
	const user = locals.user;
	if (user?.role === 'admin') {
		redirect(303, '/admin');
	}

	return {
		currentUser: user ?? null
	};
};
