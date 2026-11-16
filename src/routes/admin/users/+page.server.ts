import { error, fail } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';
import { auth } from '#lib/server/auth.ts';
import { superValidate, message, type ErrorStatus } from 'sveltekit-superforms';
import { zod4 } from 'sveltekit-superforms/adapters';
import { createAdminSchema } from './schema.ts';
import { isAPIError } from 'better-auth/api';

export const load: PageServerLoad = async ({ request, url }) => {
	const query: NonNullable<Parameters<typeof auth.api.listUsers>[0]>['query'] = {
		limit: 10,
		offset: ((Number(url.searchParams.get('page')) || 1) - 1) * 10,
		sortBy: 'createdAt',
		sortDirection: 'desc'
	};

	const sort = url.searchParams.get('sort');
	if (sort) {
		const [sortBy, sortDirection] = sort.split('-');
		query.sortBy = sortBy === 'name' ? 'name' : 'createdAt';
		query.sortDirection = sortDirection === 'asc' ? 'asc' : 'desc';
	}

	const search = url.searchParams.get('search');
	if (search) {
		query.searchValue = search;
		query.searchField = search.includes('@') ? 'email' : 'name';
		query.searchOperator = 'contains';
	}

	const role = url.searchParams.get('role');
	if (role === 'admin' || role === 'user') {
		query.filterField = 'role';
		query.filterOperator = 'eq';
		query.filterValue = role;
	}

	const users = await auth.api
		.listUsers({
			query,
			headers: request.headers
		})
		.catch((err) => {
			throw error(err.statusCode, err.message);
		});

	if (users.total === 0) {
		throw error(404, 'Tidak ada data pengguna yang ditemukan');
	}

	return {
		users,
		form: await superValidate(zod4(createAdminSchema))
	};
};

export const actions: Actions = {
	default: async ({ request }) => {
		const form = await superValidate(request, zod4(createAdminSchema));
		if (!form.valid) {
			return fail(400, { form });
		}

		try {
			await auth.api.createUser({
				body: {
					name: form.data.name,
					email: form.data.email,
					role: 'admin',
					data: {
						emailVerified: true
					}
				},
				headers: request.headers
			});
		} catch (err) {
			if (isAPIError(err)) {
				return message(form, err.message, {
					status: err.statusCode as ErrorStatus
				});
			}

			return message(form, 'Gagal menambahkan admin baru.', { status: 500 });
		}

		return message(form, 'Admin baru berhasil ditambahkan.');
	}
};
