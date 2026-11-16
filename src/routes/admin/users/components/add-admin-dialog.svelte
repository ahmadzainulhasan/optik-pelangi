<script lang="ts">
	import { page } from '$app/state';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { zod4, zod4Client } from 'sveltekit-superforms/adapters';
	import { createAdminSchema } from '../schema.ts';
	import * as Dialog from '#lib/components/ui/dialog/index.js';
	import * as Field from '#lib/components/ui/field/index.js';
	import { Button, buttonVariants } from '#lib/components/ui/button/index.js';
	import { Input } from '#lib/components/ui/input/index.js';
	import { toast } from '#lib/components/ui/sonner/index.js';
	import UserPlus from '@lucide/svelte/icons/user-plus';
	import { cn } from '#lib/utils.js';

	let open = $state(false);
	const form = superForm(page.data.form ?? defaults(zod4(createAdminSchema)), {
		validators: zod4Client(createAdminSchema),
		resetForm: true,
		onUpdated({ form }) {
			if (form.valid && form.message) {
				toast.success(form.message);
				open = false;
			} else if (form.message) {
				toast.error(form.message);
			}
		}
	});

	const { form: formData, errors, constraints, enhance, submitting } = form;
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger class={cn(buttonVariants({ variant: 'default' }), 'animate-in fade-in slide-in-from-right-3 duration-400 ease-out')}>
		<UserPlus class="size-4" />
		<span>Tambah Admin</span>
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-md">
		<form id="create-admin-form" method="POST" action="/admin/users" use:enhance>
			<Dialog.Header>
				<Dialog.Title>Tambah Admin Baru</Dialog.Title>
				<Dialog.Description>
					Daftarkan staf baru sebagai admin. Akun akan otomatis terhubung saat mereka masuk
					menggunakan akun Google dengan email ini.
				</Dialog.Description>
			</Dialog.Header>

			<Field.Group class="my-4">
				<Field.Field data-invalid={Boolean($errors.name)}>
					<Field.Label for="admin-name">Nama Lengkap</Field.Label>
					{#if $errors.name}
						<Field.Error>{$errors.name}</Field.Error>
					{/if}
					<Input
						id="admin-name"
						name="name"
						placeholder="Contoh: Budi Santoso"
						bind:value={$formData.name}
						disabled={$submitting}
						{...$constraints.name}
					/>
				</Field.Field>
				<Field.Field data-invalid={Boolean($errors.email)}>
					<Field.Label for="admin-email">Email Google</Field.Label>
					{#if $errors.email}
						<Field.Error>{$errors.email}</Field.Error>
					{/if}
					<Input
						id="admin-email"
						name="email"
						type="email"
						placeholder="nama@gmail.com"
						bind:value={$formData.email}
						disabled={$submitting}
						{...$constraints.email}
					/>

					<Field.Description>
						Pastikan alamat email Google benar dan aktif untuk mendukung keamanan website.
					</Field.Description>
				</Field.Field>
			</Field.Group>

			<Dialog.Footer>
				<Dialog.Close type="button" class={buttonVariants({ variant: 'outline' })}>
					Batal
				</Dialog.Close>
				<Dialog.Root>
					<Dialog.Trigger
						class={buttonVariants({ variant: 'default' })}
						disabled={!$formData.name ||
							!$formData.email ||
							Boolean($errors.name) ||
							Boolean($errors.email)}
					>
						Tambah Admin
					</Dialog.Trigger>
					<Dialog.Content class="sm:max-w-md">
						<Dialog.Header>
							<Dialog.Title>Konfirmasi Penambahan Admin</Dialog.Title>
							<Dialog.Description>
								Mohon tinjau kembali data staf berikut sebelum menyimpan data admin.
							</Dialog.Description>
						</Dialog.Header>

						<div class="my-2 flex flex-col gap-2 text-sm">
							<div class="flex items-center justify-between">
								<span class="text-muted-foreground">Nama Lengkap</span>
								<span class="font-medium">{$formData.name}</span>
							</div>
							<div class="flex items-center justify-between">
								<span class="text-muted-foreground">Email Google</span>
								<span class="font-medium">{$formData.email}</span>
							</div>
						</div>

						<Dialog.Footer>
							<Dialog.Close class={buttonVariants({ variant: 'outline' })} disabled={$submitting}>
								Periksa Kembali
							</Dialog.Close>
							<Button type="submit" form="create-admin-form" disabled={$submitting}>
								{$submitting ? 'Menyimpan...' : 'Lanjutkan'}
							</Button>
						</Dialog.Footer>
					</Dialog.Content>
				</Dialog.Root>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
