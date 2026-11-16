import type { PageServerLoad } from './$types';
import { auth } from '#lib/server/auth.ts';
import { error } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ request, locals }) => {
	return {
		sessions: await auth.api
			.listSessions({ headers: request.headers })
			.catch((err) => {
				throw error(err.statusCode, err.message);
			}),
		currentSession: locals.session
	};
};
