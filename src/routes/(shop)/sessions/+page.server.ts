import { redirect, error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { auth } from '#lib/server/auth.ts';

export const load: PageServerLoad = async ({ request, locals }) => {
	if (!locals.user) {
		redirect(303, '/');
	}

	return {
		sessions: await auth.api
			.listSessions({ headers: request.headers })
			.catch((err) => {
				throw error(err.statusCode, err.message);
			}),
		currentSession: locals.session
	};
};
