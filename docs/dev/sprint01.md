# Sprint 01: Operational Core (Kasir Harian) 🛠️

**Fokus Utama**: Melengkapi fitur dasar kasir agar bisa digunakan untuk operasional harian yang wajar (bisa lihat riwayat, bisa hitung kembalian).

## 🎯 Objectives
1.  Kasir bisa melihat riwayat transaksi yang sudah terjadi (hari ini atau filter tanggal).
2.  Kasir bisa melihat detail item dari sebuah transaksi (re-print struk).
3.  Kasir tidak perlu menghitung kembalian manual.
4.  Sistem mencatat metode pembayaran (Cash vs QRIS).

## 📋 Task List

### 1. Feature: Transaction History (Riwayat)
*   [ ] **Endpoint List**: `GET /api/v1/transactions`
    *   Query Param: `?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD`
    *   Response: List of Transactions (ID, Time, Total Amount, Total Item).
*   [ ] **Endpoint Detail**: `GET /api/v1/transactions/{id}`
    *   Response: Transaction Header + Details (Items).

### 2. Feature: Payment Logic (Pembayaran)
*   [ ] **Update Model**: Tambah field `paid_amount`, `change`, `payment_method` di tabel `transactions`.
*   [ ] **Update Checkout API**:
    *   Validate `req.paid_amount >= total_amount`.
    *   Calculate `change = paid_amount - total_amount`.
    *   Save to DB.
*   [ ] **Update Response**: Return info `change` ke frontend.

### 3. Testing & Verification
*   [ ] Unit Test untuk logic hitung kembalian.
*   [ ] Manual Test flow checkout dengan uang pas dan uang lebih.

## ⏱️ Estimasi Waktu
*   Backend Implementation: ~2-3 Jam
*   Testing: ~1 Jam
