import { auth } from '#lib/server/auth.ts';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, request }) => {
	return {
		user: await auth.api
			.getUser({
				query: { id: params.id },
				headers: request.headers
			})
			.catch((err) => {
				throw error(err.statusCode, err.message);
			})
	};
};
