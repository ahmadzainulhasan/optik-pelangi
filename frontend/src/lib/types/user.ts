export type UserRole = 'customer' | 'admin';

export type User = {
	id: string;
	name: string;
	email: string;
	picture: string;
	role: UserRole;
};



