import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals }) => {
	const user = locals.user;

	if (!user || user.role !== 'admin' || user.banned) {
		redirect(303, '/');
	}

	return {
		currentUser: user
	};
};
