import { defineEnvVars } from '@sveltejs/kit/env';
import { building } from '$app/env';
import { z } from 'zod';

export const variables = defineEnvVars({
	API_ORIGIN: {
		public: true,
		schema: building ? z.string().optional() : z.string()
	}
});
