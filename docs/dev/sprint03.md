# Sprint 03: Inventory & Optimization (Manajemen Stok) 📦

**Fokus Utama**: Memudahkan manajemen stok (barang masuk) dan optimasi performa query.

## 🎯 Objectives
1.  Memudahkan proses tambah stok (restock) tanpa kalkulasi manual.
2.  Mencegah aplikasi lemot saat data produk menumpuk (Pagination).
3.  Memudahkan pencarian produk (Sorting & Advanced Filter).

## 📋 Task List

### 1. Feature: Restock Endpoint
*   [ ] **Endpoint**: `POST /api/v1/products/{id}/restock`
    *   Input: `qty_added` (int)
    *   Process: `current_stock + qty_added` (Atomic Update).
    *   Response: Stok terbaru.
*   [ ] **(Opsional) Inventory Log**: Tabel baru mencatat histori keluar masuk barang.

### 2. Feature: Pagination & Filter
*   [ ] **Update List Product API**:
    *   Support query param `?page=1&limit=10`.
    *   Support `?sort_by=price&order=asc`.
*   [ ] **Backend Logic**: Implement `LIMIT` dan `OFFSET` di query database.

### 3. Refactoring & Cleanup
*   [ ] Review ulang struktur kode.
*   [ ] Hapus fungsi-fungsi helper yang tidak terpakai.
*   [ ] Update Swagger Documentation agar sesuai dengan perubahan Sprint 1, 2, 3.

## ⏱️ Estimasi Waktu
*   Implementation: ~2-3 Jam
*   Testing: ~1 Jam
