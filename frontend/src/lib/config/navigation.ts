export type NavCategory = {
	title: string;
	href: string;
	description: string;
};

export type NavItem = {
	title: string;
	href: string;
	items: NavCategory[];
};

export const mainNav: NavItem[] = [
	{
		title: 'Beranda',
		href: '/',
		items: [
			{
				title: 'Kategori Produk',
				href: '/#categories',
				description: 'Telusuri koleksi kacamata berdasarkan kategori kebutuhan.'
			},
			{
				title: 'Model Terlaris',
				href: '/#featured',
				description: 'Pilihan frame paling populer dan diminati di toko kami.'
			},
			{
				title: 'Layanan Optik',
				href: '/#services',
				description: 'Pemeriksaan refraksi mata komprehensif dan perakitan lensa faset presisi.'
			},
			{
				title: 'Pengalaman Pelanggan',
				href: '/#testimonials',
				description: 'Testimoni nyata dari pelanggan setia Optik Pelangi.'
			}
		]
	},
	{
		title: 'Produk',
		href: '/products',
		items: [
			{
				title: 'Kacamata Minus',
				href: '/products?category=minus',
				description: 'Koleksi frame dan lensa khusus untuk koreksi mata miopi jarak jauh.'
			},
			{
				title: 'Kacamata Baca',
				href: '/products?category=baca',
				description: 'Solusi presbiopi untuk jarak dekat yang nyaman dan stylish.'
			},
			{
				title: 'Kacamata Hitam',
				href: '/products?category=sunglasses',
				description: 'Perlindungan UV maksimal dengan gaya modern.'
			},
			{
				title: 'Frame Kacamata',
				href: '/products?category=frame',
				description: 'Berbagai pilihan bingkai berkualitas.'
			}
		]
	},
	{
		title: 'Tentang',
		href: '/about',
		items: [
			{
				title: 'Kisah Optik Pelangi',
				href: '/about#tentang',
				description: 'Sejarah panjang kami melayani masyarakat Kabupaten Pati sejak 2009.'
			},
			{
				title: 'Layanan Kami',
				href: '/about#layanan',
				description: 'Fasilitas dan pelayanan lengkap yang kami sediakan untuk Anda.'
			},
			{
				title: 'Lokasi & Jam Buka',
				href: '/about#lokasi',
				description: 'Petunjuk arah menuju toko fisik kami di Jl. Kiai Saleh No. 38.'
			},
			{
				title: 'Jejak Langkah',
				href: '/about#perjalanan',
				description: 'Tonggak penting pencapaian Optik Pelangi dari tahun ke tahun.'
			}
		]
	},
	{
		title: 'Bantuan',
		href: '/faq',
		items: [
			{
				title: 'Produk & Pembelian',
				href: '/faq#produk-pembelian',
				description: 'Pertanyaan umum mengenai ketersediaan, jenis lensa, dan harga.'
			},
			{
				title: 'Pemeriksaan & Layanan',
				href: '/faq#layanan-toko-pemeriksaan',
				description: 'Info tentang prosedur cek mata dan pemotongan faset lensa.'
			},
			{
				title: 'Garansi & Servis',
				href: '/faq#garansi-servis',
				description: 'Kebijakan retur, garansi produk, dan layanan cuci serta penyetelan.'
			},
			{
				title: 'Konsultasi Online',
				href: 'https://wa.me/6281249902200',
				description: 'Hubungi tim optisi kami secara langsung via obrolan WhatsApp.'
			}
		]
	}
];
