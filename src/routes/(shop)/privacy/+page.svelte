<script lang="ts">
	import { Separator } from '#lib/components/ui/separator/index.js';
</script>

<svelte:head>
	<title>Kebijakan Privasi — Optik Pelangi</title>
</svelte:head>

<article class="mx-auto max-w-3xl flex flex-col gap-10 py-8 animate-in fade-in slide-in-from-bottom-4 duration-500 ease-out">
	<h1 class="text-center font-serif text-3xl md:text-4xl lg:text-5xl">
		Kebijakan Privasi
	</h1>

	<section class="flex flex-col gap-3">
		<h2 class="font-serif text-xl font-medium md:text-2xl">1. Informasi yang Kami Kumpulkan</h2>
		<p class="leading-relaxed text-muted-foreground">
			Kami mengumpulkan informasi pribadi yang Anda berikan secara langsung kepada kami, seperti nama, alamat email, nomor telepon, dan data resep kacamata Anda saat Anda mendaftar akun atau melakukan pemesanan.
		</p>
	</section>
	<Separator />

	<section class="flex flex-col gap-3">
		<h2 class="font-serif text-xl font-medium md:text-2xl">2. Penggunaan Informasi</h2>
		<p class="leading-relaxed text-muted-foreground">
			Informasi yang kami kumpulkan digunakan untuk memproses pesanan Anda, membuat lensa sesuai resep, mengirimkan notifikasi status pesanan, dan memberikan dukungan pelanggan yang optimal.
		</p>
		<p class="leading-relaxed text-muted-foreground">
			Kami juga dapat menggunakan alamat email Anda untuk mengirimkan informasi promosi, namun Anda selalu memiliki opsi untuk berhenti berlangganan.
		</p>
	</section>
	<Separator />

	<section class="flex flex-col gap-3">
		<h2 class="font-serif text-xl font-medium md:text-2xl">3. Perlindungan Data (Data Security)</h2>
		<p class="leading-relaxed text-muted-foreground">
			Kami menerapkan langkah-langkah keamanan yang dirancang untuk melindungi informasi pribadi Anda dari akses, perubahan, pengungkapan, atau penghancuran yang tidak sah. Data medis (resep mata) disimpan secara rahasia.
		</p>
	</section>
	<Separator />

	<section class="flex flex-col gap-3">
		<h2 class="font-serif text-xl font-medium md:text-2xl">4. Berbagi dengan Pihak Ketiga</h2>
		<p class="leading-relaxed text-muted-foreground">
			Kami tidak menjual, memperdagangkan, atau menyewakan informasi identifikasi pribadi Anda kepada pihak lain. Kami hanya dapat membagikan data logistik dasar dengan layanan kurir tepercaya untuk keperluan pengiriman barang.
		</p>
	</section>
	<Separator />

	<section class="flex flex-col gap-3">
		<h2 class="font-serif text-xl font-medium md:text-2xl">5. Hak Pengguna</h2>
		<p class="leading-relaxed text-muted-foreground">
			Anda berhak untuk mengakses, memperbarui, atau menghapus informasi pribadi yang kami simpan dengan menghubungi kami melalui saluran komunikasi resmi Optik Pelangi.
		</p>
	</section>
</article>
