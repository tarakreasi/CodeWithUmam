# Sprint 02: Security Layer (Keamanan & Akses) 🔐

**Fokus Utama**: Mengamankan API dan membatasi akses berdasarkan peran (Admin vs Kasir).

## 🎯 Objectives
1.  Hanya user terdaftar yang bisa login.
2.  Tidak sembarang orang bisa akses API (harus bawa Token).
3.  Pemisahan hak akses: Kasir tidak boleh hapus produk/kategori.

## 📋 Task List

### 1. Feature: User Management (Auth)
*   [ ] **Database**: Buat tabel `users` (id, username, password_hash, role).
*   [ ] **Seed Data**: Buat 1 user Admin default.
*   [ ] **Library**: Install `bcrypt` (hashing) dan `golang-jwt` (token).

### 2. Feature: Login System
*   [ ] **Endpoint**: `POST /api/v1/login`
    *   Input: `username`, `password`
    *   Process: Cek DB, Verify Hash, Generate JWT.
    *   Response: `token`, `expired_at`.

### 3. Feature: Middleware (Satpam)
*   [ ] **Auth Middleware**:
    *   Intercept setiap request.
    *   Validasi Header `Authorization: Bearer <token>`.
    *   Deny 401 jika token invalid.
*   [ ] **Apply Middleware**: Pasang di router (kecuali endpoint Login & Public Health).

### 4. Feature: Role Based Access Control (RBAC)
*   [ ] **Admin Only Middleware**:
    *   Hanya role `admin` yang boleh akses method POST/PUT/DELETE di `/products` dan `/categories`.
*   [ ] **Cashier Access**:
    *   Kasir bebas akses `/checkout`, `/products` (GET), `/transactions`.

## ⏱️ Estimasi Waktu
*   Implementation: ~4-5 Jam (Cukup kompleks di setup awal)
*   Testing: ~1-2 Jam
