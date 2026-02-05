# Developer Collaboration Contract 🤝

Dokumen ini adalah kesepakatan teknis dan metode kerja antara **Owner/Developer** (CodeWithUmam) dan **AI Assistant** (Antigravity).

Tujuan utama proyek ini bukan hanya membuat aplikasi yang jalan, tapi juga sebagai **Media Belajar**.

---

## 1. Technology Stack 🛠️

Kita sepakat menggunakan stack berikut untuk menjaga kesederhanaan namun tetap robust:

*   **Language**: Go (Golang) versi terbaru.
*   **Database**: SQLite (via `github.com/mattn/go-sqlite3`). Ringan, tanpa setup server ribet.
*   **Architecture**: Layered Architecture (Handler -> Service -> Repository -> Model).
*   **Configuration**: Viper (Manage `.env` dan config files).
*   **Documentation**: Swagger / OpenAPI (via `swaggo`).
*   **Deployment**: Linux VPS (Systemd service).

---

## 2. Development Principles 📜

### A. TDD First (Test Driven Development)
*   **Aturan**: Sebelum atau bersamaan dengan menulis kode fitur baru, WAJIB menyertakan **Unit Test**.
*   **Goal**: Memastikan setiap fungsi berjalan benar (Business Logic aman) dan refactoring di masa depan tidak merusak fitur lama.

### B. Educational Code (Komentar Edukasi)
*   **Aturan**: Setiap kode yang "njelimet" atau mengandung konsep baru WAJIB diberi komentar `//`.
*   **Format**: Komentar harus menjelaskan **"Kenapa (Why)"** dan **"Bagaimana (How)"**, bukan sekadar translate kode.
*   **Bahasa**: Bahasa Indonesia yang santai dan mudah dimengerti.
*   **Contoh**:
    ```go
    // Kita pakai Transaction (tx) di sini karena...
    // Jangan lupa defer Rollback() untuk safety net jika...
    ```

### C. "Slow Aje" Workflow
*   **Plan First**: Jangan langsung coding. Buat plan/roadmap dulu.
*   **Review**: Minta persetujuan user sebelum eksekusi besar.
*   **Step-by-Step**: Jangan rapel semua fitur dalam satu commit/task. Pecah jadi potongan kecil yang bisa dicerna.

---

## 3. Sprint & Roadmap 🗂️

Kita akan mengerjakan fitur berdasarkan Sprint yang sudah didefinisikan:

*   **[Sprint 01: Operational Core](sprint01.md)**
    *   Fokus: History Transaksi & Logika Pembayaran (Kembalian).
*   **[Sprint 02: Security Layer](sprint02.md)**
    *   Fokus: Login, JWT, dan Role Admin/Kasir.
*   **[Sprint 03: Inventory & Optimization](sprint03.md)**
    *   Fokus: Restock Barang & Pagination.

---

## 4. Tanda Tangan Digital ✍️

Disepakati pada: **Februari 2026**

**Owner**
*(User)*

**AI Assistant**
*(Antigravity)*
