import { createAuthClient } from 'better-auth/svelte';
import { adminClient, inferAdditionalFields } from 'better-auth/client/plugins';
import { i18nClient } from '@better-auth/i18n/client';
import type { auth } from '#lib/server/auth.ts';

export const authClient = createAuthClient({
	plugins: [adminClient(), i18nClient(), inferAdditionalFields<typeof auth>()]
});

export type Session = typeof authClient.$Infer.Session;
export type User = typeof authClient.$Infer.Session.user;
