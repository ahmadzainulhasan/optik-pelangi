import type { PageLoad } from './$types';

export const load: PageLoad = () => {
	return {
		title: 'Syarat & Ketentuan',
		lastUpdated: '1 Januari 2026',
		sections: [
			{
				title: '1. Penerimaan Syarat',
				content: [
					'Dengan mengakses dan menggunakan situs web Optik Pelangi, Anda menyetujui untuk terikat oleh Syarat dan Ketentuan ini. Jika Anda tidak setuju dengan bagian mana pun dari syarat ini, Anda tidak diperkenankan menggunakan layanan kami.'
				]
			},
			{
				title: '2. Produk dan Layanan',
				content: [
					'Kami berupaya menampilkan warna dan gambar produk kami seakurat mungkin. Namun, kami tidak dapat menjamin bahwa tampilan monitor perangkat Anda akan sepenuhnya akurat.',
					'Semua deskripsi produk dan harga dapat berubah sewaktu-waktu tanpa pemberitahuan sebelumnya, atas kebijakan kami semata.'
				]
			},
			{
				title: '3. Pemeriksaan Refraksi & Lensa Resep',
				content: [
					'Layanan pemeriksaan mata (refraksi) disediakan sebagai panduan pembuatan kacamata di Optik Pelangi dan bukan merupakan diagnosis medis komprehensif untuk penyakit mata.',
					'Jika Anda menggunakan resep dari dokter mata eksternal, Anda bertanggung jawab penuh atas keakuratan resep yang diberikan kepada kami.'
				]
			},
			{
				title: '4. Pembayaran dan Pengembalian',
				content: [
					'Pesanan kacamata resep khusus (custom lenses) tidak dapat dibatalkan atau dikembalikan setelah proses pemotongan lensa dimulai.',
					'Klaim garansi kerusakan pabrik (frame/lensa) berlaku selama jangka waktu yang ditentukan pada nota pembelian fisik Anda.'
				]
			},
			{
				title: '5. Perubahan Syarat',
				content: [
					'Kami berhak untuk memperbarui, mengubah, atau mengganti bagian mana pun dari Syarat dan Ketentuan ini dengan memposting pembaruan atau perubahan di situs web kami.'
				]
			}
		]
	};
};

