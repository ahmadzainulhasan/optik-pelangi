import {
	BETTER_AUTH_URL,
	BETTER_AUTH_SECRET,
	GOOGLE_CLIENT_ID,
	GOOGLE_CLIENT_SECRET
} from '$app/env/private';
import { betterAuth } from 'better-auth';
import { drizzleAdapter } from 'better-auth/adapters/drizzle';
import { admin } from 'better-auth/plugins/admin';
import { i18n, locales } from '@better-auth/i18n';
import { db } from '#lib/server/db/index.ts';
import * as schema from '#lib/server/db/schema.ts';
import { getRequestEvent } from '$app/server';
import { sveltekitCookies } from 'better-auth/svelte-kit';
import { logger } from '#lib/server/logger.ts';

export const auth = betterAuth({
	baseURL: BETTER_AUTH_URL,
	secret: BETTER_AUTH_SECRET,
	logger: {
		level: 'warn',
		log: (level, message, ...args) => logger[level](message, ...args)
	},
	rateLimit: {
		enabled: true,
		window: 60,
		max: 100
	},
	socialProviders: {
		google: {
			clientId: GOOGLE_CLIENT_ID,
			clientSecret: GOOGLE_CLIENT_SECRET,
			scope: ['email', 'profile']
		}
	},
	database: drizzleAdapter(db, {
		provider: 'pg',
		schema,
		usePlural: true
	}),
	session: {
		expiresIn: 60 * 60 * 24 * 7,
		updateAge: 60 * 60 * 24,
		freshAge: 0
	},
	account: {
		encryptOAuthTokens: true,
		accountLinking: {
			enabled: true,
			trustedProviders: ['google']
		}
	},
	advanced: {
		useSecureCookies: true,
		cookies: {
			session_token: {
				name: 'session'
			}
		}
	},
	plugins: [
		admin({
			bannedUserMessage:
				'Akun Anda telah ditangguhkan dari aplikasi ini. Silakan hubungi admin jika Anda yakin ini adalah kesalahan.'
		}),
		i18n({
			defaultLocale: 'id',
			translations: {
				id: locales.id
			}
		}),
		sveltekitCookies(getRequestEvent)
	],
	disabledPaths: [
		'/sign-up/email',
		'/sign-in/email',
		'/reset-password',
		'/verify-password',
		'/verify-email',
		'/send-verification-email',
		'/change-email',
		'/change-password',
		'/request-password-reset',
		'/reset-password/:token',
		'/admin/set-user-password',
		'/admin/impersonate-user',
		'/admin/stop-impersonating'
	]
});
