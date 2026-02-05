# Roadmap: Menuju Aplikasi Kasir Utuh 🏪

Saat ini aplikasi sudah memiliki fondasi kuat (Layered Arch, Database, Transaction). Namun untuk disebut "Aplikasi Kasir (POS) Siap Pakai", masih ada beberapa fitur krusial yang kurang.

Berikut adalah rencana pengembangan berdasarkan prioritas kebutuhan operasional:

## 1. Authentication & Authorization (Keamanan) 🔐
**Status**: 🔴 Belum Ada (Critical)
**Masalah**: Saat ini API "telanjang". Siapa saja bisa menghapus produk atau mengintip laporan omset jika tahu URL-nya.
**Rencana Implementasi**:
- **Tabel Users**: Menyimpan username, password (di-hash dgn bcrypt), dan role (`admin` vs `cashier`).
- **JWT (JSON Web Token)**: Mekanisme login standard untuk API.
- **Middleware**: Penjaga pintu untuk mengecek token di setiap request.
- **Role-Based Access**: 
  - `Admin`: Bisa CRUD Produk & Kategori.
  - `Cashier`: Hanya bisa Checkout & Lihat History.

## 2. Manajemen Riwayat Transaksi (Transaction History) 📜
**Status**: 🔴 Belum Ada (Critical)
**Masalah**: Data tersimpan di DB, tapi Kasir tidak bisa melihat "Tadi saya scan apa saja?" atau "Cetak ulang struk transaksi #123".
**Rencana Implementasi**:
- **Endpoint List**: `GET /api/transactions` (Melihat daftar struk).
- **Endpoint Detail**: `GET /api/transactions/{id}` (Melihat detail item dalam 1 struk).
- **Filter**: Filter by date range (misal: Transaksi tanggal 1-30).

## 3. Logika Pembayaran (Payment & Change) 💰
**Status**: 🟡 Parsial (Nice to Have)
**Masalah**: Kita hanya hitung `Total`. Belum ada fitur hitung "Uang Pembeli" dan "Kembalian".
**Rencana Implementasi**:
- **Update Checkout API**: Tambah input `paid_amount` (nominal uang user) dan `payment_method` (Cash/QRIS).
- **Backend Logic**: 
  - Validasi: `paid_amount` >= `total_amount`.
  - Hitung `change` (kembalian).
- **Response**: Mengembalikan info `change` ke frontend untuk ditampilkan.

## 4. Manajemen Stok Masuk (Restock) 📦
**Status**: 🟡 Manual (Nice to Have)
**Masalah**: Saat ini untuk tambah stok, user harus "Edit Produk" dan menimpa angka stok lama. Ini rawan salah hitung.
**Rencana Implementasi**:
- **Endpoint Restock**: `POST /api/products/{id}/restock`.
- **Logic**: Input `+10` akan otomatis `stok_lama + 10`. Tidak perlu hitung manual.
- **Log History**: (Opsional) Tabel `inventory_logs` untuk mencatat "Siapa yang nambah stok & kapan".

## 5. Pagination & Filtering 🔍
**Status**: 🟢 Optimization
**Masalah**: Jika produk mencapai 1.000 item, `GET /products` akan sangat lambat dan memakan kuota data.
**Rencana Implementasi**:
- **Pagination**: `?page=1&limit=10`.
- **Sorting**: `?sort=price_asc` (Termurah) atau `?sort=newest`.

---

## 🎯 Rekomendasi Fase Berikutnya (Session 4?)

Saya menyarankan kita membagi eksekusi menjadi beberapa tahap agar tidak overwhelming:

**Phase 1: Fitur Kasir (Operational Core)**
Fokus pada kelengkapan operasional kasir sehari-hari.
1. **Transaction History**: Agar bisa cetak struk/lihat riwayat.
2. **Payment Logic**: Agar aplikasi bisa hitung kembalian.

**Phase 2: Keamanan (Security Layer)**
Fokus mengamankan sistem sebelum dipakai banyak orang.
1. **Login & JWT**.
2. **Middleware**.

Bagaimana menurut Bos? Mau kita "masak" **Phase 1** dulu biar kasirnya makin canggih, atau **Phase 2** dulu biar aman?
