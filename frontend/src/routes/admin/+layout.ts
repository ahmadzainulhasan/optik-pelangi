import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ parent }) => {
	const { user } = await parent();

	if (!user) {
		redirect(307, '/auth/login');
	}

	if (user.role !== 'admin') {
		redirect(307, '/');
	}

	return { user };
};
